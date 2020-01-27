package templates

// DefsTemplate is the template
// used to generate the definition file.
var DefsTemplate = `
<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>
	package dic

	import (
		"errors"

		"github.com/sarulabs/di/v2"
		"github.com/sarulabs/dingo/v4"
<<< range $pkg, $alias := .Imports >>>
		<<< $alias >>> "<<< $pkg >>>"<<< end >>>
	)

	var _ = errors.New("always import errors")

	func getDiDefs(provider dingo.Provider) []di.Def {
		return []di.Def{
			<<<- range $index, $def := .Defs ->>>
				<<< template "definition" $def >>>
			<<<- end >>>
		}
	}
<<< end >>>


<<</* #############################
###### DEFINITION
############################# */>>>

<<< define "definition" >>>
	{
		Name: "<<< .Name >>>",
		Scope: "<<< .Scope >>>",
		Build: func(ctn di.Container) (interface{}, error) <<< template "buildBody" . >>>,
		Close: func(obj interface{}) error <<< template "closeBody" . >>>,
	},
<<<- end >>>


<<</* #############################
###### BUILD BODY
############################# */>>>

<<< define "buildBody" ->>>
	{
	<<<- if .BuildDependsOnRawDef >>>
		d, err := provider.Get("<<< .Name >>>")
		if err != nil {
			return <<< .EmptyObject >>>, err
		}
	<<<- end >>>
	<<<- range $index, $param := .Params >>>
		<<< template "buildParam" $param >>>
	<<<- end >>>
	<<<- if .BuildIsFunc >>>
		<<< template "objectFunc" . >>>
	<<<- else >>>
		<<< template "objectNew" . >>>
	<<<- end >>>
	}
<<<- end >>>


<<</* #############################
###### BUILD PARAM
############################# */>>>

<<< define "buildParam" >>>
	<<<- if .UndefinedStructParam ->>>
		p<<< .Index >>> := <<< .Empty >>>
	<<<- else ->>>
		<<< if ne .ServiceName "" ->>>
			pi<<< .Index >>>, err := ctn.SafeGet("<<< .ServiceName >>>")
			if err != nil {
				return <<< .Def.EmptyObject >>>, err
			}
		<<< else ->>>
			pi<<< .Index >>>, ok := d.Params["<<< .Name >>>"]
			if !ok {
				return <<< .Def.EmptyObject >>>, errors.New("could not find parameter <<< .Name >>>")
			}
		<<< end ->>>
		p<<< .Index >>>, ok := pi<<< .Index >>>.(<<< .Type >>>)
		if !ok {
			return <<< .Def.EmptyObject >>>, errors.New("could not cast parameter <<< .Name >>> to <<< .Type >>>")
		}
	<<<- end ->>>
<<< end >>>


<<</* #############################
###### OBJECT FUNC
############################# */>>>

<<< define "objectFunc" ->>>
	b, ok := d.Build.(<<< .BuildType >>>)
	if !ok {
		return <<< .EmptyObject >>>, errors.New("could not cast build function to <<< .BuildType >>>")
	}
	return b(<<< .ParamsString >>>)
<<<- end >>>


<<</* #############################
###### OBJECT NEW
############################# */>>>

<<< define "objectNew" ->>>
	return &<<< .BuildType >>>{
		<<< .ParamsString >>>}, nil
<<<- end >>>


<<</* #############################
###### CLOSE BODY
############################# */>>>

<<< define "closeBody" >>>
	<<<- if eq .CloseType "" ->>>
		{
			return nil
		}
	<<<- else ->>>
		{
			d, err := provider.Get("<<< .Name >>>")
			if err != nil {
				return err
			}
			c, ok := d.Close.(<<< .CloseType >>>)
			if !ok {
				return errors.New("could not cast close function to '<<< .CloseType >>>'")
			}
			o, ok := obj.(<<< .ObjectType >>>)
			if !ok {
				return errors.New("could not cast object to '<<< .ObjectType >>>'")
			}
			return c(o)
		}
	<<<- end ->>>
<<< end >>>
`
