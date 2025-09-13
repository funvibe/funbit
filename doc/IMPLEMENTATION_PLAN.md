# План реализации библиотеки funbit

## Обзор

Данный документ содержит подробный план реализации библиотеки `funbit` - высококачественной реализации спецификации битового синтаксиса Erlang на языке Go. Библиотека предназначена для публикации на GitHub и может использоваться как самостоятельно, так и в составе Funterm.

## Архитектурные принципы

### 1. Фокус на спецификации
- Точная реализация Erlang bit syntax без избыточной расширяемости
- Соответствие техническому заданию из `BITSTRING_LIBRARY_TZ.md`
- Проход всех тестов из `BIT_SYNTAX_TESTS.md`

### 2. Производительность
- Целевые показатели: ≥ 1,000,000 операций конструирования/сек
- Целевые показатели: ≥ 500,000 операций сопоставления/сек
- Накладные расходы на память: ≤ 20%
- Минимизация аллокаций и эффективное переиспользование буферов

### 3. Качество кода
- Идиоматичный Go код
- Чистые API с fluent интерфейсами
- Комплексная обработка ошибок
- Подробная документация

### 4. Интеграция
- Легкая интеграция с Funterm
- Возможность использования как standalone библиотеки
- Четкие точки интеграции и API

## Структура проекта

```
funbit/
├── cmd/                    # Командные утилиты (примеры, бенчмарки)
│   ├── examples/
│   └── benchmarks/
├── internal/               # Внутренняя реализация
│   ├── bitstring/         # Core реализации
│   │   ├── bitstring.go
│   │   ├── segment.go
│   │   ├── bitstring_test.go
│   │   ├── size.go
│   │   ├── size_test.go
│   │   └── nested.go
│   ├── builder/           # Construction API
│   │   ├── builder.go
│   │   ├── builder_test.go
│   │   └── dynamic.go
│   ├── matcher/           # Pattern matching API
│   │   ├── matcher.go
│   │   ├── matcher_test.go
│   │   ├── dynamic.go
│   │   └── rest.go
│   ├── types/             # Обработчики типов
│   │   ├── integer.go
│   │   ├── integer_test.go
│   │   ├── float.go
│   │   ├── float_test.go
│   │   ├── binary.go
│   │   ├── binary_test.go
│   │   ├── signed.go
│   │   ├── signed_test.go
│   │   └── unit.go
│   ├── endianness/        # Endianness поддержка
│   │   ├── endianness.go
│   │   └── endianness_test.go
│   ├── utf/               # UTF кодирование
│   │   ├── utf.go
│   │   ├── utf8.go
│   │   ├── utf16.go
│   │   ├── utf32.go
│   │   └── utf_test.go
│   ├── utils/             # Утилиты
│   │   ├── utils.go
│   │   └── utils_test.go
│   ├── errors/            # Обработка ошибок
│   │   ├── errors.go
│   │   ├── validation.go
│   │   └── errors_test.go
│   └── stream/            # Streaming операции
│       ├── reader.go
│       ├── writer.go
│       └── stream_test.go
├── pkg/                   # Публичный API
│   └── funbit/            # Основной пакет
│       ├── funbit.go
│       └── funbit_test.go
├── testdata/              # Тестовые данные
│   ├── erlang_compatibility/
│   ├── protocol_samples/
│   └── edge_cases/
├── acceptancetests/       # Приемочные тесты
│   ├── bitstring_basic_test.go
│   ├── bitstring_sizes_test.go
│   ├── bitstring_types_test.go
│   ├── bitstring_endianness_test.go
│   ├── bitstring_signedness_test.go
│   ├── bitstring_utf_test.go
│   ├── bitstring_variable_size_test.go
│   ├── bitstring_rest_test.go
│   ├── bitstring_units_test.go
│   ├── bitstring_protocols_test.go
│   ├── bitstring_dynamic_test.go
│   ├── bitstring_errors_test.go
│   └── bitstring_nested_test.go
├── benchmarks/           # Бенчмарки
│   ├── construction_bench.go
│   ├── matching_bench.go
│   └── memory_bench.go
├── integration/           # Интеграционные тесты
│   ├── erlang_compatibility_test.go
│   └── funterm_integration_test.go
├── examples/              # Примеры использования
│   ├── ipv4.go
│   ├── tcp.go
│   ├── png.go
│   ├── dynamic.go
│   └── nested.go
├── docs/                  # Документация
└── go.mod                 # Модуль
```

