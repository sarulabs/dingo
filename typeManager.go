package dingo

import (
	"errors"
	"path"
	"reflect"
	"strconv"
	"strings"
)

var reservedPkgNames = map[string]struct{}{
	"dingo":       {},
	"di":          {},
	"providerPkg": {},
	"errors":      {},
	"fmt":         {},
	"http":        {},
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
func (tm *TypeManager) Register(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Invalid:
		return "", errors.New("invalid type")
	case reflect.Bool:
		return tm.registerBasicType(t)
	case reflect.Int:
		return tm.registerBasicType(t)
	case reflect.Int8:
		return tm.registerBasicType(t)
	case reflect.Int16:
		return tm.registerBasicType(t)
	case reflect.Int32:
		return tm.registerBasicType(t)
	case reflect.Int64:
		return tm.registerBasicType(t)
	case reflect.Uint:
		return tm.registerBasicType(t)
	case reflect.Uint8:
		return tm.registerBasicType(t)
	case reflect.Uint16:
		return tm.registerBasicType(t)
	case reflect.Uint32:
		return tm.registerBasicType(t)
	case reflect.Uint64:
		return tm.registerBasicType(t)
	case reflect.Uintptr:
		return "", errors.New("Uintptr is not supported")
	case reflect.Float32:
		return tm.registerBasicType(t)
	case reflect.Float64:
		return tm.registerBasicType(t)
	case reflect.Complex64:
		return tm.registerBasicType(t)
	case reflect.Complex128:
		return tm.registerBasicType(t)
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
		return tm.registerBasicType(t)
	case reflect.Struct:
		return tm.registerStruct(t)
	case reflect.UnsafePointer:
		return "", errors.New("UnsafePointer is not supported")
	default:
		return "", errors.New("type is not supported")
	}
}

func (tm *TypeManager) registerArray(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}
	eltType, err := tm.Register(t.Elem())
	if err != nil {
		return "", err
	}
	return "[" + strconv.Itoa(t.Len()) + "]" + eltType, nil
}

func (tm *TypeManager) registerBasicType(t reflect.Type) (string, error) {
	// Check if it is really a basic type ...
	if t.Kind().String() == t.Name() {
		return t.Name(), nil
	}
	// ... or if it is a type based on one.
	return tm.registerNamedType(t)
}

func (tm *TypeManager) registerNamedType(t reflect.Type) (string, error) {
	if alias := tm.addImport(t.PkgPath()); alias != "" {
		return alias + "." + t.Name(), nil
	}
	return t.Name(), nil
}

func (tm *TypeManager) registerChan(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}
	eltType, err := tm.Register(t.Elem())
	if err != nil {
		return "", err
	}
	return t.ChanDir().String() + " " + eltType, nil
}

func (tm *TypeManager) registerFunc(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}

	inTypes := make([]string, 0, t.NumIn())

	for i := 0; i < t.NumIn(); i++ {
		eltType, err := tm.Register(t.In(i))
		if err != nil {
			return "", err
		}
		inTypes = append(inTypes, eltType)
	}

	outTypes := make([]string, 0, t.NumOut())

	for i := 0; i < t.NumOut(); i++ {
		eltType, err := tm.Register(t.Out(i))
		if err != nil {
			return "", err
		}
		outTypes = append(outTypes, eltType)
	}

	switch len(outTypes) {
	case 0:
		return "func(" + strings.Join(inTypes, ", ") + ")", nil
	case 1:
		return "func(" + strings.Join(inTypes, ", ") + ") " + outTypes[0], nil
	default:
		return "func(" + strings.Join(inTypes, ", ") + ") (" + strings.Join(outTypes, ", ") + ")", nil
	}
}

func (tm *TypeManager) registerInterface(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}
	return "interface{}", nil
}

func (tm *TypeManager) registerMap(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}
	keyType, err := tm.Register(t.Key())
	if err != nil {
		return "", err
	}
	eltType, err := tm.Register(t.Elem())
	if err != nil {
		return "", err
	}
	return "map[" + keyType + "]" + eltType, nil
}

func (tm *TypeManager) registerPtr(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}
	eltType, err := tm.Register(t.Elem())
	if err != nil {
		return "", err
	}
	return "*" + eltType, nil
}

func (tm *TypeManager) registerSlice(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}
	eltType, err := tm.Register(t.Elem())
	if err != nil {
		return "", err
	}
	return "[]" + eltType, nil
}

func (tm *TypeManager) registerStruct(t reflect.Type) (string, error) {
	if t.Name() != "" {
		return tm.registerNamedType(t)
	}
	return "struct{}", nil
}

func (tm *TypeManager) addImport(pkg string) string {
	if pkg == "" {
		return ""
	}

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

	if _, ok := reservedPkgNames[name]; ok {
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
