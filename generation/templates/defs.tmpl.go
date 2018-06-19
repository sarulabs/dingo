<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>

package dic

import (
	dingoerrors "errors"

	dingodi "<<< .DiPkg >>>"
	dingodependencies "<<< .DependenciesPkg >>>"
<<< range $pkg, $alias := .Imports >>>
	<<< $alias >>> "<<< $pkg >>>"<<< end >>>
)

func getDefinitions(provider *dingodependencies.Provider) []dingodi.Definition {
	return []dingodi.Definition{
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
			Build: func(ctn dingodi.Container) (interface{}, error) {
<<< template "buildBody" . >>>
			},
			Close: func(obj interface{}) <<< template "closeBody" . >>>,
		},
<<<- end >>>


<<</* #############################
###### BUILD BODY
############################# */>>>

<<< define "buildBody" ->>>
<<< if .BuildDependsOnRawDef >>>
				d, err := provider.Get("<<< .Name >>>")
				if err != nil {
					return <<< .EmptyObject >>>, err
				}
<<<- end >>>
<<<- range $index, $param := .Params >>><<< template "buildParam" $param >>><<< end ->>>
<<< if .BuildIsFunc >>><<<- template "objectFunc" . ->>><<< else >>><<<- template "objectNew" . ->>><<< end >>>
<<<- end >>>


<<</* #############################
###### BUILD PARAM
############################# */>>>

<<< define "buildParam" >>>
<<< if .UndefinedStructParam >>>
				p<<< .Index >>> := <<< .Empty >>>
<<< else >>>
	<<< if ne .ServiceName "" >>>
				pi<<< .Index >>>, err := ctn.SafeGet("<<< .ServiceName >>>")
				if err != nil {
					return <<< .Def.EmptyObject >>>, err
				}
	<<< else >>>
				pi<<< .Index >>>, ok := d.Params["<<< .Name >>>"]
				if !ok {
					return <<< .Def.EmptyObject >>>, dingoerrors.New("could not find parameter <<< .Name >>> of <<< .Def.Name >>>")
				}
	<<< end >>>
				p<<< .Index >>>, ok := pi<<< .Index >>>.(<<< .Type >>>)
				if !ok {
					return <<< .Def.EmptyObject >>>, dingoerrors.New("could not cast parameter <<< .Name >>> of <<< .Def.Name >>> to <<< .Type >>>")
				}
<<< end >>>
<<< end >>>


<<</* #############################
###### OBJECT FUNC
############################# */>>>

<<< define "objectFunc" >>>
				b, ok := d.Build.(<<< .BuildType >>>)
				if !ok {
					return <<< .EmptyObject >>>, dingoerrors.New("could not cast build of <<< .Name >>> to <<< .BuildType >>>")
				}

				return b(<<< .ParamsString >>>)
<<<- end >>>


<<</* #############################
###### OBJECT NEW
############################# */>>>

<<< define "objectNew" >>>
				return &<<< .BuildType >>>{<<< .ParamsString >>>				}, nil
<<<- end >>>


<<</* #############################
###### CLOSE BODY
############################# */>>>

<<< define "closeBody" >>>
<<<- if eq .CloseType "" ->>>
{}
<<<- else ->>>
{
				d, err := provider.Get("<<< .Name >>>")
				if err != nil {
					return
				}

				c, ok := d.Close.(<<< .CloseType >>>)
				if !ok {
					return
				}

				o, ok := obj.(<<< .ObjectType >>>)
				if !ok {
					return
				}

				c(o)
			}
<<<- end ->>>
<<< end >>>
