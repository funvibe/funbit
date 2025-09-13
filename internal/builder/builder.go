package builder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/funvibe/funbit/internal/bitstring"
)

// Builder provides a fluent interface for constructing bitstrings
type Builder struct {
	segments []*bitstring.Segment
	buffer   *bytes.Buffer
}

// NewBuilder creates a new builder instance
func NewBuilder() *Builder {
	return &Builder{
		segments: []*bitstring.Segment{},
		buffer:   &bytes.Buffer{},
	}
}

// AddInteger adds an integer segment to the builder
func (b *Builder) AddInteger(value interface{}, options ...bitstring.SegmentOption) *Builder {
	segment := bitstring.NewSegment(value, options...)
	segment.Type = bitstring.TypeInteger

	// Set default size if not specified
	if segment.Size == nil {
		defaultSize := uint(bitstring.DefaultSizeInteger)
		segment.Size = &defaultSize
	}

	// Set default unit
	segment.Unit = bitstring.DefaultUnitInteger

	b.segments = append(b.segments, segment)
	return b
}

// Build constructs the final bitstring from all segments
func (b *Builder) Build() (*bitstring.BitString, error) {
	b.buffer.Reset()

	for _, segment := range b.segments {
		if err := b.encodeSegment(segment); err != nil {
			return nil, err
		}
	}

	// Create bitstring from the accumulated data
	data := b.buffer.Bytes()
	if len(data) == 0 {
		return bitstring.NewBitString(), nil
	}

	// Calculate total bit length
	totalBits := uint(len(data)) * 8
	return bitstring.NewBitStringFromBits(data, totalBits), nil
}

// encodeSegment encodes a single segment into the buffer
func (b *Builder) encodeSegment(segment *bitstring.Segment) error {
	if err := bitstring.ValidateSegment(segment); err != nil {
		return err
	}

	switch segment.Type {
	case bitstring.TypeInteger:
		return b.encodeInteger(segment)
	default:
		return fmt.Errorf("unsupported segment type: %s", segment.Type)
	}
}

// encodeInteger encodes an integer value into the buffer
func (b *Builder) encodeInteger(segment *bitstring.Segment) error {
	if segment.Size == nil {
		return errors.New("integer segment must have size specified")
	}

	size := *segment.Size
	if size == 0 || size > 64 {
		return fmt.Errorf("invalid integer size: %d bits", size)
	}

	// Convert value to int64
	var value int64
	switch v := segment.Value.(type) {
	case int:
		value = int64(v)
	case int8:
		value = int64(v)
	case int16:
		value = int64(v)
	case int32:
		value = int64(v)
	case int64:
		value = v
	case uint:
		value = int64(v)
	case uint8:
		value = int64(v)
	case uint16:
		value = int64(v)
	case uint32:
		value = int64(v)
	case uint64:
		if v > uint64(math.MaxInt64) {
			return fmt.Errorf("value %d too large for int64", v)
		}
		value = int64(v)
	default:
		return fmt.Errorf("unsupported integer value type: %T", segment.Value)
	}

	// Check if value fits in the specified size
	maxValue := int64(1)<<size - 1
	if value < 0 || value > maxValue {
		// Truncate to fit (as per Erlang behavior)
		value = value & maxValue
	}

	// Encode the integer
	return b.writeInteger(value, size, segment.Endianness)
}

// writeInteger writes an integer value with the specified size and endianness
func (b *Builder) writeInteger(value int64, size uint, endianness string) error {
	bytesNeeded := (size + 7) / 8

	switch endianness {
	case bitstring.EndiannessBig, "":
		return b.writeIntegerBigEndian(value, bytesNeeded)
	case bitstring.EndiannessLittle:
		return b.writeIntegerLittleEndian(value, bytesNeeded)
	case bitstring.EndiannessNative:
		// For now, default to big-endian for native
		return b.writeIntegerBigEndian(value, bytesNeeded)
	default:
		return fmt.Errorf("unsupported endianness: %s", endianness)
	}
}

// writeIntegerBigEndian writes an integer in big-endian format
func (b *Builder) writeIntegerBigEndian(value int64, bytesNeeded uint) error {
	// Create buffer for the integer
	buf := make([]byte, bytesNeeded)

	// Write from most significant to least significant
	for i := uint(0); i < bytesNeeded; i++ {
		shift := (bytesNeeded - 1 - i) * 8
		buf[i] = byte((value >> shift) & 0xFF)
	}

	_, err := b.buffer.Write(buf)
	return err
}

// writeIntegerLittleEndian writes an integer in little-endian format
func (b *Builder) writeIntegerLittleEndian(value int64, bytesNeeded uint) error {
	// Create buffer for the integer
	buf := make([]byte, bytesNeeded)

	// Write from least significant to most significant
	for i := uint(0); i < bytesNeeded; i++ {
		shift := i * 8
		buf[i] = byte((value >> shift) & 0xFF)
	}

	_, err := b.buffer.Write(buf)
	return err
}

// writeIntegerNative writes an integer in native endianness format
func (b *Builder) writeIntegerNative(value int64, bytesNeeded uint) error {
	// Use binary.Write to handle native endianness
	buf := make([]byte, bytesNeeded)

	switch bytesNeeded {
	case 1:
		buf[0] = byte(value)
	case 2:
		binary.LittleEndian.PutUint16(buf, uint16(value))
	case 4:
		binary.LittleEndian.PutUint32(buf, uint32(value))
	case 8:
		binary.LittleEndian.PutUint64(buf, uint64(value))
	default:
		// For unusual sizes, fall back to little-endian
		return b.writeIntegerLittleEndian(value, bytesNeeded)
	}

	_, err := b.buffer.Write(buf)
	return err
}
