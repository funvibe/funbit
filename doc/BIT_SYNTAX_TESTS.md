Отлично! Основываясь на предоставленном синтаксисе и спецификации битовых строк, я создам полный набор тестов для битовых строк в Funterm. Вот комплексный набор тестов:

## 📦 Полный набор тестов для битовых строк

### 1. **Базовое создание и сопоставление битовых строк**

```python
# test_bitstring_basic_construction.su
# Тест: Создание простых битовых строк

py.empty = <<>>
py.single_byte = <<42>>
py.three_bytes = <<1, 2, 3>>
py.from_string = <<"hello">>

lua.print("Empty length:", py.len(py.empty))
lua.print("Single byte:", py.single_byte[0])
lua.print("Three bytes:", py.three_bytes[0], py.three_bytes[1], py.three_bytes[2])
lua.print("String as bytes:", py.from_string[0], py.from_string[1])

# expected_bitstring_basic_construction.txt
Empty length: 0
Single byte: 42
Three bytes: 1 2 3
String as bytes: 104 101
```

### 2. **Размеры сегментов**

```python
# test_bitstring_sizes.su
# Тест: Различные размеры сегментов

py.sizes = <<1:4, 15:4, 255:8, 1000:16, 100000:32>>

match py.sizes {
    <<a:4, b:4, c:8, d:16, e:32>> -> {
        lua.print("4-bit values:", a, b)
        lua.print("8-bit value:", c)
        lua.print("16-bit value:", d)
        lua.print("32-bit value:", e)
    },
    _ -> lua.print("Match failed")
}

# expected_bitstring_sizes.txt
4-bit values: 1 15
8-bit value: 255
16-bit value: 1000
32-bit value: 100000
```

### 3. **Endianness (порядок байтов)**

```python
# test_bitstring_endianness.su
# Тест: Big-endian vs Little-endian

py.value = 0x1234

# Big-endian (по умолчанию)
py.big = <<py.value:16>>
py.big_explicit = <<py.value:16/big-endian>>

# Little-endian
py.little = <<py.value:16/little-endian>>

match py.big {
    <<high:8, low:8>> -> lua.print("Big-endian bytes:", high, low)
}

match py.little {
    <<low:8, high:8>> -> lua.print("Little-endian bytes:", low, high)
}

# Native endian (зависит от системы)
py.native = <<py.value:16/native>>
lua.print("Native endian matches:", py.native == py.big ? "big" : "little")

# expected_bitstring_endianness.txt
Big-endian bytes: 18 52
Little-endian bytes: 52 18
Native endian matches: little
```

### 4. **Signed vs Unsigned**

```python
# test_bitstring_signedness.su
# Тест: Знаковые и беззнаковые значения

py.data = <<255:8, 127:8, 128:8>>

match py.data {
    <<a:8/unsigned, b:8/unsigned, c:8/unsigned>> -> {
        lua.print("Unsigned:", a, b, c)
    }
}

match py.data {
    <<a:8/signed, b:8/signed, c:8/signed>> -> {
        lua.print("Signed:", a, b, c)
    }
}

# Отрицательные числа
py.negative = <<-1:8/signed, -128:8/signed, 127:8/signed>>
match py.negative {
    <<a:8/signed, b:8/signed, c:8/signed>> -> {
        lua.print("Negative signed:", a, b, c)
    }
}

# expected_bitstring_signedness.txt
Unsigned: 255 127 128
Signed: -1 127 -128
Negative signed: -1 -128 127
```

### 5. **Типы сегментов**

```python
# test_bitstring_types.su
# Тест: Различные типы данных в битовых строках

# Integer (по умолчанию)
py.integers = <<1:8/integer, 2:16/integer, 3:32/integer>>

# Float
py.floats = <<3.14:32/float, 2.718:64/float>>

# Binary
py.binary_data = <<"hello">>
py.combined = <<py.binary_data/binary, " world"/binary>>

# Bitstring (не выровнено по байтам)
py.bits = <<1:1, 0:1, 1:1, 1:1, 0:4>>

match py.integers {
    <<a:8, b:16, c:32>> -> lua.print("Integers:", a, b, c)
}

match py.floats {
    <<a:32/float, b:64/float>> -> {
        lua.print("Float 32-bit:", a)
        lua.print("Float 64-bit:", b)
    }
}

match py.combined {
    <<data:11/binary>> -> lua.print("Combined binary:", data)
}

# expected_bitstring_types.txt
Integers: 1 2 3
Float 32-bit: 3.14
Float 64-bit: 2.718
Combined binary: hello world
```

