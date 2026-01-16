package utilx_test

import (
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
