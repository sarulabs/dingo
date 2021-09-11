package dingo

import (
	"encoding/json"
	"reflect"
	"sort"
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

// GenerateCommentScope returns the scope as it should be printed in the generated comments.
func (def *ScannedDef) GenerateCommentScope() string {
	if def.Scope == "" {
		return "main"
	}
	return strings.ReplaceAll(def.Scope, "\n", "")
}

// GenerateCommentDescription returns the description as it should be printed in the generated comments.
func (def *ScannedDef) GenerateCommentDescription() string {
	if def.Def.Description == "" {
		return ""
	}
	comment := ""
	for _, part := range strings.Split(def.Def.Description, "\n") {
		comment += "\t\t// " + part + "\n"
	}
	return comment + "\t\t//\n"
}

// GenerateCommentParams returns the params as they should be printed in the generated comments.
func (def *ScannedDef) GenerateCommentParams() string {
	if len(def.Params) == 0 {
		return "\t\t// \tparams: nil\n"
	}

	keys := []string{}
	for key := range def.Params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	comment := "\t\t// \tparams:\n"
	for _, key := range keys {
		p, _ := def.Params[key]
		k, _ := json.Marshal(key)

		comment += "\t\t// \t\t- " + string(k) + ": "
		if p.ServiceName != "" {
			name, _ := json.Marshal(p.ServiceName)
			comment += "Service(" + strings.ReplaceAll(p.TypeString, "\n", "") + ")"
			comment += " [" + string(name) + "]\n"
		} else {
			comment += "Value(" + strings.ReplaceAll(p.TypeString, "\n", "") + ")\n"
		}
	}

	return comment
}

// GenerateComment returns the text used in the comments of the generated code.
func (def *ScannedDef) GenerateComment() string {
	comment := def.GenerateCommentDescription()

	name, _ := json.Marshal(def.Name)
	scope, _ := json.Marshal(def.GenerateCommentScope())

	comment += "\t\t// ---------------------------------------------\n"
	comment += "\t\t// \tname: " + string(name) + "\n"
	comment += "\t\t// \ttype: " + strings.ReplaceAll(def.ObjectTypeString, "\n", "") + "\n"
	comment += "\t\t// \tscope: " + string(scope) + "\n"

	if def.BuildIsFunc {
		comment += "\t\t// \tbuild: func" + "\n"
	} else {
		comment += "\t\t// \tbuild: struct" + "\n"
	}

	comment += def.GenerateCommentParams()

	if def.Unshared {
		comment += "\t\t// \tunshared: true" + "\n"
	} else {
		comment += "\t\t// \tunshared: false" + "\n"
	}
	if def.Def.Close != nil {
		comment += "\t\t// \tclose: true" + "\n"
	} else {
		comment += "\t\t// \tclose: false" + "\n"
	}

	comment += "\t\t// ---------------------------------------------"

	return comment
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
