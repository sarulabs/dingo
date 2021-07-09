package dingo

import (
	"reflect"
	"strconv"
	"strings"
)

// Scan contains the parsed information about the service definitions.
type Scan struct {
	TypeManager          *TypeManager
	ImportsWithoutParams map[string]string
	Defs                 []*ScannedDef
	ProviderPackage      string
	ProviderName         string
}

// ScannedDef contains the parsed information about a service definition.
type ScannedDef struct {
	Def              *Def
	Name             string
	FormattedName    string
	Scope            string
	ObjectType       reflect.Type
	ObjectTypeString string
	BuildIsFunc      bool
	BuildTypeString  string
	Params           map[string]*ParamInfo
	CloseTypeString  string
	Unshared         bool
}

// ParamsString returns the parameters as they should appear
// in a structure inside a go file.
func (def *ScannedDef) ParamsString() string {
	if def.BuildIsFunc {
		params := make([]string, len(def.Params))

		for i := 0; i < len(def.Params); i++ {
			params[i] = "p" + strconv.Itoa(i)
		}

		return strings.Join(params, ", ")
	}

	params := ""

	for _, param := range def.Params {
		params += param.Name + `: p` + param.Index + ",\n"
	}

	return params
}

// BuildDependsOnRawDef returns true if the service constructor
// needs the definition contained in the Provider.
func (def *ScannedDef) BuildDependsOnRawDef() bool {
	if def.BuildIsFunc {
		return true
	}
	for _, param := range def.Params {
		if param.ServiceName == "" && !param.UndefinedStructParam {
			return true
		}
	}
	return false
}

// ParamInfo contains the parsed information about a parameter.
type ParamInfo struct {
	Name                 string
	Index                string
	ServiceName          string
	Type                 reflect.Type
	TypeString           string
	UndefinedStructParam bool
	Def                  *ScannedDef
}
