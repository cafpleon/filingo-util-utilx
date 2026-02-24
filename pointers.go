// Package utilx tiene varios utilidades
package utilx

import "database/sql"

// DerefString devuelve el valor string o "" si es nil
func DerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StringPointer devuelve el valor string o  nil
func StringPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// DerefStringWithDefault devuelve el valor string o un valor por defecto
func DerefStringWithDefault(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// DerefInt devuelve el valor int o 0 si es nil
func DerefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// DerefBool devuelve el valor bool o false si es nil
func DerefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// ToPtr convierte un valor a puntero (útil para literales)
func ToPtr[T any](v T) *T {
	return &v
}