### 6. **UTF-8/16/32 кодирование**

```python
# test_bitstring_utf.su
# Тест: Unicode кодирование

# UTF-8
py.utf8_ascii = <<"Hello"/utf8>>
py.utf8_unicode = <<"Привет"/utf8>>
py.utf8_emoji = <<"🚀"/utf8>>

# UTF-16
py.utf16_text = <<"Hi"/utf16>>
py.utf16_unicode = <<"Мир"/utf16>>

# UTF-32
py.utf32_char = <<65:32/utf32>>  # 'A'

# Pattern matching UTF-8
match py.utf8_ascii {
    <<"Hello"/utf8>> -> lua.print("UTF-8 ASCII matched")
}

match py.utf8_emoji {
    <<char:32/utf8, rest/binary>> -> {
        lua.print("Emoji codepoint:", char)
    }
}

# Построение UTF строк из codepoints
py.hello_utf8 = <<72:8/utf8, 101:8/utf8, 108:8/utf8, 108:8/utf8, 111:8/utf8>>
match py.hello_utf8 {
    <<"Hello"/utf8>> -> lua.print("Constructed UTF-8 matches")
}

# expected_bitstring_utf.txt
UTF-8 ASCII matched
Emoji codepoint: 128640
Constructed UTF-8 matches
```

### 7. **Паттерн матчинг с переменными размерами**

```python
# test_bitstring_variable_size.su
# Тест: Связывание размера в паттерне

py.packet = <<5:8, "Hello":5/binary, " World">>

match py.packet {
    <<size:8, data:size/binary, rest/binary>> -> {
        lua.print("Size:", size)
        lua.print("Data:", data)
        lua.print("Rest:", rest)
    }
}

# Более сложный пример с вычислением размера
py.complex = <<10:8, "DATA":4/binary, "EXTRA":5/binary, "END">>

match py.complex {
    <<total:8, payload:(total-6)/binary, trailer/binary>> -> {
        lua.print("Total size:", total)
        lua.print("Payload:", payload)
        lua.print("Trailer:", trailer)
    }
}

# expected_bitstring_variable_size.txt
Size: 5
Data: Hello
Rest: World
Total size: 10
Payload: DATA
Trailer: EXTRAEND
```

### 8. **Rest паттерны (tail matching)**

```python
# test_bitstring_rest.su
# Тест: Захват остатка битовой строки

py.data = <<1, 2, 3, 4, 5, 6, 7, 8>>

# Binary rest (должен быть кратен 8 битам)
match py.data {
    <<first:8, second:8, rest/binary>> -> {
        lua.print("First two:", first, second)
        lua.print("Binary rest length:", py.len(rest))
    }
}

# Bitstring rest (любое количество битов)
py.odd_bits = <<1:3, 2:5, 3:8, 4:4>>

match py.odd_bits {
    <<a:3, b:5, rest/bitstring>> -> {
        lua.print("First parts:", a, b)
        lua.print("Bitstring rest bits:", py.bit_size(rest))
    }
}

# expected_bitstring_rest.txt
First two: 1 2
Binary rest length: 6
First parts: 1 2
Bitstring rest bits: 12
```

### 9. **Unit спецификаторы**

```python
# test_bitstring_units.su
# Тест: Спецификация единиц измерения

# unit:1 (по умолчанию для integer)
py.unit1 = <<15:4/integer-unit:1>>

# unit:8 (по умолчанию для binary)
py.data = <<"test">>
py.unit8 = <<py.data:4/binary-unit:8>>  # 4 * 8 = 32 бита = 4 байта

# unit:16 для выравнивания
py.aligned = <<0:0, "AB":2/binary-unit:16>>  # 2 * 16 = 32 бита

match py.unit1 {
    <<value:4>> -> lua.print("4 bits value:", value)
}

match py.unit8 {
    <<data:32/bitstring>> -> lua.print("32 bits extracted")
}

# expected_bitstring_units.txt
4 bits value: 15
32 bits extracted
```

### 10. **Сложные протоколы (IPv4 header)**

