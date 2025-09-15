package main

import (
	"fmt"
	"log"

	"github.com/funvibe/funbit/pkg/funbit"
)

func main() {
	fmt.Println("=== Funbit Public API Example ===\n")

	// Example 1: Basic BitString Construction
	fmt.Println("1. Basic BitString Construction:")
	builder := funbit.NewBuilder()
	funbit.AddInteger(builder, 42)                     // Default 8-bit integer
	funbit.AddInteger(builder, 17, funbit.WithSize(8)) // Explicit 8-bit integer
	funbit.AddBinary(builder, []byte("hello"))         // Binary data

	bs, err := funbit.Build(builder)
	if err != nil {
		log.Fatalf("Failed to build bitstring: %v", err)
	}

	fmt.Printf("   Constructed bitstring: %v bytes\n", len(bs.ToBytes()))
	fmt.Printf("   Bit length: %d bits\n", bs.Length())
	fmt.Printf("   Is binary-aligned: %t\n", bs.IsBinary())
	fmt.Printf("   Hex dump: %s\n", funbit.ToHexDump(bs))
	fmt.Printf("   Erlang format: %s\n", funbit.ToErlangFormat(bs))
	fmt.Println()

	// Example 2: Pattern Matching
	fmt.Println("2. Pattern Matching:")
	var a, b int
	var c []byte

	matcher := funbit.NewMatcher()
	funbit.Integer(matcher, &a)
	funbit.Integer(matcher, &b, funbit.WithSize(8))
	funbit.Binary(matcher, &c)

	results, err := funbit.Match(matcher, bs)
	if err != nil {
		log.Fatalf("Failed to match pattern: %v", err)
	}

	fmt.Printf("   Matched: a=%d, b=%d, c='%s'\n", a, b, string(c))
	fmt.Printf("   Number of results: %d\n", len(results))
	fmt.Println()

	// Example 3: Working with Different Data Types
	fmt.Println("3. Different Data Types:")
	builder2 := funbit.NewBuilder()
	funbit.AddInteger(builder2, -123, funbit.WithSize(16), funbit.WithSigned(true))
	funbit.AddFloat(builder2, 3.14159, funbit.WithSize(32))
	// Add UTF-8 encoded binary data for "Привет" (Russian for "Hello")
	utf8Data, err := funbit.EncodeUTF8("Привет")
	if err != nil {
		log.Fatalf("Failed to encode UTF-8: %v", err)
	}
	funbit.AddBinary(builder2, utf8Data)

	bs2, err := funbit.Build(builder2)
	if err != nil {
		log.Fatalf("Failed to build bitstring: %v", err)
	}

	fmt.Printf("   Mixed types bitstring: %d bits\n", bs2.Length())
	fmt.Printf("   Binary representation: %s\n", funbit.ToBinaryString(bs2))
	fmt.Println()

	// Example 4: Endianness Handling
	fmt.Println("4. Endianness Handling:")
	builder3 := funbit.NewBuilder()
	funbit.AddInteger(builder3, 0x1234, funbit.WithSize(16), funbit.WithEndianness("big"))
	funbit.AddInteger(builder3, 0x1234, funbit.WithSize(16), funbit.WithEndianness("little"))

	bs3, err := funbit.Build(builder3)
	if err != nil {
		log.Fatalf("Failed to build bitstring: %v", err)
	}

	fmt.Printf("   Endianness test: %s\n", funbit.ToHexDump(bs3))
	fmt.Printf("   Native endianness: %s\n", funbit.GetNativeEndianness())
	fmt.Println()

	// Example 5: Utility Functions
	fmt.Println("5. Utility Functions:")
	data := []byte{0xFF, 0x00, 0xF0}

	// Bit manipulation
	bitCount := funbit.CountBits(data)
	fmt.Printf("   Bit count in [0xFF, 0x00, 0xF0]: %d\n", bitCount)

	firstBit, err := funbit.GetBitValue(data, 0)
	if err != nil {
		log.Fatalf("Failed to get bit value: %v", err)
	}
	fmt.Printf("   First bit value: %t\n", firstBit)

	// Type conversion
	intBits, err := funbit.IntToBits(42, 16, false)
	if err != nil {
		log.Fatalf("Failed to convert int to bits: %v", err)
	}
	fmt.Printf("   Int 42 as 16-bit bits: %v\n", intBits)

	convertedInt, err := funbit.BitsToInt(intBits, false)
	if err != nil {
		log.Fatalf("Failed to convert bits to int: %v", err)
	}
	fmt.Printf("   Converted back to int: %d\n", convertedInt)
	fmt.Println()

	// Example 6: UTF Encoding/Decoding
	fmt.Println("6. UTF Encoding/Decoding:")
	text := "Hello, 世界! 🚀"

	// UTF-8
	utf8Encoded, err := funbit.EncodeUTF8(text)
	if err != nil {
		log.Fatalf("Failed to encode UTF-8: %v", err)
	}
	utf8Decoded, err := funbit.DecodeUTF8(utf8Encoded)
	if err != nil {
		log.Fatalf("Failed to decode UTF-8: %v", err)
	}
	fmt.Printf("   UTF-8: '%s' -> %d bytes -> '%s'\n", text, len(utf8Encoded), utf8Decoded)

	// UTF-16
	utf16Encoded, err := funbit.EncodeUTF16(text, "big")
	if err != nil {
		log.Fatalf("Failed to encode UTF-16: %v", err)
	}
	utf16Decoded, err := funbit.DecodeUTF16(utf16Encoded, "big")
	if err != nil {
		log.Fatalf("Failed to decode UTF-16: %v", err)
	}
	fmt.Printf("   UTF-16: '%s' -> %d bytes -> '%s'\n", text, len(utf16Encoded), utf16Decoded)
	fmt.Println()

	// Example 7: Error Handling
	fmt.Println("7. Error Handling:")

	// Try to create invalid segment
	invalidSegment := funbit.NewSegment(42, funbit.WithSize(0))
	err = funbit.ValidateSegment(invalidSegment)
	if err != nil {
		fmt.Printf("   Validation error (expected): %v\n", err)
	}

	// Try to extract bits beyond data range
	_, err = funbit.ExtractBits([]byte{0xFF}, 8, 8)
	if err != nil {
		fmt.Printf("   Extraction error (expected): %v\n", err)
	}

	// Try invalid Unicode code point
	err = funbit.ValidateUnicodeCodePoint(0x110000) // Beyond Unicode range
	if err != nil {
		fmt.Printf("   Unicode validation error (expected): %v\n", err)
	}
	fmt.Println()

	// Example 8: Network Packet Example
	fmt.Println("8. Network Packet Example:")

	// Simulate IPv4-like packet
	packetBuilder := funbit.NewBuilder()
	funbit.AddInteger(packetBuilder, 4, funbit.WithSize(4))                                 // Version
	funbit.AddInteger(packetBuilder, 5, funbit.WithSize(4))                                 // IHL
	funbit.AddInteger(packetBuilder, 20, funbit.WithSize(16), funbit.WithEndianness("big")) // Total length
	funbit.AddInteger(packetBuilder, 0x1234, funbit.WithSize(16))                           // ID
	funbit.AddBinary(packetBuilder, []byte("payload data"))                                 // Payload

	packet, err := funbit.Build(packetBuilder)
	if err != nil {
		log.Fatalf("Failed to build packet: %v", err)
	}

	// Parse the packet back
	var version, ihl uint8
	var totalLength, id uint16
	var payload []byte

	packetMatcher := funbit.NewMatcher()
	funbit.Integer(packetMatcher, &version, funbit.WithSize(4))
	funbit.Integer(packetMatcher, &ihl, funbit.WithSize(4))
	funbit.Integer(packetMatcher, &totalLength, funbit.WithSize(16), funbit.WithEndianness("big"))
	funbit.Integer(packetMatcher, &id, funbit.WithSize(16))
	funbit.Binary(packetMatcher, &payload)

	packetResults, err := funbit.Match(packetMatcher, packet)
	if err != nil {
		log.Fatalf("Failed to parse packet: %v", err)
	}

	fmt.Printf("   Parsed packet: version=%d, ihl=%d, totalLength=%d, id=0x%04X, payload='%s'\n",
		version, ihl, totalLength, id, string(payload))
	fmt.Printf("   Parse results: %d segments matched\n", len(packetResults))
	fmt.Println()

	fmt.Println("=== Example completed successfully! ===")
}
