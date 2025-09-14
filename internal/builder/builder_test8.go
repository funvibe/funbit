package builder

import (
	"testing"

	"github.com/funvibe/funbit/internal/bitstring"
)

// TestBuilder_AddFloat_SizeSpecifiedFalse tests scenarios where SizeSpecified might be false
func TestBuilder_AddFloat_SizeSpecifiedFalse(t *testing.T) {
	t.Run("Float with custom type that results in SizeSpecified=false", func(t *testing.T) {
		b := NewBuilder()

		// Create a custom type that might result in SizeSpecified=false
		type CustomFloat struct {
			val float64
		}

		customValue := CustomFloat{}
		result := b.AddFloat(customValue)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeFloat {
			t.Error("Expected segment type to be float")
		}

		// Log the actual SizeSpecified value to understand behavior
		t.Logf("CustomFloat - Size: %d, SizeSpecified: %v", b.segments[0].Size, b.segments[0].SizeSpecified)

		// If SizeSpecified is false, then the default size should be set
		if !b.segments[0].SizeSpecified {
			if b.segments[0].Size != bitstring.DefaultSizeFloat {
				t.Errorf("Expected default size %d, got %d", bitstring.DefaultSizeFloat, b.segments[0].Size)
			}
		}
	})

	t.Run("Float with map value to test unusual type handling", func(t *testing.T) {
		b := NewBuilder()

		// Test with map value (unusual type that might result in SizeSpecified=false)
		value := map[string]interface{}{"value": 3.14}
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeFloat {
			t.Error("Expected segment type to be float")
		}

		// Log the actual behavior
		t.Logf("Map value - Size: %d, SizeSpecified: %v", b.segments[0].Size, b.segments[0].SizeSpecified)
	})

	t.Run("Float with slice value to test array type handling", func(t *testing.T) {
		b := NewBuilder()

		// Test with slice value
		value := []float64{3.14, 2.718}
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeFloat {
			t.Error("Expected segment type to be float")
		}

		// Log the actual behavior
		t.Logf("Slice value - Size: %d, SizeSpecified: %v", b.segments[0].Size, b.segments[0].SizeSpecified)
	})

	t.Run("Float with function value to test function type handling", func(t *testing.T) {
		b := NewBuilder()

		// Test with function value
		value := func() float64 { return 3.14 }
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeFloat {
			t.Error("Expected segment type to be float")
		}

		// Log the actual behavior
		t.Logf("Function value - Size: %d, SizeSpecified: %v", b.segments[0].Size, b.segments[0].SizeSpecified)
	})

	t.Run("Float with channel value to test channel type handling", func(t *testing.T) {
		b := NewBuilder()

		// Test with channel value
		value := make(chan float64)
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeFloat {
			t.Error("Expected segment type to be float")
		}

		// Log the actual behavior
		t.Logf("Channel value - Size: %d, SizeSpecified: %v", b.segments[0].Size, b.segments[0].SizeSpecified)
	})

	t.Run("Float with pointer to function to test pointer type handling", func(t *testing.T) {
		b := NewBuilder()

		// Test with pointer to function
		funcPtr := func() float64 { return 3.14 }
		value := &funcPtr
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeFloat {
			t.Error("Expected segment type to be float")
		}

		// Log the actual behavior
		t.Logf("Pointer to function - Size: %d, SizeSpecified: %v", b.segments[0].Size, b.segments[0].SizeSpecified)
	})

	t.Run("Float with array value to test array type handling", func(t *testing.T) {
		b := NewBuilder()

		// Test with array value
		value := [2]float64{3.14, 2.718}
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeFloat {
			t.Error("Expected segment type to be float")
		}

		// Log the actual behavior
		t.Logf("Array value - Size: %d, SizeSpecified: %v", b.segments[0].Size, b.segments[0].SizeSpecified)
	})
}

