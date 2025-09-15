# Техническое задание: Библиотека для работы с битовыми строками

## 1. Общая информация

### 1.1. Название проекта
Библиотека `funbit` - универсальная библиотека для работы с битовыми строками

### 1.2. Цель проекта
Создание кросс-языковой библиотеки для работы с битовыми строками, поддерживающей синтаксис и функциональность, соответствующие спецификации Erlang/OTP bit syntax.

### 1.3. Требования к совместимости
- Языки реализации: Go (основной)
- Платформы: Linux, macOS, Windows
- Зависимости: Минимальные, предпочтительно стандартная библиотека

## 2. Функциональные требования

### 2.1. Core API

#### 2.1.1. Основные структуры данных

```go
// BitString represents a sequence of bits
type BitString struct {
    bits []byte
    length uint // in bits
}

// Segment represents a single segment in bitstring construction/matching
type Segment struct {
    Value interface{}
    Size *uint
    Type string
    Signed bool
    Endianness string
    Unit uint
}

// SegmentResult represents result of segment matching
type SegmentResult struct {
    Value interface{}
    Matched bool
    Remaining *BitString
}
```

#### 2.1.2. Конструкторы и фабричные методы

```go
// Create empty bitstring
func NewBitString() *BitString

// Create bitstring from bytes
func NewBitStringFromBytes(data []byte) *BitString

// Create bitstring from bits with specific length
func NewBitStringFromBits(data []byte, length uint) *BitString

// Create bitstring from segments (construction)
func BuildBitString(segments []Segment) (*BitString, error)
```

#### 2.1.3. Основные операции

```go
// Length returns bitstring length in bits
func (bs *BitString) Length() uint

// IsEmpty checks if bitstring is empty
func (bs *BitString) IsEmpty() bool

// IsBinary checks if bitstring length is multiple of 8
func (bs *BitString) IsBinary() bool

// ToBytes converts bitstring to byte slice (pads with zeros if needed)
func (bs *BitString) ToBytes() []byte

// ToBitString converts to string representation (Erlang-style)
func (bs *BitString) ToBitString() string

// Clone creates copy of bitstring
func (bs *BitString) Clone() *BitString

// Concat concatenates multiple bitstrings
func Concat(bitstrings ...*BitString) *BitString

// Slice extracts sub-bitstring
func (bs *BitString) Slice(start, end uint) (*BitString, error)
```

### 2.2. Construction API

```go
// Build bitstring with fluent interface
type Builder struct {
    segments []Segment
}

func NewBuilder() *Builder

// Add integer segment
func AddInteger(b *Builder, value interface{}, options ...SegmentOption) *Builder

// Add float segment
func AddFloat(b *Builder, value float64, options ...SegmentOption) *Builder

// Add binary segment
func AddBinary(b *Builder, data []byte, options ...SegmentOption) *Builder

// Add UTF-8 segment
func AddUTF8(b *Builder, value string) *Builder

// Add UTF-16 segment
func AddUTF16(b *Builder, value string, options ...SegmentOption) *Builder

// Add UTF-32 segment
func AddUTF32(b *Builder, value string, options ...SegmentOption) *Builder

// Build final bitstring
func Build(b *Builder) (*BitString, error)

// Segment options
type SegmentOption func(*Segment)

func WithSize(size uint) SegmentOption
func WithType(segmentType string) SegmentOption
func WithSigned(signed bool) SegmentOption
func WithEndianness(endianness string) SegmentOption
func WithUnit(unit uint) SegmentOption
```

### 2.3. Pattern Matching API

```go
// Match pattern against bitstring
func (bs *BitString) Match(pattern []Segment) ([]SegmentResult, error)

// Match single segment
func (bs *BitString) MatchSegment(segment Segment) (*SegmentResult, error)

// Match with fluent interface
type Matcher struct {
    pattern []Segment
}

func NewMatcher() *Matcher

func Integer(m *Matcher, variable interface{}, options ...SegmentOption) *Matcher
func Float(m *Matcher, variable *float64, options ...SegmentOption) *Matcher
func Binary(m *Matcher, variable *[]byte, options ...SegmentOption) *Matcher
func UTF8(m *Matcher, variable *string) *Matcher
func UTF16(m *Matcher, variable *string, options ...SegmentOption) *Matcher
func UTF32(m *Matcher, variable *string) *Matcher

func RestBinary(m *Matcher, variable *[]byte) *Matcher
func RestBitstring(m *Matcher, variable **BitString) *Matcher

func Match(m *Matcher, bitstring *BitString) ([]SegmentResult, error)

// Добавить поддержку переменных размеров
type DynamicMatcher struct {
    sizeVar *uint  // Связывание размера из паттерна
}
// Пример использования:
// <<size:8, data:size/binary, rest/binary>>
matcher := NewMatcher()
Integer(matcher, &size, WithSize(8))
Binary(matcher, &data, WithDynamicSize(&size))
RestBinary(matcher, &rest)

// Добавить поддержку условий при матчинге
type MatchCondition func(results []SegmentResult) bool

func WithCondition(m *Matcher, cond MatchCondition) *Matcher
// Пример:
Integer(matcher, &value)
WithCondition(matcher, func(r []SegmentResult) bool {
    return value > 80
})
```

