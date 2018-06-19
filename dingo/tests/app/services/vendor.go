package services

import (
	anotherpkgAlias "anotherpkg"
	"otherpkg"

	"github.com/sarulabs/dingo"
	"github.com/sarulabs/dingo/dingo/tests/app/pkg"
)

// VendorDecls is used in the tests.
var VendorDecls = []dingo.Def{
	{
		Name:  "test_vendor_1",
		Build: (*pkg.VendorTest)(nil),
	},
	{
		Name:  "test_vendor_2",
		Build: (*otherpkg.StructOtherPkg)(nil),
		Params: dingo.Params{
			"Value": "OK",
		},
	},
	{
		Name:  "test_vendor_3",
		Build: (*anotherpkgAlias.StructAnotherPkg)(nil),
		Params: dingo.Params{
			"Value": "OK",
		},
	},
}
