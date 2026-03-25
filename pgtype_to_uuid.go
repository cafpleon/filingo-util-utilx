package utilx

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgtypeToUUID(pgUUID pgtype.UUID) (uuid.UUID, error) {
	if pgUUID.Valid {
		// pgtype.UUID.Bytes es un [16]byte, que es exactamente lo que necesita uuid.FromBytes
		return uuid.FromBytes(pgUUID.Bytes[:])
	} else {
		return uuid.Nil, fmt.Errorf("UUID no presente (NULL)")
	}
}

// IsValidUUID verifica si un UUID es válido (no es nil/zero)
func IsValidUUID(id uuid.UUID) bool {
	return id != uuid.Nil
}

// IsValidUUIDString verifica si un string es un UUID válido
func IsValidUUIDString(s string) bool {
	_, err := uuid.FromString(s)
	return err == nil
}

// ParseUUID convierte string a UUID con validación
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.FromString(s)
}

// ParseUUIDOrNil convierte string a UUID, retorna Nil si es inválido
func ParseUUIDOrNil(s string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// MustParseUUID convierte string a UUID, panic si es inválido (solo para tests/init)
func MustParseUUID(s string) uuid.UUID {
	return uuid.Must(uuid.FromString(s))
}

// IsZeroOrNil verifica si el UUID es cero o nil (sinónimo de IsValidUUID)
func IsZeroOrNil(id uuid.UUID) bool {
	return id == uuid.Nil
}

// StringOrEmpty devuelve el string del UUID o vacío si es inválido
func StringOrEmpty(id uuid.UUID) string {
	if id == uuid.Nil {
		return ""
	}
	return id.String()
}

// NewUUID genera un nuevo UUID v7 con validación de error
func NewUUID() (uuid.UUID, error) {
	return uuid.NewV7()
}

// MustNewUUID genera un nuevo UUID v7, panic si falla (solo para tests/init)
func MustNewUUID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

// FormatUUID formatea un UUID para uso en logs o mensajes
func FormatUUID(id uuid.UUID) string {
	if id == uuid.Nil {
		return "<nil-uuid>"
	}
	return id.String()
}

// CompareUUIDs compara dos UUIDs y retorna true si son iguales
func CompareUUIDs(a, b uuid.UUID) bool {
	return a == b
}
