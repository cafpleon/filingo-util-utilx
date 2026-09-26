package utilx_test

import (
	"reflect"
	"testing"

	utilx "github.com/cafpleon/filingo-util-utilx"
)

func TestDerefString(t *testing.T) {
	tests := []struct {
		name  string
		input *string
		want  string
	}{
		{"nil", nil, ""},
		{"empty", utilx.ToPtr(""), ""},
		{"value", utilx.ToPtr("hello"), "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utilx.DerefString(tt.input)
			if got != tt.want {
				t.Errorf("DerefString() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestEqualJSONMaps(t *testing.T) {
	// Table-driven tests
	tests := []struct {
		name string
		a, b map[string]any
		want bool
	}{
		{
			name: "Ambos mapas vacíos",
			a:    map[string]any{},
			b:    map[string]any{},
			want: true,
		},
		{
			name: "Primitivos idénticos (Estilo JSON)",
			a:    map[string]any{"id": float64(123), "name": "Google", "active": true, "email": nil},
			b:    map[string]any{"id": float64(123), "name": "Google", "active": true, "email": nil},
			want: true,
		},
		{
			name: "Diferente longitud",
			a:    map[string]any{"id": float64(1)},
			b:    map[string]any{"id": float64(1), "extra": "data"},
			want: false,
		},
		{
			name: "Mismas llaves, diferentes valores",
			a:    map[string]any{"name": "Google"},
			b:    map[string]any{"name": "Apple"},
			want: false,
		},
		{
			name: "Objetos anidados idénticos (Fallback a reflect)",
			a:    map[string]any{"profile": map[string]any{"role": "admin", "age": float64(30)}},
			b:    map[string]any{"profile": map[string]any{"role": "admin", "age": float64(30)}},
			want: true,
		},
		{
			name: "Objetos anidados diferentes",
			a:    map[string]any{"profile": map[string]any{"role": "admin"}},
			b:    map[string]any{"profile": map[string]any{"role": "user"}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utilx.EqualJSONMaps(tt.a, tt.b); got != tt.want {
				t.Errorf("EqualJSONMaps() = %v, queremos %v", got, tt.want)
			}
		})
	}
}

func TestCloneJSONMap(t *testing.T) {
	t.Run("Retorna nil si el origen es nil", func(t *testing.T) {
		if got := utilx.CloneJSONMap(nil); got != nil {
			t.Errorf("CloneJSONMap(nil) = %v, esperábamos nil", got)
		}
	})

	t.Run("Garantizar copia profunda (Deep Copy)", func(t *testing.T) {
		// 1. Preparamos un payload original con datos anidados
		original := map[string]any{
			"provider": "google",
			"emails":   []any{"test@gmail.com", "work@gmail.com"},
			"profile": map[string]any{
				"verified": true,
				"picture":  "url_to_pic",
			},
		}

		// 2. Ejecutamos el clon
		clone := utilx.CloneJSONMap(original)

		// 3. MUTAMOS EL CLON (Cambiamos cosas en los anidados)
		clone["provider"] = "apple"

		// Modificamos el slice
		cloneEmails := clone["emails"].([]any)
		cloneEmails[0] = "hacked@gmail.com"

		// Modificamos el mapa anidado
		cloneProfile := clone["profile"].(map[string]any)
		cloneProfile["verified"] = false

		// 4. VERIFICAMOS EL ORIGINAL (No debió cambiar en absoluto)
		if original["provider"] == "apple" {
			t.Error("Copia superficial detectada: 'provider' mutó en el original")
		}

		origEmails := original["emails"].([]any)
		if origEmails[0] == "hacked@gmail.com" {
			t.Error("Copia superficial detectada: El slice 'emails' mutó en el original")
		}

		origProfile := original["profile"].(map[string]any)
		if origProfile["verified"] == false {
			t.Error("Copia superficial detectada: El mapa anidado 'profile' mutó en el original")
		}

		// 5. Verificamos estructuralmente que ya no son iguales con DeepEqual
		if reflect.DeepEqual(original, clone) {
			t.Error("El clon modificado sigue siendo igual al original")
		}
	})
}
