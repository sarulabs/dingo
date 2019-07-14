package tools

import "github.com/sarulabs/dingo/v3"

// Provider is an interface that contains the public methods
// of the generated Provider. It is used in the Scanner
// because it avoids generating a Scanner from a template
// to include the generated Provider.
type Provider interface {
	Names() []string
	Get(name string) (*dingo.Def, error)
}
