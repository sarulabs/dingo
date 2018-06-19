package generation

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ErrParserUnknownDecl is the Parser error for unknown declarations.
var ErrParserUnknownDecl = errors.New("definition files can only contain package, import, var and func declarations")

// ErrParserVarDecl is the Parser error for invalid var declarations.
var ErrParserVarDecl = errors.New("var declaration contains an error")

// NewParser is the Parser constructor.
func NewParser(loc *Locator) *Parser {
	return &Parser{Locator: loc}
}

// Parser can find the definition declarations in the source files.
type Parser struct {
	Locator *Locator
}

// Parse finds the definition declarations in the source files.
// It may also return an error if the package name of the source file
// does not match the name of its directory.
func (p *Parser) Parse() (*DeclarationList, error) {
	list := NewDeclarationList()

	for _, f := range p.Locator.SourceFiles() {
		fileList, err := p.parseFile(f)
		if err != nil {
			return nil, errors.New("could not parse definition file (" + f + "): " + err.Error())
		}

		if err := list.Merge(fileList); err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (p *Parser) parseFile(filename string) (*DeclarationList, error) {
	f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}

	if err := p.checkPkg(f); err != nil {
		return nil, err
	}

	list := NewDeclarationList()

	for _, decl := range f.Decls {
		if err := p.parseDecl(decl, list); err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (p *Parser) checkPkg(f *ast.File) error {
	if f.Name == nil {
		return errors.New("could not find package definition")
	}

	if f.Name.Name != p.Locator.SourceDirName() {
		return errors.New("expected package `" + p.Locator.SourceDirName() + "` but got `" + f.Name.Name + "`")
	}

	return nil
}

func (p *Parser) parseDecl(decl ast.Decl, list *DeclarationList) error {
	switch d := decl.(type) {
	case *ast.GenDecl:
		return p.parseGenDecl(d, list)
	case *ast.FuncDecl:
		return p.parseFuncDecl(d, list)
	default:
		return ErrParserUnknownDecl
	}
}

func (p *Parser) parseGenDecl(decl *ast.GenDecl, list *DeclarationList) error {
	if decl.Tok == token.IMPORT {
		return p.checkImportSpecs(decl.Specs)
	}

	if decl.Tok != token.VAR {
		return ErrParserUnknownDecl
	}

	return p.parseVarSpecs(decl.Specs, list)
}

func (p *Parser) checkImportSpecs(specs []ast.Spec) error {
	for _, spec := range specs {
		s, ok := spec.(*ast.ImportSpec)
		if !ok || s.Path == nil {
			return ErrParserUnknownDecl
		}

		if p.importIsRelative(s.Path.Value) {
			return errors.New("relative imports are not supported: " + s.Path.Value)
		}
	}
	return nil
}

func (p *Parser) importIsRelative(path string) bool {
	return strings.HasPrefix(path, "\".")
}

func (p *Parser) parseVarSpecs(specs []ast.Spec, list *DeclarationList) error {
	if len(specs) != 1 {
		return ErrParserVarDecl
	}

	spec, ok := specs[0].(*ast.ValueSpec)
	if !ok {
		return ErrParserVarDecl
	}

	for _, ident := range spec.Names {
		if err := p.addIdent(ident, list); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) addIdent(ident *ast.Ident, list *DeclarationList) error {
	if ident == nil {
		return nil
	}

	return list.Add(ident.Name)
}

func (p *Parser) parseFuncDecl(decl *ast.FuncDecl, list *DeclarationList) error {
	if decl == nil {
		return errors.New("functions should have a name")
	}

	if decl.Recv != nil {
		return errors.New("only functions allowed, not methods")
	}

	return p.addIdent(decl.Name, list)
}
