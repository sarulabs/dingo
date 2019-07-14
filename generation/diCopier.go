package generation

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
)

// NewDICopier is the DICopier constructor.
func NewDICopier(loc *Locator, subdir string, box packr.Box) *DICopier {
	return &DICopier{
		Output: loc.DestDir() + subdir,
		Box:    box,
	}
}

// DICopier is able to copy the github.com/sarulabs/di repository.
// The files are in github.com/sarulabs/dingo/v3/generation/di.
type DICopier struct {
	// Output is where the repository should be copied.
	Output string
	// Box contains the di files.
	Box packr.Box
}

// Copy copies github.com/sarulabs/di in the destination directory.
func (c *DICopier) Copy() error {
	if err := os.MkdirAll(c.Output, 0775); err != nil {
		return errors.New("could not create destination directory: " + err.Error())
	}

	for _, f := range c.Box.List() {
		if !strings.HasPrefix(f, "di/") {
			continue
		}
		if err := ioutil.WriteFile(c.Output+"/"+f[3:], c.Box.Bytes(f), 0664); err != nil {
			return errors.New("could not copy di file: " + err.Error())
		}
	}

	return nil
}