## Стратегия тестирования

### Общий подход
Библиотека разрабатывается с использованием методологии Test-Driven Development (TDD):
- **источник истины**: [`doc/BIT_SYNTAX_SPEC.md`](doc/BIT_SYNTAX_SPEC.md)
- **приемочные тесты**: [`doc/BIT_SYNTAX_TESTS.md`](doc/BIT_SYNTAX_TESTS.md)
- **технические требования**: [`doc/BITSTRING_LIBRARY_TZ.md`](doc/BITSTRING_LIBRARY_TZ.md)

### Процесс разработки для каждого этапа

#### Шаг 1: Анализ и создание тестов (RED)
- Изучить соответствующие разделы `BIT_SYNTAX_SPEC.md`
- Проанализировать примеры из `BIT_SYNTAX_TESTS.md`
- Создать тесты на планируемую функциональность
- **Все тесты должны быть красными** (фейлящими)
- Убедиться, что тесты предыдущих этапов остаются зелеными

#### Шаг 2: Реализация функциональности (GREEN)
- Реализовать минимально необходимый код для прохождения тестов
- **Все тесты должны стать зелеными**
- Никаких регрессий - тесты предыдущих этапов остаются зелеными

#### Шаг 3: Рефакторинг (REFACTOR)
- Оптимизировать код, сохраняя все тесты зелеными
- Улучшить производительность, если необходимо
- Обновить документацию

#### Шаг 4: Валидация
- При появлении регрессий - немедленное исправление
- Переход к следующему этапу только при всех зеленых тестах

### Тестирование для каждого этапа

#### Этап 1: Базовые операции
**Тестовые файлы:**
- `acceptancetests/bitstring_basic_test.go`
- `internal/bitstring/bitstring_test.go`
- `internal/builder/builder_test.go`
- `internal/matcher/matcher_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Sections 1-3 (Basic bit string construction and matching)
- `BIT_SYNTAX_TESTS.md`: Тесты 1-3 (Basic construction, Sizes, Endianness basics)

#### Этап 2: Поддержка размеров
**Тестовые файлы:**
- `acceptancetests/bitstring_sizes_test.go`
- `internal/bitstring/size_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Section "Integer segments", "Binary segments"
- `BIT_SYNTAX_TESTS.md`: Тест 2 (Sizes)

#### Этап 3: Типы данных
**Тестовые файлы:**
- `acceptancetests/bitstring_types_test.go`
- `internal/types/integer_test.go`
- `internal/types/float_test.go`
- `internal/types/binary_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Sections "Integer segments", "Float segments", "Binary segments"
- `BIT_SYNTAX_TESTS.md`: Тест 5 (Types)

#### Этап 4: Endianness
**Тестовые файлы:**
- `acceptancetests/bitstring_endianness_test.go`
- `internal/endianness/endianness_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Section "Endianness"
- `BIT_SYNTAX_TESTS.md`: Тест 3 (Endianness)

#### Этап 5: Signed/Unsigned
**Тестовые файлы:**
- `acceptancetests/bitstring_signedness_test.go`
- `internal/types/signed_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Section "Signedness"
- `BIT_SYNTAX_TESTS.md`: Тест 4 (Signed vs Unsigned)

#### Этап 6: UTF кодирование
**Тестовые файлы:**
- `acceptancetests/bitstring_utf_test.go`
- `internal/utf/utf_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Section "Unicode segments"
- `BIT_SYNTAX_TESTS.md`: Тест 6 (UTF-8/16/32)

#### Этап 7: Переменные размеры
**Тестовые файлы:**
- `acceptancetests/bitstring_variable_size_test.go`
- `internal/matcher/dynamic_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Section "Binding and Using a Size Variable"
- `BIT_SYNTAX_TESTS.md`: Тест 7 (Variable size patterns)

#### Этап 8: Rest паттерны
**Тестовые файлы:**
- `acceptancetests/bitstring_rest_test.go`
- `internal/matcher/rest_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Section "Getting the Rest of the Binary or Bitstring"
- `BIT_SYNTAX_TESTS.md`: Тест 8 (Rest patterns)

#### Этап 9: Unit спецификаторы
**Тестовые файлы:**
- `acceptancetests/bitstring_units_test.go`
- `internal/types/unit_test.go`

