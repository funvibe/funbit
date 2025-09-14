package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/funvibe/funbit/internal/bitstring"
	"github.com/funvibe/funbit/internal/builder"
	"github.com/funvibe/funbit/internal/matcher"
)

// PNGSignature представляет сигнатуру PNG файла
var PNGSignature = []byte{137, 80, 78, 71, 13, 10, 26, 10}

// PNGChunk представляет структуру PNG чанка
type PNGChunk struct {
	Length uint32 // Длина данных чанка (4 байта)
	Type   []byte // Тип чанка (4 байта)
	Data   []byte // Данные чанка (переменная длина)
	CRC    uint32 // Контрольная сумма CRC32 (4 байта)
}

// PNGHeader представляет структуру PNG заголовка
type PNGHeader struct {
	Signature []byte     // Сигнатура PNG (8 байт)
	Chunks    []PNGChunk // Чанки PNG
}

// PNGIHDRChunk представляет IHDR чанк PNG (заголовок изображения)
type PNGIHDRChunk struct {
	Width       uint32 // Ширина изображения (4 байта)
	Height      uint32 // Высота изображения (4 байта)
	BitDepth    uint8  // Глубина цвета (1 байт)
	ColorType   uint8  // Тип цвета (1 байт)
	Compression uint8  // Метод сжатия (1 байт)
	Filter      uint8  // Метод фильтрации (1 байт)
	Interlace   uint8  // Метод чередования (1 байт)
}

// BuildPNGHeader создает PNG заголовок из структуры
func BuildPNGHeader(header PNGHeader) (*bitstring.BitString, error) {
	// Проверяем сигнатуру
	if !bytes.Equal(header.Signature, PNGSignature) {
		return nil, errors.New("invalid PNG signature")
	}

	// Начинаем с сигнатуры
	b := builder.NewBuilder().AddBinary(header.Signature)

	// Добавляем чанки
	for _, chunk := range header.Chunks {
		chunkData, err := BuildPNGChunk(chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to build chunk %s: %v", string(chunk.Type), err)
		}
		b.AddBinary(chunkData)
	}

	return b.Build()
}

// BuildPNGChunk создает PNG чанк из структуры
func BuildPNGChunk(chunk PNGChunk) ([]byte, error) {
	if len(chunk.Type) != 4 {
		return nil, errors.New("chunk type must be 4 bytes")
	}

	// Создаем буфер для чанка
	var buf bytes.Buffer

	// Записываем длину данных (big-endian)
	if err := binary.Write(&buf, binary.BigEndian, chunk.Length); err != nil {
		return nil, fmt.Errorf("failed to write chunk length: %v", err)
	}

	// Записываем тип чанка
	if _, err := buf.Write(chunk.Type); err != nil {
		return nil, fmt.Errorf("failed to write chunk type: %v", err)
	}

	// Записываем данные
	if _, err := buf.Write(chunk.Data); err != nil {
		return nil, fmt.Errorf("failed to write chunk data: %v", err)
	}

	// Записываем CRC (big-endian)
	if err := binary.Write(&buf, binary.BigEndian, chunk.CRC); err != nil {
		return nil, fmt.Errorf("failed to write chunk CRC: %v", err)
	}

	return buf.Bytes(), nil
}

// BuildPNGIHDRChunk создает IHDR чанк из структуры
func BuildPNGIHDRChunk(ihdr PNGIHDRChunk) (PNGChunk, error) {
	// Валидация IHDR чанка
	if err := ValidatePNGIHDRChunk(&ihdr); err != nil {
		return PNGChunk{}, err
	}

	// Создаем данные IHDR чанка (13 байт)
	data := make([]byte, 13)

	// Ширина и высота (big-endian)
	binary.BigEndian.PutUint32(data[0:4], ihdr.Width)
	binary.BigEndian.PutUint32(data[4:8], ihdr.Height)

	// Остальные поля
	data[8] = ihdr.BitDepth
	data[9] = ihdr.ColorType
	data[10] = ihdr.Compression
	data[11] = ihdr.Filter
	data[12] = ihdr.Interlace

	// Вычисляем CRC (для простоты используем 0)
	crc := uint32(0) // В реальном коде здесь был бы расчет CRC32

	return PNGChunk{
		Length: 13,
		Type:   []byte("IHDR"),
		Data:   data,
		CRC:    crc,
	}, nil
}

