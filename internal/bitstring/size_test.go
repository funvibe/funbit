package bitstring

import (
	"testing"
)

func TestSizeHandling_ValidateSize(t *testing.T) {
	// Тесты для валидации размеров
	testCases := []struct {
		name        string
		size        uint
		unit        uint
		expectError bool
	}{
		{"Valid size 1", 1, 1, false},
		{"Valid size 64", 64, 1, false},
		{"Invalid size 0", 0, 1, true},
		{"Invalid size 65", 65, 1, true},
		{"Valid size with unit", 4, 16, false},  // 4 * 16 = 64 bits
		{"Invalid size with unit", 5, 16, true}, // 5 * 16 = 80 bits
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSize(tc.size, tc.unit)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSizeHandling_CalculateTotalSize(t *testing.T) {
	// Тесты для расчета общего размера
	testCases := []struct {
		name     string
		segment  Segment
		expected uint
		error    bool
	}{
		{
			name: "Integer size 8",
			segment: Segment{
				Value:         255,
				Size:          8,
				SizeSpecified: true,
				Type:          "integer",
			},
			expected: 8,
			error:    false,
		},
		{
			name: "Binary with unit",
			segment: Segment{
				Value:         []byte{0x12, 0x34},
				Size:          2,
				SizeSpecified: true,
				Unit:          8,
				Type:          "binary",
			},
			expected: 16,
			error:    false,
		},
		{
			name: "Invalid size",
			segment: Segment{
				Value:         1,
				Size:          0,
				SizeSpecified: true,
				Type:          "integer",
			},
			expected: 0,
			error:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := CalculateTotalSize(tc.segment)

			if tc.error {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if size != tc.expected {
					t.Errorf("Expected size %d, got %d", tc.expected, size)
				}
			}
		})
	}
}

func TestSizeHandling_ExtractBits(t *testing.T) {
	// Тесты для извлечения бит
	testCases := []struct {
		name     string
		data     []byte
		start    uint
		length   uint
		expected []byte
		error    bool
	}{
		{
			name:     "Extract first 4 bits",
			data:     []byte{0xF0, 0x0F}, // 11110000 00001111
			start:    0,
			length:   4,
			expected: []byte{0xF0}, // 11110000
			error:    false,
		},
		{
			name:     "Extract middle 8 bits",
			data:     []byte{0xF0, 0x0F, 0xFF}, // 11110000 00001111 11111111
			start:    4,
			length:   8,
			expected: []byte{0x00}, // 00000000 (8 бит начиная с 4-й позиции: 0000 + 0000)
			error:    false,
		},
		{
			name:     "Extract last 3 bits",
			data:     []byte{0xE0}, // 11100000
			start:    5,
			length:   3,
			expected: []byte{0x00}, // 00000000 (последние 3 бита из 11100000 это 000)
			error:    false,
		},
		{
			name:     "Invalid start position",
			data:     []byte{0xFF},
			start:    8,
			length:   1,
			expected: nil,
			error:    true,
		},
		{
			name:     "Length too large",
			data:     []byte{0xFF},
			start:    0,
			length:   9,
			expected: nil,
			error:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ExtractBits(tc.data, tc.start, tc.length)

			if tc.error {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if len(result) != len(tc.expected) {
					t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
					return
				}
				for i := range result {
					if result[i] != tc.expected[i] {
						t.Errorf("At position %d: expected 0x%02X, got 0x%02X", i, tc.expected[i], result[i])
						break
					}
				}
			}
		})
	}
}