**Источники тестов:**
- `BIT_SYNTAX_SPEC.md`: Section "Unit"
- `BIT_SYNTAX_TESTS.md`: Тест 9 (Unit specifiers)

#### Этап 10: Реальные протоколы
**Тестовые файлы:**
- `acceptancetests/bitstring_protocols_test.go`
- `examples/ipv4_test.go`
- `examples/tcp_test.go`
- `examples/png_test.go`

**Источники тестов:**
- `BIT_SYNTAX_TESTS.md`: Тест 10 (Complex protocols)

#### Этап 11: Динамическое построение
**Тестовые файлы:**
- `acceptancetests/bitstring_dynamic_test.go`
- `internal/builder/dynamic_test.go`

**Источники тестов:**
- `BIT_SYNTAX_TESTS.md`: Тест 12 (Dynamic construction)

#### Этап 12: Обработка ошибок
**Тестовые файлы:**
- `acceptancetests/bitstring_errors_test.go`
- `internal/errors/errors_test.go`

**Источники тестов:**
- `BIT_SYNTAX_TESTS.md`: Тест 13 (Errors and edge cases)

#### Этап 13: Вложенные битовые строки
**Тестовые файлы:**
- `acceptancetests/bitstring_nested_test.go`
- `internal/bitstring/nested_test.go`

**Источники тестов:**
- `BIT_SYNTAX_TESTS.md`: Тест 14 (Nested bitstrings)

### Бенчмарки

#### Производительность
```go
func BenchmarkConstruction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bs, err := NewBuilder().
            AddInteger(1).
            AddInteger(17).
            AddInteger(42).
            Build()
        if err != nil {
            b.Fatal(err)
        }
        _ = bs
    }
}

func BenchmarkMatching(b *testing.B) {
    bs, _ := NewBuilder().
        AddInteger(1).
        AddInteger(17).
        AddInteger(42).
        Build()
    
    var a, b, c int
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        results, err := NewMatcher().
            Integer(&a).
            Integer(&b).
            Integer(&c).
            Match(bs)
        if err != nil {
            b.Fatal(err)
        }
        _ = results
    }
}
```

#### Целевые показатели
- **Construction**: ≥ 1,000,000 операций/сек
- **Matching**: ≥ 500,000 операций/сек
- **Memory overhead**: ≤ 20%

### Интеграционные тесты

#### Совместимость с Erlang
```go
func TestErlangCompatibility(t *testing.T) {
    // Тесты на соответствие поведению Erlang
    testCases := loadErlangTestCases("testdata/erlang_compatibility/")
    
    for _, tc := range testCases {
        // Конструирование как в Erlang
        bs, err := BuildFromErlangSyntax(tc.ErlangSyntax)
        if err != nil {
            t.Errorf("Failed to build %s: %v", tc.ErlangSyntax, err)
            continue
        }
        
        // Сопоставление как в Erlang
        result, err := MatchWithErlangPattern(bs, tc.ErlangPattern)
        if err != nil {
            t.Errorf("Failed to match %s: %v", tc.ErlangPattern, err)
            continue
        }
        
        if !reflect.DeepEqual(result, tc.ExpectedResult) {
            t.Errorf("For %s expected %v, got %v",
                tc.ErlangSyntax, tc.ExpectedResult, result)
        }
    }
}
```

### Запуск тестов

```bash
# Запуск всех тестов
go test ./...

# Запуск только приемочных тестов
go test ./acceptancetests/

# Запуск бенчмарков
go test -bench=. ./benchmarks/

# Запуск с покрытием
go test -cover ./...

# Запуск тестов для конкретного этапа
go test -run="TestBitstringBasic" ./acceptancetests/
```

### Правила тестирования

1. **Красные тесты перед реализацией**
   - Все новые тесты должны фейлиться перед реализацией
   - Тесты предыдущих этапов должны оставаться зелеными

2. **Зеленые тесты после реализации**
   - Все тесты должны проходить после реализации
   - Никаких регрессий не допускается

3. **Немедленное исправление регрессий**
   - При появлении любых regressions - немедленная остановка и исправление
   - Переход к следующему этапу только при всех зеленых тестах

4. **Покрытие**
   - Целевое покрытие ≥ 95%
   - Каждый публичный метод должен иметь тесты
   - Все граничные случаи должны быть протестированы

