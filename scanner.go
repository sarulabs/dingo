package dingo

import (
	"errors"
	"fmt"
	"reflect"
)

// Scanner analyzes the definitions provided by a Provider.
type Scanner struct {
	Provider     Provider
	ParamScanner ParamScanner
	scan         *Scan
}

// Scan creates the Scan for the Scanner definitions.
func (s *Scanner) Scan() (*Scan, error) {
	s.scan = &Scan{
		TypeManager: &TypeManager{},
		Defs:        []*ScannedDef{},
	}

	if err := s.loadProvider(); err != nil {
		return nil, err
	}

	for _, name := range s.Provider.Names() {
		def, err := s.Provider.Get(name)
		if err != nil {
			return nil, err
		}

		if err := s.scanDef(def); err != nil {
			return nil, errors.New("could not scan definition " + def.Name + ": " + err.Error())
		}
	}

	s.scan.ImportsWithoutParams = s.scan.TypeManager.Imports()

	if err := s.ParamScanner.Scan(s.scan); err != nil {
		return nil, err
	}

	return s.scan, nil
}

func (s *Scanner) loadProvider() error {
	providerType := reflect.TypeOf(s.Provider)

	if providerType.Kind() != reflect.Ptr || providerType.Elem().Kind() != reflect.Struct {
		return errors.New("the Provider must be a pointer to a structure")
	}

	s.Provider = reflect.New(providerType.Elem()).Interface().(Provider)

	s.scan.ProviderPackage = providerType.Elem().PkgPath()
	s.scan.ProviderName = providerType.Elem().Name()

	if err := s.Provider.Load(); err != nil {
		return fmt.Errorf("could not load definitions with Provider.Load(): %v", err)
	}

	return nil
}

func (s *Scanner) scanDef(def *Def) error {
	sDef := &ScannedDef{
		Def:           def,
		Name:          def.Name,
		FormattedName: FormatDefName(def.Name),
		Scope:         def.Scope,
		Unshared:      def.Unshared,
	}

	if err := DefNameIsAllowed(sDef.FormattedName); err != nil {
		return err
	}

	if err := s.scanBuild(def, sDef); err != nil {
		return err
	}

	if err := s.scanClose(def, sDef); err != nil {
		return err
	}

	s.scan.Defs = append(s.scan.Defs, sDef)

	return nil
}

func (s *Scanner) scanBuild(def *Def, scannedDef *ScannedDef) error {
	t := reflect.TypeOf(def.Build)

	if t.Kind() == reflect.Func {
		return s.scanBuildFunc(def, scannedDef, t)
	}

	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		return s.scanBuildStruct(def, scannedDef, t)
	}

	return errors.New("Build should be a function or a pointer to a structure")
}

func (s *Scanner) scanBuildFunc(def *Def, scannedDef *ScannedDef, buildT reflect.Type) error {
	if err := s.checkBuildFunc(buildT); err != nil {
		return err
	}

	BuildTypeString, err := s.scan.TypeManager.Register(buildT)
	if err != nil {
		return err
	}

	objType, err := s.scan.TypeManager.Register(buildT.Out(0))
	if err != nil {
		return err
	}

	scannedDef.ObjectType = buildT.Out(0)
	scannedDef.ObjectTypeString = objType
	scannedDef.BuildIsFunc = true
	scannedDef.BuildTypeString = BuildTypeString

	return nil
}

func (s *Scanner) checkBuildFunc(t reflect.Type) error {
	if t.IsVariadic() {
		return errors.New("variadic Build functions are not supported")
	}

	if t.NumOut() != 2 {
		return errors.New("Build function must have 2 output parameters")
	}

	if !t.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return errors.New("Build function second output parameter should be an error")
	}

	return nil
}

func (s *Scanner) scanBuildStruct(def *Def, scannedDef *ScannedDef, buildT reflect.Type) error {
	BuildTypeString, err := s.scan.TypeManager.Register(buildT.Elem())
	if err != nil {
		return err
	}

	objType, err := s.scan.TypeManager.Register(buildT)
	if err != nil {
		return err
	}

	scannedDef.ObjectType = buildT
	scannedDef.ObjectTypeString = objType
	scannedDef.BuildIsFunc = false
	scannedDef.BuildTypeString = BuildTypeString

	return nil
}

func (s *Scanner) scanClose(def *Def, scannedDef *ScannedDef) error {
	if def.Close == nil {
		return nil
	}

	t := reflect.TypeOf(def.Close)

	if t.Kind() != reflect.Func {
		return errors.New("Close should be a function")
	}

	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	if t.NumOut() != 1 || !t.Out(0).Implements(errorInterface) {
		return errors.New("Close should return an error")
	}

	if t.NumIn() != 1 {
		return errors.New("Close should have exactly one input parameter")
	}

	return s.scanCloseParameter(def, scannedDef, t)
}

func (s *Scanner) scanCloseParameter(def *Def, scannedDef *ScannedDef, closeT reflect.Type) error {
	fType, err := s.scan.TypeManager.Register(closeT)
	if err != nil {
		return err
	}

	pType, err := s.scan.TypeManager.Register(closeT.In(0))
	if err != nil {
		return err
	}

	if pType != scannedDef.ObjectTypeString {
		return errors.New("object type is " + scannedDef.ObjectTypeString + " but " + pType + " is used is Close")
	}

	scannedDef.CloseTypeString = fType

	return nil
}
