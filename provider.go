package dingo

import (
	"errors"
	"reflect"
	"sort"
)

// Provider is the interface used to store the definitions.
// The provider is used while generating the dependency
// injection container, but also while executing the
// code of the container.
type Provider interface {
	Load() error
	Names() []string
	Get(name string) (*Def, error)
}

// BaseProvider implements the Provider interface.
// It contains no definition, but you can use this
// to create your own Provider by redefining the Load method.
type BaseProvider struct {
	defs map[string]*Def
}

// Load registers the service definitions.
// You need to override this method to add the service definitions.
// You can use the Add method to add a definition inside the provider.
func (p *BaseProvider) Load() error {
	return nil
}

// Get returns the definition for a given service.
// If the definition does not exist, an error is returned.
func (p *BaseProvider) Get(name string) (*Def, error) {
	if def, ok := p.defs[name]; ok {
		return def, nil
	}
	return nil, errors.New("could not find definition " + name)
}

// Names returns the names of the definitions.
// The names are sorted by alphabetical order.
func (p *BaseProvider) Names() []string {
	names := make([]string, 0, len(p.defs))

	for name := range p.defs {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Add adds definitions inside the Provider.
// Allowed types are:
// dingo.Def, *dingo.Def, []dingo.Def, []*dingo.Def
// func() dingo.Def, func() *dingo.Def, func() []dingo.Def, func() []*dingo.Def
func (p *BaseProvider) Add(i interface{}) error {
	switch v := i.(type) {
	case Def:
		return p.AddDef(v)
	case *Def:
		return p.AddDefPtr(v)
	case []Def:
		return p.AddDefSlice(v)
	case []*Def:
		return p.AddDefPtrSlice(v)
	case func() Def:
		return p.AddDef(v())
	case func() *Def:
		return p.AddDefPtr(v())
	case func() []Def:
		return p.AddDefSlice(v())
	case func() []*Def:
		return p.AddDefPtrSlice(v())
	default:
		return errors.New("could not load definition with type " + reflect.TypeOf(i).String() +
			" (allowed types: dingo.Def, *dingo.Def, []dingo.Def, []*dingo.Def," +
			" func() dingo.Def, func() *dingo.Def, func() []dingo.Def, func() []*dingo.Def)")
	}
}

// AddDef is the same as Add, but only for Def.
func (p *BaseProvider) AddDef(def Def) error {
	if p.defs == nil {
		p.defs = map[string]*Def{}
	}
	if _, ok := p.defs[def.Name]; ok {
		return errors.New("could not add definition: " + def.Name + " is already defined")
	}
	p.defs[def.Name] = &def
	return nil
}

// AddDefPtr is the same as Add, but only for *Def.
func (p *BaseProvider) AddDefPtr(def *Def) error {
	return p.AddDef(*def)
}

// AddDefSlice is the same as Add, but only for []Def.
func (p *BaseProvider) AddDefSlice(defs []Def) error {
	for _, def := range defs {
		if err := p.AddDef(def); err != nil {
			return err
		}
	}
	return nil
}

// AddDefPtrSlice is the same as Add, but only for []*Def.
func (p *BaseProvider) AddDefPtrSlice(defs []*Def) error {
	for _, def := range defs {
		if err := p.AddDefPtr(def); err != nil {
			return err
		}
	}
	return nil
}
