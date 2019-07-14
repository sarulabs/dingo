package pkg

// BuildStructTestA is a structure used in the tests.
type BuildStructTestA struct {
	P1 string
	P2 *BuildStructTestB
	P3 *BuildStructTestC
}

// BuildStructTestB is a structure used in the tests.
type BuildStructTestB struct {
	P1 string
	P2 *BuildStructTestC
}

// BuildStructTestC is a structure used in the tests.
type BuildStructTestC struct {
	P1 string
}