## Порядок реализации

### Этап 1: Базовые операции (Создание, Сопоставление)

#### 1.1 Core Data Structures
**Файлы:** `internal/bitstring/bitstring.go`, `internal/bitstring/segment.go`

**Задачи:**
- Реализовать структуру `BitString`
  ```go
  type BitString struct {
      data     []byte
      bitLen   uint
      capacity uint
  }
  ```
- Реализовать структуру `Segment`
  ```go
  type Segment struct {
      Value       interface{}
      Size        *uint
      Type        string
      Signed      bool
      Endianness  string
      Unit        uint
  }
  ```
- Реализовать структуру `SegmentResult`
  ```go
  type SegmentResult struct {
      Value     interface{}
      Matched   bool
      Remaining *BitString
  }
  ```

**Методы BitString:**
- `NewBitString() *BitString`
- `NewBitStringFromBytes(data []byte) *BitString`
- `NewBitStringFromBits(data []byte, length uint) *BitString`
- `Length() uint`
- `IsEmpty() bool`
- `IsBinary() bool`
- `ToBytes() []byte`
- `Clone() *BitString`

**Критерии приемки:**
- Все базовые операции работают корректно
- Память управляется эффективно
- Нет утечек памяти
- Базовые тесты проходят

#### 1.2 Basic Construction
**Файлы:** `internal/builder/builder.go`

**Задачи:**
- Реализовать базовый `Builder`
  ```go
  type Builder struct {
      segments []Segment
      buffer   *bytes.Buffer
  }
  ```
- Реализовать базовые методы:
  - `NewBuilder() *Builder`
  - `AddInteger(value interface{}, opts ...SegmentOption) *Builder`
  - `Build() (*BitString, error)`

**Segment Options:**
- `WithSize(size uint) SegmentOption`
- `WithType(segmentType string) SegmentOption`

**Критерии приемки:**
- Простые битовые строки создаются корректно
- Базовые опции работают
- Ошибки обрабатываются правильно

#### 1.3 Basic Pattern Matching
**Файлы:** `internal/matcher/matcher.go`

**Задачи:**
- Реализовать базовый `Matcher`
  ```go
  type Matcher struct {
      pattern []Segment
  }
  ```
- Реализовать базовые методы:
  - `NewMatcher() *Matcher`
  - `Integer(variable interface{}, opts ...SegmentOption) *Matcher`
  - `Match(bitstring *BitString) ([]SegmentResult, error)`

**Критерии приемки:**
- Простые паттерны сопоставляются корректно
- Переменные связываются правильно
- Ошибки сопоставления обрабатываются

---

### Этап 2: Поддержка размеров (1-64 бита, произвольные размеры)

#### 2.1 Size Handling
**Файлы:** `internal/bitstring/size.go`, `internal/types/size_handler.go`

**Задачи:**
- Реализовать обработку размеров от 1 до 64 бит
- Поддержать произвольные размеры для бинарных данных
- Реализовать выравнивание и паддинг
- Добавить валидацию размеров

**Методы:**
- `ValidateSize(size uint, unit uint) error`
- `CalculateTotalSize(segment Segment) (uint, error)`
- `ExtractBits(data []byte, start, length uint) ([]byte, error)`
- `SetBits(target, data []byte, start uint) error`

**Критерии приемки:**
- Все размеры от 1 до 64 бит работают корректно
- Произвольные размеры для бинарных данных поддерживаются
- Выравнивание работает правильно
- Граничные случаи обрабатываются

---

### Этап 3: Типы данных (integer, float, binary, bitstring)

#### 3.1 Integer Type Handler
**Файлы:** `internal/types/integer.go`

**Задачи:**
- Реализовать `IntegerTypeHandler`
- Поддержать все размеры целых чисел
- Реализовать кодирование и декодирование
- Обработать переполнение и усечение

**Методы:**
- `EncodeInteger(value interface{}, size uint, signed bool, endianness string) ([]byte, error)`
- `DecodeInteger(data []byte, size uint, signed bool, endianness string) (interface{}, error)`
- `ValidateInteger(value interface{}, size uint) error`

#### 3.2 Float Type Handler
**Файлы:** `internal/types/float.go`