// ParsePNGHeader парсит PNG заголовок из битовой строки
func ParsePNGHeader(data *bitstring.BitString) (*PNGHeader, error) {
	var header PNGHeader

	// Парсим сигнатуру
	_, err := matcher.NewMatcher().
		Binary(&header.Signature, bitstring.WithSize(8)).
		Match(data)

	if err != nil {
		return nil, fmt.Errorf("failed to parse PNG signature: %v", err)
	}

	// Проверяем сигнатуру
	if !bytes.Equal(header.Signature, PNGSignature) {
		return nil, errors.New("invalid PNG signature")
	}

	// Получаем оставшиеся данные для парсинга чанков
	remainingData := data.ToBytes()[8:]

	// Парсим чанки
	for len(remainingData) >= 12 { // Минимальный размер чанка: 4 (len) + 4 (type) + 4 (crc)
		chunk, chunkSize, err := ParsePNGChunk(remainingData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PNG chunk: %v", err)
		}

		header.Chunks = append(header.Chunks, chunk)
		remainingData = remainingData[chunkSize:]

		// Для простоты парсим только первый чанк (IHDR)
		break
	}

	return &header, nil
}

// ParsePNGChunk парсит PNG чанк из байтов
func ParsePNGChunk(data []byte) (PNGChunk, uint, error) {
	if len(data) < 12 {
		return PNGChunk{}, 0, errors.New("insufficient data for PNG chunk")
	}

	// Читаем длину (big-endian)
	length := binary.BigEndian.Uint32(data[0:4])

	// Читаем тип
	chunkType := data[4:8]

	// Проверяем, достаточно ли данных
	totalChunkSize := 12 + int(length) // 4 (len) + 4 (type) + length + 4 (crc)
	if len(data) < totalChunkSize {
		return PNGChunk{}, 0, fmt.Errorf("insufficient data for chunk %s, need %d bytes, have %d",
			string(chunkType), totalChunkSize, len(data))
	}

	// Читаем данные
	chunkData := data[8 : 8+length]

	// Читаем CRC (big-endian)
	crc := binary.BigEndian.Uint32(data[8+length : 12+length])

	chunk := PNGChunk{
		Length: length,
		Type:   chunkType,
		Data:   chunkData,
		CRC:    crc,
	}

	return chunk, uint(totalChunkSize), nil
}

// ParsePNGIHDRChunk парсит IHDR чанк из данных
func ParsePNGIHDRChunk(data []byte) (*PNGIHDRChunk, error) {
	if len(data) < 13 {
		return nil, errors.New("insufficient data for IHDR chunk, need 13 bytes")
	}

	ihdr := &PNGIHDRChunk{
		Width:       binary.BigEndian.Uint32(data[0:4]),
		Height:      binary.BigEndian.Uint32(data[4:8]),
		BitDepth:    data[8],
		ColorType:   data[9],
		Compression: data[10],
		Filter:      data[11],
		Interlace:   data[12],
	}

	// Валидация
	if err := ValidatePNGIHDRChunk(ihdr); err != nil {
		return nil, err
	}

	return ihdr, nil
}

// ValidatePNGHeader выполняет валидацию PNG заголовка
func ValidatePNGHeader(header *PNGHeader) error {
	if header == nil {
		return errors.New("header is nil")
	}

	// Проверяем сигнатуру
	if !bytes.Equal(header.Signature, PNGSignature) {
		return errors.New("invalid PNG signature")
	}

	// Проверяем наличие чанков
	if len(header.Chunks) == 0 {
		return errors.New("PNG must have at least one chunk")
	}

	// Первый чанк должен быть IHDR
	if string(header.Chunks[0].Type) != "IHDR" {
		return errors.New("first PNG chunk must be IHDR")
	}

	// Валидация каждого чанка
	for i, chunk := range header.Chunks {
		if err := ValidatePNGChunk(&chunk); err != nil {
			return fmt.Errorf("chunk %d (%s) validation failed: %v", i, string(chunk.Type), err)
		}
	}

	return nil
}