// TestBuilder_AddFloat_SizeAlreadySpecified tests the scenario where size is already specified in options
func TestBuilder_AddFloat_SizeAlreadySpecified(t *testing.T) {
	builder := NewBuilder()

	// Test case where size is already specified in options
	result := builder.AddFloat(3.14, bitstring.WithSize(64))

	// Verify that the segment was added
	if len(result.segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(result.segments))
	}

	segment := result.segments[0]
	if segment.Type != bitstring.TypeFloat {
		t.Errorf("Expected segment type to be float, got %s", segment.Type)
	}

	// When size is specified in options, SizeSpecified should be true
	// and Size should be the specified value (not default)
	if segment.Size != 64 {
		t.Errorf("Expected segment size to be 64, got %d", segment.Size)
	}

	if !segment.SizeSpecified {
		t.Errorf("Expected SizeSpecified to be true when size is provided in options")
	}

	// Test with different size
	builder2 := NewBuilder()
	result2 := builder2.AddFloat(2.71, bitstring.WithSize(32))

	if len(result2.segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(result2.segments))
	}

	segment2 := result2.segments[0]
	if segment2.Size != 32 {
		t.Errorf("Expected segment size to be 32, got %d", segment2.Size)
	}

	if !segment2.SizeSpecified {
		t.Errorf("Expected SizeSpecified to be true when size is provided in options")
	}
}

// TestBuilder_AddFloat_SizeNotSpecified tests the case where size is not specified
func TestBuilder_AddFloat_SizeNotSpecified(t *testing.T) {
	t.Run("Float32 without size specified", func(t *testing.T) {
		b := NewBuilder()
		value := float32(3.14159)
		result := b.AddFloat(value) // No size specified

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		if len(b.segments) != 1 {
			t.Fatalf("Expected 1 segment, got %d", len(b.segments))
		}

		segment := b.segments[0]
		if segment.Type != bitstring.TypeFloat {
			t.Errorf("Expected segment type 'float', got '%s'", segment.Type)
		}

		// When size is not specified, it should use default size
		if segment.Size != 32 {
			t.Errorf("Expected segment size 32, got %d", segment.Size)
		}

		// SizeSpecified should be false when using default
		if segment.SizeSpecified {
			t.Error("Expected SizeSpecified to be false when using default size")
		}
	})

	t.Run("Float64 without size specified", func(t *testing.T) {
		b := NewBuilder()
		value := float64(2.718281828459045)
		result := b.AddFloat(value) // No size specified

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Type != bitstring.TypeFloat {
			t.Errorf("Expected segment type 'float', got '%s'", segment.Type)
		}

		// When size is not specified, it should use default size for float64
		if segment.Size != 64 {
			t.Errorf("Expected segment size 64, got %d", segment.Size)
		}

		// SizeSpecified should be false when using default
		if segment.SizeSpecified {
			t.Error("Expected SizeSpecified to be false when using default size")
		}
	})

	t.Run("Float with interface value without size specified", func(t *testing.T) {
		b := NewBuilder()
		var value interface{} = float32(1.618)
		result := b.AddFloat(value) // No size specified

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Type != bitstring.TypeFloat {
			t.Errorf("Expected segment type 'float', got '%s'", segment.Type)
		}

		// When size is not specified, it should use default size for float32
		if segment.Size != 32 {
			t.Errorf("Expected segment size 32, got %d", segment.Size)
		}

		// SizeSpecified should be false when using default
		if segment.SizeSpecified {
			t.Error("Expected SizeSpecified to be false when using default size")
		}
	})
}

// TestBuilder_AddInteger_MissingReflectionPaths tests specific reflection paths that may be missing coverage
func TestBuilder_AddInteger_MissingReflectionPaths(t *testing.T) {
	t.Run("Integer with complex interface type", func(t *testing.T) {
		b := NewBuilder()

		// Use the existing CustomInt interface and MyInt struct defined at package level
		var customInt CustomInt = MyInt{val: 42}
		result := b.AddInteger(customInt)

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeInteger {
			t.Error("Expected segment type to be integer")
		}

		// Log the actual behavior to understand reflection path
		t.Logf("Custom interface - Size: %d, SizeSpecified: %v, Signed: %v",
			b.segments[0].Size, b.segments[0].SizeSpecified, b.segments[0].Signed)
	})

	t.Run("Integer with pointer to interface", func(t *testing.T) {
		b := NewBuilder()

		// Test with pointer to interface
		type Number interface{}
		var num Number = 42
		result := b.AddInteger(&num)

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeInteger {
			t.Error("Expected segment type to be integer")
		}

		// Log the actual behavior
		t.Logf("Pointer to interface - Size: %d, SizeSpecified: %v, Signed: %v",
			b.segments[0].Size, b.segments[0].SizeSpecified, b.segments[0].Signed)
	})

	t.Run("Integer with nil interface", func(t *testing.T) {
		b := NewBuilder()

		// Test with nil interface
		var num interface{} = nil
		result := b.AddInteger(num)

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		// Verify segment was added
		if len(b.segments) != 1 {
			t.Error("Expected 1 segment to be added")
		}
		if b.segments[0].Type != bitstring.TypeInteger {
			t.Error("Expected segment type to be integer")
		}

		// Log the actual behavior for nil interface
		t.Logf("Nil interface - Size: %d, SizeSpecified: %v, Signed: %v",
			b.segments[0].Size, b.segments[0].SizeSpecified, b.segments[0].Signed)
	})
}

