# Funbit - Erlang/OTP Bit Syntax Library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/funvibe/funbit.svg)](https://pkg.go.dev/github.com/funvibe/funbit)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

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
    "github.com/funvibe/funbit/internal/builder"
    "github.com/funvibe/funbit/internal/bitstring"
)

func main() {
    // Create a simple bitstring
    bs, err := builder.NewBuilder().
        AddInteger(1, builder.WithSize(4)).          // 4-bit integer
        AddInteger(17, builder.WithSize(12)).        // 12-bit integer
        AddFloat(3.14, builder.WithSize(32)).        // 32-bit float
        AddBinary([]byte("hello")).                  // Binary data
        Build()
    
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
    "github.com/funvibe/funbit/internal/matcher"
    "github.com/funvibe/funbit/internal/bitstring"
)

func main() {
    // Create a bitstring to match against
    bs := bitstring.NewBitStringFromBytes([]byte{0x12, 0x34, 0x56, 0x78})
    
    var a, b uint8
    var c uint16
    
    // Match pattern against bitstring
    results, err := matcher.NewMatcher().
        Integer(&a, matcher.WithSize(8)).           // Match 8-bit integer
        Integer(&b, matcher.WithSize(8)).           // Match 8-bit integer
        Integer(&c, matcher.WithSize(16)).          // Match 16-bit integer
        Match(bs)
    
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
    "github.com/funvibe/funbit/internal/builder"
    "github.com/funvibe/funbit/internal/matcher"
    "github.com/funvibe/funbit/internal/bitstring"
)

func main() {
    // Construction with dynamic sizing
    size := uint(5)
    data := []byte{1, 2, 3, 4, 5}
    
    bs, err := builder.NewBuilder().
        AddInteger(size, builder.WithSize(8)).               // Size field
        AddBinary(data, builder.WithDynamicSize(&size)).     // Data with dynamic size
        Build()
    
    // Matching with dynamic sizing
    var matchedSize uint
    var matchedData []byte
    
    results, err := matcher.NewMatcher().
        Integer(&matchedSize, matcher.WithSize(8)).
        RegisterVariable("size", &matchedSize).
        Binary(&matchedData, matcher.WithDynamicSizeExpression("size")).
        Match(bs)
    
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
    "github.com/funvibe/funbit/internal/builder"
    "github.com/funvibe/funbit/internal/matcher"
    "github.com/funvibe/funbit/internal/bitstring"
)

func main() {
    // Construct a simple IP-like packet
    packet, err := builder.NewBuilder().
        AddInteger(4, builder.WithSize(4)).                    // Version
        AddInteger(5, builder.WithSize(4)).                    // Header length
        AddInteger(20, builder.WithSize(16), builder.WithEndianness("big")). // Total length
        AddInteger(0x1234, builder.WithSize(16)).              // ID
        AddBinary([]byte("payload data")).                     // Payload
        Build()
    
    // Parse the packet
    var version, headerLen uint8
    var totalLength, id uint16
    var payload []byte
    
    results, err := matcher.NewMatcher().
        Integer(&version, matcher.WithSize(4)).
        Integer(&headerLen, matcher.WithSize(4)).
        Integer(&totalLength, matcher.WithSize(16), matcher.WithEndianness("big")).
        Integer(&id, matcher.WithSize(16)).
        Binary(&payload).
        Match(packet)
    
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
    "github.com/funvibe/funbit/internal/builder"
    "github.com/funvibe/funbit/internal/matcher"
    "github.com/funvibe/funbit/internal/bitstring"
)

func main() {
    // UTF-8 encoding
    text := "Hello, 世界!"
    
    bs, err := builder.NewBuilder().
        AddUTF(text, builder.WithType("utf8")).
        Build()
    
    // UTF-8 decoding
    var decodedText string
    
    results, err := matcher.NewMatcher().
        UTF(&decodedText, matcher.WithType("utf8")).
        Match(bs)
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("UTF-8: '%s' -> %v -> '%s'\n", text, bs.ToBytes(), decodedText)
}
```

## API Reference

### Builder API

#### Construction Methods
- `NewBuilder()` - Create a new builder instance
- `AddInteger(value interface{}, options ...SegmentOption)` - Add integer segment
- `AddFloat(value float64, options ...SegmentOption)` - Add float segment
- `AddBinary(data []byte, options ...SegmentOption)` - Add binary segment
- `AddUTF(value string, options ...SegmentOption)` - Add UTF segment
- `AddBitstring(value *BitString, options ...SegmentOption)` - Add bitstring segment
- `AddSegment(segment *Segment)` - Add custom segment
- `Build()` - Build the final bitstring

#### Segment Options
- `WithSize(size uint)` - Set segment size
- `WithType(segmentType string)` - Set segment type
- `WithSigned(signed bool)` - Set signedness
- `WithEndianness(endianness string)` - Set endianness
- `WithUnit(unit uint)` - Set unit size
- `WithDynamicSize(sizeVar *uint)` - Set dynamic size from variable
- `WithDynamicSizeExpression(expr string)` - Set dynamic size from expression

### Matcher API

#### Matching Methods
- `NewMatcher()` - Create a new matcher instance
- `Integer(variable interface{}, options ...SegmentOption)` - Match integer
- `Float(variable *float64, options ...SegmentOption)` - Match float
- `Binary(variable *[]byte, options ...SegmentOption)` - Match binary
- `UTF(variable *string, options ...SegmentOption)` - Match UTF
- `Bitstring(variable **BitString, options ...SegmentOption)` - Match bitstring
- `RestBinary(variable *[]byte)` - Match remaining binary data
- `RestBitstring(variable **BitString)` - Match remaining bitstring data
- `RegisterVariable(name string, variable interface{})` - Register variable for dynamic sizing
- `Match(bitstring *BitString)` - Execute pattern matching

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
- `INVALID_SIZE` - Invalid segment size
- `INVALID_TYPE` - Unsupported segment type
- `INSUFFICIENT_BITS` - Not enough bits for operation
- `TYPE_MISMATCH` - Type conversion error

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

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Erlang/OTP's bit syntax implementation
- Designed for compatibility with existing Erlang binary protocols
- Optimized for Go's performance characteristics