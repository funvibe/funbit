package protocols

import (
	"errors"

	"github.com/funvibe/funbit/internal/bitstring"
	"github.com/funvibe/funbit/internal/builder"
	"github.com/funvibe/funbit/internal/matcher"
)

// IPv4Header представляет структуру IPv4 заголовка
type IPv4Header struct {
	Version        uint // Версия (4 бита)
	HeaderLength   uint // Длина заголовка (4 бита)
	ServiceType    uint // Тип сервиса (8 бит)
	TotalLength    uint // Общая длина (16 бит)
	Identification uint // Идентификация (16 бит)
	Flags          uint // Флаги (3 бита)
	FragmentOffset uint // Смещение фрагмента (13 бит)
	TTL            uint // Время жизни (8 бит)
	Protocol       uint // Протокол (8 бит)
	Checksum       uint // Контрольная сумма (16 бит)
	SourceIP       uint // IP-адрес источника (32 бита)
	DestinationIP  uint // IP-адрес назначения (32 бита)
}

// BuildIPv4Header создает IPv4 заголовок из структуры
func BuildIPv4Header(header IPv4Header) (*bitstring.BitString, error) {
	return builder.NewBuilder().
		AddInteger(header.Version, bitstring.WithSize(4)).
		AddInteger(header.HeaderLength, bitstring.WithSize(4)).
		AddInteger(header.ServiceType, bitstring.WithSize(8)).
		AddInteger(header.TotalLength, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		AddInteger(header.Identification, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		AddInteger(header.Flags, bitstring.WithSize(3)).
		AddInteger(header.FragmentOffset, bitstring.WithSize(13)).
		AddInteger(header.TTL, bitstring.WithSize(8)).
		AddInteger(header.Protocol, bitstring.WithSize(8)).
		AddInteger(header.Checksum, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		AddInteger(header.SourceIP, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		AddInteger(header.DestinationIP, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		Build()
}

// ParseIPv4Header парсит IPv4 заголовок из битовой строки
func ParseIPv4Header(data *bitstring.BitString) (*IPv4Header, error) {
	var header IPv4Header

	_, err := matcher.NewMatcher().
		Integer(&header.Version, bitstring.WithSize(4)).
		Integer(&header.HeaderLength, bitstring.WithSize(4)).
		Integer(&header.ServiceType, bitstring.WithSize(8)).
		Integer(&header.TotalLength, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Integer(&header.Identification, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Integer(&header.Flags, bitstring.WithSize(3)).
		Integer(&header.FragmentOffset, bitstring.WithSize(13)).
		Integer(&header.TTL, bitstring.WithSize(8)).
		Integer(&header.Protocol, bitstring.WithSize(8)).
		Integer(&header.Checksum, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Integer(&header.SourceIP, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		Integer(&header.DestinationIP, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		Match(data)

	if err != nil {
		return nil, err
	}

	// Базовая валидация IPv4 заголовка
	if header.Version != 4 {
		return nil, errors.New("invalid IPv4 version")
	}

	if header.HeaderLength < 5 {
		return nil, errors.New("invalid IPv4 header length")
	}

	return &header, nil
}

// ValidateIPv4Header выполняет полную валидацию IPv4 заголовка
func ValidateIPv4Header(header *IPv4Header) error {
	if header == nil {
		return errors.New("header is nil")
	}

	if header.Version != 4 {
		return errors.New("invalid IPv4 version, must be 4")
	}

	if header.HeaderLength < 5 || header.HeaderLength > 15 {
		return errors.New("invalid IPv4 header length, must be between 5 and 15")
	}

	if header.TotalLength < 20 {
		return errors.New("invalid total length, must be at least 20 bytes")
	}

	// Проверка зарезервированных бит флагов
	if header.Flags&0xE0 != 0 {
		return errors.New("reserved bits in flags must be zero")
	}

	return nil
}

// GetIPv4HeaderLength возвращает длину заголовка в байтах
func GetIPv4HeaderLength(header *IPv4Header) uint {
	return header.HeaderLength * 4
}

// GetIPv4PayloadLength возвращает длину полезной нагрузки в байтах
func GetIPv4PayloadLength(header *IPv4Header) uint {
	return header.TotalLength - GetIPv4HeaderLength(header)
}
