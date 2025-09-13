package matcher

import (
	"testing"
)

func TestMatcher_NewMatcher(t *testing.T) {
	m := NewMatcher()

	if m == nil {
		t.Fatal("Expected NewMatcher() to return non-nil")
	}
}

func TestMatcher_Integer(t *testing.T) {
	m := NewMatcher()

	var x int

	// Test that Integer returns the matcher for chaining
	result := m.Integer(&x)
	if result != m {
		t.Error("Expected Integer() to return the same matcher instance")
	}

	// Test multiple additions for chaining
	var a, b, c int
	m2 := m.
		Integer(&a).
		Integer(&b).
		Integer(&c)

	if m2 != m {
		t.Error("Expected chaining to work correctly")
	}
}

func TestMatcher_Match(t *testing.T) {
	// This test will fail until we implement the bitstring package
	// For now, we'll test the basic structure

	var a, b, c int

	m := NewMatcher().
		Integer(&a).
		Integer(&b).
		Integer(&c)

	// This should fail because we don't have a bitstring yet
	// but we can test the method exists and returns appropriate error types
	results, err := m.Match(nil)

	if err == nil {
		t.Error("Expected error when matching nil bitstring")
	}

	if results != nil {
		t.Error("Expected nil results when matching nil bitstring")
	}
}

func TestMatcher_MatchVariables(t *testing.T) {
	// Test that variables are properly bound
	var x, y, z int

	m := NewMatcher().
		Integer(&x).
		Integer(&y).
		Integer(&z)

	// Initial values should be zero
	if x != 0 || y != 0 || z != 0 {
		t.Errorf("Expected initial values to be zero, got x=%d, y=%d, z=%d", x, y, z)
	}

	// After matching with nil, values should still be zero (no change on error)
	_, err := m.Match(nil)
	if err == nil {
		t.Error("Expected error when matching nil bitstring")
	}

	if x != 0 || y != 0 || z != 0 {
		t.Errorf("Expected values to remain zero after failed match, got x=%d, y=%d, z=%d", x, y, z)
	}
}

func TestMatcher_EmptyMatcher(t *testing.T) {
	m := NewMatcher()

	// Empty matcher should match empty bitstring
	// This will fail until we implement bitstring package
	results, err := m.Match(nil)

	if err == nil {
		t.Error("Expected error when matching nil bitstring with empty matcher")
	}

	if results != nil {
		t.Error("Expected nil results when matching nil bitstring")
	}
}

func TestMatcher_Chaining(t *testing.T) {
	var a, b, c int

	// Test that chaining works properly
	m := NewMatcher()
	result1 := m.Integer(&a)
	result2 := result1.Integer(&b)
	result3 := result2.Integer(&c)

	// All should return the same instance
	if result1 != m || result2 != m || result3 != m {
		t.Error("Expected all chained methods to return the same matcher instance")
	}

	// Variables should be properly set up
	if a != 0 || b != 0 || c != 0 {
		t.Errorf("Expected initial values to be zero, got a=%d, b=%d, c=%d", a, b, c)
	}
}

func TestMatcher_MultipleVariables(t *testing.T) {
	// Test different types of variables (for future extensibility)
	var intVar int
	var int8Var int8
	var int16Var int16
	var int32Var int32
	var int64Var int64
	var uintVar uint
	var uint8Var uint8
	var uint16Var uint16
	var uint32Var uint32
	var uint64Var uint64

	m := NewMatcher().
		Integer(&intVar).
		Integer(&int8Var).
		Integer(&int16Var).
		Integer(&int32Var).
		Integer(&int64Var).
		Integer(&uintVar).
		Integer(&uint8Var).
		Integer(&uint16Var).
		Integer(&uint32Var).
		Integer(&uint64Var)

	// Should be able to create matcher with different integer types
	if m == nil {
		t.Fatal("Expected matcher to be created successfully")
	}

	// Test with nil bitstring - should error
	_, err := m.Match(nil)
	if err == nil {
		t.Error("Expected error when matching nil bitstring")
	}
}
