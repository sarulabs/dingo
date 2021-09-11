package templates

// ContainerTemplate is the template
// used to generate the container file.
var ContainerTemplate = `
<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>
	package <<< .PkgName >>>

	import (
		"errors"
		"fmt"
		"net/http"

		"github.com/sarulabs/di/v2"
		"github.com/sarulabs/dingo/v4"

		providerPkg "<<< .ProviderPackage >>>"
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
		r, ok := i.(*http.Request)
		if !ok {
			panic("could not get the container with dic.C()")
		}
		c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
		if !ok {
			panic("could not get the container from the given *http.Request in dic.C()")
		}
		return c
	}

	type builder struct {
		builder *di.Builder
	}

	// NewBuilder creates a builder that can be used to create a Container.
	// You probably should use NewContainer to create the container directly.
	// But using NewBuilder allows you to redefine some di services.
	// This can be used for testing.
	// But this behavior is not safe, so be sure to know what you are doing.
	func NewBuilder(scopes ...string) (*builder, error) {
		if len(scopes) == 0 {
			scopes = []string{di.App, di.Request, di.SubRequest}
		}
		b, err := di.NewBuilder(scopes...)
		if err != nil {
			return nil, fmt.Errorf("could not create di.Builder: %v", err)
		}
		provider := &providerPkg.<<< .ProviderName >>>{}
		if err := provider.Load(); err != nil {
			return nil, fmt.Errorf("could not load definitions with the Provider (<<< .ProviderName >>> from <<< .ProviderPackage >>>): %v", err)
		}
		for _, d := range getDiDefs(provider) {
			if err := b.Add(d); err != nil {
				return nil, fmt.Errorf("could not add di.Def in di.Builder: %v", err)
			}
		}
		return &builder{builder: b}, nil
	}

	// Add adds one or more definitions in the Builder.
	// It returns an error if a definition can not be added.
	func (b *builder) Add(defs ...di.Def) error {
		return b.builder.Add(defs...)
	}

	// Set is a shortcut to add a definition for an already built object.
	func (b *builder) Set(name string, obj interface{}) error {
		return b.builder.Set(name, obj)
	}

	// Build creates a Container in the most generic scope.
	func (b *builder) Build() *Container {
		return &Container{ctn: b.builder.Build()}
	}

	// NewContainer creates a new Container.
	// If no scope is provided, di.App, di.Request and di.SubRequest are used.
	// The returned Container has the most generic scope (di.App).
	// The SubContainer() method should be called to get a Container in a more specific scope.
	func NewContainer(scopes ...string) (*Container, error) {
		b, err := NewBuilder(scopes...)
		if err != nil {
			return nil, err
		}
		return b.Build(), nil
	}

	// Container represents a generated dependency injection container.
	// It is a wrapper around a di.Container.
	//
	// A Container has a scope and may have a parent in a more generic scope
	// and children in a more specific scope.
	// Objects can be retrieved from the Container.
	// If the requested object does not already exist in the Container,
	// it is built thanks to the object definition.
	// The following attempts to get this object will return the same object.
	type Container struct {
		ctn di.Container
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
		// SafeGet<<< $def.FormattedName >>> retrieves the "<<< $def.Name >>>" object from the <<< $def.GenerateCommentScope >>> scope.
		//
<<< $def.GenerateComment >>>
		//
		// If the object can not be retrieved, it returns an error.
		func (c *Container) SafeGet<<< $def.FormattedName >>>() (<<< $def.ObjectTypeString >>>, error) {
			i, err := c.ctn.SafeGet("<<< $def.Name >>>")
			if err != nil {
				var eo <<< $def.ObjectTypeString >>>
				return eo, err
			}
			o, ok := i.(<<< $def.ObjectTypeString >>>)
			if !ok {
				return o, errors.New("could get '<<< $def.Name >>>' because the object could not be cast to <<< $def.ObjectTypeString >>>")
			}
			return o, nil
		}

		// Get<<< $def.FormattedName >>> retrieves the "<<< $def.Name >>>" object from the <<< $def.GenerateCommentScope >>> scope.
		//
<<< $def.GenerateComment >>>
		//
		// If the object can not be retrieved, it panics.
		func (c *Container) Get<<< $def.FormattedName >>>() <<< $def.ObjectTypeString >>> {
			o, err := c.SafeGet<<< $def.FormattedName >>>()
			if err != nil {
				panic(err)
			}
			return o
		}

		// UnscopedSafeGet<<< $def.FormattedName >>> retrieves the "<<< $def.Name >>>" object from the <<< $def.GenerateCommentScope >>> scope.
		//
<<< $def.GenerateComment >>>
		//
		// This method can be called even if <<< $def.GenerateCommentScope >>> is a sub-scope of the container.
		// If the object can not be retrieved, it returns an error.
		func (c *Container) UnscopedSafeGet<<< $def.FormattedName >>>() (<<< $def.ObjectTypeString >>>, error) {
			i, err := c.ctn.UnscopedSafeGet("<<< $def.Name >>>")
			if err != nil {
				var eo <<< $def.ObjectTypeString >>>
				return eo, err
			}
			o, ok := i.(<<< $def.ObjectTypeString >>>)
			if !ok {
				return o, errors.New("could get '<<< $def.Name >>>' because the object could not be cast to <<< $def.ObjectTypeString >>>")
			}
			return o, nil
		}

		// UnscopedGet<<< $def.FormattedName >>> retrieves the "<<< $def.Name >>>" object from the <<< $def.GenerateCommentScope >>> scope.
		//
<<< $def.GenerateComment >>>
		//
		// This method can be called even if <<< $def.GenerateCommentScope >>> is a sub-scope of the container.
		// If the object can not be retrieved, it panics.
		func (c *Container) UnscopedGet<<< $def.FormattedName >>>() <<< $def.ObjectTypeString >>> {
			o, err := c.UnscopedSafeGet<<< $def.FormattedName >>>()
			if err != nil {
				panic(err)
			}
			return o
		}

		// <<< $def.FormattedName >>> retrieves the "<<< $def.Name >>>" object from the <<< $def.GenerateCommentScope >>> scope.
		//
<<< $def.GenerateComment >>>
		//
		// It tries to find the container with the C method and the given interface.
		// If the container can be retrieved, it calls the Get<<< $def.FormattedName >>> method.
		// If the container can not be retrieved, it panics.
		func <<< $def.FormattedName >>>(i interface{}) <<< $def.ObjectTypeString >>> {
			return C(i).Get<<< $def.FormattedName >>>()
		}
	<<< end >>>
<<< end >>>
`
