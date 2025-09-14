package bitstring

// This file contains additional segment-related utilities and constants

// Segment types constants
const (
	TypeInteger       = "integer"
	TypeFloat         = "float"
	TypeBinary        = "binary"
	TypeBitstring     = "bitstring"
	TypeUTF           = "utf"   // Generic UTF type
	TypeUTF8          = "utf8"  // UTF-8 specific
	TypeUTF16         = "utf16" // UTF-16 specific
	TypeUTF32         = "utf32" // UTF-32 specific
	TypeRestBinary    = "rest_binary"
	TypeRestBitstring = "rest_bitstring"
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
	DefaultUnitUTF       = 1
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
		s.Size = size
		s.SizeSpecified = true
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
		s.UnitSpecified = true
	}
}

// WithDynamicSize sets the size for a segment using a variable reference
func WithDynamicSize(sizeVar *uint) SegmentOption {
	return func(s *Segment) {
		s.DynamicSize = sizeVar
		s.IsDynamic = true
		s.SizeSpecified = false // Dynamic size overrides explicit size
	}
}

// WithDynamicSizeExpression sets the size for a segment using an expression
func WithDynamicSizeExpression(expr string) SegmentOption {
	return func(s *Segment) {
		s.DynamicExpr = expr
		s.IsDynamic = true
		s.SizeSpecified = false // Dynamic size overrides explicit size
	}
}

// NewSegment creates a new segment with the given value and options
func NewSegment(value interface{}, options ...SegmentOption) *Segment {
	segment := &Segment{
		Value:      value,
		Type:       TypeInteger,   // default type
		Signed:     Unsigned,      // default signedness
		Endianness: EndiannessBig, // default endianness
		Unit:       0,             // start with 0 to detect if unit was set
		IsDynamic:  false,         // default to static size
	}

	for _, option := range options {
		option(segment)
	}

	// Set default size based on type if not specified
	if !segment.SizeSpecified {
		segment.Size = getDefaultSizeForType(segment.Type)
		segment.SizeSpecified = false
	}

	return segment
}

// getDefaultSizeForType returns the default size for a given type
func getDefaultSizeForType(segmentType string) uint {
	switch segmentType {
	case TypeInteger:
		return DefaultSizeInteger
	case TypeFloat:
		return DefaultSizeFloat
	case TypeUTF8:
		return 8 // UTF-8 uses 8-bit encoding
	case TypeUTF16:
		return 16 // UTF-16 uses 16-bit encoding
	case TypeUTF32:
		return 32 // UTF-32 uses 32-bit encoding
	default:
		return 0 // no default size for binary/bitstring types
	}
}

// getDefaultUnitForType returns the default unit for a given type
func getDefaultUnitForType(segmentType string) uint {
	switch segmentType {
	case TypeInteger, TypeFloat, TypeBitstring:
		return DefaultUnitInteger
	case TypeBinary:
		return DefaultUnitBinary
	case TypeUTF, TypeUTF8, TypeUTF16, TypeUTF32:
		return DefaultUnitUTF
	default:
		return DefaultUnitInteger // default for unknown types
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

	// Validate unit - check if explicitly set to invalid value
	if segment.UnitSpecified {
		if segment.Unit < 1 || segment.Unit > 256 {
			return &BitStringError{
				Code:    "INVALID_UNIT",
				Message: "unit must be between 1 and 256",
			}
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
		if segment.SizeSpecified && (segment.Size != 16 && segment.Size != 32 && segment.Size != 64) {
			return &BitStringError{
				Code:    "INVALID_FLOAT_SIZE",
				Message: "float size must be 16, 32, or 64 bits",
			}
		}
	case TypeUTF8, TypeUTF16, TypeUTF32:
		if segment.SizeSpecified {
			return &BitStringError{
				Code:    "UTF_SIZE_SPECIFIED",
				Message: "UTF types cannot have size specified",
			}
		}
		// For UTF types, unit can only be set to the default value (1)
		// but only if it was explicitly specified
		if segment.UnitSpecified && segment.Unit != getDefaultUnitForType(segment.Type) {
			return &BitStringError{
				Code:    "UTF_UNIT_MODIFIED",
				Message: "UTF types cannot have unit modified from default value",
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
