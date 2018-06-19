package services

import (
	"github.com/sarulabs/dingo"
	"github.com/sarulabs/dingo/dingo/tests/app/pkg"
)

// StructDecl is used in the tests.
var StructDecl = dingo.Def{
	Name:  "test_decl_type_0",
	Build: (*pkg.DeclTypeTest)(nil),
}

// PtrDecl is used in the tests.
var PtrDecl = &dingo.Def{
	Name:  "test_decl_type_1",
	Build: (*pkg.DeclTypeTest)(nil),
}

// StructSliceDecls is used in the tests.
var StructSliceDecls = []dingo.Def{
	{
		Name:  "test_decl_type_2",
		Build: (*pkg.DeclTypeTest)(nil),
	},
	{
		Name:  "test_decl_type_3",
		Build: (*pkg.DeclTypeTest)(nil),
	},
}

// PtrSliceDecls is used in the tests.
var PtrSliceDecls = []*dingo.Def{
	{
		Name:  "test_decl_type_4",
		Build: (*pkg.DeclTypeTest)(nil),
	},
	{
		Name:  "test_decl_type_5",
		Build: (*pkg.DeclTypeTest)(nil),
	},
}

// StructFuncDecl is used in the tests.
var StructFuncDecl = func() dingo.Def {
	return dingo.Def{
		Name:  "test_decl_type_6",
		Build: (*pkg.DeclTypeTest)(nil),
	}
}

// PtrFuncDecl is used in the tests.
var PtrFuncDecl = func() *dingo.Def {
	return &dingo.Def{
		Name:  "test_decl_type_7",
		Build: (*pkg.DeclTypeTest)(nil),
	}
}

// StructSliceFuncDecls is used in the tests.
var StructSliceFuncDecls = func() []dingo.Def {
	return []dingo.Def{
		{
			Name:  "test_decl_type_8",
			Build: (*pkg.DeclTypeTest)(nil),
		},
		{
			Name:  "test_decl_type_9",
			Build: (*pkg.DeclTypeTest)(nil),
		},
	}
}

// PtrSliceFuncDecls is used in the tests.
var PtrSliceFuncDecls = func() []*dingo.Def {
	return []*dingo.Def{
		{
			Name:  "test_decl_type_10",
			Build: (*pkg.DeclTypeTest)(nil),
		},
		{
			Name:  "test_decl_type_11",
			Build: (*pkg.DeclTypeTest)(nil),
		},
	}
}
