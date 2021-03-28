package models

// BuildFuncTestA is a structure used in the tests.
type BuildFuncTestA struct {
	P1 string
	P2 BuildFuncTestB
	P3 *BuildFuncTestC
}

// BuildFuncTestB is a structure used in the tests.
type BuildFuncTestB struct {
	P1 string
	P2 *BuildFuncTestC
}

// BuildFuncTestC is a structure used in the tests.
type BuildFuncTestC struct {
	P1 string
}

// TypeBasedOnBasicType is a type used in the tests.
type TypeBasedOnBasicType int64

// TypeBasedOnSliceOfBasicType is a type used in the tests.
type TypeBasedOnSliceOfBasicType []byte
