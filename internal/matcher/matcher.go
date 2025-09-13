package matcher

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	bitstringpkg "github.com/funvibe/funbit/internal/bitstring"
)

// Matcher provides a fluent interface for pattern matching against bitstrings
type Matcher struct {
	pattern []*bitstringpkg.Segment
}

// NewMatcher creates a new matcher instance
func NewMatcher() *Matcher {
	return &Matcher{
		pattern: []*bitstringpkg.Segment{},
	}
}

// Integer adds an integer segment to the matching pattern
func (m *Matcher) Integer(variable interface{}, options ...bitstringpkg.SegmentOption) *Matcher {
	segment := bitstringpkg.NewSegment(variable, options...)
	segment.Type = bitstringpkg.TypeInteger

	// Set default size if not specified
	if segment.Size == nil {
		defaultSize := uint(bitstringpkg.DefaultSizeInteger)
		segment.Size = &defaultSize
	}

	// Set default unit
	segment.Unit = bitstringpkg.DefaultUnitInteger

	m.pattern = append(m.pattern, segment)
	return m
}

// Match attempts to match the pattern against the provided bitstring
func (m *Matcher) Match(bitstring *bitstringpkg.BitString) ([]bitstringpkg.SegmentResult, error) {
	if bitstring == nil {
		return nil, errors.New("bitstring cannot be nil")
	}

	results := make([]bitstringpkg.SegmentResult, len(m.pattern))
	currentOffset := uint(0)

	for i, segment := range m.pattern {
		if err := bitstringpkg.ValidateSegment(segment); err != nil {
			return nil, fmt.Errorf("invalid segment %d: %v", i, err)
		}

		result, newOffset, err := m.matchSegment(segment, bitstring, currentOffset)
		if err != nil {
			return nil, fmt.Errorf("failed to match segment %d: %v", i, err)
		}

		results[i] = *result
		currentOffset = newOffset
	}

	return results, nil
}

// matchSegment matches a single segment against the bitstring at the given offset
func (m *Matcher) matchSegment(segment *bitstringpkg.Segment, bs *bitstringpkg.BitString, offset uint) (*bitstringpkg.SegmentResult, uint, error) {
	switch segment.Type {
	case bitstringpkg.TypeInteger:
		return m.matchInteger(segment, bs, offset)
	default:
		return nil, 0, fmt.Errorf("unsupported segment type: %s", segment.Type)
	}
}

// matchInteger matches an integer segment against the bitstring
func (m *Matcher) matchInteger(segment *bitstringpkg.Segment, bs *bitstringpkg.BitString, offset uint) (*bitstringpkg.SegmentResult, uint, error) {
	if segment.Size == nil {
		return nil, 0, errors.New("integer segment must have size specified")
	}

	size := *segment.Size
	if size == 0 || size > 64 {
		return nil, 0, fmt.Errorf("invalid integer size: %d bits", size)
	}

	// Check if we have enough bits remaining
	if offset+size > bs.Length() {
		return nil, 0, fmt.Errorf("insufficient bits: need %d, have %d", size, bs.Length()-offset)
	}

	// Extract the integer value
	value, err := m.extractInteger(bs, offset, size, segment.Endianness)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to extract integer: %v", err)
	}

	// Bind the value to the variable
	if err := m.bindValue(segment.Value, value); err != nil {
		return nil, 0, fmt.Errorf("failed to bind value: %v", err)
	}

	// Create remaining bitstring
	var remaining *bitstringpkg.BitString
	if offset+size < bs.Length() {
		// Extract remaining bits
		remainingData := bs.ToBytes()
		remainingOffset := (offset + size) / 8
		remainingBitOffset := (offset + size) % 8

		if remainingBitOffset == 0 {
			// Aligned to byte boundary
			remaining = bitstringpkg.NewBitStringFromBytes(remainingData[remainingOffset:])
		} else {
			// Not aligned - need bit-level extraction
			remaining = m.extractRemainingBits(bs, offset+size)
		}
	} else {
		// No remaining bits
		remaining = bitstringpkg.NewBitString()
	}

	result := &bitstringpkg.SegmentResult{
		Value:     value,
		Matched:   true,
		Remaining: remaining,
	}

	return result, offset + size, nil
}

// extractInteger extracts an integer value from the bitstring
func (m *Matcher) extractInteger(bs *bitstringpkg.BitString, offset, size uint, endianness string) (int64, error) {
	data := bs.ToBytes()
	byteOffset := offset / 8
	bitOffset := offset % 8

	// Handle bit-level extraction
	if bitOffset != 0 || size%8 != 0 {
		return m.extractIntegerBits(data, byteOffset, bitOffset, size, endianness)
	}

	// Handle byte-aligned extraction
	bytesNeeded := size / 8
	if byteOffset+bytesNeeded > uint(len(data)) {
		return 0, fmt.Errorf("insufficient data for extraction")
	}

	extractedData := data[byteOffset : byteOffset+bytesNeeded]

	switch endianness {
	case bitstringpkg.EndiannessBig, "":
		return m.bytesToInt64BigEndian(extractedData)
	case bitstringpkg.EndiannessLittle:
		return m.bytesToInt64LittleEndian(extractedData)
	case bitstringpkg.EndiannessNative:
		return m.bytesToInt64Native(extractedData)
	default:
		return 0, fmt.Errorf("unsupported endianness: %s", endianness)
	}
}

