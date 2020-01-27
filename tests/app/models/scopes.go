package models

import "math/rand"

// ScopeTest is a structure used in the tests.
type ScopeTest struct {
	CreatedAt int
}

// NewScopeTest creates a new ScopeTest.
func NewScopeTest() *ScopeTest {
	return &ScopeTest{
		CreatedAt: rand.Int(),
	}
}