**Задачи:**
- Реализовать `FloatTypeHandler`
- Поддержать размеры 32, 64 бита
- Реализовать кодирование и декодирование IEEE 754
- Обработать специальные значения (NaN, Inf)

**Методы:**
- `EncodeFloat(value float64, size uint, endianness string) ([]byte, error)`
- `DecodeFloat(data []byte, size uint, endianness string) (float64, error)`

#### 3.3 Binary/Bitstring Type Handler
**Файлы:** `internal/types/binary.go`

**Задачи:**
- Реализовать `BinaryTypeHandler` и `BitstringTypeHandler`
- Обработать выравнивание по байтам для binary
- Поддержать произвольную длину для bitstring
- Реализовать эффективное копирование данных

**Методы:**
- `EncodeBinary(data []byte, unit uint) ([]byte, error)`
- `DecodeBinary(data []byte, size uint, unit uint) ([]byte, error)`
- `EncodeBitstring(bs *BitString) ([]byte, error)`
- `DecodeBitstring(data []byte, size uint) (*BitString, error)`

**Критерии приемки:**
- Все типы данных работают корректно
- Кодирование/декодирование соответствует Erlang
- Специальные значения обрабатываются правильно
- Производительность соответствует требованиям

---

### Этап 4: Endianness (big, little, native)

#### 4.1 Endianness Support
**Файлы:** `internal/endianness/endianness.go`

**Задачи:**
- Реализовать поддержку big-endian, little-endian, native-endian
- Определить нативный endianness системы
- Реализовать конвертацию endianness
- Оптимизировать операции с учетом endianness

**Константы:**
```go
const (
    EndiannessBig    = "big"
    EndiannessLittle = "little"
    EndiannessNative = "native"
)
```

**Методы:**
- `GetNativeEndianness() string`
- `ConvertEndianness(data []byte, from, to string, size uint) ([]byte, error)`
- `EncodeWithEndianness(value interface{}, size uint, endianness string) ([]byte, error)`
- `DecodeWithEndianness(data []byte, size uint, endianness string) (interface{}, error)`

**Критерии приемки:**
- Все типы endianness работают корректно
- Native endianness определяется правильно
- Конвертация работает для всех поддерживаемых размеров
- Производительность оптимизирована

---

### Этап 5: Signed/Unsigned поддержка

#### 5.1 Signedness Support
**Файлы:** `internal/types/signed.go`

**Задачи:**
- Реализовать поддержку знаковых и беззнаковых целых чисел
- Обработать преобразование знаковых представлений
- Реализовать правильное декодирование отрицательных чисел
- Добавить валидацию диапазонов

**Методы:**
- `EncodeSignedInteger(value int64, size uint) ([]byte, error)`
- `EncodeUnsignedInteger(value uint64, size uint) ([]byte, error)`
- `DecodeSignedInteger(data []byte, size uint) (int64, error)`
- `DecodeUnsignedInteger(data []byte, size uint) (uint64, error)`
- `ValidateSignedRange(value int64, size uint) error`
- `ValidateUnsignedRange(value uint64, size uint) error`

**Критерии приемки:**
- Знаковые и беззнаковые числа работают корректно
- Отрицательные числа обрабатываются правильно
- Переполнение и диапазоны проверяются
- Соответствие поведению Erlang

---

### Этап 6: UTF кодирование (UTF-8, UTF-16, UTF-32)

#### 6.1 UTF Support
**Файлы:** `internal/utf/utf.go`, `internal/utf/utf8.go`, `internal/utf/utf16.go`, `internal/utf/utf32.go`

**Задачи:**
- Реализовать поддержку UTF-8, UTF-16, UTF-32
- Обработать кодирование и декодирование Unicode
- Поддержать endianness для UTF-16/UTF-32
- Валидировать Unicode code points

**Методы:**
- `EncodeUTF8(value string) ([]byte, error)`
- `DecodeUTF8(data []byte) (string, error)`
- `EncodeUTF16(value string, endianness string) ([]byte, error)`
- `DecodeUTF16(data []byte, endianness string) (string, error)`
- `EncodeUTF32(value string, endianness string) ([]byte, error)`
- `DecodeUTF32(data []byte, endianness string) (string, error)`
- `ValidateUnicodeCodePoint(codePoint int) error`

**Критерии приемки:**
- Все UTF кодировки работают корректно
- Unicode символы обрабатываются правильно
- Endianness для UTF-16/UTF-32 поддерживается
- Невалидные последовательности обрабатываются