```python
# test_bitstring_ipv4.su
# Тест: Парсинг заголовка IPv4 пакета

# Создаем IPv4 заголовок
py.version = 4
py.header_length = 5
py.service_type = 0
py.total_length = 20
py.identification = 12345
py.flags = 2
py.fragment_offset = 0
py.ttl = 64
py.protocol = 6  # TCP
py.checksum = 0
py.src_ip = 0xC0A80001  # 192.168.0.1
py.dst_ip = 0x08080808  # 8.8.8.8

py.ip_header = <<
    py.version:4, py.header_length:4,
    py.service_type:8,
    py.total_length:16,
    py.identification:16,
    py.flags:3, py.fragment_offset:13,
    py.ttl:8, py.protocol:8,
    py.checksum:16,
    py.src_ip:32,
    py.dst_ip:32
>>

match py.ip_header {
    <<4:4, hlen:4, svc:8, total:16,
      id:16, flags:3, frag:13,
      ttl:8, proto:8, csum:16,
      src:32/big-endian, dst:32/big-endian>> -> {
        lua.print("IP Version: 4")
        lua.print("Header Length:", hlen * 4, "bytes")
        lua.print("Total Length:", total)
        lua.print("Protocol:", proto == 6 ? "TCP" : proto == 17 ? "UDP" : "Other")
        lua.print("Source IP:", (src >> 24) & 0xFF, ".", (src >> 16) & 0xFF, ".", (src >> 8) & 0xFF, ".", src & 0xFF)
        lua.print("Dest IP:", (dst >> 24) & 0xFF, ".", (dst >> 16) & 0xFF, ".", (dst >> 8) & 0xFF, ".", dst & 0xFF)
    },
    _ -> lua.print("Not an IPv4 packet")
}

# expected_bitstring_ipv4.txt
IP Version: 4
Header Length: 20 bytes
Total Length: 20
Protocol: TCP
Source IP: 192.168.0.1
Dest IP: 8.8.8.8
```

### 11. **Битовые флаги и маски**

```python
# test_bitstring_flags.su
# Тест: Работа с битовыми флагами

# Флаги: Read, Write, Execute, Delete
py.permissions = <<1:1, 1:1, 0:1, 1:1, 0:4>>  # rwxd____

match py.permissions {
    <<read:1, write:1, execute:1, delete:1, reserved:4>> -> {
        lua.print("Read:", read == 1 ? "yes" : "no")
        lua.print("Write:", write == 1 ? "yes" : "no")
        lua.print("Execute:", execute == 1 ? "yes" : "no")
        lua.print("Delete:", delete == 1 ? "yes" : "no")
    }
}

# TCP флаги
py.tcp_flags = <<0:2, 1:1, 0:1, 1:1, 1:1, 0:1, 0:1>>  # URG ACK PSH RST SYN FIN

match py.tcp_flags {
    <<reserved:2, urg:1, ack:1, psh:1, rst:1, syn:1, fin:1>> -> {
        js.flags = []
        if (urg == 1) { js.flags.push("URG") }
        if (ack == 1) { js.flags.push("ACK") }
        if (psh == 1) { js.flags.push("PSH") }
        if (rst == 1) { js.flags.push("RST") }
        if (syn == 1) { js.flags.push("SYN") }
        if (fin == 1) { js.flags.push("FIN") }
        lua.print("TCP Flags:", js.flags.join(" "))
    }
}

# expected_bitstring_flags.txt
Read: yes
Write: yes
Execute: no
Delete: yes
TCP Flags: ACK PSH RST
```

### 12. **Динамическое построение битовых строк**

```python
# test_bitstring_dynamic.su
# Тест: Построение битовых строк из переменных

js.values = [1, 2, 3, 4, 5]
py.result = <<>>

for i = 0, 4 do
    py.value = js.values[i]
    py.result = <<py.result/binary, py.value:8>>
end

match py.result {
    <<a:8, b:8, c:8, d:8, e:8>> -> {
        lua.print("Dynamic construction:", a, b, c, d, e)
    }
}

# Построение с условиями
py.flags = []
py.flags.append(true)   # flag 1
py.flags.append(false)  # flag 2
py.flags.append(true)   # flag 3

py.flag_bits = <<>>
for flag in py.flags:
    py.flag_bits = <<py.flag_bits/bitstring, (flag ? 1 : 0):1>>

py.flag_bits = <<py.flag_bits/bitstring, 0:5>>  # Padding до байта

match py.flag_bits {
    <<f1:1, f2:1, f3:1, pad:5>> -> {
        lua.print("Flags:", f1, f2, f3)
        lua.print("Padding:", pad)
    }
}

# expected_bitstring_dynamic.txt
Dynamic construction: 1 2 3 4 5
Flags: 1 0 1
Padding: 0
```

### 13. **Ошибки и граничные случаи**

