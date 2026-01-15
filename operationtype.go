package utilx

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OperationType define las operaciones CRUD básicas de un sistema
type OperationType int

const (
	Create OperationType = iota // 0
	Update                      // 1
	Delete                      // 2
	// NOTA: Puedes agregar más valores aquí si necesitas (ej: Read, Patch)
)

// String devuelve la representación en string del OperationType
func (ot OperationType) String() string {
	names := [...]string{"CREATE", "UPDATE", "DELETE"}
	if ot < Create || ot > Delete {
		return fmt.Sprintf("OperationType(%d)", int(ot))
	}
	return names[ot]
}

// IsValid verifica si el valor del OperationType es válido
func (ot OperationType) IsValid() bool {
	return ot >= Create && ot <= Delete
}

// Parse convierte un string en OperationType (case-insensitive)
func Parse(str string) (OperationType, error) {
	switch strings.ToUpper(strings.TrimSpace(str)) {
	case "CREATE":
		return Create, nil
	case "UPDATE":
		return Update, nil
	case "DELETE":
		return Delete, nil
	default:
		return -1, fmt.Errorf("operationtype: valor inválido %q", str)
	}
}

// MustParse convierte un string en OperationType, panic si es inválido
func MustParse(str string) OperationType {
	ot, err := Parse(str)
	if err != nil {
		panic(err)
	}
	return ot
}

// FromInt convierte un int en OperationType (útil para bases de datos)
func FromInt(i int) (OperationType, error) {
	ot := OperationType(i)
	if !ot.IsValid() {
		return -1, fmt.Errorf("operationtype: valor int inválido %d", i)
	}
	return ot, nil
}

// Int devuelve el valor int subyacente
func (ot OperationType) Int() int {
	return int(ot)
}

// IsDestructive indica si la operación es destructiva (borra datos)
func (ot OperationType) IsDestructive() bool {
	return ot == Delete
}

// IsMutation indica si la operación modifica datos (no es solo lectura)
func (ot OperationType) IsMutation() bool {
	return ot == Create || ot == Update || ot == Delete
}

// MarshalJSON implementa json.Marshaler para serialización automática
func (ot OperationType) MarshalJSON() ([]byte, error) {
	if !ot.IsValid() {
		return nil, fmt.Errorf("operationtype: no se puede serializar valor inválido %d", ot)
	}
	return json.Marshal(ot.String())
}

// UnmarshalJSON implementa json.Unmarshaler para deserialización automática
func (ot *OperationType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, err := Parse(str)
	if err != nil {
		return fmt.Errorf("operationtype: %w", err)
	}
	*ot = parsed
	return nil
}

// Scan implementa sql.Scanner para leer desde bases de datos
func (ot *OperationType) Scan(value interface{}) error {
	switch v := value.(type) {
	case int64:
		parsed, err := FromInt(int(v))
		if err != nil {
			return err
		}
		*ot = parsed
		return nil
	case string:
		parsed, err := Parse(v)
		if err != nil {
			return err
		}
		*ot = parsed
		return nil
	default:
		return fmt.Errorf("operationtype: tipo no soportado para Scan: %T", value)
	}
}

// Value implementa driver.Valuer para escribir a bases de datos
func (ot OperationType) Value() (interface{}, error) {
	return ot.Int(), nil
}
