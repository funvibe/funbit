package matcher

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
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
	if !segment.SizeSpecified {
		segment.Size = bitstringpkg.DefaultSizeInteger
		segment.SizeSpecified = false
	}

	// Set default unit
	segment.Unit = bitstringpkg.DefaultUnitInteger

	m.pattern = append(m.pattern, segment)
	return m
}

// Float adds a float segment to the matching pattern
func (m *Matcher) Float(variable interface{}, options ...bitstringpkg.SegmentOption) *Matcher {
	segment := bitstringpkg.NewSegment(variable, options...)
	segment.Type = bitstringpkg.TypeFloat

	// Set default size if not specified
	if !segment.SizeSpecified {
		segment.Size = bitstringpkg.DefaultSizeFloat
		segment.SizeSpecified = false
	}

	// Set default unit
	segment.Unit = bitstringpkg.DefaultUnitFloat

	m.pattern = append(m.pattern, segment)
	return m
}

// Binary adds a binary segment to the matching pattern
func (m *Matcher) Binary(variable interface{}, options ...bitstringpkg.SegmentOption) *Matcher {
	segment := bitstringpkg.NewSegment(variable, options...)
	segment.Type = bitstringpkg.TypeBinary

	// Set default size if not specified
	if !segment.SizeSpecified {
		// For binary, size will be determined by the data length
		// We'll handle this during matching
	}

	// Set default unit
	segment.Unit = bitstringpkg.DefaultUnitBinary

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
	case bitstringpkg.TypeFloat:
		return m.matchFloat(segment, bs, offset)
	case bitstringpkg.TypeBinary:
		return m.matchBinary(segment, bs, offset)
	case bitstringpkg.TypeBitstring:
		return m.matchBitstring(segment, bs, offset)
	default:
		return nil, 0, fmt.Errorf("unsupported segment type: %s", segment.Type)
	}
}

// matchInteger matches an integer segment against the bitstring
func (m *Matcher) matchInteger(segment *bitstringpkg.Segment, bs *bitstringpkg.BitString, offset uint) (*bitstringpkg.SegmentResult, uint, error) {
	if !segment.SizeSpecified {
		return nil, 0, errors.New("integer segment must have size specified")
	}

	size := segment.Size
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

// matchFloat matches a float segment against the bitstring
func (m *Matcher) matchFloat(segment *bitstringpkg.Segment, bs *bitstringpkg.BitString, offset uint) (*bitstringpkg.SegmentResult, uint, error) {
	if !segment.SizeSpecified {
		return nil, 0, errors.New("float segment must have size specified")
	}

	size := segment.Size
	if size != 16 && size != 32 && size != 64 {
		return nil, 0, fmt.Errorf("invalid float size: %d bits (must be 16, 32, or 64)", size)
	}

	// Check if we have enough bits remaining
	if offset+size > bs.Length() {
		return nil, 0, fmt.Errorf("insufficient bits: need %d, have %d", size, bs.Length()-offset)
	}

	// Extract the float value
	value, err := m.extractFloat(bs, offset, size, segment.Endianness)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to extract float: %v", err)
	}

	// Bind the value to the variable
	if err := m.bindFloatValue(segment.Value, value); err != nil {
		return nil, 0, fmt.Errorf("failed to bind float value: %v", err)
	}

	// Create remaining bitstring
	var remaining *bitstringpkg.BitString
	if offset+size < bs.Length() {
		remaining = m.extractRemainingBits(bs, offset+size)
	} else {
		remaining = bitstringpkg.NewBitString()
	}

	result := &bitstringpkg.SegmentResult{
		Value:     value,
		Matched:   true,
		Remaining: remaining,
	}

	return result, offset + size, nil
}

