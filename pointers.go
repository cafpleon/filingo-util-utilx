// Package utilx tiene varios utilidades
package utilx

import (
	"database/sql"
	"reflect"
)

// DerefString (dereference string) devuelve el valor string o "" si es nil
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

// NullStringToPointer devuelve el valor string o  nil
func NullStringToPointer(ns sql.NullString) *string {
	return StringPointer(ns)
}

// DerefStringWithDefault (dereference string) devuelve el valor string o un valor por defecto
func DerefStringWithDefault(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

// DerefInt (dereference integer) devuelve el valor int o 0 si es nil
func DerefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// DerefBool (dereference booleano) devuelve el valor bool o false si es nil
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

// CompareStringPointers para comparar strings pointers
func CompareStringPointers(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// EqualJSONMaps
// esta comparación está pensada para mapas que provienen de un JSON (justificando por qué se espera un float64).
func EqualJSONMaps(a, b map[string]any) bool {
	// 1. Protección contra pánico y atajo de punteros
	//if a == b {
	//	return true // Son el mismo puntero, o ambos son nil
	//}
	if a == nil || b == nil {
		return false // Uno es nil y el otro no
	}

	// 2. Desreferenciamos los punteros para trabajar con los mapas
	//mA := *a
	//mB := *b

	//if len(mA) != len(mB) {
	//	return false
	//}

	if len(a) != len(b) {
		return false
	}

	for k, vA := range a { //mA {
		vB, ok := b[k] // mB[k]
		if !ok {
			return false
		}

		// 3. Comparación rápida de tipos (Ajustada para JSON)
		switch valA := vA.(type) {
		case string:
			if valB, ok := vB.(string); !ok || valA != valB {
				return false
			}
		case float64:
			// JSON SIEMPRE decodifica los números como float64 en los map[string]any
			if valB, ok := vB.(float64); !ok || valA != valB {
				return false
			}
		case bool:
			if valB, ok := vB.(bool); !ok || valA != valB {
				return false
			}
		case nil:
			// Maneja los null de JSON explícitamente sin usar reflect
			if vB != nil {
				return false
			}
		default:
			// Fallback a DeepEqual para arreglos ([]any) u objetos anidados (map[string]any)
			if !reflect.DeepEqual(vA, vB) {
				return false
			}
		}
	}
	return true
}

// CloneJSONMap crea una copia profunda de un payload JSON.
func CloneJSONMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}

	dst := make(map[string]any, len(src))
	for k, v := range src {
		// Evaluamos si el valor es un tipo compuesto (anidado) de JSON
		switch val := v.(type) {
		case map[string]any:
			dst[k] = CloneJSONMap(val) // Llamada recursiva para objetos anidados
		case []any:
			// Clonamos arreglos/slices
			newSlice := make([]any, len(val))
			for i, item := range val {
				// (Opcional) Si esperas arreglos de arreglos, podrías necesitar más lógica,
				// pero para JSON estándar con 1 nivel de arreglos, una copia simple sirve.
				newSlice[i] = item
			}
			dst[k] = newSlice
		default:
			// Primitivos (string, float64, bool, nil) se copian por valor sin riesgo
			dst[k] = v
		}
	}
	return dst
}
