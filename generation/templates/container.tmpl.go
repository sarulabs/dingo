<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>

package dic

import (
	dingoerrors "errors"
	dingohttp "net/http"

	dingo "github.com/sarulabs/dingo/v3"
	dingodi "<<< .DiPkg >>>"
	dingodependencies "<<< .DependenciesPkg >>>"
<<< range $pkg, $alias := .Imports >>>
	<<< $alias >>> "<<< $pkg >>>"<<< end >>>
)

// C retrieves a Container from an interface.
// The function panics if the Container can not be retrieved.
//
// The interface can be :
// - a *Container
// - an *http.Request containing a *Container in its context.Context
//   for the dingo.ContainerKey("dingo") key.
//
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}

	r, ok := i.(*dingohttp.Request)
	if !ok {
		panic("could not get the container with C()")
	}

	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok {
		panic("could not get the container from the given *http.Request")
	}

	return c
}

// NewContainer creates a new Container.
// If no scope is provided, dingo.App, dingo.Request and dingo.SubRequest are used.
// The returned Container has the most generic scope (dingo.App).
// The SubContainer() method should be called to get a Container in a more specific scope.
func NewContainer(scopes ...string) (*Container, error) {
	if dingo.Version != "1" {
		return nil, dingoerrors.New("The generated code requires github.com/sarulabs/dingo/v3 at version 1, but got version " + dingo.Version)
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
		if err := b.Add(d); err != nil {
			return nil, err
		}
	}

	return &Container{ctn: b.Build()}, nil
}

// Container represents a dependency injection container.
// To create a Container, you should use a Builder or another Container.
//
// A Container has a scope and may have a parent in a more generic scope
// and children in a more specific scope.
// Objects can be retrieved from the Container.
// If the requested object does not already exist in the Container,
// it is built thanks to the object definition.
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

// SubScopes returns the list of scopes that are more specific than the Container scope.
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

// SubContainer creates a new Container in the next sub-scope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object has to belong to this scope or a more generic one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// Instead it panics.
func (c *Container) Get(name string) interface{} {
	return c.ctn.Get(name)
}

// Fill is similar to SafeGet but it does not return the object.
// Instead it fills the provided object with the value returned by SafeGet.
// The provided object must be a pointer to the value returned by SafeGet.
func (c *Container) Fill(name string, dst interface{}) error {
	return c.ctn.Fill(name, dst)
}

// Put places destination object to the container by key
func (c *Container) Put(name string, dst interface{}) error {
	return c.ctn.Put(name, dst)
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a more specific scope.
// To do so, UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGet(name string) interface{} {
	return c.ctn.UnscopedGet(name)
}

// UnscopedFill is similar to UnscopedSafeGet but copies the object in dst instead of returning it.
func (c *Container) UnscopedFill(name string, dst interface{}) error {
	return c.ctn.UnscopedFill(name, dst)
}

// Clean deletes the sub-container created by UnscopedSafeGet, UnscopedGet or UnscopedFill.
func (c *Container) Clean() error {
	return c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls the Close function of their Definition on them.
// It will also call DeleteWithSubContainers on each child and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
// The sub-containers are deleted even if they are still used in other goroutines.
// It can cause errors. You may want to use the Delete method instead.
func (c *Container) DeleteWithSubContainers() error {
	return c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers if the Container does not have any child.
// But if the Container has sub-containers, it will not be deleted right away.
// The deletion only occurs when all the sub-containers have been deleted manually.
// So you have to call Delete or DeleteWithSubContainers on all the sub-containers.
func (c *Container) Delete() error {
	return c.ctn.Delete()
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
		return <<< $def.EmptyObject >>>, dingoerrors.New("could get `<<< $def.Name >>>` because the object could not be cast to <<< $def.ObjectType >>>")
	}

	return o, nil
}

// Get<<< $def.FormattedName >>> is similar to SafeGet<<< $def.FormattedName >>> but it does not return the error.
// Instead it panics.
func (c *Container) Get<<< $def.FormattedName >>>() <<< $def.ObjectType >>> {
	o, err := c.SafeGet<<< $def.FormattedName >>>()
	if err != nil {
		panic(err)
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
		return <<< $def.EmptyObject >>>, dingoerrors.New("could get `<<< $def.Name >>>` because the object could not be cast to <<< $def.ObjectType >>>")
	}

	return o, nil
}

// UnscopedGet<<< $def.FormattedName >>> is similar to UnscopedSafeGet<<< $def.FormattedName >>> but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGet<<< $def.FormattedName >>>() <<< $def.ObjectType >>> {
	o, err := c.UnscopedSafeGet<<< $def.FormattedName >>>()
	if err != nil {
		panic(err)
	}
	return o
}

// <<< $def.FormattedName >>> is similar to Get<<< $def.FormattedName >>>.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the Get<<< $def.FormattedName >>> method.
// If the container can not be retrieved, it panics.
func <<< $def.FormattedName >>>(i interface{}) <<< $def.ObjectType >>> {
	return C(i).Get<<< $def.FormattedName >>>()
}

<<< end >>>

<<< end >>>
