package tools

import (
	"errors"
	"path"
	"reflect"
	"strconv"
	"strings"
)

// RegisteredType is the representation of a type.
// It contains the type name and how to write
// its empty value in a go file.
type RegisteredType struct {
	Type       string
	EmptyValue string
}

// NewRegisteredType is the RegisteredType constructor.
func NewRegisteredType(typ, empty string) *RegisteredType {
	return &RegisteredType{Type: typ, EmptyValue: empty}
}

// TypeManager maintains a list of all the import paths
// that are used in the types that it has registered.
// It associates a unique alias to all the import paths.
type TypeManager struct {
	imports map[string]string
	aliases map[string]int
}

// Imports returns a map with all the imports that are used
// in the registered types. The key is the import path
// and the value is the alias that has been given by the TypeManager.
func (tm *TypeManager) Imports() map[string]string {
	m := map[string]string{}

	for k, v := range tm.imports {
		m[k] = v
	}

	return m
}

// Register adds a new type in the TypeManager.
func (tm *TypeManager) Register(t reflect.Type) (*RegisteredType, error) {
	switch t.Kind() {
	case reflect.Invalid:
		return nil, errors.New("invalid type")
	case reflect.Bool:
		return NewRegisteredType("bool", "false"), nil
	case reflect.Int:
		return NewRegisteredType("int", "0"), nil
	case reflect.Int8:
		return NewRegisteredType("int8", "0"), nil
	case reflect.Int16:
		return NewRegisteredType("int16", "0"), nil
	case reflect.Int32:
		return NewRegisteredType("int32", "0"), nil
	case reflect.Int64:
		return NewRegisteredType("int64", "0"), nil
	case reflect.Uint:
		return NewRegisteredType("uint", "0"), nil
	case reflect.Uint8:
		return NewRegisteredType("uint8", "0"), nil
	case reflect.Uint16:
		return NewRegisteredType("uint16", "0"), nil
	case reflect.Uint32:
		return NewRegisteredType("uint32", "0"), nil
	case reflect.Uint64:
		return NewRegisteredType("uint64", "0"), nil
	case reflect.Uintptr:
		return nil, errors.New("Uintptr is not supported")
	case reflect.Float32:
		return NewRegisteredType("float32", "0"), nil
	case reflect.Float64:
		return NewRegisteredType("float64", "0"), nil
	case reflect.Complex64:
		return NewRegisteredType("complex64", "0"), nil
	case reflect.Complex128:
		return NewRegisteredType("complex128", "0"), nil
	case reflect.Array:
		return tm.registerArray(t)
	case reflect.Chan:
		return tm.registerChan(t)
	case reflect.Func:
		return tm.registerFunc(t)
	case reflect.Interface:
		return tm.registerInterface(t)
	case reflect.Map:
		return tm.registerMap(t)
	case reflect.Ptr:
		return tm.registerPtr(t)
	case reflect.Slice:
		return tm.registerSlice(t)
	case reflect.String:
		return NewRegisteredType("string", "\"\""), nil
	case reflect.Struct:
		return tm.registerStruct(t)
	case reflect.UnsafePointer:
		return nil, errors.New("UnsafePointer is not supported")
	default:
		return nil, errors.New("type is not supported")
	}
}

func (tm *TypeManager) registerArray(t reflect.Type) (*RegisteredType, error) {
	elt, err := tm.Register(t.Elem())
	if err != nil {
		return nil, err
	}
	return NewRegisteredType("["+strconv.Itoa(t.Len())+"]"+elt.Type, "nil"), nil
}

func (tm *TypeManager) registerChan(t reflect.Type) (*RegisteredType, error) {
	elt, err := tm.Register(t.Elem())
	if err != nil {
		return nil, err
	}
	return NewRegisteredType(t.ChanDir().String()+" "+elt.Type, "nil"), nil
}

func (tm *TypeManager) registerFunc(t reflect.Type) (*RegisteredType, error) {
	in := make([]string, 0, t.NumIn())

	for i := 0; i < t.NumIn(); i++ {
		elt, err := tm.Register(t.In(i))
		if err != nil {
			return nil, err
		}
		in = append(in, elt.Type)
	}

	out := make([]string, 0, t.NumOut())

	for i := 0; i < t.NumOut(); i++ {
		elt, err := tm.Register(t.Out(i))
		if err != nil {
			return nil, err
		}
		out = append(out, elt.Type)
	}

	switch len(out) {
	case 0:
		return NewRegisteredType("func("+strings.Join(in, ", ")+")", "nil"), nil
	case 1:
		return NewRegisteredType("func("+strings.Join(in, ", ")+") "+out[0], "nil"), nil
	default:
		return NewRegisteredType("func("+strings.Join(in, ", ")+") ("+strings.Join(out, ", ")+")", "nil"), nil
	}
}

func (tm *TypeManager) registerInterface(t reflect.Type) (*RegisteredType, error) {
	if alias := tm.addImport(t.PkgPath()); alias != "" {
		return NewRegisteredType(alias+"."+t.Name(), "nil"), nil
	}
	if t.Name() == "" {
		return NewRegisteredType("interface{}", "nil"), nil
	}
	return NewRegisteredType(t.Name(), "nil"), nil
}

func (tm *TypeManager) registerMap(t reflect.Type) (*RegisteredType, error) {
	key, err := tm.Register(t.Key())
	if err != nil {
		return nil, err
	}
	elt, err := tm.Register(t.Elem())
	if err != nil {
		return nil, err
	}
	return NewRegisteredType("map["+key.Type+"]"+elt.Type, "nil"), nil
}

func (tm *TypeManager) registerPtr(t reflect.Type) (*RegisteredType, error) {
	elt, err := tm.Register(t.Elem())
	if err != nil {
		return nil, err
	}
	return NewRegisteredType("*"+elt.Type, "nil"), nil
}

func (tm *TypeManager) registerSlice(t reflect.Type) (*RegisteredType, error) {
	elt, err := tm.Register(t.Elem())
	if err != nil {
		return nil, err
	}
	return NewRegisteredType("[]"+elt.Type, "nil"), nil
}

func (tm *TypeManager) registerStruct(t reflect.Type) (*RegisteredType, error) {
	if alias := tm.addImport(t.PkgPath()); alias != "" {
		return NewRegisteredType(alias+"."+t.Name(), alias+"."+t.Name()+"{}"), nil
	}
	return NewRegisteredType(t.Name(), t.Name()+"{}"), nil
}

func (tm *TypeManager) addImport(pkg string) string {
	if pkg == "" {
		return ""
	}

	pkg = tm.removeVendor(pkg)

	if alias, ok := tm.imports[pkg]; ok {
		return alias
	}

	if tm.imports == nil {
		tm.imports = map[string]string{}
	}

	alias := tm.createAlias(pkg)
	tm.imports[pkg] = alias

	return alias
}

func (tm *TypeManager) createAlias(pkg string) string {
	name := FormatPkgName(path.Base(pkg))

	if strings.HasPrefix(name, "dingo") {
		name = "alias" + name
	}

	if tm.aliases == nil {
		tm.aliases = map[string]int{}
	}

	counter := tm.aliases[name]

	tm.aliases[name] = counter + 1

	if counter > 0 {
		return name + strconv.Itoa(counter)
	}

	return name
}

func (tm *TypeManager) removeVendor(pkg string) string {
	parts := strings.Split(pkg, "/vendor/")

	if len(parts) < 2 {
		return pkg
	}

	return parts[len(parts)-1]
}
