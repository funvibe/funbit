package matcher

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/funvibe/funbit/internal/bitstring"
	bitstringpkg "github.com/funvibe/funbit/internal/bitstring"
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

func TestMatcher_Float(t *testing.T) {
	t.Run("Float 32-bit big endian", func(t *testing.T) {
		var result float64
		m := NewMatcher()
		returnedMatcher := m.Float(&result, bitstring.WithSize(32), bitstring.WithEndianness("big"))

		if returnedMatcher != m {
			t.Error("Expected Float() to return the same matcher instance")
		}

		// Check that pattern was added
		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Size != 32 {
			t.Errorf("Expected size 32, got %d", segment.Size)
		}

		if segment.Type != "float" {
			t.Errorf("Expected type 'float', got '%s'", segment.Type)
		}

		if segment.Endianness != "big" {
			t.Errorf("Expected endianness 'big', got '%s'", segment.Endianness)
		}
	})

	t.Run("Float 64-bit little endian", func(t *testing.T) {
		var result float64
		m := NewMatcher()
		returnedMatcher := m.Float(&result, bitstring.WithSize(64), bitstring.WithEndianness("little"))

		if returnedMatcher != m {
			t.Error("Expected Float() to return the same matcher instance")
		}

		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Size != 64 {
			t.Errorf("Expected size 64, got %d", segment.Size)
		}

		if segment.Endianness != "little" {
			t.Errorf("Expected endianness 'little', got '%s'", segment.Endianness)
		}
	})

	t.Run("Float native endianness", func(t *testing.T) {
		var result float64
		m := NewMatcher()
		returnedMatcher := m.Float(&result, bitstring.WithSize(32), bitstring.WithEndianness("native"))

		if returnedMatcher != m {
			t.Error("Expected Float() to return the same matcher instance")
		}

		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Endianness != "native" {
			t.Errorf("Expected endianness 'native', got '%s'", segment.Endianness)
		}
	})

	t.Run("Float with options", func(t *testing.T) {
		var result float64
		m := NewMatcher()
		returnedMatcher := m.Float(&result, bitstring.WithSize(32), bitstring.WithEndianness("big"), bitstring.WithSigned(true))

		if returnedMatcher != m {
			t.Error("Expected Float() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if !segment.Signed {
			t.Error("Expected signed to be true")
		}
	})

	t.Run("Float multiple options", func(t *testing.T) {
		var result float64
		m := NewMatcher()
		returnedMatcher := m.Float(&result,
			bitstring.WithSize(64),
			bitstring.WithEndianness("little"),
			bitstring.WithSigned(true),
			bitstring.WithUnit(8), // unit is uint, not string
		)

		if returnedMatcher != m {
			t.Error("Expected Float() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if !segment.Signed {
			t.Error("Expected signed to be true")
		}

		if segment.Unit != 8 {
			t.Errorf("Expected unit 8, got %d", segment.Unit)
		}
	})
}

func TestMatcher_Binary(t *testing.T) {
	t.Run("Binary 16-bit", func(t *testing.T) {
		var result []byte
		m := NewMatcher()
		returnedMatcher := m.Binary(&result, bitstring.WithSize(16))

		if returnedMatcher != m {
			t.Error("Expected Binary() to return the same matcher instance")
		}

		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Size != 16 {
			t.Errorf("Expected size 16, got %d", segment.Size)
		}

		if segment.Type != "binary" {
			t.Errorf("Expected type 'binary', got '%s'", segment.Type)
		}
	})

	t.Run("Binary 32-bit with options", func(t *testing.T) {
		var result []byte
		m := NewMatcher()
		returnedMatcher := m.Binary(&result, bitstring.WithSize(32), bitstring.WithEndianness("little"))

		if returnedMatcher != m {
			t.Error("Expected Binary() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if segment.Endianness != "little" {
			t.Errorf("Expected endianness 'little', got '%s'", segment.Endianness)
		}
	})

	t.Run("Binary with multiple options", func(t *testing.T) {
		var result []byte
		m := NewMatcher()
		returnedMatcher := m.Binary(&result,
			bitstring.WithSize(64),
			bitstring.WithSigned(true),
			bitstring.WithEndianness("big"),
			bitstring.WithUnit(8), // unit is uint, not string
		)

		if returnedMatcher != m {
			t.Error("Expected Binary() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if !segment.Signed {
			t.Error("Expected signed to be true")
		}

		if segment.Endianness != "big" {
			t.Errorf("Expected endianness 'big', got '%s'", segment.Endianness)
		}

		if segment.Unit != 8 {
			t.Errorf("Expected unit 8, got %d", segment.Unit)
		}
	})
}

func TestMatcher_UTF(t *testing.T) {
	t.Run("UTF-8", func(t *testing.T) {
		var result string
		m := NewMatcher()
		returnedMatcher := m.UTF(&result, bitstring.WithEndianness("utf-8"))

		if returnedMatcher != m {
			t.Error("Expected UTF() to return the same matcher instance")
		}

		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Type != "integer" {
			t.Errorf("Expected type 'integer', got '%s'", segment.Type)
		}

		if segment.Endianness != "utf-8" {
			t.Errorf("Expected endianness 'utf-8', got '%s'", segment.Endianness)
		}
	})

	t.Run("UTF-16 big endian", func(t *testing.T) {
		var result string
		m := NewMatcher()
		returnedMatcher := m.UTF(&result, bitstring.WithEndianness("utf-16be"))

		if returnedMatcher != m {
			t.Error("Expected UTF() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if segment.Endianness != "utf-16be" {
			t.Errorf("Expected endianness 'utf-16be', got '%s'", segment.Endianness)
		}
	})

	t.Run("UTF-32 little endian", func(t *testing.T) {
		var result string
		m := NewMatcher()
		returnedMatcher := m.UTF(&result, bitstring.WithEndianness("utf-32le"))

		if returnedMatcher != m {
			t.Error("Expected UTF() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if segment.Endianness != "utf-32le" {
			t.Errorf("Expected endianness 'utf-32le', got '%s'", segment.Endianness)
		}
	})

	t.Run("UTF with size option", func(t *testing.T) {
		var result string
		m := NewMatcher()
		returnedMatcher := m.UTF(&result, bitstring.WithSize(16), bitstring.WithEndianness("utf-8"))

		if returnedMatcher != m {
			t.Error("Expected UTF() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if segment.Size != 16 {
			t.Errorf("Expected size 16, got %d", segment.Size)
		}
	})

	t.Run("UTF with multiple options", func(t *testing.T) {
		var result string
		m := NewMatcher()
		returnedMatcher := m.UTF(&result,
			bitstring.WithSize(32),
			bitstring.WithEndianness("utf-16"),
			bitstring.WithUnit(1), // UTF unit is always 1
		)

		if returnedMatcher != m {
			t.Error("Expected UTF() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if segment.Size != 32 {
			t.Errorf("Expected size 32, got %d", segment.Size)
		}

		if segment.Unit != 1 {
			t.Errorf("Expected unit 1, got %d", segment.Unit)
		}
	})
}

func TestMatcher_Bitstring(t *testing.T) {
	t.Run("Bitstring 8-bit", func(t *testing.T) {
		var result *bitstringpkg.BitString
		m := NewMatcher()
		returnedMatcher := m.Bitstring(&result, bitstring.WithSize(8))

		if returnedMatcher != m {
			t.Error("Expected Bitstring() to return the same matcher instance")
		}

		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Size != 8 {
			t.Errorf("Expected size 8, got %d", segment.Size)
		}

		if segment.Type != "bitstring" {
			t.Errorf("Expected type 'bitstring', got '%s'", segment.Type)
		}
	})

	t.Run("Bitstring with options", func(t *testing.T) {
		var result *bitstringpkg.BitString
		m := NewMatcher()
		returnedMatcher := m.Bitstring(&result, bitstring.WithSize(16), bitstring.WithSigned(true))

		if returnedMatcher != m {
			t.Error("Expected Bitstring() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if !segment.Signed {
			t.Error("Expected signed to be true")
		}
	})

	t.Run("Bitstring with multiple options", func(t *testing.T) {
		var result *bitstringpkg.BitString
		m := NewMatcher()
		returnedMatcher := m.Bitstring(&result,
			bitstring.WithSize(32),
			bitstring.WithSigned(true),
			bitstring.WithUnit(1), // unit is uint, not string
		)

		if returnedMatcher != m {
			t.Error("Expected Bitstring() to return the same matcher instance")
		}

		segment := m.pattern[0]
		if !segment.Signed {
			t.Error("Expected signed to be true")
		}

		if segment.Unit != 1 {
			t.Errorf("Expected unit 1, got %d", segment.Unit)
		}
	})
}

func TestMatcher_RegisterVariable(t *testing.T) {
	t.Run("Register variable", func(t *testing.T) {
		var testVar uint = 42
		m := NewMatcher()
		returnedMatcher := m.RegisterVariable("test_var", &testVar)

		if returnedMatcher != m {
			t.Error("Expected RegisterVariable() to return the same matcher instance")
		}

		// Check that variable was registered
		if val, exists := m.variables["test_var"]; !exists || val != &testVar {
			t.Error("Expected variable to be registered in variables map")
		}
	})

	t.Run("Register variable with options", func(t *testing.T) {
		var anotherVar uint = 24
		m := NewMatcher()
		returnedMatcher := m.RegisterVariable("another_var", &anotherVar)

		if returnedMatcher != m {
			t.Error("Expected RegisterVariable() to return the same matcher instance")
		}

		// Check that variable was registered
		if val, exists := m.variables["another_var"]; !exists || val != &anotherVar {
			t.Error("Expected variable to be registered in variables map")
		}
	})
}

func TestMatcher_RestBinary(t *testing.T) {
	t.Run("Rest binary", func(t *testing.T) {
		var result []byte
		m := NewMatcher()
		returnedMatcher := m.RestBinary(&result)

		if returnedMatcher != m {
			t.Error("Expected RestBinary() to return the same matcher instance")
		}

		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Type != "rest_binary" {
			t.Errorf("Expected type 'rest_binary', got '%s'", segment.Type)
		}

		// Rest binary doesn't set unit by default, it's 0
		if segment.Unit != 0 {
			t.Errorf("Expected unit 0, got %d", segment.Unit)
		}
	})

	t.Run("Rest binary with options", func(t *testing.T) {
		var result []byte
		m := NewMatcher()
		returnedMatcher := m.RestBinary(&result)

		if returnedMatcher != m {
			t.Error("Expected RestBinary() to return the same matcher instance")
		}

		segment := m.pattern[0]
		// Rest binary doesn't have signed property set by default
		if segment.Signed {
			t.Error("Expected signed to be false")
		}
	})
}

func TestMatcher_RestBitstring(t *testing.T) {
	t.Run("Rest bitstring", func(t *testing.T) {
		var result *bitstringpkg.BitString
		m := NewMatcher()
		returnedMatcher := m.RestBitstring(&result)

		if returnedMatcher != m {
			t.Error("Expected RestBitstring() to return the same matcher instance")
		}

		if len(m.pattern) != 1 {
			t.Errorf("Expected 1 segment in pattern, got %d", len(m.pattern))
		}

		segment := m.pattern[0]
		if segment.Type != "rest_bitstring" {
			t.Errorf("Expected type 'rest_bitstring', got '%s'", segment.Type)
		}

		// Rest bitstring doesn't set unit by default, it's 0
		if segment.Unit != 0 {
			t.Errorf("Expected unit 0, got %d", segment.Unit)
		}
	})

	t.Run("Rest bitstring with options", func(t *testing.T) {
		var result *bitstringpkg.BitString
		m := NewMatcher()
		returnedMatcher := m.RestBitstring(&result)

		if returnedMatcher != m {
			t.Error("Expected RestBitstring() to return the same matcher instance")
		}

		segment := m.pattern[0]
		// Rest bitstring doesn't have signed property set by default
		if segment.Signed {
			t.Error("Expected signed to be false")
		}
	})
}

func TestMatcher_MatchFunctions(t *testing.T) {
	t.Run("Match with simple integer pattern", func(t *testing.T) {
		var result int
		m := NewMatcher()
		m.Integer(&result, bitstring.WithSize(16))

		// Create a bitstring with 16-bit value 0x1234
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0x12, 0x34})

		matcherResults, err := m.Match(bs)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(matcherResults) != 1 {
			t.Errorf("Expected 1 result, got %d", len(matcherResults))
		}

		if !matcherResults[0].Matched {
			t.Error("Expected match to be successful")
		}

		if result != 0x1234 {
			t.Errorf("Expected result 0x1234, got 0x%X", result)
		}
	})

	t.Run("Match with pattern chaining", func(t *testing.T) {
		var result1, result2 int
		m := NewMatcher()
		m.Integer(&result1, bitstring.WithSize(8)).
			Integer(&result2, bitstring.WithSize(8))

		// Create a bitstring with two 8-bit values
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0xAB, 0xCD})

		matcherResults, err := m.Match(bs)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(matcherResults) != 2 {
			t.Errorf("Expected 2 results, got %d", len(matcherResults))
		}

		if !matcherResults[0].Matched || !matcherResults[1].Matched {
			t.Error("Expected all matches to be successful")
		}

		if result1 != 0xAB || result2 != 0xCD {
			t.Errorf("Expected results 0xAB, 0xCD, got 0x%X, 0x%X", result1, result2)
		}
	})

	t.Run("Match with insufficient data", func(t *testing.T) {
		var result int
		m := NewMatcher()
		m.Integer(&result, bitstring.WithSize(16))

		// Create a bitstring with only 8 bits
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0x12})

		_, err := m.Match(bs)

		if err == nil {
			t.Error("Expected error for insufficient data")
		}
	})

	t.Run("Match with float", func(t *testing.T) {
		var result float64
		m := NewMatcher()
		m.Float(&result, bitstring.WithSize(32))

		// Create a bitstring with 32-bit float value 1.5
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, math.Float32bits(1.5))
		bs := bitstringpkg.NewBitStringFromBytes(bytes)

		matcherResults, err := m.Match(bs)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(matcherResults) != 1 {
			t.Errorf("Expected 1 result, got %d", len(matcherResults))
		}

		if !matcherResults[0].Matched {
			t.Error("Expected match to be successful")
		}

		if math.Abs(result-1.5) > 0.001 {
			t.Errorf("Expected result 1.5, got %f", result)
		}
	})

	t.Run("Match with binary", func(t *testing.T) {
		var result []byte
		m := NewMatcher()
		m.Binary(&result, bitstring.WithSize(2)) // 2 bytes

		// Create a bitstring with 16-bit binary data (2 bytes)
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0xAB, 0xCD})

		matcherResults, err := m.Match(bs)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(matcherResults) != 1 {
			t.Errorf("Expected 1 result, got %d", len(matcherResults))
		}

		if !matcherResults[0].Matched {
			t.Error("Expected match to be successful")
		}

		if len(result) != 2 || result[0] != 0xAB || result[1] != 0xCD {
			t.Errorf("Expected result [0xAB, 0xCD], got %v", result)
		}
	})

	t.Run("Match with UTF-8", func(t *testing.T) {
		var result string
		m := NewMatcher()
		m.UTF(&result, bitstring.WithType("utf8")) // Use type instead of endianness for UTF

		// Create a bitstring with UTF-8 character 'A'
		bs := bitstringpkg.NewBitStringFromBytes([]byte{'A'})

		matcherResults, err := m.Match(bs)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(matcherResults) != 1 {
			t.Errorf("Expected 1 result, got %d", len(matcherResults))
		}

		if !matcherResults[0].Matched {
			t.Error("Expected match to be successful")
		}

		if result != "A" {
			t.Errorf("Expected result 'A', got '%s'", result)
		}
	})

	t.Run("Match with nil bitstring", func(t *testing.T) {
		var result int
		m := NewMatcher()
		m.Integer(&result, bitstring.WithSize(8))

		_, err := m.Match(nil)

		if err == nil {
			t.Error("Expected error for nil bitstring")
		}
	})

	t.Run("Match with rest binary", func(t *testing.T) {
		var result []byte
		m := NewMatcher()
		m.RestBinary(&result)

		// Create a bitstring with some data
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0xAB, 0xCD, 0xEF})

		matcherResults, err := m.Match(bs)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(matcherResults) != 1 {
			t.Errorf("Expected 1 result, got %d", len(matcherResults))
		}

		if !matcherResults[0].Matched {
			t.Error("Expected match to be successful")
		}

		expected := []byte{0xAB, 0xCD, 0xEF}
		if len(result) != len(expected) || !bytesEqual(result, expected) {
			t.Errorf("Expected result %v, got %v", expected, result)
		}
	})

	t.Run("Match with rest bitstring", func(t *testing.T) {
		var result *bitstringpkg.BitString
		m := NewMatcher()
		m.RestBitstring(&result)

		// Create a bitstring with some data
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0xAB, 0xCD, 0xEF})

		matcherResults, err := m.Match(bs)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(matcherResults) != 1 {
			t.Errorf("Expected 1 result, got %d", len(matcherResults))
		}

		if !matcherResults[0].Matched {
			t.Error("Expected match to be successful")
		}

		if result == nil || result.Length() != 24 {
			t.Errorf("Expected result with 24 bits, got %d bits", result.Length())
		}
	})
}

