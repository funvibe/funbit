package protocols

import (
	"fmt"
	"strings"

	"github.com/funvibe/funbit/internal/bitstring"
	"github.com/funvibe/funbit/internal/builder"
	"github.com/funvibe/funbit/internal/matcher"
)

// TCPFlags представляет структуру TCP флагов
type TCPFlags struct {
	Reserved uint // Зарезервированные биты (2 бита)
	URG      uint // Urgent flag (1 бит)
	ACK      uint // Acknowledgment flag (1 бит)
	PSH      uint // Push flag (1 бит)
	RST      uint // Reset flag (1 бит)
	SYN      uint // Synchronize flag (1 бит)
	FIN      uint // Finish flag (1 бит)
}

// TCPHeader представляет полный TCP заголовок
type TCPHeader struct {
	SourcePort      uint     // Порт источника (16 бит)
	DestinationPort uint     // Порт назначения (16 бит)
	SequenceNumber  uint     // Номер последовательности (32 бита)
	Acknowledgment  uint     // Номер подтверждения (32 бита)
	DataOffset      uint     // Смещение данных (4 бита)
	Reserved        uint     // Зарезервировано (3 бита)
	Flags           TCPFlags // Флаги (9 бит: 2 зарезервированных + 7 флагов)
	Window          uint     // Размер окна (16 бит)
	Checksum        uint     // Контрольная сумма (16 бит)
	UrgentPointer   uint     // Указатель срочности (16 бит)
}

// BuildTCPFlags создает TCP флаги из структуры
func BuildTCPFlags(flags TCPFlags) (*bitstring.BitString, error) {
	return builder.NewBuilder().
		AddInteger(flags.Reserved, bitstring.WithSize(2)).
		AddInteger(flags.URG, bitstring.WithSize(1)).
		AddInteger(flags.ACK, bitstring.WithSize(1)).
		AddInteger(flags.PSH, bitstring.WithSize(1)).
		AddInteger(flags.RST, bitstring.WithSize(1)).
		AddInteger(flags.SYN, bitstring.WithSize(1)).
		AddInteger(flags.FIN, bitstring.WithSize(1)).
		Build()
}

// ParseTCPFlags парсит TCP флаги из битовой строки
func ParseTCPFlags(data *bitstring.BitString) (*TCPFlags, error) {
	var flags TCPFlags

	_, err := matcher.NewMatcher().
		Integer(&flags.Reserved, bitstring.WithSize(2)).
		Integer(&flags.URG, bitstring.WithSize(1)).
		Integer(&flags.ACK, bitstring.WithSize(1)).
		Integer(&flags.PSH, bitstring.WithSize(1)).
		Integer(&flags.RST, bitstring.WithSize(1)).
		Integer(&flags.SYN, bitstring.WithSize(1)).
		Integer(&flags.FIN, bitstring.WithSize(1)).
		Match(data)

	if err != nil {
		return nil, err
	}

	return &flags, nil
}