// TestBuilder_Build_MissingAlignmentCases tests specific alignment scenarios that may be missing coverage
func TestBuilder_Build_MissingAlignmentCases(t *testing.T) {
	t.Run("Build with segment that causes validation error before alignment", func(t *testing.T) {
		b := NewBuilder()
		// Add a segment that will fail validation
		b.AddInteger(42, bitstring.WithSize(0)) // Invalid size

		_, err := b.Build()
		if err == nil {
			t.Error("Expected error for invalid segment")
		}

		// Should be a BitStringError
		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeInvalidSize {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeInvalidSize, bitStringErr.Code)
			}
		}
	})

	t.Run("Build with mixed segment types and alignment", func(t *testing.T) {
		b := NewBuilder()
		// Test complex alignment scenario with different segment types
		b.AddInteger(0b101, bitstring.WithSize(3))   // 3 bits
		b.AddBinary([]byte{0xAB})                    // Should trigger alignment
		b.AddInteger(0x1234, bitstring.WithSize(16)) // 16 bits

		bs, err := b.Build()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if bs == nil {
			t.Fatal("Expected non-nil bitstring")
		}

		// Should be 3 + 5 (padding) + 8 + 16 = 32 bits
		if bs.Length() == 0 {
			t.Error("Expected non-empty bitstring")
		}
		t.Logf("Mixed alignment bitstring length: %d", bs.Length())
	})

	t.Run("Build with exact byte boundary alignment", func(t *testing.T) {
		b := NewBuilder()
		// Test case where total bits exactly align to byte boundary
		b.AddInteger(0b1, bitstring.WithSize(1))       // 1 bit
		b.AddInteger(0b1111111, bitstring.WithSize(7)) // 7 bits = 1 byte total
		b.AddInteger(0xFF, bitstring.WithSize(8))      // 8 bits = 1 byte

		bs, err := b.Build()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Should be 1 + 7 + 8 = 16 bits (exactly 2 bytes, no padding needed)
		if bs.Length() != 16 {
			t.Errorf("Expected bitstring length 16, got %d", bs.Length())
		}
	})

	t.Run("Build with multiple alignment scenarios", func(t *testing.T) {
		b := NewBuilder()
		// Test multiple alignment scenarios in sequence
		b.AddInteger(0b1, bitstring.WithSize(1))      // 1 bit
		b.AddInteger(0b1, bitstring.WithSize(1))      // 1 bit (total 2 bits)
		b.AddInteger(0b111111, bitstring.WithSize(6)) // 6 bits (total 8 bits = 1 byte)
		b.AddInteger(0xAB, bitstring.WithSize(8))     // 8 bits (should align automatically)

		bs, err := b.Build()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Should be 1 + 1 + 6 + 8 = 16 bits
		if bs.Length() != 16 {
			t.Errorf("Expected bitstring length 16, got %d", bs.Length())
		}
	})
}

