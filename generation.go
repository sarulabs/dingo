package dingo

import (
	"fmt"
	"os"

	"github.com/sarulabs/dingo/v4/templates"
)

// GenerateContainer generates a depedency injection container.
// The definitions are loaded from the Provider.
// The code is generated in the outputDirectory.
func GenerateContainer(provider Provider, outputDirectory string) error {
	return GenerateContainerWithCustomPkgName(provider, outputDirectory, "dic")
}

// GenerateContainerWithCustomPkgName works like GenerateContainer
// but let you customize the package name in which files are generated.
// The package name must be a valid package name or the generation may fail unexpectedly.
func GenerateContainerWithCustomPkgName(provider Provider, outputDirectory, pkgName string) error {
	scan, err := scanDefs(provider)
	if err != nil {
		return err
	}

	return writeScan(scan, outputDirectory, pkgName)
}

func scanDefs(provider Provider) (*Scan, error) {
	scanner := &Scanner{Provider: provider}

	scan, err := scanner.Scan()
	if err != nil {
		return nil, fmt.Errorf("could not scan definitions: %v", err)
	}

	return scan, nil
}

func writeScan(scan *Scan, outputDirectory, pkgName string) error {
	err := os.RemoveAll(outputDirectory + "/" + pkgName)
	if err != nil {
		return fmt.Errorf("could not remove destination directory: %v", err)
	}

	err = templates.WriteTemplate(
		outputDirectory+"/"+pkgName+"/defs.go",
		templates.DefsTemplate,
		map[string]interface{}{
			"PkgName": pkgName,
			"Imports": scan.TypeManager.Imports(),
			"Defs":    scan.Defs,
		},
	)
	if err != nil {
		return fmt.Errorf("could not generate definition file: %v", err)
	}

	err = templates.WriteTemplate(
		outputDirectory+"/"+pkgName+"/container.go",
		templates.ContainerTemplate,
		map[string]interface{}{
			"PkgName":         pkgName,
			"Imports":         scan.ImportsWithoutParams,
			"Defs":            scan.Defs,
			"ProviderPackage": scan.ProviderPackage,
			"ProviderName":    scan.ProviderName,
		},
	)
	if err != nil {
		return fmt.Errorf("could not generate container file: %v", err)
	}

	return nil
}