// BuildTCPHeader создает полный TCP заголовок из структуры
func BuildTCPHeader(header TCPHeader) (*bitstring.BitString, error) {
	return builder.NewBuilder().
		AddInteger(header.SourcePort, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		AddInteger(header.DestinationPort, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		AddInteger(header.SequenceNumber, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		AddInteger(header.Acknowledgment, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		AddInteger(header.DataOffset, bitstring.WithSize(4)).
		AddInteger(header.Reserved, bitstring.WithSize(3)).
		AddInteger(header.Flags.Reserved, bitstring.WithSize(2)).
		AddInteger(header.Flags.URG, bitstring.WithSize(1)).
		AddInteger(header.Flags.ACK, bitstring.WithSize(1)).
		AddInteger(header.Flags.PSH, bitstring.WithSize(1)).
		AddInteger(header.Flags.RST, bitstring.WithSize(1)).
		AddInteger(header.Flags.SYN, bitstring.WithSize(1)).
		AddInteger(header.Flags.FIN, bitstring.WithSize(1)).
		AddInteger(header.Window, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		AddInteger(header.Checksum, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		AddInteger(header.UrgentPointer, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Build()
}

// ParseTCPHeader парсит полный TCP заголовок из битовой строки
func ParseTCPHeader(data *bitstring.BitString) (*TCPHeader, error) {
	var header TCPHeader

	_, err := matcher.NewMatcher().
		Integer(&header.SourcePort, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Integer(&header.DestinationPort, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Integer(&header.SequenceNumber, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		Integer(&header.Acknowledgment, bitstring.WithSize(32), bitstring.WithEndianness("big")).
		Integer(&header.DataOffset, bitstring.WithSize(4)).
		Integer(&header.Reserved, bitstring.WithSize(3)).
		Integer(&header.Flags.Reserved, bitstring.WithSize(2)).
		Integer(&header.Flags.URG, bitstring.WithSize(1)).
		Integer(&header.Flags.ACK, bitstring.WithSize(1)).
		Integer(&header.Flags.PSH, bitstring.WithSize(1)).
		Integer(&header.Flags.RST, bitstring.WithSize(1)).
		Integer(&header.Flags.SYN, bitstring.WithSize(1)).
		Integer(&header.Flags.FIN, bitstring.WithSize(1)).
		Integer(&header.Window, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Integer(&header.Checksum, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Integer(&header.UrgentPointer, bitstring.WithSize(16), bitstring.WithEndianness("big")).
		Match(data)

	if err != nil {
		return nil, err
	}

	return &header, nil
}

// ValidateTCPFlags выполняет валидацию TCP флагов
func ValidateTCPFlags(flags *TCPFlags) error {
	if flags == nil {
		return fmt.Errorf("flags is nil")
	}

	// Зарезервированные биты должны быть 0
	if flags.Reserved != 0 {
		return fmt.Errorf("reserved bits must be zero")
	}

	// Проверка значений флагов (должны быть 0 или 1)
	flagValues := []uint{flags.URG, flags.ACK, flags.PSH, flags.RST, flags.SYN, flags.FIN}
	for i, value := range flagValues {
		if value > 1 {
			flagNames := []string{"URG", "ACK", "PSH", "RST", "SYN", "FIN"}
			return fmt.Errorf("%s flag must be 0 or 1", flagNames[i])
		}
	}

	return nil
}

// ValidateTCPHeader выполняет полную валидацию TCP заголовка
func ValidateTCPHeader(header *TCPHeader) error {
	if header == nil {
		return fmt.Errorf("header is nil")
	}

	// Валидация портов
	if header.SourcePort > 65535 {
		return fmt.Errorf("invalid source port")
	}

	if header.DestinationPort > 65535 {
		return fmt.Errorf("invalid destination port")
	}

	// Валидация смещения данных
	if header.DataOffset < 5 || header.DataOffset > 15 {
		return fmt.Errorf("invalid data offset, must be between 5 and 15")
	}

	// Валидация зарезервированных битов
	if header.Reserved != 0 {
		return fmt.Errorf("reserved bits must be zero")
	}

	// Валидация флагов
	if err := ValidateTCPFlags(&header.Flags); err != nil {
		return err
	}

	// Валидация размера окна
	if header.Window > 65535 {
		return fmt.Errorf("invalid window size")
	}

	return nil
}

// GetTCPHeaderLength возвращает длину TCP заголовка в байтах
func GetTCPHeaderLength(header *TCPHeader) uint {
	return header.DataOffset * 4
}

// GetTCPFlagsString возвращает строковое представление установленных флагов
func GetTCPFlagsString(flags TCPFlags) string {
	var flagStrings []string

	if flags.URG == 1 {
		flagStrings = append(flagStrings, "URG")
	}
	if flags.ACK == 1 {
		flagStrings = append(flagStrings, "ACK")
	}
	if flags.PSH == 1 {
		flagStrings = append(flagStrings, "PSH")
	}
	if flags.RST == 1 {
		flagStrings = append(flagStrings, "RST")
	}
	if flags.SYN == 1 {
		flagStrings = append(flagStrings, "SYN")
	}
	if flags.FIN == 1 {
		flagStrings = append(flagStrings, "FIN")
	}

	if len(flagStrings) == 0 {
		return "NONE"
	}

	return strings.Join(flagStrings, "|")
}

// IsTCPConnectionEstablishment проверяет, является ли пакет установлением соединения (SYN)
func IsTCPConnectionEstablishment(flags TCPFlags) bool {
	return flags.SYN == 1 && flags.ACK == 0
}

// IsTCPConnectionEstablished проверяет, является ли пакет подтверждением установки соединения (SYN+ACK)
func IsTCPConnectionEstablished(flags TCPFlags) bool {
	return flags.SYN == 1 && flags.ACK == 1
}

// IsTCPConnectionTermination проверяет, является ли пакет завершением соединения (FIN или RST)
func IsTCPConnectionTermination(flags TCPFlags) bool {
	return flags.FIN == 1 || flags.RST == 1
}

// IsTCPDataPacket проверяет, содержит ли пакет данные (PSH или ACK)
func IsTCPDataPacket(flags TCPFlags) bool {
	return flags.PSH == 1 && flags.ACK == 1
}
