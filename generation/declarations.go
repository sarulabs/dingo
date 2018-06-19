package generation

import (
	"errors"
	"strings"
)

// NewDeclarationList is the DeclarationList constructor.
func NewDeclarationList() *DeclarationList {
	return &DeclarationList{
		decls: map[string]struct{}{},
	}
}

// DeclarationList is the structure
// used to store the name of the declarations in the source files.
// A declaration is the name of the variable containing one or more dingo.Def.
type DeclarationList struct {
	list  []string
	decls map[string]struct{}
}

// Add adds a new declaration in the DeclarationList.
func (dl *DeclarationList) Add(name string) error {
	if name == "" {
		return errors.New("declaration name can not be empty")
	}

	if name != strings.Title(name) {
		return errors.New("declaration `" + name + "` must be exported")
	}

	if _, ok := dl.decls[name]; ok {
		return errors.New("declaration `" + name + "` is defined more than once")
	}

	dl.list = append(dl.list, name)
	dl.decls[name] = struct{}{}

	return nil
}

// Merge adds the declarations of the declList to the DeclarationList.
func (dl *DeclarationList) Merge(declList *DeclarationList) error {
	for _, name := range declList.list {
		if err := dl.Add(name); err != nil {
			return err
		}
	}

	return nil
}

// String turns the declarations into a string.
// It works like strings.Join on the declaration names,
// but it also allows to add a prefix before each element.
func (dl *DeclarationList) String(prefix, sep string) string {
	decls := make([]string, 0, len(dl.list))

	for _, decl := range dl.list {
		decls = append(decls, prefix+decl)
	}

	return strings.Join(decls, sep)
}
