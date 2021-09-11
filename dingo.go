package dingo

// Def is the structure containing a service definition.
type Def struct {
	// Name is the key that is used to retrieve the object from the container.
	Name string
	// Scope determines in which container the object is stored.
	// Typical scopes are "app" and "request".
	Scope string
	// NotForAutoFill should be set to true if you
	// do not want to use this service automatically
	// as a dependency in other services.
	NotForAutoFill bool
	// Build defines the service constructor. It can be either:
	// - a pointer to a structure: (*MyStruct)(nil)
	// - a factory function: func(any, any, ...) (any, error)
	Build interface{}
	// Params are used to assist the service constructor.
	Params Params
	// Close should be a function: func(any) error.
	// With any being the type of the service.
	Close interface{}
	// Unshared is false by default. That means that the object is only created once in a given container.
	// They are singleton and the same instance will be returned each time "Get", "SafeGet" or "Fill" is called.
	// If you want to retrieve a new object every time, "Unshared" needs to be set to true.
	Unshared bool
	// Description is a text that describes the service.
	// If provided, the description is used in the comments of the generated code.
	Description string
}

// Params are used to assist the service constructor.
// If the Def.Build field is a pointer to a structure, the keys of the map should be
// among the names of the structure fields. These fields will be filled with
// the associated value in the map.
// If the Def.Build field is a function, it works the same way.
// But the key of the map should be the index
// of the function parameters (e.g.: "0", "1", ...).
//
// key=fieldName¦paramIndex value=any¦dingo.Service|dingo.AutoFill
type Params map[string]interface{}

// Service can be used as Params value.
// It means that the field (or parameter) should be replaced
// by an other service. This service should be retrieved from the container.
type Service string

// AutoFill can be used as Params value to avoid autofill (default is Autofill(true)).
// If a structure field is not in the Params map, the container will try to use
// a service from the container that has the same type.
// Setting the entry to AutoFill(false) will let the field empty in the structure.
type AutoFill bool

// ContainerKey is a type that can be used as key in a context.Context.
// For example it can be use if you want to store
// a container in the Context of an http.Request.
// It is used in the generated C function.
type ContainerKey string
