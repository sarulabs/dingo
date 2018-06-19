<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>

package dic

import (
	dingoerrors "errors"
	dingolog "log"
	dingohttp "net/http"

	dingo "github.com/sarulabs/dingo"
	dingodi "<<< .DiPkg >>>"
	dingodependencies "<<< .DependenciesPkg >>>"
<<< range $pkg, $alias := .Imports >>>
	<<< $alias >>> "<<< $pkg >>>"<<< end >>>
)

// C retrieves a Container from an interface.
// The interface can be :
// - a *Container
// - an *http.Request containing a *Container in its Context
//   for the dingo.ContainerKey("dingo") key.
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}

	r, ok := i.(*dingohttp.Request)
	if !ok && ErrorCallback != nil {
		ErrorCallback(dingoerrors.New("could not get container with C()"))
		return nil
	}
	
	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok && ErrorCallback != nil {
		ErrorCallback(dingoerrors.New("could not get container from *http.Request"))
		return nil
	}

	return c
}

// ErrorCallback is a function that is called
// when there is an error while retrieving an object
// with the Get method (and its derivatives).
// The function can be changed to match the needs of your application.
var ErrorCallback = func(err error) {
	dingolog.Println(err.Error())
}

// NewContainer creates a new Container.
// If no scope is provided, dingo.App, dingo.Request and dingo.SubRequest are used.
// The returned Container has the wider scope (dingo.App).
// The SubContainer() method should be called to get a Container in a narrower scope.
func NewContainer(scopes ...string) (*Container, error) {
	if dingo.Version != "1" {
		return nil, dingoerrors.New("The generated code requires github.com/sarulabs/dingo at version 1, but got version " + dingo.Version)
	}

	if len(scopes) == 0 {
		scopes = []string{dingo.App, dingo.Request, dingo.SubRequest}
	}

	b, err := dingodi.NewBuilder(scopes...)
	if err != nil {
		return nil, err
	}

	provider := dingodependencies.NewProvider()
	if err := provider.Load(); err != nil {
		return nil, err
	}

	for _, d := range getDefinitions(provider) {
		if err := b.AddDefinition(d); err != nil {
			return nil, err
		}
	}

	return &Container{ctn: b.Build()}, nil
}

// Container represents a dependency injection container.
// A Container has a scope and may have a parent with a wider scope
// and children with a narrower scope.
// Objects can be retrieved from the Container.
// If the desired object does not already exist in the Container,
// it is built thanks to the object Definition.
// The following attempts to get this object will return the same object.
type Container struct {
	ctn dingodi.Container
}

// Scope returns the Container scope.
func (c *Container) Scope() string {
	return c.ctn.Scope()
}

// Scopes returns the list of available scopes.
func (c *Container) Scopes() []string {
	return c.ctn.Scopes()
}

// ParentScopes returns the list of scopes wider than the Container scope.
func (c *Container) ParentScopes() []string {
	return c.ctn.ParentScopes()
}

// SubScopes returns the list of scopes narrower than the Container scope.
func (c *Container) SubScopes() []string {
	return c.ctn.SubScopes()
}

// Parent returns the parent Container.
func (c *Container) Parent() *Container {
	if p := c.ctn.Parent(); p != nil {
		return &Container{ctn: p}
	}
	return nil
}

// SubContainer creates a new Container in the next subscope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object needs to belong to this scope or a wider one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) Get(name string) interface{} {
	o, err := c.ctn.SafeGet(name)
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a narrower scope.
// To do so UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to Delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
func (c *Container) UnscopedGet(name string) interface{} {
	o, err := c.ctn.UnscopedSafeGet(name)
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// Clean deletes the sub-container created by UnscopedSafeGet or UnscopedGet.
func (c *Container) Clean() {
	c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls their Close function if it exists.
// It will also call DeleteWithSubContainers on each child Container
// and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
func (c *Container) DeleteWithSubContainers() {
	c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers but do not delete the subcontainers.
// If the Container has subcontainers, it will not be deleted right away.
// The deletion only occurs when all the subcontainers have been deleted.
func (c *Container) Delete() {
	c.ctn.Delete()
}

// IsClosed returns true if the Container has been deleted.
func (c *Container) IsClosed() bool {
	return c.ctn.IsClosed()
}

<<< range $index, $def := .Defs ->>>
// SafeGet<<< $def.FormattedName >>> works like SafeGet but only for <<< $def.FormattedName >>>.
// It does not return an interface but a <<< $def.ObjectType >>>.
func (c *Container) SafeGet<<< $def.FormattedName >>>() (<<< $def.ObjectType >>>, error) {
	i, err := c.ctn.SafeGet("<<< $def.Name >>>")
	if err != nil {
		return <<< $def.EmptyObject >>>, err
	}

	o, ok := i.(<<< $def.ObjectType >>>)
	if !ok {
		return <<< $def.EmptyObject >>>, dingoerrors.New("could not cast object to <<< $def.ObjectType >>>")
	}

	return o, nil
}

// Get<<< $def.FormattedName >>> is similar to SafeGet<<< $def.FormattedName >>> but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) Get<<< $def.FormattedName >>>() <<< $def.ObjectType >>> {
	o, err := c.SafeGet<<< $def.FormattedName >>>()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// UnscopedSafeGet<<< $def.FormattedName >>> works like UnscopedSafeGet but only for <<< $def.FormattedName >>>.
// It does not return an interface but a <<< $def.ObjectType >>>.
func (c *Container) UnscopedSafeGet<<< $def.FormattedName >>>() (<<< $def.ObjectType >>>, error) {
	i, err := c.ctn.UnscopedSafeGet("<<< $def.Name >>>")
	if err != nil {
		return <<< $def.EmptyObject >>>, err
	}

	o, ok := i.(<<< $def.ObjectType >>>)
	if !ok {
		return <<< $def.EmptyObject >>>, dingoerrors.New("could not cast object to <<< $def.ObjectType >>>")
	}

	return o, nil
}

// UnscopedGet<<< $def.FormattedName >>> is similar to UnscopedSafeGet<<< $def.FormattedName >>> but it does not return the error.
// The error is handled by ErrorCallback.
func (c *Container) UnscopedGet<<< $def.FormattedName >>>() <<< $def.ObjectType >>> {
	o, err := c.UnscopedSafeGet<<< $def.FormattedName >>>()
	if err != nil && ErrorCallback != nil {
		ErrorCallback(err)
	}
	return o
}

// <<< $def.FormattedName >>> is similar to Get<<< $def.FormattedName >>>.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the Get<<< $def.FormattedName >>> method.
// If the container can not be retrieved, it returns the default value for the returned type.
func <<< $def.FormattedName >>>(i interface{}) <<< $def.ObjectType >>> {
	c := C(i)
	if c == nil {
		return <<< $def.EmptyObject >>>
	}
	return c.Get<<< $def.FormattedName >>>()
}

<<< end >>>

<<< end >>>
