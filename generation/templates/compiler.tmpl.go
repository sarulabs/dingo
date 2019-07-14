<<</* #############################
###### BASE
############################# */>>>

<<< define "base" ->>>

package main

import (
   	"fmt"

	"github.com/sarulabs/dingo/v2"
	"<<< .TemplatesPkg >>>"
	"<<< .ToolsPkg >>>"
	"<<< .TmpDependenciesPkg >>>"
)

func main() {
	if dingo.Version != "1" {
		fmt.Println("This command requires having github.com/sarulabs/dingo/v2 at version 1, but got version " + dingo.Version)
		return
	}

	// Load provider.
	provider := dependencies.NewProvider()

	err := provider.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Scan definitions.
	scanner := &tools.Scanner{Provider: provider}

	scan, err := scanner.Scan()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write definitions.
	err = tools.WriteTemplate("<<< .DicDir >>>/defs.go", templates.Defs, map[string]interface{}{
		"DependenciesPkg": "<<< .DependenciesPkg >>>",
		"DiPkg": "<<< .DiPkg >>>",
		"Imports": scan.TypeManager.Imports(),
		"Defs": scan.Defs,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate container.
	err = tools.WriteTemplate("<<< .DicDir >>>/container.go", templates.Container, map[string]interface{}{
		"DependenciesPkg": "<<< .DependenciesPkg >>>",
		"DiPkg": "<<< .DiPkg >>>",
		"Imports": scan.ImportsWithoutParams,
		"Defs": scan.Defs,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("dingoCompilerSuccess")
}
<<< end >>>
