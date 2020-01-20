package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"

	"github.com/fatih/color"
	"github.com/sarulabs/dingo/v3/generation"
	"github.com/sarulabs/dingo/v3/generation/tools"
)

const version = "3.1.0"

type cmdFlags struct {
	src     *string
	dest    *string
	destPkg *string
}

func main() {
	flags, err := parseFlags()
	if err != nil {
		printError(err.Error())
		printTitle()
		printUsage()
		return
	}

	// Create Locator.
	loc, err := generation.NewLocator(*flags.src, *flags.dest, *flags.destPkg)
	handleError(err)

	// generate container
	generateContainer(loc, packr.NewBox("../generation"))

	printSuccess("code has been successfully generated in " + loc.DestDir())
}

func generateContainer(loc *generation.Locator, box packr.Box) {
	// Clear the temporary destination directory.
	err := os.RemoveAll(loc.DestDir() + "/dictmp")
	handleError(errors.Wrap(err, "could not remove temporary destination directory"))

	// Find declarations in the definition files.
	// Also check package name in files.
	declarations, err := generation.NewParser(loc).Parse()
	handleError(err)

	// Copy files from the definition directory.
	// Only the copies will be used from now on.
	err = generation.NewDefinitionCopier(loc, "dictmp/dependencies/dingorawdefs").Copy()
	handleError(err)

	// Generate a Provider that will allow access to the definitions.
	err = tools.WriteTemplate(loc.DestDir()+"/dictmp/dependencies/provider.go", box.String("templates/provider.tmpl.go"), map[string]interface{}{
		"DumpedDefsPkg": loc.DestPkg() + "/dic/dependencies/dingorawdefs",
		"Decls":         declarations.String("dingorawdefs.", ", "),
	})
	handleError(err)

	// Copy github.com/sarulabs/di.
	err = generation.NewDICopier(loc, "/dictmp/dependencies/di", box).Copy()
	handleError(err)

	// Generate and execute the code that will scan the definitions
	// and generate the dependency injection container.
	err = generation.NewGenerator(loc, box, declarations).Run()
	handleError(err)

	// Code has been generated in the temporary directory.
	// Move it the the destination directory.
	err = os.RemoveAll(loc.DestDir() + "/dic")
	handleError(errors.Wrap(err, "could not remove destination directory"))
	err = os.Rename(loc.DestDir()+"/dictmp", loc.DestDir()+"/dic")
	handleError(errors.Wrap(err, "could not move files to destination directory"))
}

func parseFlags() (cmdFlags, error) {
	fs := flag.NewFlagSet("dingo", flag.ContinueOnError)
	fs.Usage = func() {}
	fs.SetOutput(ioutil.Discard)

	flags := cmdFlags{
		src:     fs.String("src", "", ""),
		dest:    fs.String("dest", "", ""),
		destPkg: fs.String("destPkg", "", ""),
	}

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return flags, err
	}

	if *flags.src == "" || *flags.dest == "" {
		return flags, errors.New("-src and -dest flags can not be empty")
	}

	return flags, nil
}

func printTitle() {
	fmt.Println(color.BlueString(`  __                                
 /\ \  __                           
 \_\ \/\_\    ___      __     ___   
 /'_`+"`"+` \/\ \ /' _ `+"`"+`\  /'_ `+"`"+`\  / __`+"`"+`\ 
/\ \L\ \ \ \/\ \/\ \/\ \L\ \/\ \L\ \
\ \___,_\ \_\ \_\ \_\ \____ \ \____/
 \/__,_ /\/_/\/_/\/_/\/___L\ \/___/ 
                       /\____/      
       `) + color.YellowString("v"+version) + color.BlueString(`          \_/__/
`))
}

func printUsage() {
	fmt.Println(`dingo is a code generator that creates a dependency injection container for your application.

` + color.GreenString(`----------------------------------
usage is simple:

dingo -src=path/to/definition/directory -dest=path/to/generated/code
----------------------------------`) + `

• The generated code will be in the path/to/generated/code directory.
• If the destination directory previously existed, it will be deleted and replaced by a new one.
`)
}

func printError(msg string) {
	fmt.Println(color.RedString("[KO]") + " " + toTitle(msg))
}

func printSuccess(msg string) {
	fmt.Println(color.GreenString("[OK]") + " " + toTitle(msg))
}

func handleError(err error) {
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}

func toTitle(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