func TestSizeHandling_SetBits(t *testing.T) {
	// Тесты для установки бит
	testCases := []struct {
		name     string
		target   []byte
		data     []byte
		start    uint
		expected []byte
		error    bool
	}{
		{
			name:     "Set first 4 bits",
			target:   []byte{0x00, 0xFF},
			data:     []byte{0xF0}, // 11110000
			start:    0,
			expected: []byte{0xF0, 0xFF}, // 11110000 11111111
			error:    false,
		},
		{
			name:     "Set middle 8 bits",
			target:   []byte{0xFF, 0xFF, 0xFF},
			data:     []byte{0x00}, // 00000000
			start:    4,
			expected: []byte{0xF0, 0x0F, 0xFF}, // 11110000 00001111 11111111
			error:    false,
		},
		{
			name:     "Set bits crossing byte boundary",
			target:   []byte{0x00, 0x00, 0x00}, // 3 байта = 24 бита
			data:     []byte{0xFF, 0xC0},       // 11111111 11000000 (12 бит)
			start:    2,
			expected: []byte{0x3F, 0xF0, 0x00}, // 00111111 11110000 00000000
			error:    false,
		},
		{
			name:     "Invalid start position",
			target:   []byte{0xFF},
			data:     []byte{0x01},
			start:    8,
			expected: nil,
			error:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			targetCopy := make([]byte, len(tc.target))
			copy(targetCopy, tc.target)

			err := SetBits(targetCopy, tc.data, tc.start)

			if tc.error {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if len(targetCopy) != len(tc.expected) {
					t.Errorf("Expected length %d, got %d", len(tc.expected), len(targetCopy))
					return
				}
				for i := range targetCopy {
					if targetCopy[i] != tc.expected[i] {
						t.Errorf("At position %d: expected 0x%02X, got 0x%02X", i, tc.expected[i], targetCopy[i])
						break
					}
				}
			}
		})
	}
}

func TestSizeHandling_Alignment(t *testing.T) {
	// Тесты для выравнивания
	testCases := []struct {
		name        string
		data        []byte
		offset      uint
		alignment   uint
		expected    []byte
		expectError bool
	}{
		{
			name:      "Align 3 bits to byte boundary",
			data:      []byte{0xE0}, // 11100000 (3 бита данных)
			offset:    3,
			alignment: 8,
			expected:  []byte{0xE0, 0x00}, // 11100000 00000000 (добавлен байт выравнивания)
		},
		{
			name:      "Align 12 bits to 16-bit boundary",
			data:      []byte{0x12, 0x30}, // 00010010 00110000 (12 бит данных)
			offset:    12,
			alignment: 16,
			expected:  []byte{0x12, 0x30, 0x00}, // добавить 4 бита паддинга
		},
		{
			name:        "Invalid alignment",
			data:        []byte{0xFF},
			offset:      8,
			alignment:   0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := AlignData(tc.data, tc.offset, tc.alignment)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if len(result) != len(tc.expected) {
					t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
					return
				}
				for i := range result {
					if result[i] != tc.expected[i] {
						t.Errorf("At position %d: expected 0x%02X, got 0x%02X", i, tc.expected[i], result[i])
						break
					}
				}
			}
		})
	}
}

func TestSizeHandling_Padding(t *testing.T) {
	// Тесты для паддинга
	testCases := []struct {
		name     string
		data     []byte
		bitLen   uint
		target   uint
		expected []byte
	}{
		{
			name:     "Pad 3 bits to byte",
			data:     []byte{0xE0}, // 11100000 (3 бита данных)
			bitLen:   3,
			target:   8,
			expected: []byte{0xE0, 0x00}, // добавить 5 бит паддинга (1 байт)
		},
		{
			name:     "Pad 12 bits to 16 bits",
			data:     []byte{0x12, 0x30}, // 00010010 00110000 (12 бит)
			bitLen:   12,
			target:   16,
			expected: []byte{0x12, 0x30, 0x00}, // добавить 4 бита паддинга (1 байт)
		},
		{
			name:     "Pad 9 bits to 16 bits",
			data:     []byte{0x80, 0x00}, // 10000000 00000000 (9 бит)
			bitLen:   9,
			target:   16,
			expected: []byte{0x80, 0x00, 0x00}, // добавить 7 бит паддинга (1 байт)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PadData(tc.data, tc.bitLen, tc.target)

			if len(result) != len(tc.expected) {
				t.Errorf("Expected length %d, got %d", len(tc.expected), len(result))
				return
			}

			for i := range result {
				if result[i] != tc.expected[i] {
					t.Errorf("At position %d: expected 0x%02X, got 0x%02X", i, tc.expected[i], result[i])
					break
				}
			}
		})
	}
}

// Вспомогательные функции
func uintPtr(val uint) *uint {
	return &val
}
