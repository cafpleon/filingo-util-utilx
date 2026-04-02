package utilx

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// Definimos un tipo privado para la clave del contexto.
// Esto evita colisiones con otras librerías que usen strings simples.
type contextKey string

const ActorIDKey contextKey = "actor_id"

// GetActorID extrae el UUID del actor del contexto de forma segura.
func GetActorID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(ActorIDKey)
	if val == nil {
		return uuid.Nil, fmt.Errorf("actor_id no encontrado en el contexto")
	}

	// Intentamos el casteo. El middleware debe asegurarse de guardarlo como uuid.UUID
	actorID, ok := val.(uuid.UUID)
	if !ok {
		// Opcional: Si el middleware lo guarda como string, intentamos parsearlo
		if strVal, ok := val.(string); ok {
			return uuid.FromString(strVal)
		}
		return uuid.Nil, fmt.Errorf("el valor de actor_id en el contexto no es un UUID válido")
	}

	return actorID, nil
}
