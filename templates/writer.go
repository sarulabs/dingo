package templates

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/tools/imports"
)

// WriteTemplate executes the given templates with the given data.
// Then it writes the result in the output file.
// gofmt is used to format the output file.
func WriteTemplate(filename string, tmpl string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0775); err != nil {
		return fmt.Errorf("mkdir failed: %v", err)
	}

	content, err := ExecuteTemplate(tmpl, data)
	if err != nil {
		return err
	}

	content, err = imports.Process(filename, content, nil)
	if err != nil {
		return fmt.Errorf("formatting file failed: %v", err)
	}

	err = ioutil.WriteFile(filename, content, 0664)
	if err != nil {
		return fmt.Errorf("writing file failed: %v", err)
	}

	return nil
}

// ExecuteTemplate renders the given template.
func ExecuteTemplate(tmpl string, data interface{}) ([]byte, error) {
	t, err := template.New("").Delims("<<<", ">>>").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parsing template failed: %v", err)
	}

	buf := bytes.NewBuffer(nil)

	err = t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return nil, fmt.Errorf("executing template failed: %v", err)
	}

	return buf.Bytes(), nil
}
