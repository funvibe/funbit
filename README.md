# Funbit - Erlang/OTP Bit Syntax Library for Go

Funbit is a comprehensive Go library that provides Erlang/OTP bit syntax compatibility for working with bitstrings and binary data. It offers a fluent interface for both constructing and pattern matching bitstrings with full support for dynamic sizing, various data types, endianness handling, and advanced bit manipulation operations.

## Features

- **Erlang/OTP Bit Syntax Compatibility**: Full support for Erlang's bit syntax expressions for construction and matching
- **True Bit-Level Operations**: The builder operates as a true bit stream, allowing for unaligned data construction
- **Multiple Data Types**: Integer, float, binary, bitstring, UTF-8/16/32 encoding
- **Dynamic Sizing**: Support for variable-sized segments using expressions in the Matcher
- **Unit Specifiers**: Advanced size control with customizable unit values (1-256 bits)
- **Endianness Support**: Big, little, and native endianness handling
- **Bit-Level Manipulation**: Extract, manipulate, and convert individual bits
- **Protocol Support**: Ready for parsing real-world protocols (IPv4, TCP, PNG, etc.)
- **UTF Encoding**: Full Unicode support with UTF-8/16/32 encoding/decoding
- **Fluent Interface**: Clean, chainable API for both construction and matching

## Installation

```bash
go get github.com/funvibe/funbit
```

## Core Concepts: How Funbit Thinks

To use Funbit effectively, it's crucial to understand these core principles, which may not be immediately obvious.

### 1. The Builder is a True Bit Stream

The `funbit.Builder` is not a segment assembler; it operates as a **true bit stream writer**. When you add a segment, its bits are appended to the stream, regardless of byte boundaries.

- `builder.AddInteger(1, funbit.WithSize(1))` adds **exactly one bit**.
- `builder.AddInteger(0, funbit.WithSize(3))` then adds **exactly three bits**.

The builder does not automatically pad segments to the nearest byte. This allows for the creation of tightly packed, unaligned binary data.

### 2. Concatenation Requires Re-Building

Bitstrings in Funbit are immutable. You cannot "append" bits to an existing `BitString` object. To achieve concatenation, as in the Erlang expression `New = <<Old/bitstring, ...>>`, you must create a new builder, add the old bitstring as the first segment, and then add the new segments.

```go
// Correct way to concatenate
builder := funbit.NewBuilder()
builder.AddBitstring(oldBitstring) // Add the old bitstring first
builder.AddInteger(newValue, funbit.WithSize(4))
newBitstring, _ := builder.Build()
```

### 3. Endianness Applies Only to Byte-Aligned Segments

The `WithEndianness` option affects segments whose total size in bits is a multiple of 8 (e.g., 16, 24, 32, 64 bits). For segments with unaligned sizes (e.g., `12` bits), the endianness specifier is ignored, and the bits are written in a big-endian fashion.

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/funvibe/funbit/pkg/funbit"
)

func main() {
	// Basic Construction
	builder := funbit.NewBuilder()
	funbit.AddInteger(builder, 42)                     // 8-bit integer (default)
	funbit.AddInteger(builder, 17, funbit.WithSize(8)) // Explicit 8-bit integer
	funbit.AddBinary(builder, []byte("hello"))         // Binary data

	bitstring, err := funbit.Build(builder)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Built: %d bytes, %d bits\n", len(bitstring.ToBytes()), bitstring.Length())

	// Pattern Matching
	var a, b int
	var data []byte

	matcher := funbit.NewMatcher()
	funbit.Integer(matcher, &a)
	funbit.Integer(matcher, &b, funbit.WithSize(8))
	funbit.Binary(matcher, &data)

	results, err := funbit.Match(matcher, bitstring)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched: a=%d, b=%d, data='%s' (%d segments)\n", a, b, string(data), len(results))
}
```

## Advanced Examples

### Unit Specifiers - Advanced Size Control

```go
// Unit specifiers control how Size * Unit = TotalBits
builder := funbit.NewBuilder()
funbit.AddInteger(builder, 15, funbit.WithSize(4), funbit.WithUnit(1))  // 4*1 = 4 bits
funbit.AddInteger(builder, 1, funbit.WithSize(8), funbit.WithUnit(1))   // 8*1 = 8 bits
bitstring, _ := funbit.Build(builder)

