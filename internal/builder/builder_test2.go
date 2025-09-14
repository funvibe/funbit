package builder

import (
	"bytes"
	"testing"

	"github.com/funvibe/funbit/internal/bitstring"
)

func TestBuilder_AddInteger(t *testing.T) {
	t.Run("Chaining and basic functionality", func(t *testing.T) {
		b := NewBuilder()

		// Test that AddInteger returns the builder for chaining
		result := b.AddInteger(42, bitstring.WithSize(8))
		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		// Test multiple additions for chaining
		b2 := b.
			AddInteger(1, bitstring.WithSize(8)).
			AddInteger(17, bitstring.WithSize(8)).
			AddInteger(42, bitstring.WithSize(8))

		if b2 != b {
			t.Error("Expected chaining to work correctly")
		}
	})

	t.Run("Auto-detect signedness for negative values", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(-42, bitstring.WithSize(8))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if !segment.Signed {
			t.Error("Expected segment to be marked as signed for negative value")
		}
	})

	t.Run("Auto-detect signedness for positive values", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(42, bitstring.WithSize(8))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Signed {
			t.Error("Expected segment to be marked as unsigned for positive value")
		}
	})

	t.Run("Explicit signedness", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(42, bitstring.WithSize(8), bitstring.WithSigned(true))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if !segment.Signed {
			t.Error("Expected segment to be marked as signed when explicitly set")
		}
	})

	t.Run("With type specified", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(42, bitstring.WithSize(8), bitstring.WithType("custom"))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Type != "custom" {
			t.Errorf("Expected segment type 'custom', got '%s'", segment.Type)
		}
	})
}

