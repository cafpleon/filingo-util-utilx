package utilx

import (
	"context"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetActorID(t *testing.T) {
	// Generamos UUIDs de prueba válidos usando v5 (o v7 en producción)
	validUUID, _ := uuid.NewV4()

	t.Run("Éxito con tipo uuid.UUID nativo", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ActorIDKey, validUUID)

		result, err := GetActorID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, validUUID, result)
	})

	t.Run("Éxito con tipo string (Fallback Parsing)", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ActorIDKey, validUUID.String())

		result, err := GetActorID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, validUUID, result)
	})

	t.Run("Error cuando la clave no existe en el contexto", func(t *testing.T) {
		ctx := context.Background()

		result, err := GetActorID(ctx)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, result)
		assert.Contains(t, err.Error(), "actor_id no encontrado")
	})

	t.Run("Error cuando el valor es un string pero no es un UUID válido", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ActorIDKey, "un-string-cualquiera")

		result, err := GetActorID(ctx)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, result)
		assert.Contains(t, err.Error(), "no es un UUID válido")
	})
}

func TestGetTenantID(t *testing.T) {
	validTenantUUID, _ := uuid.NewV4()

	t.Run("Éxito con tipo uuid.UUID nativo", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, validTenantUUID)

		result, err := GetTenantID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, validTenantUUID, result)
	})

	t.Run("Éxito con tipo string (Fallback Parsing)", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, validTenantUUID.String())

		result, err := GetTenantID(ctx)
		assert.NoError(t, err)
		assert.Equal(t, validTenantUUID, result)
	})

	t.Run("Error cuando no se encuentra el tenant_id", func(t *testing.T) {
		ctx := context.Background()

		result, err := GetTenantID(ctx)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, result)
	})
}

func TestGetActorNameAndTenantCode(t *testing.T) {
	t.Run("Extraer Actor Name exitosamente", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ActorNameKey, "Carlos Fonseca")
		name, err := GetActorName(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "Carlos Fonseca", name)
	})

	t.Run("Extraer Tenant Code exitosamente", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantCodeKey, "ACME_LOGISTICA")
		code, err := GetTenantCode(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "ACME_LOGISTICA", code)
	})
}

func TestGetRoles(t *testing.T) {
	expectedRoles := []string{"ROLE_SUPER_ADMIN", "ROLE_ANALYST"}

	t.Run("Éxito al obtener slice de roles", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), RolesKey, expectedRoles)
		roles, err := GetRoles(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
	})

	t.Run("Error cuando los roles no existen", func(t *testing.T) {
		ctx := context.Background()
		_, err := GetRoles(ctx)
		assert.Error(t, err)
	})
}

func TestGetEnvironment(t *testing.T) {
	t.Run("Devuelve el entorno configurado", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), EnvironmentKey, "staging")
		env, err := GetEnvironment(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "staging", env)
	})

	t.Run("Devuelve entorno de producción por defecto si está vacío", func(t *testing.T) {
		ctx := context.Background()
		env, err := GetEnvironment(ctx)
		assert.NoError(t, err) // No debe dar error
		assert.Equal(t, "production", env)
	})
}