// TestBuilder_encodeSegment_MissingErrorPaths tests error paths in encodeSegment that may be missing coverage
func TestBuilder_encodeSegment_MissingErrorPaths(t *testing.T) {
	t.Run("Encode segment with nil value", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         nil,
			Type:          bitstring.TypeInteger,
			Size:          8,
			SizeSpecified: true,
		}

		err := encodeSegment(w, segment)
		if err == nil {
			t.Error("Expected error for nil value")
		}
		t.Logf("Expected error for nil value: %v", err)
	})

	t.Run("Encode segment with invalid UTF type", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:      65,
			Type:       "utf64", // Invalid UTF type
			Endianness: "big",
		}

		err := encodeSegment(w, segment)
		if err == nil {
			t.Error("Expected error for invalid UTF type")
		}

		if err.Error() != "unsupported segment type: utf64" {
			t.Errorf("Expected 'unsupported segment type: utf64', got %v", err)
		}
	})

	t.Run("Encode segment with float type and invalid size", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         float32(3.14),
			Type:          bitstring.TypeFloat,
			Size:          16, // Invalid float size
			SizeSpecified: true,
		}

		err := encodeSegment(w, segment)
		if err == nil {
			t.Error("Expected error for invalid float size")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeInvalidFloatSize {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeInvalidFloatSize, bitStringErr.Code)
			}
		}
	})

	t.Run("Encode segment with binary type and size not specified", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         []byte{0xAB},
			Type:          bitstring.TypeBinary,
			SizeSpecified: false, // Size not specified for binary
		}

		err := encodeSegment(w, segment)
		if err == nil {
			t.Error("Expected error for binary with size not specified")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeBinarySizeRequired {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeBinarySizeRequired, bitStringErr.Code)
			}
		}
	})
}

// TestBuilder_encodeInteger_MissingBitstringPaths tests bitstring type paths in encodeInteger
func TestBuilder_encodeInteger_MissingBitstringPaths(t *testing.T) {
	t.Run("Encode integer with bitstring type and exact size match", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         int64(0xAB),
			Type:          bitstring.TypeBitstring,
			Size:          8, // Exact match for integer value
			SizeSpecified: true,
		}

		err := encodeInteger(w, segment)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		data, totalBits := w.final()
		if totalBits != 8 {
			t.Errorf("Expected totalBits 8, got %d", totalBits)
		}

		if len(data) != 1 || data[0] != 0xAB {
			t.Errorf("Expected byte [0xAB], got %v", data)
		}
	})

	t.Run("Encode integer with bitstring type and slice value", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         []byte{0xAB, 0xCD},
			Type:          bitstring.TypeBitstring,
			Size:          16, // Exact match for slice length
			SizeSpecified: true,
		}

		err := encodeInteger(w, segment)
		if err == nil {
			t.Error("Expected error for unsupported integer type, got nil")
		} else if err.Error() != "unsupported integer type for bitstring value: []uint8" {
			t.Errorf("Expected 'unsupported integer type for bitstring value: []uint8', got %v", err)
		}

		// Verify that no data was written due to the error
		data, totalBits := w.final()
		if totalBits != 0 {
			t.Errorf("Expected totalBits 0 due to error, got %d", totalBits)
		}

		if len(data) != 0 {
			t.Errorf("Expected 0 bytes due to error, got %d", len(data))
		}
	})

	t.Run("Encode integer with bitstring type and insufficient slice data", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         []byte{0xAB}, // Only 1 byte (8 bits)
			Type:          bitstring.TypeBitstring,
			Size:          16, // Requesting 16 bits
			SizeSpecified: true,
		}

		err := encodeInteger(w, segment)
		if err == nil {
			t.Error("Expected error for insufficient bits")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeInsufficientBits {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeInsufficientBits, bitStringErr.Code)
			}
		}
	})

	t.Run("Encode integer with bitstring type and non-byte slice", func(t *testing.T) {
		w := newBitWriter()
		// Create a slice that's not []byte to test different reflection path
		intSlice := []int{0xAB, 0xCD}
		segment := &bitstring.Segment{
			Value:         intSlice,
			Type:          bitstring.TypeBitstring,
			Size:          16,
			SizeSpecified: true,
		}

		err := encodeInteger(w, segment)
		// This should fail because it's not a []byte slice
		if err == nil {
			t.Error("Expected error for non-byte slice")
		}
		t.Logf("Expected error for non-byte slice: %v", err)
	})
}

