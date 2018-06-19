package tools

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/sarulabs/dingo"
)

// Scan contains the parsed information about the service definitions.
type Scan struct {
	TypeManager          *TypeManager
	ImportsWithoutParams map[string]string
	Defs                 []*ScannedDef
}

// ScannedDef contains the parsed information about a service definition.
type ScannedDef struct {
	Def           *dingo.Def
	Name          string
	FormattedName string
	Scope         string
	ObjectType    string
	EmptyObject   string
	BuildIsFunc   bool
	BuildType     string
	Params        map[string]*ParamInfo
	CloseType     string
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
		params += `
					` + param.Name + `: p` + param.Index + `,
`
	}

	return params
}

// BuildDependsOnRawDef returns true if the service contructor
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
	Type                 string
	Empty                string
	UndefinedStructParam bool
	Def                  *ScannedDef
}

var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digits = []byte("0123456789")

// FormatDefName is the function used to turn the definition name
// into something that can be used as a method name
// for the generated container.
func FormatDefName(name string) string {
	formatted := strings.Builder{}
	start := true

	for _, c := range name {
		if !bytes.ContainsRune(chars, c) {
			start = true
			continue
		}

		if !start {
			formatted.WriteRune(c)
			continue
		}

		formatted.WriteRune(unicode.ToUpper(c))
		start = false
	}

	return formatted.String()
}

// DefNameIsAllowed returns an error if the definition name is not allowed.
func DefNameIsAllowed(name string) error {
	names := []string{"", "C", "ErrorCallback", "Container", "NewContainer"}

	formatted := FormatDefName(name)

	for _, n := range names {
		if n == formatted {
			return errors.New("DefName '" + name + "' is not allowed (reserved key word)")
		}
	}

	if bytes.ContainsRune(digits, rune(formatted[0])) {
		return errors.New("DefName '" + name + "' is not allowed (first char is a digit)")
	}

	return nil
}

// FormatPkgName formats a package name by keeping only the letters.
func FormatPkgName(name string) string {
	formatted := strings.Builder{}

	for _, c := range []byte(name) {
		if bytes.Contains(letters, []byte{c}) {
			formatted.WriteByte(c)
		}
	}

	return formatted.String()
}
