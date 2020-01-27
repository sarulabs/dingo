package services

import (
	"github.com/sarulabs/dingo/v4"
	"github.com/sarulabs/dingo/v4/tests/app/models"
)

// StructDecl is used in the tests.
var StructDecl = dingo.Def{
	Name:  "test_decl_type_0",
	Build: (*models.DeclTypeTest)(nil),
}

// PtrDecl is used in the tests.
var PtrDecl = &dingo.Def{
	Name:  "test_decl_type_1",
	Build: (*models.DeclTypeTest)(nil),
}

// StructSliceDecls is used in the tests.
var StructSliceDecls = []dingo.Def{
	{
		Name:  "test_decl_type_2",
		Build: (*models.DeclTypeTest)(nil),
	},
	{
		Name:  "test_decl_type_3",
		Build: (*models.DeclTypeTest)(nil),
	},
}

// PtrSliceDecls is used in the tests.
var PtrSliceDecls = []*dingo.Def{
	{
		Name:  "test_decl_type_4",
		Build: (*models.DeclTypeTest)(nil),
	},
	{
		Name:  "test_decl_type_5",
		Build: (*models.DeclTypeTest)(nil),
	},
}

// StructFuncDecl is used in the tests.
var StructFuncDecl = func() dingo.Def {
	return dingo.Def{
		Name:  "test_decl_type_6",
		Build: (*models.DeclTypeTest)(nil),
	}
}

// PtrFuncDecl is used in the tests.
var PtrFuncDecl = func() *dingo.Def {
	return &dingo.Def{
		Name:  "test_decl_type_7",
		Build: (*models.DeclTypeTest)(nil),
	}
}

// StructSliceFuncDecls is used in the tests.
var StructSliceFuncDecls = func() []dingo.Def {
	return []dingo.Def{
		{
			Name:  "test_decl_type_8",
			Build: (*models.DeclTypeTest)(nil),
		},
		{
			Name:  "test_decl_type_9",
			Build: (*models.DeclTypeTest)(nil),
		},
	}
}

// PtrSliceFuncDecls is used in the tests.
var PtrSliceFuncDecls = func() []*dingo.Def {
	return []*dingo.Def{
		{
			Name:  "test_decl_type_10",
			Build: (*models.DeclTypeTest)(nil),
		},
		{
			Name:  "test_decl_type_11",
			Build: (*models.DeclTypeTest)(nil),
		},
	}
}
