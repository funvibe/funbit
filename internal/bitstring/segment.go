package bitstring

// This file contains additional segment-related utilities and constants

// Segment types constants
const (
	TypeInteger   = "integer"
	TypeFloat     = "float"
	TypeBinary    = "binary"
	TypeBitstring = "bitstring"
	TypeUTF8      = "utf8"
	TypeUTF16     = "utf16"
	TypeUTF32     = "utf32"
)

// Endianness constants
const (
	EndiannessBig    = "big"
	EndiannessLittle = "little"
	EndiannessNative = "native"
)

// Signedness constants
const (
	Signed   = true
	Unsigned = false
)

// Default unit values for different types
const (
	DefaultUnitInteger   = 1
	DefaultUnitFloat     = 1
	DefaultUnitBinary    = 8
	DefaultUnitBitstring = 1
)

// Default size values for different types (in bits)
const (
	DefaultSizeInteger = 8
	DefaultSizeFloat   = 64
)

// SegmentOption is a function type for configuring segments
type SegmentOption func(*Segment)

// WithSize sets the size for a segment
func WithSize(size uint) SegmentOption {
	return func(s *Segment) {
		s.Size = &size
	}
}

// WithType sets the type for a segment
func WithType(segmentType string) SegmentOption {
	return func(s *Segment) {
		s.Type = segmentType
	}
}

// WithSigned sets the signedness for a segment
func WithSigned(signed bool) SegmentOption {
	return func(s *Segment) {
		s.Signed = signed
	}
}

// WithEndianness sets the endianness for a segment
func WithEndianness(endianness string) SegmentOption {
	return func(s *Segment) {
		s.Endianness = endianness
	}
}

// WithUnit sets the unit for a segment
func WithUnit(unit uint) SegmentOption {
	return func(s *Segment) {
		s.Unit = unit
	}
}

// NewSegment creates a new segment with the given value and options
func NewSegment(value interface{}, options ...SegmentOption) *Segment {
	segment := &Segment{
		Value:      value,
		Type:       TypeInteger,        // default type
		Signed:     Unsigned,           // default signedness
		Endianness: EndiannessBig,      // default endianness
		Unit:       DefaultUnitInteger, // default unit
	}

	for _, option := range options {
		option(segment)
	}

	// Set default size based on type if not specified
	if segment.Size == nil {
		defaultSize := getDefaultSizeForType(segment.Type)
		segment.Size = &defaultSize
	}

	// Set default unit based on type
	segment.Unit = getDefaultUnitForType(segment.Type)

	return segment
}

// getDefaultSizeForType returns the default size for a given type
func getDefaultSizeForType(segmentType string) uint {
	switch segmentType {
	case TypeInteger:
		return DefaultSizeInteger
	case TypeFloat:
		return DefaultSizeFloat
	default:
		return 0 // no default size for binary/bitstring/utf types
	}
}

// getDefaultUnitForType returns the default unit for a given type
func getDefaultUnitForType(segmentType string) uint {
	switch segmentType {
	case TypeInteger, TypeFloat, TypeBitstring:
		return DefaultUnitInteger
	case TypeBinary:
		return DefaultUnitBinary
	default:
		return DefaultUnitInteger // default for utf types
	}
}

// ValidateSegment checks if a segment has valid configuration
func ValidateSegment(segment *Segment) error {
	if segment == nil {
		return &BitStringError{
			Code:    "INVALID_SEGMENT",
			Message: "segment cannot be nil",
		}
	}

	// Validate type
	if segment.Type == "" {
		segment.Type = TypeInteger // default to integer
	}

	// Validate unit
	if segment.Unit == 0 {
		segment.Unit = getDefaultUnitForType(segment.Type)
	}

	if segment.Unit < 1 || segment.Unit > 256 {
		return &BitStringError{
			Code:    "INVALID_UNIT",
			Message: "unit must be between 1 and 256",
		}
	}

	// Validate endianness
	if segment.Endianness == "" {
		segment.Endianness = EndiannessBig // default to big
	}

	if segment.Endianness != EndiannessBig &&
		segment.Endianness != EndiannessLittle &&
		segment.Endianness != EndiannessNative {
		return &BitStringError{
			Code:    "INVALID_ENDIANNESS",
			Message: "endianness must be 'big', 'little', or 'native'",
		}
	}

	// Type-specific validations
	switch segment.Type {
	case TypeFloat:
		if segment.Size != nil && (*segment.Size != 16 && *segment.Size != 32 && *segment.Size != 64) {
			return &BitStringError{
				Code:    "INVALID_FLOAT_SIZE",
				Message: "float size must be 16, 32, or 64 bits",
			}
		}
	case TypeUTF8, TypeUTF16, TypeUTF32:
		if segment.Size != nil {
			return &BitStringError{
				Code:    "UTF_SIZE_SPECIFIED",
				Message: "UTF types cannot have size specified",
			}
		}
		if segment.Unit != getDefaultUnitForType(segment.Type) {
			return &BitStringError{
				Code:    "UTF_UNIT_MODIFIED",
				Message: "UTF types cannot have unit modified",
			}
		}
	}

	return nil
}

// BitStringError represents an error in bitstring operations
type BitStringError struct {
	Code    string
	Message string
	Context interface{}
}

func (e *BitStringError) Error() string {
	if e.Context != nil {
		return e.Message + " (context: " + string(e.Context.(string)) + ")"
	}
	return e.Message
}