---

### Этап 7: Переменные размеры в паттернах

#### 7.1 Dynamic Size Binding
**Файлы:** `internal/matcher/dynamic.go`

**Задачи:**
- Реализовать связывание размеров в паттернах
- Поддержать выражения для вычисления размеров
- Реализовать валидацию зависимостей
- Оптимизировать обработку динамических размеров

**Методы:**
- `BindSizeVariable(name string, value uint)`
- `EvaluateSizeExpression(expr string, context map[string]uint) (uint, error)`
- `ValidateSizeDependencies(pattern []Segment) error`
- `OptimizeDynamicPattern(pattern []Segment) ([]Segment, error)`

**Пример использования:**
```go
// <<size:8, data:size/binary, rest/binary>>
matcher := NewMatcher().
    Integer(&size, WithSize(8)).
    Binary(&data, WithDynamicSize("size")).
    RestBinary(&rest)
```

**Критерии приемки:**
- Динамические размеры работают корректно
- Выражения вычисляются правильно
- Зависимости валидируются
- Производительность оптимизирована

---

### Этап 8: Rest паттерны (binary/bitstring остатки)

#### 8.1 Rest Pattern Support
**Файлы:** `internal/matcher/rest.go`

**Задачи:**
- Реализовать захват остатка битовой строки
- Поддержать `rest/binary` (выравнено по байтам)
- Поддержать `rest/bitstring` (произвольная длина)
- Обработать валидацию остатков

**Методы:**
- `MatchRestBinary(bs *BitString, offset uint) ([]byte, error)`
- `MatchRestBitstring(bs *BitString, offset uint) (*BitString, error)`
- `ValidateRestPattern(pattern []Segment) error`
- `ExtractRestSegment(pattern []Segment) (int, error)`

**Пример использования:**
```go
matcher := NewMatcher().
    Integer(&header, WithSize(16)).
    RestBinary(&payload)    // Для binary
matcher := NewMatcher().
    Integer(&flags, WithSize(4)).
    RestBitstring(&data)    // Для bitstring
```

**Критерии приемки:**
- Rest паттерны работают корректно
- Binary остатки выравнены по байтам
- Bitstring остатки поддерживают произвольную длину
- Валидация работает правильно

---

### Этап 9: Unit спецификаторы и выравнивание

#### 9.1 Unit Specifiers
**Файлы:** `internal/types/unit.go`

**Задачи:**
- Реализовать поддержку unit спецификаторов
- Поддержать диапазон 1-256 для unit
- Реализовать выравнивание данных
- Оптимизировать операции с учетом unit

**Методы:**
- `ValidateUnit(unit uint) error`
- `CalculateAlignedSize(size uint, unit uint) uint`
- `AlignData(data []byte, unit uint) ([]byte, error)`
- `EncodeWithUnit(value interface{}, size uint, unit uint) ([]byte, error)`
- `DecodeWithUnit(data []byte, size uint, unit uint) (interface{}, error)`

**Пример использования:**
```go
// unit:16 для выравнивания по 16 бит
segment := Segment{
    Value:      data,
    Size:       uintPtr(4),
    Unit:       16,  // 4 * 16 = 64 бита
    Type:       "binary",
}
```

**Критерии приемки:**
- Unit спецификаторы работают корректно
- Выравнивание выполняется правильно
- Диапазон 1-256 поддерживается
- Производительность оптимизирована

---

### Этап 10: Реальные протоколы (IPv4, TCP флаги, PNG)

#### 10.1 Protocol Examples
**Файлы:** `examples/ipv4.go`, `examples/tcp.go`, `examples/png.go`

**Задачи:**
- Реализовать парсинг IPv4 заголовков
- Реализовать обработку TCP флагов
- Реализовать парсинг PNG заголовков
- Создать комплексные примеры использования

**IPv4 Header Example:**
```go
func ParseIPv4Header(data *BitString) (*IPv4Header, error) {
    matcher := NewMatcher().
        Integer(&version, WithSize(4)).
        Integer(&headerLength, WithSize(4)).
        Integer(&serviceType, WithSize(8)).
        Integer(&totalLength, WithSize(16), WithEndianness("big")).
        Integer(&identification, WithSize(16)).
        Integer(&flags, WithSize(3)).
        Integer(&fragmentOffset, WithSize(13)).
        Integer(&ttl, WithSize(8)).
        Integer(&protocol, WithSize(8)).
        Integer(&checksum, WithSize(16)).
        Integer(&srcIP, WithSize(32)).
        Integer(&dstIP, WithSize(32))
    
    results, err := matcher.Match(data)
    // ...
}
```