// extractIntegerBits extracts an integer value from non-byte-aligned bits
func (m *Matcher) extractIntegerBits(data []byte, byteOffset, bitOffset, size uint, endianness string) (int64, error) {
	var value int64 = 0

	remainingBits := size

	for remainingBits > 0 {
		if byteOffset >= uint(len(data)) {
			return 0, fmt.Errorf("insufficient data for bit extraction")
		}

		bitsAvailable := 8 - bitOffset
		bitsToExtract := remainingBits
		if bitsToExtract > bitsAvailable {
			bitsToExtract = bitsAvailable
		}

		// Extract bits from current byte
		byteVal := data[byteOffset]
		mask := byte((1 << bitsToExtract) - 1)
		extractedBits := (byteVal >> (bitsAvailable - bitsToExtract)) & mask

		value = (value << bitsToExtract) | int64(extractedBits)

		remainingBits -= bitsToExtract
		byteOffset++
		bitOffset = 0
	}

	return value, nil
}

// bytesToInt64BigEndian converts bytes to int64 in big-endian format
func (m *Matcher) bytesToInt64BigEndian(data []byte) (int64, error) {
	var result int64 = 0

	for _, b := range data {
		result = (result << 8) | int64(b)
	}

	return result, nil
}

// bytesToInt64LittleEndian converts bytes to int64 in little-endian format
func (m *Matcher) bytesToInt64LittleEndian(data []byte) (int64, error) {
	var result int64 = 0

	for i := len(data) - 1; i >= 0; i-- {
		result = (result << 8) | int64(data[i])
	}

	return result, nil
}

// bytesToInt64Native converts bytes to int64 in native endianness format
func (m *Matcher) bytesToInt64Native(data []byte) (int64, error) {
	switch len(data) {
	case 1:
		return int64(data[0]), nil
	case 2:
		return int64(binary.LittleEndian.Uint16(data)), nil
	case 4:
		return int64(binary.LittleEndian.Uint32(data)), nil
	case 8:
		return int64(binary.LittleEndian.Uint64(data)), nil
	default:
		// Fall back to little-endian for unusual sizes
		return m.bytesToInt64LittleEndian(data)
	}
}

// bindValue binds the extracted value to the variable
func (m *Matcher) bindValue(variable interface{}, value int64) error {
	if variable == nil {
		return errors.New("variable cannot be nil")
	}

	// Use reflection to set the value
	val := reflect.ValueOf(variable)

	// Check if it's a pointer
	if val.Kind() != reflect.Ptr {
		return errors.New("variable must be a pointer")
	}

	// Dereference the pointer
	val = val.Elem()

	// Check if it's settable
	if !val.CanSet() {
		return errors.New("variable is not settable")
	}

	// Set the value based on the type
	switch val.Kind() {
	case reflect.Int:
		val.SetInt(value)
	case reflect.Int8:
		val.SetInt(value)
	case reflect.Int16:
		val.SetInt(value)
	case reflect.Int32:
		val.SetInt(value)
	case reflect.Int64:
		val.SetInt(value)
	case reflect.Uint:
		val.SetUint(uint64(value))
	case reflect.Uint8:
		val.SetUint(uint64(value))
	case reflect.Uint16:
		val.SetUint(uint64(value))
	case reflect.Uint32:
		val.SetUint(uint64(value))
	case reflect.Uint64:
		val.SetUint(uint64(value))
	default:
		return fmt.Errorf("unsupported variable type: %v", val.Kind())
	}

	return nil
}

// extractRemainingBits extracts the remaining bits from a bitstring after a given offset
func (m *Matcher) extractRemainingBits(bs *bitstringpkg.BitString, offset uint) *bitstringpkg.BitString {
	if offset >= bs.Length() {
		return bitstringpkg.NewBitString()
	}

	remainingSize := bs.Length() - offset
	data := bs.ToBytes()
	byteOffset := offset / 8
	bitOffset := offset % 8

	if bitOffset == 0 {
		// Aligned case
		return bitstringpkg.NewBitStringFromBytes(data[byteOffset:])
	}

	// Non-aligned case - need to shift bits
	resultData := make([]byte, (remainingSize+7)/8)
	dataIndex := byteOffset
	resultIndex := uint(0)
	remainingBits := remainingSize

	// Handle first partial byte
	if bitOffset > 0 && remainingBits > 0 {
		bitsFromFirstByte := 8 - bitOffset
		if bitsFromFirstByte > remainingBits {
			bitsFromFirstByte = remainingBits
		}

		firstByte := data[dataIndex]
		mask := byte((1 << bitsFromFirstByte) - 1)
		extractedBits := (firstByte >> (8 - bitOffset - bitsFromFirstByte)) & mask

		resultData[resultIndex] = extractedBits << (8 - bitsFromFirstByte)

		dataIndex++
		remainingBits -= bitsFromFirstByte

		if remainingBits > 0 {
			resultIndex++
		}
	}

	// Handle full bytes
	for remainingBits >= 8 && dataIndex < uint(len(data)) {
		resultData[resultIndex] = data[dataIndex]
		dataIndex++
		resultIndex++
		remainingBits -= 8
	}

	// Handle last partial byte
	if remainingBits > 0 && dataIndex < uint(len(data)) {
		lastByte := data[dataIndex]
		mask := byte((1 << remainingBits) - 1)
		extractedBits := (lastByte >> (8 - remainingBits)) & mask

		if resultIndex < uint(len(resultData)) {
			resultData[resultIndex] = extractedBits << (8 - remainingBits)
		}
	}

	return bitstringpkg.NewBitStringFromBits(resultData, remainingSize)
}