func TestMatcher_ExtractFunctions(t *testing.T) {
	t.Run("extractBinaryBits", func(t *testing.T) {
		m := NewMatcher()

		// Test extracting 4 bits from the middle of a byte
		data := []byte{0b11110000}                     // 0xF0
		result, err := m.extractBinaryBits(data, 2, 4) // Extract 4 bits starting at position 2

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Should extract bits 1100 (positions 2-5) and left-align them: 11000000
		expected := []byte{0b11000000}
		if len(result) != 1 || result[0] != expected[0] {
			t.Errorf("Expected result [0b%08b], got [0b%08b]", expected[0], result[0])
		}
	})

	t.Run("extractBinaryBits full byte", func(t *testing.T) {
		m := NewMatcher()

		data := []byte{0xAB}
		result, err := m.extractBinaryBits(data, 0, 8)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(result) != 1 || result[0] != 0xAB {
			t.Errorf("Expected result [0xAB], got %v", result)
		}
	})

	t.Run("extractBinaryBits empty result", func(t *testing.T) {
		m := NewMatcher()

		data := []byte{0xAB}
		result, err := m.extractBinaryBits(data, 0, 0)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(result) != 0 {
			t.Errorf("Expected empty result, got %v", result)
		}
	})

	t.Run("extractNestedBitstring", func(t *testing.T) {
		m := NewMatcher()

		// Create a bitstring with 24 bits
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0x12, 0x34, 0x56})

		// Extract 8 bits starting from bit 8 (second byte)
		result, err := m.extractNestedBitstring(bs, 8, 8)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Length() != 8 {
			t.Errorf("Expected result with 8 bits, got %d bits", result.Length())
		}

		// Check the extracted value
		extractedBytes := result.ToBytes()
		if len(extractedBytes) != 1 || extractedBytes[0] != 0x34 {
			t.Errorf("Expected extracted value 0x34, got %v", extractedBytes)
		}
	})

	t.Run("extractNestedBitstring full bitstring", func(t *testing.T) {
		m := NewMatcher()

		bs := bitstringpkg.NewBitStringFromBytes([]byte{0x12, 0x34})
		result, err := m.extractNestedBitstring(bs, 0, 16)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Length() != 16 {
			t.Errorf("Expected result with 16 bits, got %d bits", result.Length())
		}

		extractedBytes := result.ToBytes()
		expected := []byte{0x12, 0x34}
		if !bytesEqual(extractedBytes, expected) {
			t.Errorf("Expected extracted value %v, got %v", expected, extractedBytes)
		}
	})

	t.Run("validateExtractionBounds", func(t *testing.T) {
		bs := bitstringpkg.NewBitStringFromBytes([]byte{0x12, 0x34}) // 16 bits

		// Valid extraction
		err := validateExtractionBounds(bs, 0, 8)
		if err != nil {
			t.Errorf("Expected no error for valid extraction, got %v", err)
		}

		// Invalid extraction - beyond bounds
		err = validateExtractionBounds(bs, 8, 16) // Starting at bit 8, need 16 bits, only have 8
		if err == nil {
			t.Error("Expected error for extraction beyond bounds")
		}
	})

	t.Run("truncateBitstring", func(t *testing.T) {
		m := NewMatcher()

		bs := bitstringpkg.NewBitStringFromBytes([]byte{0x12, 0x34}) // 16 bits

		// Truncate to 8 bits
		result, err := m.truncateBitstring(bs, 8)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Length() != 8 {
			t.Errorf("Expected result with 8 bits, got %d bits", result.Length())
		}

		extractedBytes := result.ToBytes()
		if len(extractedBytes) != 1 || extractedBytes[0] != 0x12 {
			t.Errorf("Expected truncated value 0x12, got %v", extractedBytes)
		}
	})

	t.Run("truncateBitstring to zero", func(t *testing.T) {
		m := NewMatcher()

		bs := bitstringpkg.NewBitStringFromBytes([]byte{0x12, 0x34})
		result, err := m.truncateBitstring(bs, 0)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result.Length() != 0 {
			t.Errorf("Expected empty bitstring, got %d bits", result.Length())
		}
	})

	t.Run("extractBitsFromByte", func(t *testing.T) {
		m := NewMatcher()

		// Extract 4 bits from a byte
		result, err := m.extractBitsFromByte(0xF0, 4)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Should extract 1111 and left-align: 11110000
		expected := byte(0xF0)
		if result != expected {
			t.Errorf("Expected result 0x%02X, got 0x%02X", expected, result)
		}
	})

	t.Run("extractBitsFromByte too many bits", func(t *testing.T) {
		m := NewMatcher()

		_, err := m.extractBitsFromByte(0x12, 9)

		if err == nil {
			t.Error("Expected error for extracting more than 8 bits")
		}
	})

	t.Run("extractIntegerBits non-byte-aligned", func(t *testing.T) {
		// Test the package-level extractIntegerBits function
		data := []byte{0b11001100} // 0xCC

		// Extract 4 bits starting at position 2 (should get 0011)
		result, err := extractIntegerBits(data, 2, 4, false)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != 0x3 {
			t.Errorf("Expected result 0x3, got 0x%X", result)
		}
	})

	t.Run("extractIntegerBits signed negative", func(t *testing.T) {
		data := []byte{0b11110000} // 0xF0

		// Extract 4 bits starting at position 0 as signed (should get -8 in two's complement 4-bit)
		result, err := extractIntegerBits(data, 0, 4, true)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// In 4-bit two's complement, 1111 is -1
		if result != -1 {
			t.Errorf("Expected result -1, got %d", result)
		}
	})

	t.Run("extractUTF16 basic", func(t *testing.T) {
		m := NewMatcher()

		// Test UTF-16 BE encoding of 'A' (0x0041)
		data := []byte{0x00, 0x41}
		result, bytesConsumed, err := m.extractUTF16(data, "big")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != "A" {
			t.Errorf("Expected result 'A', got '%s'", result)
		}

		if bytesConsumed != 2 {
			t.Errorf("Expected 2 bytes consumed, got %d", bytesConsumed)
		}
	})

	t.Run("extractUTF32 basic", func(t *testing.T) {
		m := NewMatcher()

		// Test UTF-32 BE encoding of 'A' (0x00000041)
		data := []byte{0x00, 0x00, 0x00, 0x41}
		result, bytesConsumed, err := m.extractUTF32(data, "big")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if result != "A" {
			t.Errorf("Expected result 'A', got '%s'", result)
		}

		if bytesConsumed != 4 {
			t.Errorf("Expected 4 bytes consumed, got %d", bytesConsumed)
		}
	})

	t.Run("extractUTF8 invalid sequence", func(t *testing.T) {
		m := NewMatcher()

		// Invalid UTF-8 sequence (0xFF is not a valid start byte)
		data := []byte{0xFF}
		_, _, err := m.extractUTF8(data)

		if err == nil {
			t.Error("Expected error for invalid UTF-8 sequence")
		}
	})

	t.Run("extractUTF8 incomplete sequence", func(t *testing.T) {
		m := NewMatcher()

		// Incomplete 2-byte UTF-8 sequence
		data := []byte{0xC3} // Missing second byte
		_, _, err := m.extractUTF8(data)

		if err == nil {
			t.Error("Expected error for incomplete UTF-8 sequence")
		}
	})
}