// TestBuilder_encodeBinary_MissingSizeValidation tests size validation edge cases in encodeBinary
func TestBuilder_encodeBinary_MissingSizeValidation(t *testing.T) {
	t.Run("Encode binary with size specified as zero", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         []byte{0xAB},
			Type:          bitstring.TypeBinary,
			Size:          0, // Zero size
			SizeSpecified: true,
		}

		err := encodeBinary(w, segment)
		if err == nil {
			t.Error("Expected error for zero size")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeInvalidSize {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeInvalidSize, bitStringErr.Code)
			}
		}
	})

	t.Run("Encode binary with size larger than data", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         []byte{0xAB}, // 1 byte
			Type:          bitstring.TypeBinary,
			Size:          2, // Requesting 2 bytes
			SizeSpecified: true,
		}

		err := encodeBinary(w, segment)
		if err == nil {
			t.Error("Expected error for size mismatch")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeBinarySizeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeBinarySizeMismatch, bitStringErr.Code)
			}
		}
	})

	t.Run("Encode binary with size smaller than data", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         []byte{0xAB, 0xCD}, // 2 bytes
			Type:          bitstring.TypeBinary,
			Size:          1, // Requesting 1 byte
			SizeSpecified: true,
		}

		err := encodeBinary(w, segment)
		if err == nil {
			t.Error("Expected error for size mismatch")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeBinarySizeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeBinarySizeMismatch, bitStringErr.Code)
			}
		}
	})

	t.Run("Encode binary with empty data and size specified", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         []byte{}, // Empty data
			Type:          bitstring.TypeBinary,
			Size:          1, // Size specified but no data
			SizeSpecified: true,
		}

		err := encodeBinary(w, segment)
		if err == nil {
			t.Error("Expected error for size mismatch with empty data")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeBinarySizeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeBinarySizeMismatch, bitStringErr.Code)
			}
		}
	})
}

// TestBuilder_validateBitstringValue_MissingCases tests missing cases in validateBitstringValue
func TestBuilder_validateBitstringValue_MissingCases(t *testing.T) {
	t.Run("Validate bitstring with non-bitstring pointer", func(t *testing.T) {
		segment := &bitstring.Segment{
			Value:         "not a bitstring",
			Type:          bitstring.TypeBitstring,
			Size:          8,
			SizeSpecified: true,
		}

		_, err := validateBitstringValue(segment)
		if err == nil {
			t.Error("Expected error for non-bitstring value")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeTypeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeTypeMismatch, bitStringErr.Code)
			}
		}
	})

	t.Run("Validate bitstring with nil pointer", func(t *testing.T) {
		var bs *bitstring.BitString = nil
		segment := &bitstring.Segment{
			Value:         bs,
			Type:          bitstring.TypeBitstring,
			Size:          8,
			SizeSpecified: true,
		}

		_, err := validateBitstringValue(segment)
		if err == nil {
			t.Error("Expected error for nil bitstring pointer")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeInvalidSegment {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeInvalidSegment, bitStringErr.Code)
			}
		}
	})

	t.Run("Validate bitstring with integer value", func(t *testing.T) {
		segment := &bitstring.Segment{
			Value:         42,
			Type:          bitstring.TypeBitstring,
			Size:          8,
			SizeSpecified: true,
		}

		_, err := validateBitstringValue(segment)
		if err == nil {
			t.Error("Expected error for integer value")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeTypeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeTypeMismatch, bitStringErr.Code)
			}
		}
	})

	t.Run("Validate bitstring with slice value", func(t *testing.T) {
		segment := &bitstring.Segment{
			Value:         []byte{0xAB},
			Type:          bitstring.TypeBitstring,
			Size:          8,
			SizeSpecified: true,
		}

		_, err := validateBitstringValue(segment)
		if err == nil {
			t.Error("Expected error for slice value")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeTypeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeTypeMismatch, bitStringErr.Code)
			}
		}
	})
}

