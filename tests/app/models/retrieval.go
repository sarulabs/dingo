package models

import "math/rand"

// RetrievalTest is a structure used in the tests.
type RetrievalTest struct {
	CreatedAt int
}

// NewRetrievalTest creates a new RetrievalTest.
func NewRetrievalTest() *RetrievalTest {
	return &RetrievalTest{
		CreatedAt: rand.Int(),
	}
}