**TCP Flags Example:**
```go
func ParseTCPFlags(data *BitString) (*TCPFlags, error) {
    matcher := NewMatcher().
        Integer(&reserved, WithSize(2)).
        Integer(&urg, WithSize(1)).
        Integer(&ack, WithSize(1)).
        Integer(&psh, WithSize(1)).
        Integer(&rst, WithSize(1)).
        Integer(&syn, WithSize(1)).
        Integer(&fin, WithSize(1))
    
    results, err := matcher.Match(data)
    // ...
}
```

**PNG Header Example:**
```go
func ParsePNGHeader(data *BitString) (*PNGHeader, error) {
    matcher := NewMatcher().
        Binary(&signature, WithSize(64)).  // 8 bytes
        Integer(&chunkLength, WithSize(32), WithEndianness("big")).
        Binary(&chunkType, WithSize(32)).
        // IHDR data
        Integer(&width, WithSize(32), WithEndianness("big")).
        Integer(&height, WithSize(32), WithEndianness("big")).
        Integer(&bitDepth, WithSize(8)).
        Integer(&colorType, WithSize(8)).
        Integer(&compression, WithSize(8)).
        Integer(&filter, WithSize(8)).
        Integer(&interlace, WithSize(8))
    
    results, err := matcher.Match(data)
    // ...
}
```

**Критерии приемки:**
- Все примеры протоколов работают корректно
- Реальные данные парсятся правильно
- Производительность соответствует требованиям
- Код легко читается и поддерживается

---

### Этап 11: Динамическое построение (циклы, условия)

#### 11.1 Dynamic Construction
**Файлы:** `examples/dynamic.go`, `internal/builder/dynamic.go`

**Задачи:**
- Реализовать динамическое построение битовых строк
- Поддержать построение в циклах
- Реализовать условное построение
- Оптимизировать динамические операции

**Методы:**
- `AppendToBitString(target *BitString, segments ...Segment) error`
- `BuildBitStringDynamically(generator func() ([]Segment, error)) (*BitString, error)`
- `BuildConditionalBitString(condition bool, trueSegments, falseSegments []Segment) (*BitString, error)`

**Пример использования:**
```go
// Динамическое построение в цикле
func BuildPacket(values []int) (*BitString, error) {
    builder := NewBuilder()
    
    for _, value := range values {
        builder.AddInteger(value, WithSize(16))
    }
    
    return builder.Build()
}

// Условное построение
func BuildMessage(withHeader bool, payload []byte) (*BitString, error) {
    builder := NewBuilder()
    
    if withHeader {
        builder.AddInteger(0x MAGIC, WithSize(32))
    }
    
    builder.AddBinary(payload)
    return builder.Build()
}
```

**Критерии приемки:**
- Динамическое построение работает корректно
- Циклы и условия обрабатываются правильно
- Производительность оптимизирована
- Память управляется эффективно

---

### Этап 12: Обработка ошибок (переполнение, несовпадение)

#### 12.1 Error Handling
**Файлы:** `internal/errors/errors.go`, `internal/errors/validation.go`

**Задачи:**
- Реализовать комплексную систему ошибок
- Обработать переполнение при конструировании
- Обработать несовпадение при паттерн матчинге
- Предоставить детальную информацию об ошибках

**Типы ошибок:**
```go
type BitStringError struct {
    Code    string
    Message string
    Context map[string]interface{}
    Stack   []string
}

const (
    ErrInvalidSegment     = "INVALID_SEGMENT"
    ErrSizeMismatch       = "SIZE_MISMATCH"
    ErrTypeMismatch       = "TYPE_MISMATCH"
    ErrEndiannessError   = "ENDIANNESS_ERROR"
    ErrUTFError          = "UTF_ERROR"
    ErrOutOfRange        = "OUT_OF_RANGE"
    ErrInvalidBitString  = "INVALID_BITSTRING"
    ErrOverflow          = "OVERFLOW"
    ErrPatternMismatch   = "PATTERN_MISMATCH"
)
```