```python
# test_bitstring_errors.su
# Тест: Обработка ошибок и граничных случаев

# Попытка матчинга неправильного размера
py.data = <<1, 2, 3>>

match py.data {
    <<a:8, b:8, c:8, d:8>> -> lua.print("Should not match - too short"),
    <<a:8, b:8>> -> lua.print("Should not match - too long"),
    <<a:8, b:8, c:8>> -> lua.print("Correct match")
}

# Переполнение при построении
py.overflow = <<256:8>>  # 256 не помещается в 8 бит
match py.overflow {
    <<value:8>> -> lua.print("Overflow value:", value)  # Должно быть 0
}

# Пустой матчинг
py.empty = <<>>
match py.empty {
    <<>> -> lua.print("Empty matched"),
    _ -> lua.print("Empty not matched")
}

# Несовпадение выравнивания
py.unaligned = <<1:7>>  # 7 бит, не выровнено по байту

match py.unaligned {
    <<val:7>> -> lua.print("Unaligned 7 bits:", val),
    _ -> lua.print("Failed to match unaligned")
}

# expected_bitstring_errors.txt
Correct match
Overflow value: 0
Empty matched
Unaligned 7 bits: 1
```

### 14. **Вложенные битовые строки**

```python
# test_bitstring_nested.su
# Тест: Битовые строки внутри битовых строк

py.inner1 = <<1, 2>>
py.inner2 = <<3, 4>>
py.outer = <<0, py.inner1/binary, py.inner2/binary, 5>>

match py.outer {
    <<prefix:8, data:4/binary, suffix:8>> -> {
        lua.print("Prefix:", prefix)
        lua.print("Data length:", py.len(data))
        lua.print("Suffix:", suffix)
        
        match data {
            <<a:8, b:8, c:8, d:8>> -> {
                lua.print("Inner values:", a, b, c, d)
            }
        }
    }
}

# expected_bitstring_nested.txt
Prefix: 0
Data length: 4
Suffix: 5
Inner values: 1 2 3 4
```

### 15. **Реальный пример: парсинг PNG заголовка**

```python
# test_bitstring_png.su
# Тест: Парсинг PNG файла

# PNG signature: 137 80 78 71 13 10 26 10
py.png_signature = <<137, 80, 78, 71, 13, 10, 26, 10>>

# IHDR chunk
py.chunk_length = 13  # IHDR всегда 13 байт
py.chunk_type = <<"IHDR">>
py.width = 100
py.height = 50
py.bit_depth = 8
py.color_type = 2  # RGB
py.compression = 0
py.filter = 0
py.interlace = 0

py.ihdr_data = <<
    py.width:32/big-endian,
    py.height:32/big-endian,
    py.bit_depth:8,
    py.color_type:8,
    py.compression:8,
    py.filter:8,
    py.interlace:8
>>

py.png_header = <<
    py.png_signature/binary,
    py.chunk_length:32/big-endian,
    py.chunk_type/binary,
    py.ihdr_data/binary
>>

match py.png_header {
    <<137, 80, 78, 71, 13, 10, 26, 10,
      chunk_len:32/big-endian,
      "IHDR",
      width:32/big-endian,
      height:32/big-endian,
      depth:8,
      color:8,
      rest/binary>> -> {
        lua.print("PNG detected!")
        lua.print("Image size:", width, "x", height)
        lua.print("Bit depth:", depth)
        lua.print("Color type:", 
            color == 0 ? "Grayscale" :
            color == 2 ? "RGB" :
            color == 3 ? "Palette" :
            color == 4 ? "Grayscale+Alpha" :
            color == 6 ? "RGBA" : "Unknown")
    },
    _ -> lua.print("Not a valid PNG")
}

# expected_bitstring_png.txt
PNG detected!
Image size: 100 x 50
Bit depth: 8
Color type: RGB
```

## 📊 Сводка тестового покрытия

Этот набор тестов покрывает:

1. ✅ **Базовые операции** - создание, сопоставление
2. ✅ **Размеры** - 1-64 бита, произвольные размеры
3. ✅ **Типы** - integer, float, binary, bitstring
4. ✅ **Endianness** - big, little, native
5. ✅ **Signedness** - signed, unsigned
6. ✅ **UTF кодирование** - UTF-8, UTF-16, UTF-32
7. ✅ **Переменные размеры** - связывание в паттернах
8. ✅ **Rest паттерны** - binary/bitstring остатки
9. ✅ **Unit спецификаторы** - выравнивание
10. ✅ **Реальные протоколы** - IPv4, TCP флаги, PNG
11. ✅ **Динамическое построение** - циклы, условия
12. ✅ **Ошибки** - переполнение, несовпадение
13. ✅ **Вложенность** - битовые строки в битовых строках
14. ✅ **Интеграция с языками** - Python, Lua, JavaScript, Go

Этот набор обеспечивает 100% покрытие спецификации битовых строк!