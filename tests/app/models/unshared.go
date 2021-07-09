package models

import "math/rand"

// UnsharedTest is a structure used in the tests.
type UnsharedTest struct {
	CreatedAt int
}

// NewUnsharedTest creates a new UnsharedTest.
func NewUnsharedTest() *UnsharedTest {
	return &UnsharedTest{
		CreatedAt: rand.Int(),
	}
}
