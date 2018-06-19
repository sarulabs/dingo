<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>

package dependencies

import (
	"errors"
	"reflect"
	"sort"

	"github.com/sarulabs/dingo"
	"<<< .DumpedDefsPkg >>>"
)

// NewProvider is the Provider constructor.
func NewProvider() *Provider {
	return &Provider{
		defs: map[string]*dingo.Def{},
	}
}

// Provider can provide the service definitions.
type Provider struct {
	defs map[string]*dingo.Def
}

// Load registers the service definitions.
// It should be called first.
func (p *Provider) Load() error {
	decls := []interface{}{<<< .Decls >>>}

	for _, i := range decls {
		if err := p.load(i); err != nil {
			return err
		}
	}

	return nil
}

// Get returns the definition for a given service.
// If the definition does not exist, an error is returned.
func (p *Provider) Get(name string) (*dingo.Def, error) {
	if def, ok := p.defs[name]; ok {
		return def, nil
	}
	return nil, errors.New("could not find definition " + name)
}

// Names returns the names of the definitions.
// The names are sorted by alphabetical order.
func (p *Provider) Names() []string {
	names := make([]string, 0, len(p.defs))

	for name := range p.defs {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func (p *Provider) load(i interface{}) error {
	switch v := i.(type) {
	case dingo.Def:
		return p.add(&v)
	case *dingo.Def:
		return p.add(v)
	case []dingo.Def:
		return p.addSlice(v)
	case []*dingo.Def:
		return p.addPtrSlice(v)
	case func() dingo.Def:
		return p.load(v())
	case func() *dingo.Def:
		return p.load(v())
	case func() []dingo.Def:
		return p.load(v())
	case func() []*dingo.Def:
		return p.load(v())
	default:
		return errors.New("could not load definition with type " + reflect.TypeOf(i).String() +
			" (allowed types: dingo.Def, *dingo.Def, []dingo.Def, []*dingo.Def," +
			" func() dingo.Def, func() *dingo.Def, func() []dingo.Def, func() []*dingo.Def)")
	}
}

func (p *Provider) add(def *dingo.Def) error {
	if _, ok := p.defs[def.Name]; ok {
		return errors.New("could not add definition: " + def.Name + " is already defined")
	}

	d := *def

	p.defs[def.Name] = &d
	return nil
}

func (p *Provider) addSlice(defs []dingo.Def) error {
	for _, def := range defs {
		if err := p.add(&def); err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) addPtrSlice(defs []*dingo.Def) error {
	for _, def := range defs {
		if err := p.add(def); err != nil {
			return err
		}
	}
	return nil
}

<<< end >>>
