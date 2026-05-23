package utilx

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// StringValidationOpts define las reglas de negocio para un campo de texto.
type StringValidationOpts struct {
	Required  bool
	MinLength int
	MaxLength int
	Pattern   string
	// Si mañana necesitas validar correos, solo agregas:
	// IsEmail bool
}

// ValidateDomainString aplica las reglas del framework a un string.
func ValidateDomainString(fieldName, val string, opts StringValidationOpts) error {
	// 1. Obligatoriedad
	if opts.Required && val == "" {
		return fmt.Errorf("campo '%s' es obligatorio", fieldName)
	}

	// Si está vacío y no es obligatorio, saltamos el resto de validaciones
	if val == "" {
		return nil
	}

	length := utf8.RuneCountInString(val)

	// 2. Longitud Mínima
	if opts.MinLength > 0 && length < opts.MinLength {
		return fmt.Errorf("campo '%s' debe tener al menos %d caracteres", fieldName, opts.MinLength)
	}

	// 3. Longitud Máxima
	if opts.MaxLength > 0 && length > opts.MaxLength {
		return fmt.Errorf("campo '%s' no puede exceder los %d caracteres", fieldName, opts.MaxLength)
	}

	// 4. Patrón de Expresión Regular
	if opts.Pattern != "" {
		// Nota: Para optimizar aún más, en el futuro podrías precompilar las regex
		// y pasarlas como *regexp.Regexp en las opciones.
		matched, err := regexp.MatchString(opts.Pattern, val)
		if err != nil || !matched {
			return fmt.Errorf("campo '%s' tiene un formato inválido", fieldName)
		}
	}

	return nil
}

// Number es un constraint para que los genéricos acepten cualquier tipo numérico.
type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

// ValidateNumericFilter valida exclusión y rangos para cualquier número usando punteros (*T).
func ValidateNumericFilter[T Number](fieldName string, exact, min, max *T, maxAllowed T) error {
	count := 0
	if exact != nil {
		count++
	}
	if min != nil {
		count++
	}
	if max != nil {
		count++
	}

	if count > 1 {
		return fmt.Errorf("filtro '%s': use solo exact, min o max, no combinaciones", fieldName)
	}

	if min != nil && max != nil && *min > *max {
		return fmt.Errorf("filtro '%s': Min (%v) > Max (%v)", fieldName, *min, *max)
	}

	if min != nil && *min > maxAllowed {
		return fmt.Errorf("filtro '%s' Min supera el máximo permitido", fieldName)
	}

	if max != nil && *max > maxAllowed {
		return fmt.Errorf("filtro '%s' Max supera el máximo permitido", fieldName)
	}

	return nil
}

// ValidateStringFilter limpia y valida los filtros de texto de forma segura.
// Acepta punteros para poder modificar (TrimSpace) los valores originales del struct.
func ValidateStringFilter(fieldName string, exact, contains, prefix *string, minPartialLen, maxLen int) error {
	count := 0
	var activeVal string
	var isPartial bool

	// 1. Limpieza y conteo de variantes
	if exact != nil {
		*exact = strings.TrimSpace(*exact)
		if *exact != "" {
			count++
			activeVal = *exact
		}
	}

	if contains != nil {
		*contains = strings.TrimSpace(*contains)
		if *contains != "" {
			count++
			activeVal = *contains
			isPartial = true
		}
	}

	if prefix != nil {
		*prefix = strings.TrimSpace(*prefix)
		if *prefix != "" {
			count++
			activeVal = *prefix
			isPartial = true
		}
	}

	// 2. Exclusión Mutua
	if count > 1 {
		return fmt.Errorf("conflicto en filtro '%s': use solo una variante (exact, contains, prefix)", fieldName)
	}

	// 3. Validaciones de Seguridad y Dominio
	if activeVal != "" {
		// Defensa 1: Longitud Máxima (DoS por Memoria)
		if utf8.RuneCountInString(activeVal) > maxLen {
			return fmt.Errorf("filtro '%s' excede el máximo de %d caracteres", fieldName, maxLen)
		}

		// Defensa 2: Longitud Mínima (DoS por CPU/Disco en consultas LIKE)
		if isPartial && utf8.RuneCountInString(activeVal) < minPartialLen {
			return fmt.Errorf("filtro '%s' parcial requiere al menos %d caracteres", fieldName, minPartialLen)
		}

		// Defensa 3: Anti-Wildcards (SQL Injection Light)
		if strings.ContainsAny(activeVal, "%_") {
			return fmt.Errorf("filtro '%s' contiene caracteres prohibidos (%%, _)", fieldName)
		}
	}

	return nil
}

// ValidateDateFilterRange asegura que los filtros de fecha no se contradigan y no excedan el rango.
func ValidateDateFilterRange(fieldName string, exact, from, to time.Time, maxDays float64) error {
	count := 0
	if !exact.IsZero() {
		count++
	}
	if !from.IsZero() {
		count++
	}
	if !to.IsZero() {
		count++
	}

	if count > 1 {
		return fmt.Errorf("filtro '%s': use solo exact, from o to", fieldName)
	}

	if !from.IsZero() && !to.IsZero() {
		if from.After(to) {
			return fmt.Errorf("filtro '%s': From (%v) después de To (%v)", fieldName, from, to)
		}

		days := to.Sub(from).Hours() / 24
		if days > maxDays {
			return fmt.Errorf("filtro '%s': rango de fechas excede %d días", fieldName, int(maxDays))
		}
	}

	return nil
}
