# Funbit - Erlang/OTP Bit Syntax Library for Go

Funbit is a comprehensive Go library that provides Erlang/OTP bit syntax compatibility for working with bitstrings and binary data. It offers a fluent interface for both constructing and pattern matching bitstrings with full support for dynamic sizing, various data types, and endianness handling.

## Features

- **Erlang/OTP Bit Syntax Compatibility**: Full support for Erlang's bit syntax expressions for construction and matching.
- **True Bit-Level Operations**: The builder operates as a true bit stream, allowing for unaligned data construction.
- **Multiple Data Types**: Integer, float, binary, bitstring, UTF-8/16/32.
- **Dynamic Sizing**: Support for variable-sized segments using expressions in the Matcher.
- **Endianness Support**: Big, little, and native endianness handling for byte-aligned segments.
- **Fluent Interface**: Clean, chainable API for both construction and matching.

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

## Usage Examples

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
	builder.AddInteger(1, funbit.WithSize(4))          // 4-bit integer
	builder.AddInteger(17, funbit.WithSize(12))        // 12-bit integer
	builder.AddFloat(3.14, funbit.WithSize(32))        // 32-bit float
	builder.AddBinary([]byte("hello"))                 // Binary data
	
	bs, err := builder.Build()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Constructed bitstring: %x\n", bs.ToBytes())
}
```

### Advanced: Unaligned Data and Concatenation

This example demonstrates how to correctly build a bitstring by iteratively adding unaligned, 1-bit flags. This is the correct pattern for `New = <<Old/bitstring, NewBit:1>>`.

```go
package main

import (
	"fmt"
	"github.com/funvibe/funbit/pkg/funbit"
)

func main() {
	flags := []bool{true, false, true} // We want to build the bit sequence 101

	// Start with an empty bitstring
	currentBitstring, _ := funbit.NewBuilder().Build()

	fmt.Printf("Start: len=%d bits, data=%x\n", currentBitstring.Length(), currentBitstring.ToBytes())

	for i, flag := range flags {
		// To concatenate, create a new builder for each step
		builder := funbit.NewBuilder()

		// 1. Add the previous bitstring as the first segment
		builder.AddBitstring(currentBitstring)

		// 2. Add the new 1-bit flag
		var bit uint = 0
		if flag {
			bit = 1
		}
		builder.AddInteger(bit, funbit.WithSize(1))

		// 3. Build the new, longer bitstring
		currentBitstring, _ = builder.Build()
		
		fmt.Printf("Step %d: Added bit %d, new len=%d bits, data=%x\n", i+1, bit, currentBitstring.Length(), currentBitstring.ToBytes())
	}
    
    // Final result is 3 bits long, stored in a single byte: 10100000 (0xA8)
	fmt.Printf("Final: len=%d bits, data=%x\n", currentBitstring.Length(), currentBitstring.ToBytes())
}
```

### Pattern Matching with Endianness

```go
package main

import (
	"fmt"
	"github.com/funvibe/funbit/pkg/funbit"
)

func main() {
	// 1. Construct a 16-bit little-endian integer (0x1234)
	builder := funbit.NewBuilder()
	val := uint16(0x1234)
	builder.AddInteger(val, funbit.WithSize(16), funbit.WithEndianness("little"))
	bs, _ := builder.Build()

	// In memory, this will be [0x34, 0x12]
	fmt.Printf("Little-endian bytes: %x\n", bs.ToBytes())

	// 2. Match the little-endian integer
	var matchedVal uint16
	matcher := funbit.NewMatcher()
	matcher.Integer(&matchedVal, funbit.WithSize(16), funbit.WithEndianness("little"))
	
	_, err := matcher.Match(bs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Matched value: 0x%x\n", matchedVal)

	// NOTE: Endianness only applies to segments with a size that is a multiple of 8.
}
```

## API Reference

(This section can be expanded, but the examples above cover the most critical usage patterns.)

### Builder API
- `NewBuilder()`
- `AddInteger(value interface{}, options ...SegmentOption)`
- `AddFloat(value float64, options ...SegmentOption)`
- `AddBinary(data []byte, options ...SegmentOption)`
- `AddBitstring(value *BitString, options ...SegmentOption)`
- `Build()`

### Matcher API
- `NewMatcher()`
- `Integer(variable interface{}, options ...SegmentOption)`
- `Float(variable *float64, options ...SegmentOption)`
- `Binary(variable *[]byte, options ...SegmentOption)`
- `RestBinary(variable *[]byte)`
- `RestBitstring(variable **BitString)`
- `Match(bitstring *BitString)`
- `RegisterVariable(name string, variable interface{})` for dynamic sized matching.

### Segment Options
- `WithSize(size uint)`: Sets segment size **in bits** for `Integer`, `Float`, and `Bitstring`. For `Binary`, size is in multiples of the `Unit` (default 8 bits).
- `WithEndianness(endianness string)`: `big`, `little`, `native`.
- `WithSigned(signed bool)`
- `WithUnit(unit uint)`