fmt.Printf("Total: %d bits\n", bitstring.Length()) // 12 bits

// Parse back with same units
var a, b uint
matcher := funbit.NewMatcher()
funbit.Integer(matcher, &a, funbit.WithSize(4), funbit.WithUnit(1))
funbit.Integer(matcher, &b, funbit.WithSize(8), funbit.WithUnit(1))
results, _ := funbit.Match(matcher, bitstring)
```

### Dynamic Size Expressions

```go
// Create packet: <<5:8, "Hello":5/binary, "World">>
builder := funbit.NewBuilder()
funbit.AddInteger(builder, 5, funbit.WithSize(8))
funbit.AddBinary(builder, []byte("Hello"), funbit.WithSize(5))
funbit.AddBinary(builder, []byte("World"))
packet, _ := funbit.Build(builder)

// Parse with dynamic size
var size uint
var data, rest []byte
matcher := funbit.NewMatcher()
funbit.RegisterVariable(matcher, "size", &size)
funbit.Integer(matcher, &size, funbit.WithSize(8))
funbit.Binary(matcher, &data, funbit.WithDynamicSizeExpression("size"))
funbit.Binary(matcher, &rest)
results, _ := funbit.Match(matcher, packet)

fmt.Printf("Size: %d, Data: '%s', Rest: '%s'\n", size, string(data), string(rest))
```

### Bit-Level Manipulation

```go
// Extract individual bits and ranges
data := []byte{0xB4} // 10110100 in binary

bit0, _ := funbit.GetBitValue(data, 0) // LSB
bit7, _ := funbit.GetBitValue(data, 7) // MSB
fmt.Printf("Byte 0xB4: bit0=%t, bit7=%t\n", bit0, bit7)

// Extract bit range
bits3to5, _ := funbit.ExtractBits(data, 3, 3) // bits 3,4,5
fmt.Printf("Bits 3-5: %08b\n", bits3to5[0]>>5)

// Convert between int and bits
intVal := 42
bits, _ := funbit.IntToBits(int64(intVal), 8, false)
fmt.Printf("Int %d as bits: %08b\n", intVal, bits[0])
```

### Endianness Handling

```go
// Mixed endianness in single construction
builder := funbit.NewBuilder()
funbit.AddInteger(builder, 0x1234, funbit.WithSize(16), funbit.WithEndianness("big"))
funbit.AddInteger(builder, 0x5678, funbit.WithSize(16), funbit.WithEndianness("little"))
bitstring, _ := funbit.Build(builder)

fmt.Printf("Mixed endian: %s\n", funbit.ToHexDump(bitstring))
fmt.Printf("Native endianness: %s\n", funbit.GetNativeEndianness())
```

### Binary vs Bitstring Types

```go
// Binary: byte-aligned, Bitstring: bit-aligned
builder := funbit.NewBuilder()
funbit.AddInteger(builder, 42, funbit.WithSize(10)) // 10 bits - not byte aligned
bitstring10, _ := funbit.Build(builder)

builder2 := funbit.NewBuilder()
funbit.AddBinary(builder2, []byte("AB")) // 16 bits - byte aligned
binary16, _ := funbit.Build(builder2)

fmt.Printf("10-bit value: IsBinary=%t\n", bitstring10.IsBinary()) // false
fmt.Printf("16-bit binary: IsBinary=%t\n", binary16.IsBinary())   // true
```

### UTF Encoding/Decoding

```go
text := "Hello, 世界! 🚀"

// UTF-8 encoding/decoding
utf8Encoded, _ := funbit.EncodeUTF8(text)
utf8Decoded, _ := funbit.DecodeUTF8(utf8Encoded)
fmt.Printf("UTF-8: '%s' -> %d bytes -> '%s'\n", text, len(utf8Encoded), utf8Decoded)

