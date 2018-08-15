package dingo

// App is the name of the application scope.
const App = "app"

// Request is the name of the request scope.
const Request = "request"

// SubRequest is the name of the subrequest scope.
const SubRequest = "subrequest"

// Def is the structure containing a service definition.
type Def struct {
	Name  string
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

// Version is used by the dingo command to ensure that
// the github.com/sarulabs/dingo/dingo package
// is in the right version inside your project.
const Version = "1"

// ContainerKey is a type that can be used to store a container
// in the Context of an http.Request.
// By default, it is used in the generated C function.
type ContainerKey string
