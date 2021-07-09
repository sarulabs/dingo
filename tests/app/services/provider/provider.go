package provider

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/services"
)

// Provider with the test definitions.
type Provider struct {
	dingo.BaseProvider
}

// Load adds the definitions in the provider.
func (p *Provider) Load() error {
	if err := p.Add(services.StructDecl); err != nil {
		return err
	}
	if err := p.Add(services.PtrDecl); err != nil {
		return err
	}
	if err := p.Add(services.StructSliceDecls); err != nil {
		return err
	}
	if err := p.Add(services.PtrSliceDecls); err != nil {
		return err
	}
	if err := p.Add(services.StructFuncDecl); err != nil {
		return err
	}
	if err := p.Add(services.PtrFuncDecl); err != nil {
		return err
	}
	if err := p.Add(services.StructSliceFuncDecls); err != nil {
		return err
	}
	if err := p.Add(services.PtrSliceFuncDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.AutofillDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.BuildFuncDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.BuildStructDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.CloseDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.DiDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.InterfacesDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.RetrievalDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.ScopeDecls); err != nil {
		return err
	}
	if err := p.AddDefSlice(services.UnsharedDecls); err != nil {
		return err
	}
	return nil
}