func TestBuilder_AddInteger_AdditionalCoverage(t *testing.T) {
	t.Run("Integer with explicit size and type", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(42, bitstring.WithSize(16), bitstring.WithType("custom"))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		if len(b.segments) != 1 {
			t.Errorf("Expected 1 segment, got %d", len(b.segments))
		}

		segment := b.segments[0]
		if segment.Type != "custom" {
			t.Errorf("Expected segment type 'custom', got '%s'", segment.Type)
		}

		if segment.Size != 16 {
			t.Errorf("Expected segment size 16, got %d", segment.Size)
		}

		if !segment.SizeSpecified {
			t.Error("Expected SizeSpecified to be true")
		}
	})

	t.Run("Integer with explicit signedness false", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(42, bitstring.WithSize(8), bitstring.WithSigned(false))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Signed {
			t.Error("Expected segment to be marked as unsigned when explicitly set to false")
		}
	})

	t.Run("Integer with explicit signedness true", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(42, bitstring.WithSize(8), bitstring.WithSigned(true))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if !segment.Signed {
			t.Error("Expected segment to be marked as signed when explicitly set to true")
		}
	})

	t.Run("Integer with endianness", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(0xABCD, bitstring.WithSize(16), bitstring.WithEndianness(bitstring.EndiannessLittle))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Endianness != bitstring.EndiannessLittle {
			t.Errorf("Expected segment endianness %s, got %s", bitstring.EndiannessLittle, segment.Endianness)
		}
	})

	t.Run("Integer with unit", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(42, bitstring.WithSize(8), bitstring.WithUnit(16))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Unit != 16 {
			t.Errorf("Expected segment unit 16, got %d", segment.Unit)
		}
	})

	t.Run("Integer with all options", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(-42,
			bitstring.WithSize(16),
			bitstring.WithType("custom_int"),
			bitstring.WithSigned(true),
			bitstring.WithEndianness(bitstring.EndiannessLittle),
			bitstring.WithUnit(32))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Type != "custom_int" {
			t.Errorf("Expected segment type 'custom_int', got '%s'", segment.Type)
		}
		if segment.Size != 16 {
			t.Errorf("Expected segment size 16, got %d", segment.Size)
		}
		if !segment.Signed {
			t.Error("Expected segment to be marked as signed")
		}
		if segment.Endianness != bitstring.EndiannessLittle {
			t.Errorf("Expected segment endianness %s, got %s", bitstring.EndiannessLittle, segment.Endianness)
		}
		if segment.Unit != 32 {
			t.Errorf("Expected segment unit 32, got %d", segment.Unit)
		}
		if !segment.SizeSpecified {
			t.Error("Expected SizeSpecified to be true")
		}
	})

	t.Run("Zero value integer", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(0, bitstring.WithSize(8))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Signed {
			t.Error("Expected zero value to be marked as unsigned")
		}
	})

	t.Run("Uint8 value", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(uint8(255), bitstring.WithSize(8))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Signed {
			t.Error("Expected uint8 value to be marked as unsigned")
		}
	})

	t.Run("Int8 negative value", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(int8(-128), bitstring.WithSize(8))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if !segment.Signed {
			t.Error("Expected negative int8 value to be marked as signed")
		}
	})

	t.Run("Int8 positive value", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddInteger(int8(127), bitstring.WithSize(8))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Signed {
			t.Error("Expected positive int8 value to be marked as unsigned")
		}
	})

	t.Run("Large uint64 value", func(t *testing.T) {
		b := NewBuilder()
		largeValue := uint64(0xFFFFFFFFFFFFFFFF)
		result := b.AddInteger(largeValue, bitstring.WithSize(64))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Signed {
			t.Error("Expected uint64 value to be marked as unsigned")
		}
		if segment.Size != 64 {
			t.Errorf("Expected segment size 64, got %d", segment.Size)
		}
	})

	t.Run("Large int64 value", func(t *testing.T) {
		b := NewBuilder()
		largeValue := int64(-9223372036854775808)
		result := b.AddInteger(largeValue, bitstring.WithSize(64))

		if result != b {
			t.Error("Expected AddInteger() to return the same builder instance")
		}

		segment := b.segments[0]
		if !segment.Signed {
			t.Error("Expected negative int64 value to be marked as signed")
		}
		if segment.Size != 64 {
			t.Errorf("Expected segment size 64, got %d", segment.Size)
		}
	})
}

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Builder)
		wantErr bool
	}{
		{
			name: "empty builder",
			setup: func(b *Builder) {
				// No additions
			},
			wantErr: false,
		},
		{
			name: "single integer",
			setup: func(b *Builder) {
				b.AddInteger(42, bitstring.WithSize(8))
			},
			wantErr: false,
		},
		{
			name: "multiple integers",
			setup: func(b *Builder) {
				b.AddInteger(1, bitstring.WithSize(8)).AddInteger(17, bitstring.WithSize(8)).AddInteger(42, bitstring.WithSize(8))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()
			tt.setup(b)

			bs, err := b.Build()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected Build() to return an error")
				}
				if bs != nil {
					t.Error("Expected Build() to return nil bitstring on error")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if bs == nil {
				t.Fatal("Expected Build() to return non-nil bitstring")
			}

			// Basic validation of created bitstring
			if bs.Length() == 0 && tt.name != "empty builder" {
				t.Errorf("Expected non-zero length for %s", tt.name)
			}

			// Should be binary since we're adding integers (default 8 bits each)
			if tt.name != "empty builder" && !bs.IsBinary() {
				t.Error("Expected bitstring to be binary")
			}
		})
	}
}

func TestBuilder_BuildContent(t *testing.T) {
	// Test specific content generation
	b := NewBuilder().
		AddInteger(1, bitstring.WithSize(8)).
		AddInteger(17, bitstring.WithSize(8)).
		AddInteger(42, bitstring.WithSize(8))

	bs, err := b.Build()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bs.Length() != 24 { // 3 integers * 8 bits = 24 bits
		t.Errorf("Expected bitstring length 24, got %d", bs.Length())
	}

	if !bs.IsBinary() {
		t.Error("Expected bitstring to be binary")
	}

	bytes := bs.ToBytes()
	if len(bytes) != 3 {
		t.Fatalf("Expected 3 bytes, got %d", len(bytes))
	}

	if bytes[0] != 1 || bytes[1] != 17 || bytes[2] != 42 {
		t.Errorf("Expected [1, 17, 42], got %v", bytes)
	}
}

func TestBuilder_EmptyBuild(t *testing.T) {
	b := NewBuilder()
	bs, err := b.Build()

	if err != nil {
		t.Errorf("Expected no error for empty build, got %v", err)
	}

	if bs == nil {
		t.Fatal("Expected non-nil bitstring for empty build")
	}

	if bs.Length() != 0 {
		t.Errorf("Expected empty bitstring length 0, got %d", bs.Length())
	}

	if !bs.IsEmpty() {
		t.Error("Expected empty bitstring to be empty")
	}

	if !bs.IsBinary() {
		t.Error("Expected empty bitstring to be binary")
	}

	bytes := bs.ToBytes()
	if len(bytes) != 0 {
		t.Errorf("Expected empty bytes, got %v", bytes)
	}
}