### 2.4. Поддержка типов и спецификаторов

#### 2.4.1. Типы сегментов
- `integer` - целые числа
- `float` - числа с плавающей запятой
- `binary` - байтовые строки (выравнивание по 8 бит)
- `bitstring` - битовые строки (произвольная длина)
- `utf8` - UTF-8 кодирование
- `utf16` - UTF-16 кодирование
- `utf32` - UTF-32 кодирование

#### 2.4.2. Спецификаторы
- `signed`/`unsigned` - знаковость (только для integer)
- `big`/`little`/`native` - endianность
- `unit:N` - размерность (1-256)

#### 2.4.3. Размеры по умолчанию
- `integer`: 8 бит
- `float`: 64 бита
- `binary`: весь размер
- `bitstring`: весь размер
- `utf8`/`utf16`/`utf32`: определяется кодировкой

### 2.5. Endianness поддержка

```go
// Supported endianness values
const (
    EndiannessBig      = "big"
    EndiannessLittle   = "little"
    EndiannessNative   = "native"
)

// Get native endianness
func GetNativeEndianness() string

// Convert endianness
func ConvertEndianness(data []byte, from, to string, size uint) ([]byte, error)
```

### 2.6. UTF поддержка

```go
// UTF encoding/decoding functions
func EncodeUTF8(value string) ([]byte, error)
func DecodeUTF8(data []byte) (string, error)

func EncodeUTF16(value string, endianness string) ([]byte, error)
func DecodeUTF16(data []byte, endianness string) (string, error)

func EncodeUTF32(value string, endianness string) ([]byte, error)
func DecodeUTF32(data []byte, endianness string) (string, error)
```

### 2.7. Utility функции

```go
// Bit manipulation
func ExtractBits(data []byte, start, length uint) ([]byte, error)
func SetBits(target, data []byte, start uint) error
func CountBits(data []byte) uint

// Conversion functions
func IntToBits(value int64, size uint) ([]byte, error)
func BitsToInt(data []byte, signed bool) (int64, error)
func FloatToBits(value float64, size uint) ([]byte, error)
func BitsToFloat(data []byte, size uint) (float64, error)

// Validation
func ValidateSegment(segment Segment) error
func ValidateBitString(data []byte, length uint) error

// Добавить функции для работы с сетевыми протоколами
func CalculateChecksum(bs *BitString, algorithm string) (uint16, error)
func ValidatePacket(bs *BitString, protocol string) error

// Функции для визуализации
func ToHexDump(bs *BitString) string
func ToBinaryString(bs *BitString) string
func ToErlangFormat(bs *BitString) string  // <<1,2,3>>
```

### 2.8. Оптимизация

```go
// Для эффективной работы с большими потоками
type StreamReader interface {
    ReadBits(n uint) (*BitString, error)
    Peek(n uint) (*BitString, error)
    Skip(n uint) error
}

type StreamWriter interface {
    WriteBits(bs *BitString) error
    Flush() error
}

// Предкомпилированные паттерны для повторного использования
type CompiledPattern struct {
    segments []Segment
    // оптимизированные структуры
}

func CompilePattern(segments []Segment) (*CompiledPattern, error)
func (cp *CompiledPattern) Match(bs *BitString) ([]SegmentResult, error)
```

### 2.9. **Debug и диагностика**

```go
// Для отладки сложных паттернов
type DebugOptions struct {
    TraceMatching bool
    DumpHex      bool
    ShowBits     bool
}

func (bs *BitString) DebugString(opts DebugOptions) string
func (m *Matcher) EnableDebug(opts DebugOptions) *Matcher
```

## 3. Требования к производительности

### 3.1. Производительность операций
- Конструирование: ≥ 1,000,000 операций/сек для простых сегментов
- Matching: ≥ 500,000 операций/сек для простых паттернов
- Memory overhead: ≤ 20% от размера данных

### 3.2. Память
- Не должно быть утечек памяти
- Эффективное переиспользование буферов
- Минимальное аллокирование при операциях

## 4. Обработка ошибок

