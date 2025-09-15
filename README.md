# Funbit - Erlang/OTP Bit Syntax Library for Go

Funbit is a comprehensive Go library that provides Erlang/OTP bit syntax compatibility for working with bitstrings and binary data. It offers a fluent interface for both constructing and pattern matching bitstrings with full support for dynamic sizing, various data types, and endianness handling.

## Features

- **Complete Erlang/OTP Bit Syntax Compatibility**: Full support for Erlang's bit syntax expressions
- **Multiple Data Types**: Integer, float, binary, bitstring, UTF-8/16/32
- **Dynamic Sizing**: Support for variable-sized segments using expressions
- **Endianness Support**: Big, little, and native endianness handling
- **Fluent Interface**: Clean, chainable API for both construction and matching
- **Type Safety**: Strong typing with comprehensive validation
- **Performance Optimized**: Efficient bit-level operations with minimal memory overhead

## Installation

```bash
go get github.com/funvibe/funbit
```

## Quick Start

### Basic Construction

```go
package main

import (
    "fmt"
    "github.com/funvibe/funbit/pkg/funbit"
)

func main() {
    // Create a simple bitstring
    builder := funbit.NewBuilder()
    funbit.AddInteger(builder, 1, funbit.WithSize(4))          // 4-bit integer
    funbit.AddInteger(builder, 17, funbit.WithSize(12))        // 12-bit integer
    funbit.AddFloat(builder, 3.14, funbit.WithSize(32))        // 32-bit float
    funbit.AddBinary(builder, []byte("hello"))                  // Binary data
    
    bs, err := funbit.Build(builder)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Constructed bitstring: %v\n", bs.ToBytes())
}
```

### Basic Pattern Matching

```go
package main

import (
    "fmt"
    "github.com/funvibe/funbit/pkg/funbit"
)

func main() {
    // Create a bitstring to match against
    bs := funbit.NewBitStringFromBytes([]byte{0x12, 0x34, 0x56, 0x78})
    
    var a, b uint8
    var c uint16
    
    // Match pattern against bitstring
    matcher := funbit.NewMatcher()
    funbit.Integer(matcher, &a, funbit.WithSize(8))           // Match 8-bit integer
    funbit.Integer(matcher, &b, funbit.WithSize(8))           // Match 8-bit integer
    funbit.Integer(matcher, &c, funbit.WithSize(16))          // Match 16-bit integer
    
    results, err := funbit.Match(matcher, bs)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Matched: a=%d, b=%d, c=%d\n", a, b, c)
}
```

## Advanced Usage

### Dynamic Sizing

```go
package main

import (
    "fmt"
    "github.com/funvibe/funbit/pkg/funbit"
)

func main() {
    // Construction with dynamic sizing
    size := uint(5)
    data := []byte{1, 2, 3, 4, 5}
    
    builder := funbit.NewBuilder()
    funbit.AddInteger(builder, size, funbit.WithSize(8))               // Size field
    funbit.AddBinary(builder, data, funbit.WithDynamicSize(&size))     // Data with dynamic size
    
    bs, err := funbit.Build(builder)
    
    // Matching with dynamic sizing
    var matchedSize uint
    var matchedData []byte
    
    matcher := funbit.NewMatcher()
    funbit.Integer(matcher, &matchedSize, funbit.WithSize(8))
    funbit.RegisterVariable(matcher, "size", &matchedSize)
    funbit.Binary(matcher, &matchedData, funbit.WithDynamicSizeExpression("size"))
    
    results, err := funbit.Match(matcher, bs)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Dynamic match: size=%d, data=%v\n", matchedSize, matchedData)
}
```

### Network Packet Example