// TestBuilder_writeBitstringBits_MissingBoundaryConditions tests boundary conditions in writeBitstringBits
func TestBuilder_writeBitstringBits_MissingBoundaryConditions(t *testing.T) {
	t.Run("Write bitstring bits with size exactly matching bitstring length", func(t *testing.T) {
		w := newBitWriter()
		bs := bitstring.NewBitStringFromBits([]byte{0xAB, 0xCD}, 16)

		err := writeBitstringBits(w, bs, 16) // Exact match
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		data, totalBits := w.final()
		if totalBits != 16 {
			t.Errorf("Expected totalBits 16, got %d", totalBits)
		}

		if len(data) != 2 || data[0] != 0xAB || data[1] != 0xCD {
			t.Errorf("Expected bytes [0xAB, 0xCD], got %v", data)
		}
	})

	t.Run("Write bitstring bits with size one less than bitstring length", func(t *testing.T) {
		w := newBitWriter()
		bs := bitstring.NewBitStringFromBits([]byte{0xAB, 0xCD}, 16)

		err := writeBitstringBits(w, bs, 15) // One less than available
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		data, totalBits := w.final()
		if totalBits != 15 {
			t.Errorf("Expected totalBits 15, got %d", totalBits)
		}

		// Should have first 15 bits of 0xABCD
		if len(data) != 2 {
			t.Errorf("Expected 2 bytes, got %d", len(data))
		}
	})

	t.Run("Write bitstring bits with single byte boundary", func(t *testing.T) {
		w := newBitWriter()
		bs := bitstring.NewBitStringFromBits([]byte{0xAB, 0xCD}, 16)

		err := writeBitstringBits(w, bs, 8) // Exactly one byte
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		data, totalBits := w.final()
		if totalBits != 8 {
			t.Errorf("Expected totalBits 8, got %d", totalBits)
		}

		if len(data) != 1 || data[0] != 0xAB {
			t.Errorf("Expected byte [0xAB], got %v", data)
		}
	})

	t.Run("Write bitstring bits with size crossing byte boundary", func(t *testing.T) {
		w := newBitWriter()
		bs := bitstring.NewBitStringFromBits([]byte{0xAB, 0xCD}, 16)

		err := writeBitstringBits(w, bs, 12) // Crosses byte boundary
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		data, totalBits := w.final()
		if totalBits != 12 {
			t.Errorf("Expected totalBits 12, got %d", totalBits)
		}

		if len(data) != 2 {
			t.Errorf("Expected 2 bytes, got %d", len(data))
		}
	})
}

// TestBuilder_encodeFloat_MissingNativeEndianness tests native endianness paths in encodeFloat
func TestBuilder_encodeFloat_MissingNativeEndianness(t *testing.T) {
	t.Run("Encode float32 with native endianness", func(t *testing.T) {
		w := newBitWriter()
		value := float32(3.14159)
		segment := &bitstring.Segment{
			Value:         value,
			Type:          bitstring.TypeFloat,
			Size:          32,
			SizeSpecified: true,
			Endianness:    bitstring.EndiannessNative,
		}

		err := encodeFloat(w, segment)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		data, totalBits := w.final()
		if totalBits != 32 {
			t.Errorf("Expected totalBits 32, got %d", totalBits)
		}

		if len(data) != 4 {
			t.Errorf("Expected 4 bytes, got %d", len(data))
		}

		// The actual byte order depends on the native endianness
		t.Logf("Native endianness float32 result: %v", data)
	})

	t.Run("Encode float64 with native endianness", func(t *testing.T) {
		w := newBitWriter()
		value := float64(2.718281828459045)
		segment := &bitstring.Segment{
			Value:         value,
			Type:          bitstring.TypeFloat,
			Size:          64,
			SizeSpecified: true,
			Endianness:    bitstring.EndiannessNative,
		}

		err := encodeFloat(w, segment)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		data, totalBits := w.final()
		if totalBits != 64 {
			t.Errorf("Expected totalBits 64, got %d", totalBits)
		}

		if len(data) != 8 {
			t.Errorf("Expected 8 bytes, got %d", len(data))
		}

		// The actual byte order depends on the native endianness
		t.Logf("Native endianness float64 result: %v", data)
	})

	t.Run("Encode float with int value type conversion", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         42, // int value
			Type:          bitstring.TypeFloat,
			Size:          32,
			SizeSpecified: true,
		}

		err := encodeFloat(w, segment)
		// This should fail because int is not a supported float type
		if err == nil {
			t.Error("Expected error for int value type")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeTypeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeTypeMismatch, bitStringErr.Code)
			}
		}
	})

	t.Run("Encode float with string value type", func(t *testing.T) {
		w := newBitWriter()
		segment := &bitstring.Segment{
			Value:         "3.14", // string value
			Type:          bitstring.TypeFloat,
			Size:          32,
			SizeSpecified: true,
		}

		err := encodeFloat(w, segment)
		// This should fail because string is not a supported float type
		if err == nil {
			t.Error("Expected error for string value type")
		}

		if bitStringErr, ok := err.(*bitstring.BitStringError); ok {
			if bitStringErr.Code != bitstring.CodeTypeMismatch {
				t.Errorf("Expected error code %s, got %s", bitstring.CodeTypeMismatch, bitStringErr.Code)
			}
		}
	})
}
