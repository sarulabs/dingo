package generation

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// NewLocator creates an initialized Locator.
func NewLocator(src, dest, destPkg string) (*Locator, error) {
	loc := &Locator{}
	return loc, loc.Init(src, dest, destPkg)
}

// Locator provides an easy way to get information
// about the source and destination directories.
type Locator struct {
	src        string
	srcFiles   []string
	srcDirName string
	dest       string
	destPkg    string
}

// SourceDir returns the absolute path to the source directory.
// Init should be called first.
func (l *Locator) SourceDir() string {
	return l.src
}

// SourceFiles returns the name of the go files
// contained in the source directory.
// Init should be called first.
func (l *Locator) SourceFiles() []string {
	return l.srcFiles
}

// SourceDirName returns the name of the source directory.
// Init should be called first.
func (l *Locator) SourceDirName() string {
	return l.srcDirName
}

// DestDir returns the absolute path to the destination directory.
// Init should be called first.
func (l *Locator) DestDir() string {
	return l.dest
}

// DestPkg returns the name of the package
// for the files in the destination directory.
// Init should be called first.
func (l *Locator) DestPkg() string {
	return l.destPkg
}

// Init initialize the Locator instance
// with a given source and destination directories.
// dest must be in the $GOPATH.
func (l *Locator) Init(src, dest, destPkg string) error {
	if err := l.setSource(src); err != nil {
		return err
	}
	if err := l.setSourceFiles(); err != nil {
		return err
	}
	if err := l.setDest(dest); err != nil {
		return err
	}
	return l.setDestPkg(destPkg)
}

func (l *Locator) setSource(src string) error {
	dir, err := filepath.Abs(src)
	if err != nil {
		return errors.New("could not get absolute path of source (" + src + ")")
	}

	s, err := os.Stat(dir)
	if err != nil {
		return errors.New("could not read source stat (" + dir + ") :" + err.Error())
	}

	if !s.IsDir() {
		return errors.New("source should be a directory (" + dir + ")")
	}

	l.src = dir
	l.srcDirName = filepath.Base(dir)

	return nil
}

func (l *Locator) setSourceFiles() error {
	goFiles := []string{}

	files, err := ioutil.ReadDir(l.src)
	if err != nil {
		return errors.New("could not read source files: " + err.Error())
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
			goFiles = append(goFiles, l.src+"/"+f.Name())
		}
	}

	l.srcFiles = goFiles

	return nil
}

func (l *Locator) setDest(dest string) error {
	dir, err := filepath.Abs(dest)
	if err != nil {
		return errors.New("could not get absolute path of destination (" + dest + ")")
	}

	l.dest = strings.Replace(dir, "\\", "/", -1) // use / as path separator even on windows

	return nil
}

func (l *Locator) setDestPkg(destPkg string) error {
	if destPkg != "" {
		l.destPkg = destPkg
		return nil
	}

	if l.setDestPkgWithGopath() {
		return nil
	}

	if l.setDestPkgWithModFile() {
		return nil
	}

	return errors.New("could not determine destination package (destination is not in GOPATH or in a module)")
}

func (l *Locator) setDestPkgWithGopath() bool {
	if os.Getenv("GOPATH") == "" {
		return false
	}

	gopath, err := filepath.Abs(os.Getenv("GOPATH"))
	if err != nil {
		return false
	}

	if !strings.HasPrefix(l.dest, gopath) {
		return false
	}

	gopath = strings.Replace(gopath, "\\", "/", -1) // use / as path separator even on windows

	l.destPkg = ""

	for i := 5 + len(gopath); i < len(l.dest); i++ {
		l.destPkg += string(l.dest[i])
	}

	return true
}

func (l *Locator) setDestPkgWithModFile() bool {
	dir := l.dest

	for true {
		if dir == "/" {
			return false
		}

		dir = path.Dir(dir)

		if s, err := os.Stat(dir + "/go.mod"); err != nil || s.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(dir + "/go.mod")
		if err != nil {
			continue
		}

		if l.setDesPkgWithModFileContent(dir, content) {
			return true
		}
	}

	return false
}

func (l *Locator) setDesPkgWithModFileContent(dir string, content []byte) bool {
	mod := l.modulePath(content)
	if mod == "" {
		return false
	}

	l.destPkg = strings.ReplaceAll(l.dest, dir, mod)

	return true
}

// copied from https://github.com/golang/go/blob/master/src/cmd/go/internal/modfile/read.go#L837
func (l *Locator) modulePath(mod []byte) string {
	slashSlash := []byte("//")
	moduleStr := []byte("module")

	for len(mod) > 0 {
		line := mod
		mod = nil
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, mod = line[:i], line[i+1:]
		}
		if i := bytes.Index(line, slashSlash); i >= 0 {
			line = line[:i]
		}
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, moduleStr) {
			continue
		}
		line = line[len(moduleStr):]
		n := len(line)
		line = bytes.TrimSpace(line)
		if len(line) == n || len(line) == 0 {
			continue
		}

		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(string(line))
			if err != nil {
				return "" // malformed quoted string or multiline module path
			}
			return p
		}

		return string(line)
	}
	return "" // missing module path
}
