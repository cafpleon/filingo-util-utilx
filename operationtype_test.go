package utilx

import (
	"encoding/json"
	"testing"
)

func TestOperationType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    OperationType
		wantErr bool
	}{
		{"CREATE", "CREATE", Create, false},
		{"UPDATE", "UPDATE", Update, false},
		{"DELETE", "DELETE", Delete, false},
		{"case insensitive", "create", Create, false},
		{"trim spaces", "  UPDATE  ", Update, false},
		{"invalid", "INVALID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationType_JSON(t *testing.T) {
	type Request struct {
		Op OperationType `json:"op"`
	}

	// Test Marshal
	req := Request{Op: Update}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"op":"UPDATE"}` {
		t.Errorf("Marshal() = %s, want {\"op\":\"UPDATE\"}", string(data))
	}

	// Test Unmarshal
	var req2 Request
	if err := json.Unmarshal([]byte(`{"op":"DELETE"}`), &req2); err != nil {
		t.Fatal(err)
	}
	if req2.Op != Delete {
		t.Errorf("Unmarshal() = %v, want DELETE", req2.Op)
	}
}
