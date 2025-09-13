package bitstring

import (
	"reflect"
	"testing"
)

func TestBitString_NewBitString(t *testing.T) {
	bs := NewBitString()

	if bs == nil {
		t.Fatal("Expected NewBitString() to return non-nil")
	}

	if bs.Length() != 0 {
		t.Errorf("Expected empty bitstring length 0, got %d", bs.Length())
	}

	if !bs.IsEmpty() {
		t.Error("Expected empty bitstring to be empty")
	}

	if !bs.IsBinary() {
		t.Error("Expected empty bitstring to be binary (0 is multiple of 8)")
	}
}

func TestBitString_NewBitStringFromBytes(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		expectedLen    uint
		expectedBinary bool
		expectedBytes  []byte
	}{
		{
			name:           "empty bytes",
			data:           []byte{},
			expectedLen:    0,
			expectedBinary: true,
			expectedBytes:  []byte{},
		},
		{
			name:           "single byte",
			data:           []byte{42},
			expectedLen:    8,
			expectedBinary: true,
			expectedBytes:  []byte{42},
		},
		{
			name:           "multiple bytes",
			data:           []byte{1, 2, 3, 4},
			expectedLen:    32,
			expectedBinary: true,
			expectedBytes:  []byte{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitStringFromBytes(tt.data)

			if bs == nil {
				t.Fatal("Expected NewBitStringFromBytes() to return non-nil")
			}

			if bs.Length() != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, bs.Length())
			}

			if bs.IsEmpty() != (tt.expectedLen == 0) {
				t.Errorf("Expected IsEmpty() to be %v for length %d", tt.expectedLen == 0, tt.expectedLen)
			}

			if bs.IsBinary() != tt.expectedBinary {
				t.Errorf("Expected IsBinary() to be %v", tt.expectedBinary)
			}

			bytes := bs.ToBytes()
			if !reflect.DeepEqual(bytes, tt.expectedBytes) {
				t.Errorf("Expected bytes %v, got %v", tt.expectedBytes, bytes)
			}
		})
	}
}

func TestBitString_NewBitStringFromBits(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		length         uint
		expectedLen    uint
		expectedBinary bool
		wantErr        bool
	}{
		{
			name:           "zero bits",
			data:           []byte{},
			length:         0,
			expectedLen:    0,
			expectedBinary: true,
			wantErr:        false,
		},
		{
			name:           "4 bits (half byte)",
			data:           []byte{0b10110000},
			length:         4,
			expectedLen:    4,
			expectedBinary: false,
			wantErr:        false,
		},
		{
			name:           "12 bits (1.5 bytes)",
			data:           []byte{0xFF, 0xF0},
			length:         12,
			expectedLen:    12,
			expectedBinary: false,
			wantErr:        false,
		},
		{
			name:           "16 bits (exactly 2 bytes)",
			data:           []byte{0xFF, 0xAA},
			length:         16,
			expectedLen:    16,
			expectedBinary: true,
			wantErr:        false,
		},
		{
			name:           "insufficient data",
			data:           []byte{0xFF},
			length:         16,
			expectedLen:    0,
			expectedBinary: false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitStringFromBits(tt.data, tt.length)

			if tt.wantErr {
				if bs != nil {
					t.Error("Expected NewBitStringFromBits() to return nil on error")
				}
				return
			}

			if bs == nil {
				t.Fatal("Expected NewBitStringFromBits() to return non-nil")
			}

			if bs.Length() != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, bs.Length())
			}

			if bs.IsEmpty() != (tt.expectedLen == 0) {
				t.Errorf("Expected IsEmpty() to be %v for length %d", tt.expectedLen == 0, tt.expectedLen)
			}

			if bs.IsBinary() != tt.expectedBinary {
				t.Errorf("Expected IsBinary() to be %v", tt.expectedBinary)
			}
		})
	}
}

func TestBitString_Clone(t *testing.T) {
	original := NewBitStringFromBytes([]byte{1, 2, 3})

	cloned := original.Clone()

	if cloned == nil {
		t.Fatal("Expected Clone() to return non-nil")
	}

	if cloned.Length() != original.Length() {
		t.Errorf("Expected clone to have same length %d, got %d", original.Length(), cloned.Length())
	}

	if cloned.IsEmpty() != original.IsEmpty() {
		t.Errorf("Expected clone to have same IsEmpty() value %v", original.IsEmpty())
	}

	if cloned.IsBinary() != original.IsBinary() {
		t.Errorf("Expected clone to have same IsBinary() value %v", original.IsBinary())
	}

	originalBytes := original.ToBytes()
	clonedBytes := cloned.ToBytes()

	if !reflect.DeepEqual(originalBytes, clonedBytes) {
		t.Errorf("Expected clone to have same bytes %v, got %v", originalBytes, clonedBytes)
	}
}

func TestBitString_EmptyString(t *testing.T) {
	bs := NewBitString()

	if !bs.IsEmpty() {
		t.Error("Expected empty bitstring to be empty")
	}

	if bs.Length() != 0 {
		t.Errorf("Expected empty bitstring length 0, got %d", bs.Length())
	}

	if !bs.IsBinary() {
		t.Error("Expected empty bitstring to be binary (0 is multiple of 8)")
	}

	bytes := bs.ToBytes()
	if len(bytes) != 0 {
		t.Errorf("Expected empty bytes, got %v", bytes)
	}
}
