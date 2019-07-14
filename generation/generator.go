package generation

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/sarulabs/dingo/v3/generation/tools"
)

// NewGenerator is the Generator constructor.
func NewGenerator(loc *Locator, box packr.Box, decls *DeclarationList) *Generator {
	return &Generator{
		Locator:      loc,
		Box:          box,
		Declarations: decls,
	}
}

// Generator creates the compiler files and runs the compiler.
// It generates the dependency injection container.
type Generator struct {
	Locator      *Locator
	Box          packr.Box
	Declarations *DeclarationList
}

// Run executes the Generator.
func (gen *Generator) Run() error {
	tmpDir := gen.Locator.DestDir() + "/dictmp/tmp"

	defer os.RemoveAll(tmpDir)

	if err := gen.generateProvider(tmpDir); err != nil {
		return errors.New("could not create temporary provider: " + err.Error())
	}

	if err := gen.generateTemplates(tmpDir); err != nil {
		return errors.New("could not copy templates: " + err.Error())
	}

	if err := gen.generateTools(tmpDir); err != nil {
		return errors.New("could not generate tools: " + err.Error())
	}

	err := tools.WriteTemplate(tmpDir+"/main.go", gen.Box.String("templates/compiler.tmpl.go"), map[string]interface{}{
		"DiPkg":              gen.Locator.DestPkg() + "/dic/dependencies/di",
		"DicDir":             gen.Locator.DestDir() + "/dictmp",
		"TemplatesPkg":       gen.Locator.DestPkg() + "/dictmp/tmp/templates",
		"ToolsPkg":           gen.Locator.DestPkg() + "/dictmp/tmp/tools",
		"DependenciesPkg":    gen.Locator.DestPkg() + "/dic/dependencies",
		"TmpDependenciesPkg": gen.Locator.DestPkg() + "/dictmp/tmp/dependencies",
	})
	if err != nil {
		return errors.New("could not write Compiler template: " + err.Error())
	}


	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = tmpDir
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "dingoCompilerSuccess") {
		return errors.New("an error occurred while running the generated compiler:\n" + string(out))
	}

	return nil
}

func (gen *Generator) generateProvider(tmpDir string) error {
	if err := os.MkdirAll(tmpDir+"/dependencies", 0775); err != nil {
		return errors.New("could not create " + tmpDir + "/dependencies directory: " + err.Error())
	}

	// Generate a Provider that will be used only during compilation.
	return tools.WriteTemplate(
		gen.Locator.DestDir()+"/dictmp/tmp/dependencies/provider.go",
		gen.Box.String("templates/provider.tmpl.go"),
		map[string]interface{}{
			"DumpedDefsPkg": gen.Locator.DestPkg() + "/dictmp/dependencies/dingorawdefs",
			"Decls":         gen.Declarations.String("dingorawdefs.", ", "),
		},
	)
}

func (gen *Generator) generateTemplates(tmpDir string) error {
	if err := os.MkdirAll(tmpDir+"/templates", 0775); err != nil {
		return errors.New("could not create " + tmpDir + "/templates directory: " + err.Error())
	}

	tmpls := map[string]string{
		"Container": "templates/container.tmpl.go",
		"Defs":      "templates/defs.tmpl.go",
	}

	for varName, fileName := range tmpls {
		content := strings.Replace(gen.Box.String(fileName), "`", "`+\"`\"+`", -1)

		code := "package templates\n"
		code += "\n"
		code += "// " + varName + " is a generated template file.\n"
		code += "var " + varName + " = `\n"
		code += string(content)
		code += "`\n"

		if err := ioutil.WriteFile(tmpDir+"/"+fileName, []byte(code), 0664); err != nil {
			return errors.New("could not generate template file " + fileName + ": " + err.Error())
		}
	}

	return nil
}

func (gen *Generator) generateTools(tmpDir string) error {
	if err := os.MkdirAll(tmpDir+"/tools", 0775); err != nil {
		return errors.New("could not create " + tmpDir + "/templates directory: " + err.Error())
	}

	for _, f := range gen.Box.List() {
		if !strings.HasPrefix(f, "tools/") {
			continue
		}
		if err := ioutil.WriteFile(tmpDir+"/"+f, gen.Box.Bytes(f), 0664); err != nil {
			return errors.New("could not generate tool file " + f + ": " + err.Error())
		}
	}

	return nil
}