// matchBinary matches a binary segment against the bitstring
func (m *Matcher) matchBinary(segment *bitstringpkg.Segment, bs *bitstringpkg.BitString, offset uint) (*bitstringpkg.SegmentResult, uint, error) {
	// Determine size if not specified
	var size uint
	if !segment.SizeSpecified {
		// For binary without explicit size, use remaining bits by default
		size = bs.Length() - offset
		if size == 0 {
			return nil, 0, errors.New("binary size cannot be zero")
		}
	} else {
		size = segment.Size
		// For binary type, if size is specified in bytes (unit=8), convert to bits
		if segment.Unit == 8 {
			size = size * 8
		}
		if size == 0 {
			// If size is explicitly set to 0, use remaining bits
			size = bs.Length() - offset
			if size == 0 {
				return nil, 0, errors.New("binary size cannot be zero")
			}
		}
	}

	// Check if we have enough bits remaining
	if offset+size > bs.Length() {
		return nil, 0, fmt.Errorf("insufficient bits: need %d, have %d", size, bs.Length()-offset)
	}

	// Extract the binary data
	value, err := m.extractBinary(bs, offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to extract binary: %v", err)
	}

	// Bind the value to the variable
	if err := m.bindBinaryValue(segment.Value, value); err != nil {
		return nil, 0, fmt.Errorf("failed to bind binary value: %v", err)
	}

	// Create remaining bitstring
	var remaining *bitstringpkg.BitString
	if offset+size < bs.Length() {
		remaining = m.extractRemainingBits(bs, offset+size)
	} else {
		remaining = bitstringpkg.NewBitString()
	}

	result := &bitstringpkg.SegmentResult{
		Value:     value,
		Matched:   true,
		Remaining: remaining,
	}

	return result, offset + size, nil
}

// matchBitstring matches a bitstring segment against the bitstring
func (m *Matcher) matchBitstring(segment *bitstringpkg.Segment, bs *bitstringpkg.BitString, offset uint) (*bitstringpkg.SegmentResult, uint, error) {
	if !segment.SizeSpecified {
		return nil, 0, errors.New("bitstring segment must have size specified")
	}

	size := segment.Size
	if size == 0 || size > 64 {
		return nil, 0, fmt.Errorf("invalid bitstring size: %d bits", size)
	}

	// Check if we have enough bits remaining
	if offset+size > bs.Length() {
		return nil, 0, fmt.Errorf("insufficient bits: need %d, have %d", size, bs.Length()-offset)
	}

	// Bitstring segments are always extracted as big-endian unsigned integers.
	value, err := extractIntegerBits(bs.ToBytes(), offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to extract bitstring: %v", err)
	}

	// Bind the value to the variable
	if err := m.bindValue(segment.Value, value); err != nil {
		return nil, 0, fmt.Errorf("failed to bind bitstring value: %v", err)
	}

	// Create remaining bitstring
	var remaining *bitstringpkg.BitString
	if offset+size < bs.Length() {
		remaining = m.extractRemainingBits(bs, offset+size)
	} else {
		remaining = bitstringpkg.NewBitString()
	}

	result := &bitstringpkg.SegmentResult{
		Value:     value,
		Matched:   true,
		Remaining: remaining,
	}

	return result, offset + size, nil
}

// extractFloat extracts a float value from the bitstring
func (m *Matcher) extractFloat(bs *bitstringpkg.BitString, offset, size uint, endianness string) (float64, error) {
	data := bs.ToBytes()
	byteOffset := offset / 8
	bitOffset := offset % 8

	// For now, only support byte-aligned floats
	if bitOffset != 0 {
		return 0, fmt.Errorf("non-byte-aligned floats not supported yet")
	}

	bytesNeeded := size / 8
	if byteOffset+bytesNeeded > uint(len(data)) {
		return 0, fmt.Errorf("insufficient data for float extraction")
	}

	extractedData := data[byteOffset : byteOffset+bytesNeeded]

	switch size {
	case 16:
		// 16-bit float (half precision)
		var bits uint16
		switch endianness {
		case bitstringpkg.EndiannessBig, "":
			bits = binary.BigEndian.Uint16(extractedData)
		case bitstringpkg.EndiannessLittle:
			bits = binary.LittleEndian.Uint16(extractedData)
		case bitstringpkg.EndiannessNative:
			bits = binary.LittleEndian.Uint16(extractedData)
		default:
			return 0, fmt.Errorf("unsupported endianness: %s", endianness)
		}
		// Convert half precision to single precision (simplified)
		float32Bits := uint32(bits) << 16
		return float64(math.Float32frombits(float32Bits)), nil
	case 32:
		var bits uint32
		switch endianness {
		case bitstringpkg.EndiannessBig, "":
			bits = binary.BigEndian.Uint32(extractedData)
		case bitstringpkg.EndiannessLittle:
			bits = binary.LittleEndian.Uint32(extractedData)
		case bitstringpkg.EndiannessNative:
			bits = binary.LittleEndian.Uint32(extractedData)
		default:
			return 0, fmt.Errorf("unsupported endianness: %s", endianness)
		}
		return float64(math.Float32frombits(bits)), nil
	case 64:
		var bits uint64
		switch endianness {
		case bitstringpkg.EndiannessBig, "":
			bits = binary.BigEndian.Uint64(extractedData)
		case bitstringpkg.EndiannessLittle:
			bits = binary.LittleEndian.Uint64(extractedData)
		case bitstringpkg.EndiannessNative:
			bits = binary.LittleEndian.Uint64(extractedData)
		default:
			return 0, fmt.Errorf("unsupported endianness: %s", endianness)
		}
		return math.Float64frombits(bits), nil
	default:
		return 0, fmt.Errorf("unsupported float size: %d", size)
	}
}