**Методы:**
- `NewBitStringError(code, message string, context map[string]interface{}) *BitStringError`
- `ValidateSegment(segment Segment) error`
- `ValidateBitString(data []byte, length uint) error`
- `CheckOverflow(value interface{}, size uint) error`
- `CheckPatternMatch(expected, actual interface{}) error`

**Примеры обработки ошибок:**
```go
// Переполнение
_, err := builder.AddInteger(256, WithSize(8))
if err != nil {
    if errors.Is(err, ErrOverflow) {
        // Обработка переполнения
    }
}

// Несовпадение паттерна
results, err := matcher.Match(data)
if err != nil {
    if bitstringErr, ok := err.(*BitStringError); ok {
        fmt.Printf("Error %s: %s", bitstringErr.Code, bitstringErr.Message)
    }
}
```

**Критерии приемки:**
- Все типы ошибок обрабатываются корректно
- Детальная информация об ошибках предоставляется
- Переполнение обнаруживается правильно
- Несовпадение паттернов обрабатывается

---

### Этап 13: Вложенные битовые строки

#### 13.1 Nested Bitstrings
**Файлы:** `internal/bitstring/nested.go`, `examples/nested.go`

**Задачи:**
- Реализовать поддержку вложенных битовых строк
- Обработать конструирование вложенных структур
- Реализовать сопоставление вложенных паттернов
- Оптимизировать операции с вложенными структурами

**Методы:**
- `EncodeNestedBitstrings(segments []Segment) (*BitString, error)`
- `DecodeNestedBitstrings(data *BitString, pattern []Segment) ([]SegmentResult, error)`
- `ValidateNestedStructure(segments []Segment) error`
- `OptimizeNestedOperations(segments []Segment) []Segment`

**Пример использования:**
```go
// Вложенное конструирование
inner1 := NewBuilder().AddInteger(1).AddInteger(2).Build()
inner2 := NewBuilder().AddInteger(3).AddInteger(4).Build()

outer := NewBuilder().
    AddInteger(0).
    AddBitstring(inner1).
    AddBitstring(inner2).
    AddInteger(5).
    Build()

// Вложенное сопоставление
matcher := NewMatcher().
    Integer(&prefix, WithSize(8)).
    Bitstring(&inner1, WithSize(16)).
    Bitstring(&inner2, WithSize(16)).
    Integer(&suffix, WithSize(8))

results, err := matcher.Match(outer)
```

**Критерии приемки:**
- Вложенные битовые строки работают корректно
- Конструирование и сопоставление работают правильно
- Производительность оптимизирована
- Глубокая вложенность поддерживается

---

## Финальные этапы

### Тестирование
- Создать unit тесты для всех компонентов
- Реализовать интеграционные тесты
- Создать benchmark тесты для производительности
- Реализовать property-based тесты

### Документация
- Написать GoDoc для всех публичных API
- Создать README с примерами использования
- Написать руководство по миграции
- Создать performance guide

### Валидация
- Проверить соответствие техническому заданию
- Пройти все тесты из `BIT_SYNTAX_TESTS.md`
- Убедиться, что производительность соответствует требованиям
- Проверить интеграцию с Funterm

### Публикация
- Подготовить релиз
- Создать GitHub репозиторий
- Настроить CI/CD
- Опубликовать документацию

## Критерии успеха

### Функциональные требования
- [ ] Все функции из Core API реализованы
- [ ] Поддержка всех типов сегментов
- [ ] Поддержка всех спецификаторов
- [ ] Корректная работа endianness
- [ ] Поддержка UTF кодировок
- [ ] Все 13 этапов реализации завершены

### Качество кода
- [ ] Покрытие тестами ≥ 95%
- [ ] Нет memory leaks
- [ ] Производительность соответствует требованиям
- [ ] Полная документация
- [ ] Code review пройден

### Интеграция
- [ ] Совместимость с текущей версией Funterm
- [ ] Наличие layer совместимости
- [ ] Возможность постепенной миграции
- [ ] Примеры интеграции работают

### Производительность
- [ ] Конструирование: ≥ 1,000,000 операций/сек
- [ ] Matching: ≥ 500,000 операций/сек
- [ ] Memory overhead: ≤ 20% от размера данных
- [ ] Нет утечек памяти
- [ ] Бенчмарки стабильны

---