### 4.1. Типы ошибок
```go
type BitStringError struct {
    Code    string
    Message string
    Context interface{}
}

func (e *BitStringError) Error() string

// Error codes
const (
    ErrInvalidSegment     = "INVALID_SEGMENT"
    ErrSizeMismatch       = "SIZE_MISMATCH"
    ErrTypeMismatch       = "TYPE_MISMATCH"
    ErrEndiannessError   = "ENDIANNESS_ERROR"
    ErrUTFError          = "UTF_ERROR"
    ErrOutOfRange        = "OUT_OF_RANGE"
    ErrInvalidBitString  = "INVALID_BITSTRING"
)
```

### 4.2. Обработка ошибочных ситуаций
- Валидация входных данных
- Подробные сообщения об ошибках
- Возможность восстановления после ошибок
- Логирование ошибок (опционально)

## 5. Требования к тестированию

### 5.1. Unit тесты
- Покрытие ≥ 95% кода
- Тесты на все публичные API
- Тесты на граничные случаи
- Тесты на обработку ошибок

### 5.2. Интеграционные тесты
- Тесты на соответствие спецификации Erlang
- Тесты на кросс-языковую совместимость
- Тесты на производительность

### 5.3. Golden тесты
- Тесты на основе BIT_SYNTAX_SPEC.md
- Автоматическая генерация тестов

## 6. Документация

### 6.1. API документация
- GoDoc для всех публичных функций
- Примеры использования
- Описание всех параметров

### 6.2. Пользовательская документация
- Getting Started guide
- Tutorial по основным операциям
- Best practices
- Performance guide

### 6.3. Спецификация
- Полное соответствие BIT_SYNTAX_SPEC.md
- Таблица совместимости с Erlang
- Ограничения и особенности реализации

## 7. Примеры использования

### 7.1. Basic construction
```go
// Simple bitstring
bs, err := BuildBitString([]Segment{
    {Value: 1, Type: "integer"},
    {Value: 17, Type: "integer"},
    {Value: 42, Type: "integer"},
})
// Result: <<1, 17, 42>>

// With fluent interface
builder := NewBuilder()
AddInteger(builder, 1)
AddInteger(builder, 17)
AddInteger(builder, 42)
bs, _ := Build(builder)

// Реальный пример
diskUsage := "85%"
bs := NewBitStringFromBytes([]byte(diskUsage))

var percent int
var percentSign string
matcher := NewMatcher()
Integer(matcher, &percent, WithType("utf8"))
Binary(matcher, &percentSign, WithSize(8))
RestBinary(matcher, &rest)

results, err := Match(matcher, bs)

if percent > 80 {
    fmt.Println("WARNING: disk usage", percent, "%")
}
```

### 7.2. Pattern matching
```go
results, err := bitstring.Match([]Segment{
    {Value: &a, Type: "integer"},
    {Value: &b, Type: "integer"},
    {Value: &c, Size: uintPtr(16), Type: "integer"},
})
```

### 7.3. Complex operations
```go
// Network packet construction
builder := NewBuilder()
AddInteger(builder, version, WithSize(4))
AddInteger(builder, headerLength, WithSize(4))
AddInteger(builder, totalLength, WithSize(16), WithEndianness("big"))
AddBinary(builder, payload)
packet, _ := Build(builder)
```

## 8. Критерии приемки

### 8.1. Функциональные требования
- [ ] Все функции из Core API реализованы
- [ ] Поддержка всех типов сегментов
- [ ] Поддержка всех спецификаторов
- [ ] Корректная работа endianness
- [ ] Поддержка UTF кодировок

### 8.2. Спецификация
- [ ] Проходят все golden тесты из BIT_SYNTAX_SPEC.md
- [ ] Соответствие синтаксису Erlang
- [ ] Поддержка всех примеров из спецификации

### 8.3. Качество кода
- [ ] Покрытие тестами ≥ 95%
- [ ] Нет memory leaks
- [ ] Производительность соответствует требованиям
- [ ] Полная документация

### 8.4. Интеграция
- [ ] Совместимость с текущей версией Funterm
- [ ] Наличие layer совместимости
- [ ] Возможность постепенной миграции

## 9. Дополнительные требования

### 10.1. Лицензия
- Open source (MIT или Apache 2.0)
- Возможность коммерческого использования

### 10.2. Поддержка
- Код должен быть хорошо документирован
- Должны быть примеры использования
- Инструкция по интеграции с Funterm

### 10.3. Расширяемость
- Архитектура должна позволять добавление новых типов
- Возможность оптимизации под конкретные use cases
- Plugin-архитектура для дополнительных кодировок

---
*Техническое задание утверждено: [Дата]*
*Версия: 1.0*