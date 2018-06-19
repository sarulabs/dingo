package generation

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NewDefinitionCopier is the DefinitionCopier constructor.
func NewDefinitionCopier(loc *Locator, subDir string) *DefinitionCopier {
	return &DefinitionCopier{Locator: loc, SubDir: subDir}
}

// DefinitionCopier copies the definition files from
// the source directory, to the destination directory.
// It changes the package name in the files.
type DefinitionCopier struct {
	Locator *Locator
	SubDir  string
}

// Copy copies each source file in the destination directory.
func (cop *DefinitionCopier) Copy() error {
	if err := os.MkdirAll(cop.Locator.DestDir()+"/"+cop.SubDir, 0775); err != nil {
		return errors.New("could not create destination directory: " + err.Error())
	}

	for _, f := range cop.Locator.SourceFiles() {
		if err := cop.copyFile(f); err != nil {
			return errors.New("could not dump file: " + err.Error())
		}
	}

	return nil
}

func (cop *DefinitionCopier) copyFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("could not read definition file (" + filename + "): " + err.Error())
	}

	// replace package name
	content = bytes.Replace(
		content,
		[]byte(cop.Locator.SourceDirName()),
		[]byte(filepath.Base(cop.SubDir)),
		1,
	)

	_, name := filepath.Split(filename)

	dest := cop.Locator.DestDir() + "/" + cop.SubDir + "/" + name

	if err := ioutil.WriteFile(dest, content, 0664); err != nil {
		return errors.New("could not write definition file (" + dest + "): " + err.Error())
	}

	return nil
}