// ValidatePNGChunk выполняет валидацию PNG чанка
func ValidatePNGChunk(chunk *PNGChunk) error {
	if chunk == nil {
		return errors.New("chunk is nil")
	}

	if len(chunk.Type) != 4 {
		return errors.New("chunk type must be 4 bytes")
	}

	// Проверяем, что тип состоит из букв латинского алфавита
	for _, c := range chunk.Type {
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
			return fmt.Errorf("chunk type must contain only letters, got: %c", c)
		}
	}

	// Проверяем длину данных
	if uint32(len(chunk.Data)) != chunk.Length {
		return fmt.Errorf("chunk data length mismatch, expected %d, got %d",
			chunk.Length, len(chunk.Data))
	}

	return nil
}

// ValidatePNGIHDRChunk выполняет валидацию IHDR чанка
func ValidatePNGIHDRChunk(ihdr *PNGIHDRChunk) error {
	if ihdr == nil {
		return errors.New("IHDR chunk is nil")
	}

	// Проверяем ширину и высоту
	if ihdr.Width == 0 {
		return errors.New("image width must be greater than 0")
	}

	if ihdr.Height == 0 {
		return errors.New("image height must be greater than 0")
	}

	// Проверяем глубину цвета
	validBitDepths := []uint8{1, 2, 4, 8, 16}
	isValidBitDepth := false
	for _, bd := range validBitDepths {
		if ihdr.BitDepth == bd {
			isValidBitDepth = true
			break
		}
	}
	if !isValidBitDepth {
		return fmt.Errorf("invalid bit depth %d, must be one of: %v", ihdr.BitDepth, validBitDepths)
	}

	// Проверяем тип цвета
	validColorTypes := []uint8{0, 2, 3, 4, 6}
	isValidColorType := false
	for _, ct := range validColorTypes {
		if ihdr.ColorType == ct {
			isValidColorType = true
			break
		}
	}
	if !isValidColorType {
		return fmt.Errorf("invalid color type %d, must be one of: %v", ihdr.ColorType, validColorTypes)
	}

	// Проверяем совместимость глубины цвета и типа цвета
	if ihdr.ColorType == 3 && ihdr.BitDepth < 8 {
		return errors.New("indexed color (type 3) requires bit depth of 8")
	}

	// Проверяем метод сжатия (должен быть 0)
	if ihdr.Compression != 0 {
		return errors.New("compression method must be 0")
	}

	// Проверяем метод фильтрации (должен быть 0)
	if ihdr.Filter != 0 {
		return errors.New("filter method must be 0")
	}

	// Проверяем метод чередования (должен быть 0 или 1)
	if ihdr.Interlace > 1 {
		return errors.New("interlace method must be 0 or 1")
	}

	return nil
}

// GetPNGColorTypeName возвращает название типа цвета
func GetPNGColorTypeName(colorType uint8) string {
	switch colorType {
	case 0:
		return "Grayscale"
	case 2:
		return "RGB"
	case 3:
		return "Indexed"
	case 4:
		return "Grayscale with Alpha"
	case 6:
		return "RGB with Alpha"
	default:
		return "Unknown"
	}
}

// GetPNGInterlaceMethodName возвращает название метода чередования
func GetPNGInterlaceMethodName(interlace uint8) string {
	switch interlace {
	case 0:
		return "None"
	case 1:
		return "Adam7"
	default:
		return "Unknown"
	}
}

// FormatPNGIHDRInfo форматирует информацию об IHDR чанке для вывода
func FormatPNGIHDRInfo(ihdr *PNGIHDRChunk) string {
	var info strings.Builder

	info.WriteString(fmt.Sprintf("PNG Image Header:\n"))
	info.WriteString(fmt.Sprintf("  Width: %d px\n", ihdr.Width))
	info.WriteString(fmt.Sprintf("  Height: %d px\n", ihdr.Height))
	info.WriteString(fmt.Sprintf("  Bit Depth: %d bits\n", ihdr.BitDepth))
	info.WriteString(fmt.Sprintf("  Color Type: %d (%s)\n", ihdr.ColorType, GetPNGColorTypeName(ihdr.ColorType)))
	info.WriteString(fmt.Sprintf("  Compression: %d\n", ihdr.Compression))
	info.WriteString(fmt.Sprintf("  Filter: %d\n", ihdr.Filter))
	info.WriteString(fmt.Sprintf("  Interlace: %d (%s)\n", ihdr.Interlace, GetPNGInterlaceMethodName(ihdr.Interlace)))

	return info.String()
}
