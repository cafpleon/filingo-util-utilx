package utilx

import (
	"testing"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantError bool
	}{
		{"empty", "", "Nombre", true},
		{"spaces", "   ", "Email", true},
		{"valid", "Juan", "Nombre", false},
		{"tab", "\t", "Campo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Required(tt.value, tt.fieldName)
			if (err != nil) != tt.wantError {
				t.Errorf("Required() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestRequiredWithCustomMsg(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		msg       string
		args      []interface{}
		wantError string
	}{
		{"empty with default", "", "", nil, "campo '' es obligatorio"},
		{"empty with custom", "", "Error personalizado", nil, "Error personalizado"},
		{"empty with format", "", "Error: %s", []interface{}{"detalle"}, "Error: detalle"},
		{"valid", "test", "", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RequiredWithCustomMsg(tt.value, "", tt.msg, tt.args...)

			if tt.wantError == "" && err != nil {
				t.Errorf("no se esperaba error, se obtuvo: %v", err)
			}
			if tt.wantError != "" && (err == nil || err.Error() != tt.wantError) {
				t.Errorf("error = %v, se esperaba: %s", err, tt.wantError)
			}
		})
	}
}
