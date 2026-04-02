package utilx

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
)

// contextKey Definimos un tipo privado para la clave del contexto.
// Esto evita colisiones con otras librerías que usen strings simples.
// TODO mover a ultx_auth i a filingo_util_auth
type contextKey string

const (
	ActorIDKey   contextKey = "actor_id"
	ActorNameKey contextKey = "actor_name"
	ActorTypeKey contextKey = "actor_type"
	RolesKey     contextKey = "roles"

	// Llaves para Multitenancy
	TenantIDKey   contextKey = "tenant_id"
	TenantCodeKey contextKey = "tenant_code"

	// Llaves para Entorno
	EnvironmentKey contextKey = "environment"
)

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

// --- GETTERS DE IDENTIDAD ---

// GetActorName devuelve el nombre legible del actor (ej: "Juan Perez").
func GetActorName(ctx context.Context) (string, error) {
	val := ctx.Value(ActorNameKey)
	if name, ok := val.(string); ok {
		return name, nil
	}
	return "", fmt.Errorf("actor_name no encontrado o no es string")
}

// GetRoles devuelve el slice de roles del actor para validaciones de RBAC.
func GetRoles(ctx context.Context) ([]string, error) {
	val := ctx.Value(RolesKey)
	if roles, ok := val.([]string); ok {
		return roles, nil
	}
	return nil, fmt.Errorf("roles no encontrados en el contexto")
}

// --- GETTERS DE MULTITENANCY ---

// GetTenantID es CRUCIAL para el Registry de Base de Datos.
func GetTenantID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(TenantIDKey)
	if id, ok := val.(uuid.UUID); ok {
		return id, nil
	}
	// Fallback si el middleware lo inyectó como string
	if s, ok := val.(string); ok {
		return uuid.FromString(s)
	}
	return uuid.Nil, fmt.Errorf("tenant_id no encontrado o inválido")
}

// GetTenantCode devuelve el código semántico (ej: "WELCOME_DEMO").
// Útil para lógica de negocio que dependa del tipo de cliente.
func GetTenantCode(ctx context.Context) (string, error) {
	val := ctx.Value(TenantCodeKey)
	if code, ok := val.(string); ok {
		return code, nil
	}
	return "", fmt.Errorf("tenant_code no encontrado")
}

// --- GETTERS DE ENTORNO ---

// GetEnvironment devuelve "development", "staging" o "production".
func GetEnvironment(ctx context.Context) (string, error) {
	val := ctx.Value(EnvironmentKey)
	if env, ok := val.(string); ok {
		return env, nil
	}
	return "production", nil // Default seguro
}
