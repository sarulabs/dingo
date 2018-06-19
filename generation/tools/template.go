package tools

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// WriteTemplate executes the given templates with the given data.
// Then it stores the result in the output file.
func WriteTemplate(output string, tmpl string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(output), 0775); err != nil {
		return errors.New("could not create directory: " + err.Error())
	}

	t, err := template.New("").Delims("<<<", ">>>").Parse(tmpl)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)

	err = t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(output, buf.Bytes(), 0664)
}
