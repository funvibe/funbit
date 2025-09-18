package main

import (
	"fmt"
	"log"

	"github.com/funvibe/funbit/pkg/funbit"
)

func main() {
	fmt.Println("=== Funbit Public API Example ===")

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
	var a1, b1 int
	var c1 []byte

	matcher1 := funbit.NewMatcher()
	funbit.Integer(matcher1, &a1)
	funbit.Integer(matcher1, &b1, funbit.WithSize(8))
	funbit.Binary(matcher1, &c1)

	results1, err := funbit.Match(matcher1, bs)
	if err != nil {
		log.Fatalf("Failed to match pattern: %v", err)
	}

	fmt.Printf("   Matched: a=%d, b=%d, c='%s'\n", a1, b1, string(c1))
	fmt.Printf("   Number of results: %d\n", len(results1))
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
	packetBuilder8 := funbit.NewBuilder()
	funbit.AddInteger(packetBuilder8, 4, funbit.WithSize(4))                                 // Version
	funbit.AddInteger(packetBuilder8, 5, funbit.WithSize(4))                                 // IHL
	funbit.AddInteger(packetBuilder8, 20, funbit.WithSize(16), funbit.WithEndianness("big")) // Total length
	funbit.AddInteger(packetBuilder8, 0x1234, funbit.WithSize(16))                           // ID
	funbit.AddBinary(packetBuilder8, []byte("payload data"))                                 // Payload

	packet8, err := funbit.Build(packetBuilder8)
	if err != nil {
		log.Fatalf("Failed to build packet: %v", err)
	}

	// Parse the packet back
	var version8, ihl8 uint8
	var totalLength8, id8 uint16
	var payload8 []byte

	packetMatcher8 := funbit.NewMatcher()
	funbit.Integer(packetMatcher8, &version8, funbit.WithSize(4))
	funbit.Integer(packetMatcher8, &ihl8, funbit.WithSize(4))
	funbit.Integer(packetMatcher8, &totalLength8, funbit.WithSize(16), funbit.WithEndianness("big"))
	funbit.Integer(packetMatcher8, &id8, funbit.WithSize(16))
	funbit.Binary(packetMatcher8, &payload8)

	packetResults8, err := funbit.Match(packetMatcher8, packet8)
	if err != nil {
		log.Fatalf("Failed to parse packet: %v", err)
	}

	fmt.Printf("   Parsed packet: version=%d, ihl=%d, totalLength=%d, id=0x%04X, payload='%s'\n",
		version8, ihl8, totalLength8, id8, string(payload8))
	fmt.Printf("   Parse results: %d segments matched\n", len(packetResults8))
	fmt.Println()

	// Example 9: Unit Specifiers - Advanced Size Control
	fmt.Println("9. Unit Specifiers - Advanced Size Control:")
	fmt.Printf("   Unit specifiers control how Size * Unit = TotalBits\n")
	fmt.Printf("   Integer default: unit=1, Binary default: unit=8\n\n")

	// Erlang equivalent: <<15:4/unit:1, 1:8/unit:8>>
	builder9 := funbit.NewBuilder()
	funbit.AddInteger(builder9, 15, funbit.WithSize(4), funbit.WithUnit(1)) // 4*1 = 4 bits
	funbit.AddInteger(builder9, 1, funbit.WithSize(8), funbit.WithUnit(1))  // 8*1 = 8 bits (default unit for int)

	bs9, err := funbit.Build(builder9)
	if err != nil {
		log.Fatalf("Failed to build unit example: %v", err)
	}

	fmt.Printf("   Unit example: %d bits total\n", bs9.Length())
	fmt.Printf("   Hex dump: %s\n", funbit.ToHexDump(bs9))

	// Parse back with same units
	var a9, b9 uint
	matcher9 := funbit.NewMatcher()
	funbit.Integer(matcher9, &a9, funbit.WithSize(4), funbit.WithUnit(1))
	funbit.Integer(matcher9, &b9, funbit.WithSize(8), funbit.WithUnit(1))

	results9, err := funbit.Match(matcher9, bs9)
	if err != nil {
		log.Fatalf("Failed to match unit example: %v", err)
	}
	fmt.Printf("   Parsed: a=%d, b=%d (%d results)\n", a9, b9, len(results9))
	fmt.Printf("   Demonstrates: Size*Unit calculation (4*1=4 bits, 8*1=8 bits)\n")
	fmt.Println()

	// Example 10: Dynamic Size Expressions
	fmt.Println("10. Dynamic Size Expressions:")
	fmt.Printf("   Use RegisterVariable and WithDynamicSizeExpression for variable-length fields\n")

	// Create packet: <<5:8, "Hello":5/binary, "World">>
	// Erlang equivalent: <<Size:8, Data:Size/binary, Rest/binary>>
	builder10 := funbit.NewBuilder()
	funbit.AddInteger(builder10, 5, funbit.WithSize(8))
	funbit.AddBinary(builder10, []byte("Hello"), funbit.WithSize(5))
	funbit.AddBinary(builder10, []byte("World"))

	bs10, err := funbit.Build(builder10)
	if err != nil {
		log.Fatalf("Failed to build dynamic size example: %v", err)
	}

	// Parse with dynamic size
	var size10 uint
	var data10, rest10 []byte
	matcher10 := funbit.NewMatcher()
	funbit.RegisterVariable(matcher10, "size", &size10)
	funbit.Integer(matcher10, &size10, funbit.WithSize(8))
	funbit.Binary(matcher10, &data10, funbit.WithDynamicSizeExpression("size"))
	funbit.Binary(matcher10, &rest10)

	results10, err := funbit.Match(matcher10, bs10)
	_ = results10 // Avoid unused variable warning
	if err != nil {
		log.Fatalf("Failed to match dynamic size: %v", err)
	}
	fmt.Printf("   Dynamic size parsed: size=%d, data='%s', rest='%s'\n", size10, string(data10), string(rest10))

	// Example with expression: <<10:8, "DATA":4/binary, "EXTRA":5/binary, "END">>
	// Erlang equivalent: <<Total:8, Payload:(Total-6)/binary, Trailer/binary>>
	builder10b := funbit.NewBuilder()
	funbit.AddInteger(builder10b, 10, funbit.WithSize(8))
	funbit.AddBinary(builder10b, []byte("DATA"), funbit.WithSize(4))
	funbit.AddBinary(builder10b, []byte("EXTRA"), funbit.WithSize(5))
	funbit.AddBinary(builder10b, []byte("END"))

	bs10b, err := funbit.Build(builder10b)
	if err != nil {
		log.Fatalf("Failed to build expression example: %v", err)
	}

	var total10b uint
	var payload10b, trailer10b []byte
	matcher10b := funbit.NewMatcher()
	funbit.RegisterVariable(matcher10b, "total", &total10b)
	funbit.Integer(matcher10b, &total10b, funbit.WithSize(8))
	funbit.Binary(matcher10b, &payload10b, funbit.WithDynamicSizeExpression("total-6"))
	funbit.Binary(matcher10b, &trailer10b)

	results10b, err := funbit.Match(matcher10b, bs10b)
	_ = results10b // Avoid unused variable warning
	if err != nil {
		log.Fatalf("Failed to match expression: %v", err)
	}
	fmt.Printf("   Expression parsed: total=%d, payload='%s', trailer='%s'\n", total10b, string(payload10b), string(trailer10b))
	fmt.Println()

	// Example 11: Bit-Level Manipulation
	fmt.Println("11. Bit-Level Manipulation:")
	fmt.Printf("   Extract and manipulate individual bits\n")

	// Create a byte with pattern 10110100 (0xB4)
	testData := []byte{0xB4} // 10110100 in binary

	// Get individual bits
	bit0, _ := funbit.GetBitValue(testData, 0) // LSB
	bit7, _ := funbit.GetBitValue(testData, 7) // MSB
	fmt.Printf("   Byte 0xB4 (0b%b): bit0=%t, bit7=%t\n", testData[0], bit0, bit7)

	// Extract bit ranges
	bits3to5, _ := funbit.ExtractBits(testData, 3, 3) // bits 3,4,5 (101)
	fmt.Printf("   Bits 3-5: 0b%b (decimal: %d)\n", bits3to5[0]>>5, bits3to5[0]>>5)

	// Convert between int and bits
	intVal := 42
	bits, _ := funbit.IntToBits(int64(intVal), 8, false)
	fmt.Printf("   Int %d as bits: %08b\n", intVal, bits[0])

	convertedBack, _ := funbit.BitsToInt(bits, false)
	fmt.Printf("   Converted back: %d\n", convertedBack)
	fmt.Println()

	// Example 12: Binary vs Bitstring Types
	fmt.Println("12. Binary vs Bitstring Types:")
	fmt.Printf("   Binary: byte-aligned, Bitstring: bit-aligned\n")

	// Create 10-bit bitstring (not byte-aligned)
	builder12 := funbit.NewBuilder()
	funbit.AddInteger(builder12, 42, funbit.WithSize(10)) // 10 bits
	bs12, _ := funbit.Build(builder12)

	fmt.Printf("   10-bit value: %d bits, IsBinary: %t\n", bs12.Length(), bs12.IsBinary())

	// Binary segment (requires byte alignment)
	builder12b := funbit.NewBuilder()
	funbit.AddBinary(builder12b, []byte("AB")) // 16 bits, byte-aligned
	bs12b, _ := funbit.Build(builder12b)

	fmt.Printf("   Binary data: %d bits, IsBinary: %t\n", bs12b.Length(), bs12b.IsBinary())

	// Bitstring can handle non-byte-aligned data
	matcher12 := funbit.NewMatcher()
	funbit.Bitstring(matcher12, &bs12)
	results12, _ := funbit.Match(matcher12, bs12)
	_ = results12 // Avoid unused variable warning
	fmt.Printf("   Bitstring extracted: %d bits\n", bs12.Length())
	fmt.Println()

	// Example 13: Real-World Protocol - IPv4 Header
	fmt.Println("13. Real-World Protocol - IPv4 Header:")
	fmt.Printf("   Parse IPv4-like header using bit syntax patterns\n")

	// Create IPv4-like header (simplified)
	ipv4Data := []byte{
		0x45, 0x00, 0x00, 0x3C, // Version+IHL, TOS, Total Length
		0x30, 0x39, 0x00, 0x00, // ID, Flags+Fragment Offset
		0x40, 0x06, 0xAB, 0xCD, // TTL, Protocol, Checksum
		0xC0, 0xA8, 0x01, 0x01, // Source IP
		0x0A, 0x00, 0x00, 0x01, // Destination IP
	}

	ipv4Packet := funbit.NewBitStringFromBytes(ipv4Data)

	// Parse IPv4 header fields
	var version13, ihl13, tos13 uint8
	var totalLen13, id13 uint16
	var flags13, fragOff13 uint
	var ttl13, protocol13 uint8
	var checksum13 uint16
	var srcIP13, dstIP13 uint32

	ipv4Matcher := funbit.NewMatcher()
	funbit.Integer(ipv4Matcher, &version13, funbit.WithSize(4))
	funbit.Integer(ipv4Matcher, &ihl13, funbit.WithSize(4))
	funbit.Integer(ipv4Matcher, &tos13, funbit.WithSize(8))
	funbit.Integer(ipv4Matcher, &totalLen13, funbit.WithSize(16), funbit.WithEndianness("big"))
	funbit.Integer(ipv4Matcher, &id13, funbit.WithSize(16), funbit.WithEndianness("big"))
	funbit.Integer(ipv4Matcher, &flags13, funbit.WithSize(3))
	funbit.Integer(ipv4Matcher, &fragOff13, funbit.WithSize(13))
	funbit.Integer(ipv4Matcher, &ttl13, funbit.WithSize(8))
	funbit.Integer(ipv4Matcher, &protocol13, funbit.WithSize(8))
	funbit.Integer(ipv4Matcher, &checksum13, funbit.WithSize(16), funbit.WithEndianness("big"))
	funbit.Integer(ipv4Matcher, &srcIP13, funbit.WithSize(32), funbit.WithEndianness("big"))
	funbit.Integer(ipv4Matcher, &dstIP13, funbit.WithSize(32), funbit.WithEndianness("big"))

	ipv4Results, err := funbit.Match(ipv4Matcher, ipv4Packet)
	if err != nil {
		log.Fatalf("Failed to parse IPv4: %v", err)
	}

	fmt.Printf("   IPv4 Header Parsed:\n")
	fmt.Printf("     Version: %d, IHL: %d (words), TOS: %d\n", version13, ihl13, tos13)
	fmt.Printf("     Total Length: %d bytes, ID: %d\n", totalLen13, id13)
	fmt.Printf("     Flags: %d, Fragment Offset: %d\n", flags13, fragOff13)
	fmt.Printf("     TTL: %d, Protocol: %d, Checksum: 0x%04X\n", ttl13, protocol13, checksum13)
	fmt.Printf("     Source IP: %d.%d.%d.%d\n", (srcIP13>>24)&0xFF, (srcIP13>>16)&0xFF, (srcIP13>>8)&0xFF, srcIP13&0xFF)
	fmt.Printf("     Dest IP: %d.%d.%d.%d\n", (dstIP13>>24)&0xFF, (dstIP13>>16)&0xFF, (dstIP13>>8)&0xFF, dstIP13&0xFF)
	fmt.Printf("     Matched %d fields\n", len(ipv4Results))
	fmt.Println()

	// Example 14: Erlang Bit Syntax Equivalents
	fmt.Println("14. Erlang Bit Syntax Equivalents:")
	fmt.Printf("   Direct translations from Erlang bit syntax to funbit API\n")

	examples := []struct {
		erlang string
		funbit string
		result string
	}{
		{
			"<<42>>",
			"AddInteger(builder, 42)",
			"8-bit integer (default)",
		},
		{
			"<<42:16>>",
			"AddInteger(builder, 42, WithSize(16))",
			"16-bit integer",
		},
		{
			"<<42:16/little>>",
			"AddInteger(builder, 42, WithSize(16), WithEndianness(\"little\"))",
			"16-bit little-endian",
		},
		{
			"<<\"hello\">>",
			"AddBinary(builder, []byte(\"hello\"))",
			"Binary string",
		},
		{
			"<<Value:Size/binary>>",
			"Binary(matcher, &dest, WithDynamicSizeExpression(\"Size\"))",
			"Dynamic binary size",
		},
		{
			"<<Data:10/bitstring>>",
			"Bitstring(matcher, &dest, WithSize(10))",
			"10-bit bitstring",
		},
	}

	for i, ex := range examples {
		fmt.Printf("   %d. Erlang: %-25s -> Funbit: %s\n", i+1, ex.erlang, ex.funbit)
		fmt.Printf("      Result: %s\n", ex.result)
	}
	fmt.Println()

	fmt.Println("=== Comprehensive Funbit Public API Example Completed! ===")
	fmt.Printf("This example demonstrates all major bit syntax features:\n")
	fmt.Printf("• Basic construction and pattern matching\n")
	fmt.Printf("• Multiple data types (integer, float, binary, UTF)\n")
	fmt.Printf("• Endianness control (big/little/native)\n")
	fmt.Printf("• Unit specifiers for advanced size control\n")
	fmt.Printf("• Dynamic size expressions with variables\n")
	fmt.Printf("• Bit-level manipulation functions\n")
	fmt.Printf("• Binary vs Bitstring type differences\n")
	fmt.Printf("• Real-world protocol parsing (IPv4)\n")
	fmt.Printf("• Direct Erlang bit syntax equivalents\n")
}