// UTF-16 encoding/decoding
utf16Encoded, _ := funbit.EncodeUTF16(text, "big")
utf16Decoded, _ := funbit.DecodeUTF16(utf16Encoded, "big")
fmt.Printf("UTF-16: '%s' -> %d bytes -> '%s'\n", text, len(utf16Encoded), utf16Decoded)
```

### Real-World Protocol: IPv4 Header

```go
// IPv4 header data (simplified)
ipv4Data := []byte{
	0x45, 0x00, 0x00, 0x3C, // Version+IHL, TOS, Total Length
	0x30, 0x39, 0x00, 0x00, // ID, Flags+Fragment Offset
	0x40, 0x06, 0xAB, 0xCD, // TTL, Protocol, Checksum
	0xC0, 0xA8, 0x01, 0x01, // Source IP
	0x0A, 0x00, 0x00, 0x01, // Destination IP
}

packet := funbit.NewBitStringFromBytes(ipv4Data)

// Parse IPv4 header fields
var version, ihl, tos uint8
var totalLen, id uint16
var flags, fragOff uint
var ttl, protocol uint8
var checksum uint16
var srcIP, dstIP uint32

matcher := funbit.NewMatcher()
funbit.Integer(matcher, &version, funbit.WithSize(4))
funbit.Integer(matcher, &ihl, funbit.WithSize(4))
funbit.Integer(matcher, &tos, funbit.WithSize(8))
funbit.Integer(matcher, &totalLen, funbit.WithSize(16), funbit.WithEndianness("big"))
funbit.Integer(matcher, &id, funbit.WithSize(16), funbit.WithEndianness("big"))
funbit.Integer(matcher, &flags, funbit.WithSize(3))
funbit.Integer(matcher, &fragOff, funbit.WithSize(13))
funbit.Integer(matcher, &ttl, funbit.WithSize(8))
funbit.Integer(matcher, &protocol, funbit.WithSize(8))
funbit.Integer(matcher, &checksum, funbit.WithSize(16), funbit.WithEndianness("big"))
funbit.Integer(matcher, &srcIP, funbit.WithSize(32), funbit.WithEndianness("big"))
funbit.Integer(matcher, &dstIP, funbit.WithSize(32), funbit.WithEndianness("big"))