```go
package main

import (
    "fmt"
    "github.com/funvibe/funbit/pkg/funbit"
)

func main() {
    // Construct a simple IP-like packet
    builder := funbit.NewBuilder()
    funbit.AddInteger(builder, 4, funbit.WithSize(4))                    // Version
    funbit.AddInteger(builder, 5, funbit.WithSize(4))                    // Header length
    funbit.AddInteger(builder, 20, funbit.WithSize(16), funbit.WithEndianness("big")) // Total length
    funbit.AddInteger(builder, 0x1234, funbit.WithSize(16))              // ID
    funbit.AddBinary(builder, []byte("payload data"))                     // Payload
    
    packet, err := funbit.Build(builder)
    
    // Parse the packet
    var version, headerLen uint8
    var totalLength, id uint16
    var payload []byte
    
    matcher := funbit.NewMatcher()
    funbit.Integer(matcher, &version, funbit.WithSize(4))
    funbit.Integer(matcher, &headerLen, funbit.WithSize(4))
    funbit.Integer(matcher, &totalLength, funbit.WithSize(16), funbit.WithEndianness("big"))
    funbit.Integer(matcher, &id, funbit.WithSize(16))
    funbit.Binary(matcher, &payload)
    
    results, err := funbit.Match(matcher, packet)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Packet: version=%d, headerLen=%d, totalLength=%d, id=%04x, payload=%s\n",
        version, headerLen, totalLength, id, string(payload))
}
```

### UTF Encoding

```go
package main

import (
    "fmt"
    "github.com/funvibe/funbit/pkg/funbit"
)

func main() {
    // UTF-8 encoding using utility functions
    text := "Hello, 世界!"
    
    utf8Data, err := funbit.EncodeUTF8(text)
    if err != nil {
        panic(err)
    }
    
    builder := funbit.NewBuilder()
    funbit.AddBinary(builder, utf8Data)
    
    bs, err := funbit.Build(builder)
    
    // UTF-8 decoding using utility functions
    decodedText, err := funbit.DecodeUTF8(utf8Data)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("UTF-8: '%s' -> %v -> '%s'\n", text, bs.ToBytes(), decodedText)
}
```

## API Reference

### Builder API

#### Factory Functions
- `NewBuilder()` - Create a new builder instance

#### Construction Functions
- `AddInteger(builder *Builder, value interface{}, options ...SegmentOption)` - Add integer segment
- `AddFloat(builder *Builder, value float64, options ...SegmentOption)` - Add float segment
- `AddBinary(builder *Builder, data []byte, options ...SegmentOption)` - Add binary segment
- `AddUTF8(builder *Builder, value string)` - Add UTF-8 segment
- `AddUTF16(builder *Builder, value string, options ...SegmentOption)` - Add UTF-16 segment
- `AddUTF32(builder *Builder, value string, options ...SegmentOption)` - Add UTF-32 segment
- `AddBitstring(builder *Builder, value *BitString, options ...SegmentOption)` - Add bitstring segment
- `AddSegment(builder *Builder, segment Segment)` - Add custom segment
- `Build(builder *Builder)` - Build the final bitstring

#### Segment Options
- `WithSize(size uint)` - Set segment size
- `WithType(segmentType string)` - Set segment type
- `WithSigned(signed bool)` - Set signedness
- `WithEndianness(endianness string)` - Set endianness
- `WithUnit(unit uint)` - Set unit size
- `WithDynamicSize(sizeVar *uint)` - Set dynamic size from variable
- `WithDynamicSizeExpression(expr string)` - Set dynamic size from expression

### Matcher API

#### Factory Functions
- `NewMatcher()` - Create a new matcher instance

#### Matching Functions
- `Integer(matcher *Matcher, variable interface{}, options ...SegmentOption)` - Match integer
- `Float(matcher *Matcher, variable *float64, options ...SegmentOption)` - Match float
- `Binary(matcher *Matcher, variable *[]byte, options ...SegmentOption)` - Match binary
- `UTF(matcher *Matcher, variable *string, options ...SegmentOption)` - Match UTF
- `UTF8(matcher *Matcher, variable *string)` - Match UTF-8
- `UTF16(matcher *Matcher, variable *string, options ...SegmentOption)` - Match UTF-16
- `UTF32(matcher *Matcher, variable *string, options ...SegmentOption)` - Match UTF-32
- `Bitstring(matcher *Matcher, variable **BitString, options ...SegmentOption)` - Match bitstring
- `RestBinary(matcher *Matcher, variable *[]byte)` - Match remaining binary data
- `RestBitstring(matcher *Matcher, variable **BitString)` - Match remaining bitstring data
- `RegisterVariable(matcher *Matcher, name string, variable interface{})` - Register variable for dynamic sizing
- `Match(matcher *Matcher, bitstring *BitString)` - Execute pattern matching