// extractBinary extracts binary data from the bitstring
func (m *Matcher) extractBinary(bs *bitstringpkg.BitString, offset, size uint) ([]byte, error) {
	data := bs.ToBytes()
	byteOffset := offset / 8
	bitOffset := offset % 8

	// Handle bit-level extraction
	if bitOffset != 0 || size%8 != 0 {
		return m.extractBinaryBits(data, offset, size)
	}

	// Handle byte-aligned extraction
	bytesNeeded := size / 8
	if byteOffset+bytesNeeded > uint(len(data)) {
		return nil, fmt.Errorf("insufficient data for binary extraction")
	}

	// Return a copy of the extracted data
	result := make([]byte, bytesNeeded)
	copy(result, data[byteOffset:byteOffset+bytesNeeded])
	return result, nil
}

// extractBinaryBits extracts binary data with proper bit alignment for binary type
func (m *Matcher) extractBinaryBits(data []byte, start, length uint) ([]byte, error) {
	if start >= uint(len(data))*8 {
		return nil, fmt.Errorf("start position %d is beyond data length", start)
	}

	if length == 0 {
		return []byte{}, nil
	}

	if start+length > uint(len(data))*8 {
		return nil, fmt.Errorf("cannot extract %d bits from position %d", length, start)
	}

	// For binary data, we need to extract bits and pack them into bytes
	resultBytes := make([]byte, (length+7)/8)

	for i := uint(0); i < length; i++ {
		currentBitPos := start + i
		bytePos := currentBitPos / 8
		bitInByte := 7 - (currentBitPos % 8) // 0 is MSB

		bit := (data[bytePos] >> bitInByte) & 1
		resultBytePos := i / 8
		bitInResult := 7 - (i % 8)

		if bit != 0 {
			resultBytes[resultBytePos] |= (1 << bitInResult)
		}
	}

	return resultBytes, nil
}

// bindFloatValue binds the extracted float value to the variable
func (m *Matcher) bindFloatValue(variable interface{}, value float64) error {
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
	case reflect.Float32:
		val.SetFloat(value)
	case reflect.Float64:
		val.SetFloat(value)
	default:
		return fmt.Errorf("unsupported float variable type: %v", val.Kind())
	}

	return nil
}

// bindBinaryValue binds the extracted binary value to the variable
func (m *Matcher) bindBinaryValue(variable interface{}, value []byte) error {
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
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			val.SetBytes(value)
		} else {
			return fmt.Errorf("unsupported slice type: %v", val.Type())
		}
	case reflect.String:
		val.SetString(string(value))
	default:
		return fmt.Errorf("unsupported binary variable type: %v", val.Kind())
	}

	return nil
}

// extractInteger extracts an integer value from the bitstring
func (m *Matcher) extractInteger(bs *bitstringpkg.BitString, offset, size uint, endianness string) (int64, error) {
	data := bs.ToBytes()
	byteOffset := offset / 8
	bitOffset := offset % 8

	// Handle bit-level extraction
	if bitOffset != 0 || size%8 != 0 {
		return extractIntegerBits(data, offset, size)
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

// extractIntegerBits extracts an integer value from a non-byte-aligned position.
func extractIntegerBits(data []byte, startBit, numBits uint) (int64, error) {
	if startBit+numBits > uint(len(data))*8 {
		return 0, fmt.Errorf("cannot extract %d bits starting from bit %d", numBits, startBit)
	}
	if numBits == 0 {
		return 0, nil
	}
	if numBits > 64 {
		return 0, fmt.Errorf("cannot extract more than 64 bits into an int64")
	}

	var value uint64
	for i := uint(0); i < numBits; i++ {
		currentBitPos := startBit + i
		bytePos := currentBitPos / 8
		bitInByte := 7 - (currentBitPos % 8) // 0 is MSB

		bit := (data[bytePos] >> bitInByte) & 1
		value = (value << 1) | uint64(bit)
	}

	return int64(value), nil
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