results, _ := funbit.Match(matcher, packet)
fmt.Printf("IPv4 parsed: Version=%d, TTL=%d, Protocol=%d\n", version, ttl, protocol)
```

## Erlang Bit Syntax Equivalents

Funbit provides direct equivalents to Erlang's bit syntax expressions:

| Erlang Expression | Funbit Equivalent | Description |
|-------------------|-------------------|-------------|
| `<<42>>` | `AddInteger(builder, 42)` | 8-bit integer (default) |
| `<<42:16>>` | `AddInteger(builder, 42, WithSize(16))` | 16-bit integer |
| `<<42:16/little>>` | `AddInteger(builder, 42, WithSize(16), WithEndianness("little"))` | 16-bit little-endian |
| `<<"hello">>` | `AddBinary(builder, []byte("hello"))` | Binary string |
| `<<Value:Size/binary>>` | `Binary(matcher, &dest, WithDynamicSizeExpression("Size"))` | Dynamic binary size |
| `<<Data:10/bitstring>>` | `Bitstring(matcher, &dest, WithSize(10))` | 10-bit bitstring |
| `<<X:4/unit:1>>` | `AddInteger(builder, X, WithSize(4), WithUnit(1))` | Unit specifier (4*1=4 bits) |
| `<<Float:32/float>>` | `AddFloat(builder, Float, WithSize(32))` | 32-bit float |
| `<<$a, $b, $c>>` | `AddBinary(builder, []byte("abc"))` | Character bytes |

## Performance Notes

- **Builder is a true bit stream**: No automatic byte alignment - bits are appended exactly as specified
- **Immutable bitstrings**: Use new builders for concatenation rather than mutation
- **Unit specifiers**: Use unit:1 for bit-level control, unit:8 for byte-aligned operations
- **Dynamic sizing**: Use `RegisterVariable` for variable-length field parsing
- **Memory efficient**: Bitstrings store data in minimal byte arrays with length tracking

## Contributing

Funbit is designed to provide comprehensive Erlang bit syntax compatibility. Areas for contribution:

- Additional protocol parsers (TCP, UDP, HTTP, etc.)
- Performance optimizations
- Extended UTF support
- More comprehensive test coverage

## License

See LICENSE.md for licensing information.

## API Reference

### Core Types
- `Builder` - Constructs bitstrings with fluent interface
- `Matcher` - Pattern matches bitstrings with variable binding
- `BitString` - Immutable bitstring representation

### Builder Functions
- `NewBuilder() *Builder`
- `AddInteger(b *Builder, value interface{}, options ...SegmentOption) *Builder`
- `AddFloat(b *Builder, value interface{}, options ...SegmentOption) *Builder`
- `AddBinary(b *Builder, value []byte, options ...SegmentOption) *Builder`
- `AddBitstring(b *Builder, value *BitString, options ...SegmentOption) *Builder`
- `AddSegment(b *Builder, segment Segment) *Builder`
- `Build(b *Builder) (*BitString, error)`

### Matcher Functions
- `NewMatcher() *Matcher`
- `Integer(m *Matcher, dest interface{}, options ...SegmentOption)`
- `Float(m *Matcher, dest interface{}, options ...SegmentOption)`
- `Binary(m *Matcher, dest *[]byte, options ...SegmentOption)`
- `Bitstring(m *Matcher, dest **BitString, options ...SegmentOption)`
- `Match(m *Matcher, bs *BitString) ([]MatchResult, error)`
- `RegisterVariable(m *Matcher, name string, variable interface{})`

### BitString Constructors
- `NewBitString() *BitString`
- `NewBitStringFromBytes(data []byte) *BitString`
- `NewBitStringFromBits(data []byte, length uint) *BitString`

### Utility Functions
- `CountBits(data []byte) int` - Count set bits in byte array
- `ToHexDump(bs *BitString) string` - Format as hex dump
- `ToErlangFormat(bs *BitString) string` - Format as Erlang binary syntax
- `ToBinaryString(bs *BitString) string` - Format as binary string
- `GetBitValue(data []byte, bitIndex uint) (bool, error)` - Extract single bit
- `ExtractBits(data []byte, startBit, numBits uint) ([]byte, error)` - Extract bit range
- `IntToBits(value int64, size uint, signed bool) ([]byte, error)` - Convert int to bits
- `BitsToInt(bits []byte, signed bool) (int64, error)` - Convert bits to int
- `GetNativeEndianness() string` - Get system endianness

### UTF Encoding Functions
- `EncodeUTF8(text string) ([]byte, error)`
- `DecodeUTF8(data []byte) (string, error)`
- `EncodeUTF16(text string, endianness string) ([]byte, error)`
- `DecodeUTF16(data []byte, endianness string) (string, error)`
- `EncodeUTF32(text string, endianness string) ([]byte, error)`
- `DecodeUTF32(data []byte, endianness string) (string, error)`
- `ValidateUnicodeCodePoint(codePoint int) error`

### Segment Options
- `WithSize(size uint)` - Segment size in bits (integer/float) or units (binary)
- `WithType(typeStr string)` - Data type: "integer", "float", "binary", "bitstring", "utf8", "utf16", "utf32"
- `WithSigned(signed bool)` - Signed/unsigned for integer types
- `WithEndianness(endianness string)` - "big", "little", "native"
- `WithUnit(unit uint)` - Unit size (1-256), affects Size calculation
- `WithDynamicSizeExpression(expr string)` - Dynamic size using variables

### Constants
- `TypeInteger`, `TypeFloat`, `TypeBinary`, `TypeBitstring`, `TypeUTF`, `TypeUTF8`, `TypeUTF16`, `TypeUTF32`
- `EndiannessBig`, `EndiannessLittle`, `EndiannessNative`

### Error Handling
- `ConvertFunbitError(err error) error` - Convert funbit errors to standard format