### Core Types

#### BitString
- `NewBitString()` - Create empty bitstring
- `NewBitStringFromBytes(data []byte)` - Create from bytes
- `NewBitStringFromBits(data []byte, length uint)` - Create from bits with specific length
- `Length() uint` - Get length in bits
- `ToBytes() []byte` - Convert to byte slice
- `IsEmpty() bool` - Check if empty
- `IsBinary() bool` - Check if length is multiple of 8

#### Segment Types
- `TypeInteger` - Integer values (default 8 bits)
- `TypeFloat` - Floating point values (16, 32, or 64 bits)
- `TypeBinary` - Byte-aligned binary data
- `TypeBitstring` - Arbitrary bit length data
- `TypeUTF8/TypeUTF16/TypeUTF32` - Unicode encoded strings
- `TypeRestBinary/TypeRestBitstring` - Remaining data

#### Utility Functions
- `EncodeUTF8(value string) ([]byte, error)` - Encode string to UTF-8
- `DecodeUTF8(data []byte) (string, error)` - Decode UTF-8 to string
- `EncodeUTF16(value string, endianness string) ([]byte, error)` - Encode string to UTF-16
- `DecodeUTF16(data []byte, endianness string) (string, error)` - Decode UTF-16 to string
- `EncodeUTF32(value string, endianness string) ([]byte, error)` - Encode string to UTF-32
- `DecodeUTF32(data []byte, endianness string) (string, error)` - Decode UTF-32 to string
- `ExtractBits(data []byte, start, length uint) ([]byte, error)` - Extract bits from data
- `CountBits(data []byte) uint` - Count set bits in data
- `IntToBits(value int64, size uint, signed bool) ([]byte, error)` - Convert integer to bits
- `BitsToInt(data []byte, signed bool) (int64, error)` - Convert bits to integer
- `GetNativeEndianness() string` - Get system endianness

## Supported Types and Specifiers

### Data Types
- **Integer**: Signed/unsigned integers with configurable size
- **Float**: 16/32/64-bit floating point numbers
- **Binary**: Byte-aligned data (8-bit units)
- **Bitstring**: Arbitrary bit-length data (1-bit units)
- **UTF-8/16/32**: Unicode encoded strings

### Specifiers
- **Endianness**: `big`, `little`, `native`
- **Signedness**: `signed`, `unsigned` (integers only)
- **Unit**: 1-256 (multiplier for size)
- **Size**: Explicit size or dynamic sizing

## Error Handling

The library provides detailed error information through the `BitStringError` type:

```go
type BitStringError struct {
    Code    string      // Error code
    Message string      // Error message
    Context interface{} // Additional context
}
```

Common error codes:
- `ErrInvalidSize` - Invalid segment size
- `ErrInvalidType` - Unsupported segment type
- `ErrInsufficientBits` - Not enough bits for operation
- `ErrTypeMismatch` - Type conversion error
- `ErrOverflow` - Integer overflow
- `ErrSignedOverflow` - Signed integer overflow
- `ErrInvalidEndianness` - Invalid endianness specification
- `ErrBinarySizeRequired` - Binary segment requires size specification
- `ErrBinarySizeMismatch` - Binary data size doesn't match specified size
- `ErrInvalidBinaryData` - Invalid binary data type
- `ErrInvalidBitstringData` - Invalid bitstring data
- `ErrUTFSizeSpecified` - Size cannot be specified for UTF segments
- `ErrInvalidUnicodeCodepoint` - Invalid Unicode code point
- `ErrInvalidSegment` - Invalid segment configuration
- `ErrInvalidUnit` - Invalid unit value
- `ErrInvalidFloatSize` - Invalid float size (must be 32 or 64)

## Performance Considerations

- **Memory Efficiency**: Minimal allocations during operations
- **Bit-Level Operations**: Optimized for direct bit manipulation
- **Dynamic Sizing**: Efficient expression evaluation
- **Zero-Copy**: Where possible, avoid unnecessary data copying

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.

## Acknowledgments

- Inspired by Erlang/OTP's bit syntax implementation
- Designed for compatibility with existing Erlang binary protocols
- Optimized for Go's performance characteristics