// Test Builder methods with zero coverage
func TestBuilder_AddBinary(t *testing.T) {
	t.Run("Valid binary data", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddBinary([]byte{0xAB, 0xCD})

		if result != b {
			t.Error("Expected AddBinary() to return the same builder instance")
		}

		if len(b.segments) != 1 {
			t.Errorf("Expected 1 segment, got %d", len(b.segments))
		}

		segment := b.segments[0]
		if segment.Type != bitstring.TypeBinary {
			t.Errorf("Expected segment type %s, got %s", bitstring.TypeBinary, segment.Type)
		}

		if segment.Size != 2 {
			t.Errorf("Expected segment size 2, got %d", segment.Size)
		}

		if !segment.SizeSpecified {
			t.Error("Expected SizeSpecified to be true")
		}
	})

	t.Run("Empty binary data", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddBinary([]byte{})

		if result != b {
			t.Error("Expected AddBinary() to return the same builder instance")
		}

		if len(b.segments) != 0 {
			t.Errorf("Expected 0 segments for empty data, got %d", len(b.segments))
		}
	})

	t.Run("Binary data with explicit size", func(t *testing.T) {
		b := NewBuilder()
		data := []byte{0xAB, 0xCD, 0xEF}
		result := b.AddBinary(data, bitstring.WithSize(3))

		if result != b {
			t.Error("Expected AddBinary() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Size != 3 {
			t.Errorf("Expected segment size 3, got %d", segment.Size)
		}
	})

	t.Run("Binary data with unit", func(t *testing.T) {
		b := NewBuilder()
		data := []byte{0xAB}
		result := b.AddBinary(data, bitstring.WithUnit(16))

		if result != b {
			t.Error("Expected AddBinary() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Unit != 16 {
			t.Errorf("Expected segment unit 16, got %d", segment.Unit)
		}
	})
}

func TestBuilder_AddFloat(t *testing.T) {
	t.Run("Valid float32", func(t *testing.T) {
		b := NewBuilder()
		value := float32(3.14)
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		if len(b.segments) != 1 {
			t.Errorf("Expected 1 segment, got %d", len(b.segments))
		}

		segment := b.segments[0]
		if segment.Type != bitstring.TypeFloat {
			t.Errorf("Expected segment type %s, got %s", bitstring.TypeFloat, segment.Type)
		}

		if segment.Size != 32 { // Default size for float32 is 32 bits when not explicitly set
			t.Errorf("Expected segment size 32, got %d", segment.Size)
		}
	})

	t.Run("Valid float64", func(t *testing.T) {
		b := NewBuilder()
		value := float64(2.71828)
		result := b.AddFloat(value)

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Value != value {
			t.Errorf("Expected segment value %v, got %v", value, segment.Value)
		}
	})

	t.Run("Float with explicit size", func(t *testing.T) {
		b := NewBuilder()
		value := float32(1.618)
		result := b.AddFloat(value, bitstring.WithSize(32))

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Size != 32 {
			t.Errorf("Expected segment size 32, got %d", segment.Size)
		}
	})

	t.Run("Float with endianness", func(t *testing.T) {
		b := NewBuilder()
		value := float64(123.456)
		result := b.AddFloat(value, bitstring.WithEndianness(bitstring.EndiannessLittle))

		if result != b {
			t.Error("Expected AddFloat() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Endianness != bitstring.EndiannessLittle {
			t.Errorf("Expected segment endianness %s, got %s", bitstring.EndiannessLittle, segment.Endianness)
		}
	})
}

func TestBuilder_AddSegment(t *testing.T) {
	t.Run("Valid segment", func(t *testing.T) {
		b := NewBuilder()
		segment := bitstring.NewSegment(42, bitstring.WithSize(8), bitstring.WithType(bitstring.TypeInteger))
		result := b.AddSegment(*segment)

		if result != b {
			t.Error("Expected AddSegment() to return the same builder instance")
		}

		if len(b.segments) != 1 {
			t.Errorf("Expected 1 segment, got %d", len(b.segments))
		}

		addedSegment := b.segments[0]
		if addedSegment.Type != segment.Type {
			t.Errorf("Expected segment type %s, got %s", segment.Type, addedSegment.Type)
		}

		if addedSegment.Value != segment.Value {
			t.Errorf("Expected segment value %v, got %v", segment.Value, addedSegment.Value)
		}
	})

	t.Run("Binary segment with size", func(t *testing.T) {
		b := NewBuilder()
		segment := bitstring.NewSegment([]byte{0xAB}, bitstring.WithType(bitstring.TypeBinary))
		segment.Size = 1
		result := b.AddSegment(*segment)

		if result != b {
			t.Error("Expected AddSegment() to return the same builder instance")
		}

		addedSegment := b.segments[0]
		if !addedSegment.SizeSpecified {
			t.Error("Expected SizeSpecified to be true for binary segment with size > 0")
		}
	})

	t.Run("UTF segment", func(t *testing.T) {
		b := NewBuilder()
		segment := bitstring.NewSegment(65, bitstring.WithType("utf8"))
		result := b.AddSegment(*segment)

		if result != b {
			t.Error("Expected AddSegment() to return the same builder instance")
		}

		addedSegment := b.segments[0]
		if addedSegment.SizeSpecified {
			t.Error("Expected SizeSpecified to be false for UTF segment")
		}
	})

	t.Run("Non-UTF segment with zero size", func(t *testing.T) {
		b := NewBuilder()
		segment := bitstring.NewSegment(0, bitstring.WithType(bitstring.TypeInteger))
		segment.Size = 0
		result := b.AddSegment(*segment)

		if result != b {
			t.Error("Expected AddSegment() to return the same builder instance")
		}

		addedSegment := b.segments[0]
		if !addedSegment.SizeSpecified {
			t.Error("Expected SizeSpecified to be true for non-UTF segment with explicit size")
		}
	})
}

func TestBuilder_AddBitstring(t *testing.T) {
	t.Run("Valid bitstring", func(t *testing.T) {
		b := NewBuilder()
		bs := bitstring.NewBitStringFromBits([]byte{0xAB}, 8)
		result := b.AddBitstring(bs)

		if result != b {
			t.Error("Expected AddBitstring() to return the same builder instance")
		}

		if len(b.segments) != 1 {
			t.Errorf("Expected 1 segment, got %d", len(b.segments))
		}

		segment := b.segments[0]
		if segment.Type != bitstring.TypeBitstring {
			t.Errorf("Expected segment type %s, got %s", bitstring.TypeBitstring, segment.Type)
		}

		if segment.Unit != 1 {
			t.Errorf("Expected segment unit 1, got %d", segment.Unit)
		}

		if segment.Size != 8 {
			t.Errorf("Expected segment size 8, got %d", segment.Size)
		}
	})

	t.Run("Nil bitstring", func(t *testing.T) {
		b := NewBuilder()
		result := b.AddBitstring(nil)

		if result != b {
			t.Error("Expected AddBitstring() to return the same builder instance")
		}

		if len(b.segments) != 0 {
			t.Errorf("Expected 0 segments for nil bitstring, got %d", len(b.segments))
		}
	})

	t.Run("Bitstring with explicit size", func(t *testing.T) {
		b := NewBuilder()
		bs := bitstring.NewBitStringFromBits([]byte{0xAB, 0xCD}, 16)
		result := b.AddBitstring(bs, bitstring.WithSize(16))

		if result != b {
			t.Error("Expected AddBitstring() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Size != 16 {
			t.Errorf("Expected segment size 16, got %d", segment.Size)
		}
	})

	t.Run("Bitstring with unit", func(t *testing.T) {
		b := NewBuilder()
		bs := bitstring.NewBitStringFromBits([]byte{0xAB}, 8)
		result := b.AddBitstring(bs, bitstring.WithUnit(2))

		if result != b {
			t.Error("Expected AddBitstring() to return the same builder instance")
		}

		segment := b.segments[0]
		if segment.Unit != 2 {
			t.Errorf("Expected segment unit 2, got %d", segment.Unit)
		}
	})
}

// Test encode functions and helper functions
func TestBuilder_alignToByte(t *testing.T) {
	t.Run("No alignment needed", func(t *testing.T) {
		w := newBitWriter()
		w.alignToByte()

		if w.bitCount != 0 {
			t.Errorf("Expected bitCount 0, got %d", w.bitCount)
		}

		if w.acc != 0 {
			t.Errorf("Expected acc 0, got %d", w.acc)
		}
	})

	t.Run("Alignment needed", func(t *testing.T) {
		w := newBitWriter()
		w.writeBits(0b101, 3) // Write 3 bits
		w.alignToByte()

		if w.bitCount != 0 {
			t.Errorf("Expected bitCount 0 after alignment, got %d", w.bitCount)
		}

		if w.acc != 0 {
			t.Errorf("Expected acc 0 after alignment, got %d", w.acc)
		}

		// Check that buffer has one byte with 3 bits padded to 8
		bytes := w.buf.Bytes()
		if len(bytes) != 1 {
			t.Errorf("Expected 1 byte in buffer, got %d", len(bytes))
		}

		// 0b101 padded with 5 zeros becomes 0b10100000
		if bytes[0] != 0b10100000 {
			t.Errorf("Expected byte 0b10100000, got 0b%08b", bytes[0])
		}
	})

	t.Run("Already aligned", func(t *testing.T) {
		w := newBitWriter()
		w.writeBits(0xFF, 8) // Write full byte
		w.alignToByte()

		if w.bitCount != 0 {
			t.Errorf("Expected bitCount 0, got %d", w.bitCount)
		}

		// Should have one byte in buffer, no additional padding
		bufBytes := w.buf.Bytes()
		if len(bufBytes) != 1 {
			t.Errorf("Expected 1 byte in buffer, got %d", len(bufBytes))
		}
	})
}

func TestBuilder_writeBytes(t *testing.T) {
	t.Run("Write bytes to empty writer", func(t *testing.T) {
		w := newBitWriter()
		data := []byte{0xAB, 0xCD, 0xEF}
		n, err := w.writeBytes(data)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if n != len(data) {
			t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
		}

		if w.bitCount != 0 {
			t.Errorf("Expected bitCount 0 after writeBytes, got %d", w.bitCount)
		}

		bufBytes := w.buf.Bytes()
		if !bytes.Equal(bufBytes, data) {
			t.Errorf("Expected bytes %v, got %v", data, bufBytes)
		}
	})

	t.Run("Write bytes after partial byte", func(t *testing.T) {
		w := newBitWriter()
		w.writeBits(0b101, 3) // Write 3 bits first
		data := []byte{0xAB, 0xCD}
		n, err := w.writeBytes(data)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if n != len(data) {
			t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
		}

		if w.bitCount != 0 {
			t.Errorf("Expected bitCount 0 after writeBytes, got %d", w.bitCount)
		}

		// Should have 3 bytes: padded partial byte + 2 data bytes
		bytes := w.buf.Bytes()
		if len(bytes) != 3 {
			t.Errorf("Expected 3 bytes in buffer, got %d", len(bytes))
		}
	})

	t.Run("Write empty data", func(t *testing.T) {
		w := newBitWriter()
		data := []byte{}
		n, err := w.writeBytes(data)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if n != 0 {
			t.Errorf("Expected to write 0 bytes, wrote %d", n)
		}

		if w.bitCount != 0 {
			t.Errorf("Expected bitCount 0, got %d", w.bitCount)
		}
	})
}

func TestBuilder_final(t *testing.T) {
	t.Run("Empty writer", func(t *testing.T) {
		w := newBitWriter()
		data, totalBits := w.final()

		if len(data) != 0 {
			t.Errorf("Expected empty data, got %v", data)
		}

		if totalBits != 0 {
			t.Errorf("Expected totalBits 0, got %d", totalBits)
		}
	})

	t.Run("Writer with full bytes", func(t *testing.T) {
		w := newBitWriter()
		w.writeBits(0xAB, 8)
		w.writeBits(0xCD, 8)
		data, totalBits := w.final()

		if len(data) != 2 {
			t.Errorf("Expected 2 bytes, got %d", len(data))
		}

		if totalBits != 16 {
			t.Errorf("Expected totalBits 16, got %d", totalBits)
		}

		if data[0] != 0xAB || data[1] != 0xCD {
			t.Errorf("Expected bytes [0xAB, 0xCD], got %v", data)
		}
	})

	t.Run("Writer with partial byte", func(t *testing.T) {
		w := newBitWriter()
		w.writeBits(0b101, 3)
		data, totalBits := w.final()

		if len(data) != 1 {
			t.Errorf("Expected 1 byte, got %d", len(data))
		}

		if totalBits != 3 {
			t.Errorf("Expected totalBits 3, got %d", totalBits)
		}

		// 0b101 should be shifted to MSB: 0b10100000
		if data[0] != 0b10100000 {
			t.Errorf("Expected byte 0b10100000, got 0b%08b", data[0])
		}
	})
}
