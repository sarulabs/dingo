package generation

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// NewLocator creates an initialized Locator.
func NewLocator(src, dest string) (*Locator, error) {
	loc := &Locator{}
	return loc, loc.Init(src, dest)
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
func (l *Locator) Init(src, dest string) error {
	if err := l.setSource(src); err != nil {
		return err
	}
	if err := l.setSourceFiles(); err != nil {
		return err
	}
	if err := l.setDest(dest); err != nil {
		return err
	}
	return l.setDestPkg()
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

func (l *Locator) setDestPkg() error {
	gopath, err := filepath.Abs(os.Getenv("GOPATH"))
	if err != nil {
		return errors.New("could not get absolute path of GOPATH (" + os.Getenv("GOPATH") + ")")
	}

	gopath = strings.Replace(gopath, "\\", "/", -1) // use / as path separator even on windows

	l.destPkg = ""

	for i := 5 + len(gopath); i < len(l.dest); i++ {
		l.destPkg += string(l.dest[i])
	}

	return nil
}
