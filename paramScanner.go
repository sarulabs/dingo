package dingo

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ParamScanner helps the Scanner.
// It scans information about params.
type ParamScanner struct {
	scan       *Scan
	defsByName map[string]*ScannedDef
	defsByType map[string][]*ScannedDef
}

// Scan updates the given Scan with data about params.
func (s *ParamScanner) Scan(scan *Scan) error {
	s.scan = scan

	// set defsByName and defsByType
	// defsByType only contains definitions available for autofill
	s.defsByName = map[string]*ScannedDef{}
	s.defsByType = map[string][]*ScannedDef{}

	for _, def := range scan.Defs {
		s.defsByName[def.Name] = def
		if !def.Def.NotForAutoFill {
			s.defsByType[def.ObjectTypeString] = append(s.defsByType[def.ObjectTypeString], def)
		}
	}

	// scan parameters
	for _, def := range scan.Defs {
		if err := s.scanParams(def); err != nil {
			return errors.New("could not scan parameters for definition " + def.Name + ": " + err.Error())
		}
	}

	return nil
}

func (s *ParamScanner) scanParams(def *ScannedDef) error {
	params, err := s.expectedParams(def)
	if err != nil {
		return err
	}

	def.Params = params

	for name := range def.Def.Params {
		if _, ok := params[name]; !ok {
			return errors.New("definition should not have parameter " + name)
		}
	}

	for _, param := range params {
		if err := s.setParam(param, def); err != nil {
			return err
		}
	}

	return nil
}

func (s *ParamScanner) expectedParams(def *ScannedDef) (map[string]*ParamInfo, error) {
	if def.BuildIsFunc {
		return s.expectedFuncParams(def)
	}
	return s.expectedStructParams(def)
}

func (s *ParamScanner) expectedFuncParams(def *ScannedDef) (map[string]*ParamInfo, error) {
	params := map[string]*ParamInfo{}

	t := reflect.TypeOf(def.Def.Build)

	for i := 0; i < t.NumIn(); i++ {
		index := strconv.Itoa(i)

		pType, err := s.scan.TypeManager.Register(t.In(i))
		if err != nil {
			return nil, err
		}

		params[index] = &ParamInfo{
			Name:       index,
			Index:      index,
			Type:       t.In(i),
			TypeString: pType,
			Def:        def,
		}
	}

	return params, nil
}

func (s *ParamScanner) expectedStructParams(def *ScannedDef) (map[string]*ParamInfo, error) {
	params := map[string]*ParamInfo{}

	t := reflect.TypeOf(def.Def.Build).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name != strings.Title(f.Name) {
			continue
		}

		index := strconv.Itoa(i)

		pType, err := s.scan.TypeManager.Register(f.Type)
		if err != nil {
			return nil, err
		}

		params[f.Name] = &ParamInfo{
			Name:       f.Name,
			Index:      index,
			Type:       f.Type,
			TypeString: pType,
			Def:        def,
		}
	}

	return params, nil
}

func (s *ParamScanner) setParam(param *ParamInfo, def *ScannedDef) error {
	p, ok := def.Def.Params[param.Name]
	if !ok {
		return s.autofill(param, !def.BuildIsFunc)
	}

	if v, ok := p.(Service); ok {
		return s.setServiceParam(param, string(v))
	}

	autofill, ok := p.(AutoFill)
	if ok && bool(autofill) {
		return s.autofill(param, false)
	}
	if ok && !bool(autofill) && def.BuildIsFunc {
		return errors.New("definition can not have parameters with AutoFill(false) because it uses a Build function")
	}
	if ok {
		return nil
	}

	pType, err := s.scan.TypeManager.Register(reflect.TypeOf(p))
	if err != nil {
		return err
	}

	if pType != param.TypeString && !s.implementsInterface(reflect.TypeOf(p), param.Type) {
		return errors.New("param " + param.Name + " should be a " + param.TypeString + " but is a " + pType)
	}

	return nil
}

func (s *ParamScanner) autofill(param *ParamInfo, acceptNotFound bool) error {
	defs, _ := s.defsByType[param.TypeString]
	if len(defs) == 0 && acceptNotFound {
		param.UndefinedStructParam = true
		return nil
	}
	if len(defs) == 0 {
		return fmt.Errorf("autofill require exactly one %s, but found 0 definition with this type", param.TypeString)
	}
	if len(defs) > 1 {
		return fmt.Errorf("autofill require exactly one %s, but found %d definitions with this type", param.TypeString, len(defs))
	}

	param.ServiceName = defs[0].Name

	return nil
}

func (s *ParamScanner) setServiceParam(param *ParamInfo, service string) error {
	def, ok := s.defsByName[service]
	if !ok {
		return errors.New("could not find definition " + service + " for param " + param.Name)
	}

	if def.ObjectTypeString != param.TypeString && !s.implementsInterface(def.ObjectType, param.Type) {
		return errors.New("param " + param.Name + " should be a " + param.TypeString + " but is a " + def.ObjectTypeString)
	}

	param.ServiceName = service

	return nil
}

func (s *ParamScanner) implementsInterface(t, i reflect.Type) bool {
	return i.Kind() == reflect.Interface && t.Implements(i)
}
