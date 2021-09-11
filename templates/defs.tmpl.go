package templates

// DefsTemplate is the template
// used to generate the definition file.
var DefsTemplate = `
<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>
	package <<< .PkgName >>>

	import (
		"errors"

		"github.com/sarulabs/di/v2"
		"github.com/sarulabs/dingo/v4"
<<< range $pkg, $alias := .Imports >>>
		<<< $alias >>> "<<< $pkg >>>"<<< end >>>
	)

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
		<<<-  if ne .CloseTypeString "" >>>
		Close: func(obj interface{}) error <<< template "closeBody" . >>>,
		<<<- end >>>
		Unshared: <<< .Unshared >>>,
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
			var eo <<< .ObjectTypeString >>>
			return eo, err
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
		var p<<< .Index >>> <<< .TypeString >>>
	<<<- else ->>>
		<<< if ne .ServiceName "" ->>>
			pi<<< .Index >>>, err := ctn.SafeGet("<<< .ServiceName >>>")
			if err != nil {
				var eo <<< .Def.ObjectTypeString >>>
				return eo, err
			}
		<<< else ->>>
			pi<<< .Index >>>, ok := d.Params["<<< .Name >>>"]
			if !ok {
				var eo <<< .Def.ObjectTypeString >>>
				return eo, errors.New("could not find parameter <<< .Name >>>")
			}
		<<< end ->>>
		p<<< .Index >>>, ok := pi<<< .Index >>>.(<<< .TypeString >>>)
		if !ok {
			var eo <<< .Def.ObjectTypeString >>>
			return eo, errors.New("could not cast parameter <<< .Name >>> to <<< .TypeString >>>")
		}
	<<<- end ->>>
<<< end >>>


<<</* #############################
###### OBJECT FUNC
############################# */>>>

<<< define "objectFunc" ->>>
	b, ok := d.Build.(<<< .BuildTypeString >>>)
	if !ok {
		var eo <<< .ObjectTypeString >>>
		return eo, errors.New("could not cast build function to <<< .BuildTypeString >>>")
	}
	return b(<<< .ParamsString >>>)
<<<- end >>>


<<</* #############################
###### OBJECT NEW
############################# */>>>

<<< define "objectNew" ->>>
	return &<<< .BuildTypeString >>>{
		<<< .ParamsString >>>}, nil
<<<- end >>>


<<</* #############################
###### CLOSE BODY
############################# */>>>

<<< define "closeBody" >>>
	<<<- if eq .CloseTypeString "" ->>>
		{
			return nil
		}
	<<<- else ->>>
		{
			d, err := provider.Get("<<< .Name >>>")
			if err != nil {
				return err
			}
			c, ok := d.Close.(<<< .CloseTypeString >>>)
			if !ok {
				return errors.New("could not cast close function to '<<< .CloseTypeString >>>'")
			}
			o, ok := obj.(<<< .ObjectTypeString >>>)
			if !ok {
				return errors.New("could not cast object to '<<< .ObjectTypeString >>>'")
			}
			return c(o)
		}
	<<<- end ->>>
<<< end >>>
`
