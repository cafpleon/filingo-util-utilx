package utilx

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gofrs/uuid/v5"
)

// Required valida que un string no esté vacío (tras trim)
func Required(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("campo '%s' es obligatorio", fieldName)
	}
	return nil
}

// RequiredWithCustomMsg permite mensaje personalizado
func RequiredWithCustomMsg(value, fieldName string, customMsg string, args ...any) error {
	if strings.TrimSpace(value) == "" {
		if customMsg != "" {
			if len(args) > 0 {
				return fmt.Errorf(customMsg, args...)
			}
			return errors.New(customMsg)
		}
		return fmt.Errorf("campo '%s' es obligatorio", fieldName)
	}
	return nil
}

// RequiredMinLength valida longitud mínima
func RequiredMinLength(value, fieldName string, minLength int) error {
	if err := Required(value, fieldName); err != nil {
		return err
	}
	if utf8.RuneCountInString(strings.TrimSpace(value)) < minLength {
		return fmt.Errorf("campo '%s' debe tener al menos %d caracteres", fieldName, minLength)
	}
	return nil
}

// RequiredMaxLength valida longitud máxima
func RequiredMaxLength(value, fieldName string, maxLength int) error {
	if utf8.RuneCountInString(value) > maxLength {
		return fmt.Errorf("campo '%s' no puede exceder %d caracteres", fieldName, maxLength)
	}
	return nil
}

// Enum valida que un valor esté en un conjunto de strings
func Enum(value string, validValues []string, fieldName string) error {
	for _, valid := range validValues {
		if value == valid {
			return nil
		}
	}
	return fmt.Errorf("campo '%s' tiene valor inválido '%s' (valores válidos: %v)",
		fieldName, value, validValues)
}

// EnumCaseInsensitive valida ignorando mayúsculas/minúsculas
func EnumCaseInsensitive(value string, validValues []string, fieldName string) error {
	valueNorm := strings.ToUpper(strings.TrimSpace(value))
	for _, valid := range validValues {
		if valueNorm == strings.ToUpper(valid) {
			return nil
		}
	}
	return fmt.Errorf("campo '%s' tiene valor inválido '%s'", fieldName, value)
}

// ValidateUUID verifica que el string tenga un formato UUID válido (versión 4 o cualquier UUID estándar).
// Retorna error si es inválido o si es el UUID cero (opcional, según tu regla de negocio).
func ValidateUUID(s string) error {
	if s == "" {
		return fmt.Errorf("UUID no puede estar vacío")
	}
	_, err := uuid.FromString(s)
	if err != nil {
		return fmt.Errorf("formato UUID inválido: %s", s)
	}
	// Opcional: si quieres rechazar el UUID cero (00000000-0000-0000-0000-000000000000)
	// if u == uuid.Nil {
	// 	return fmt.Errorf("UUID no puede ser el valor cero")
	// }
	return nil
}
