# 🛡️ Filingo Utility: Identity Context (utilx)
Este paquete es el Componente Core de Infraestructura de Seguridad del ecosistema Filingo. Su única responsabilidad es actuar como el Puente de Identidad entre las solicitudes que entran por la red (HTTP/gRPC) y el corazón de nuestra lógica de negocio (Dominio/Servicios).

Garantiza que la información del usuario (actor_id, tenant_id, roles, etc.) viaje de forma inmutable, segura y fuertemente tipada a lo largo de toda la cadena de ejecución.

## 🧠 ¿Por qué existe este paquete? (Para Analistas Jóvenes)
En arquitecturas de microservicios, el olvido o mal manejo de la identidad causa fallas graves de seguridad. Este paquete soluciona tres problemas críticos:

Evita Colisiones globales: No usamos cadenas de texto comunes (como "actor_id") para guardar datos en el contexto de Go. Usamos un tipo privado (contextKey) que nadie fuera de este paquete puede alterar o suplantar.

Elimina la "Obsesión por los Primitivos": En la base de datos usamos UUID para los identificadores. Los tokens de red a veces envían textos (string). Este paquete centraliza la conversión automática; si el dato es un texto válido, lo transforma en un objeto uuid.UUID real de Go.

Previene Caídas de Sistema (Panics): Intentar extraer un dato del contexto haciendo un casteo manual directo (ej. ctx.Value("roles").([]string)) tumbará el servidor en producción si el dato no existe. Las funciones de este paquete controlan ese riesgo devolviendo un error controlado de Go (error).

## 🧭 Estructura del Pasaporte de Identidad
Cada vez que una petición es validada, el contexto se "hidrata" con las siguientes constantes universales:

Identidad del Actor: ActorIDKey, ActorNameKey, ActorTypeKey.

Seguridad y Permisos: RolesKey (utilizado por el motor de Casbin).

Aislamiento Organizacional (Multi-tenancy): TenantIDKey, TenantCodeKey (utilizado por el Registry dinámico de bases de datos).

Entorno de Ejecución: EnvironmentKey (development, staging, production).

## 💻 Guía de Uso Rápido (¿Cómo lo programo?)
1. En la capa de Entrada (Middlewares / Handlers)
Si estás construyendo un middleware de autenticación, debes inyectar los datos del token en el contexto estándar de Go usando las llaves del paquete:

Go


// Ejemplo de inyección en un middleware
ctx := c.Request.Context()
ctx = context.WithValue(ctx, utilx.ActorIDKey, claims.Sub) // claims.Sub es el UUID del usuario
c.Request = c.Request.WithContext(ctx)
2. En la capa de Negocio (Servicios / Use Cases / Repositorios)
REGLA DE ORO: Nunca accedas al contexto manualmente usando ctx.Value(). Llama siempre a los getters de este paquete.

Go


package services

import (
	"context"
	"fmt"
	"github.com/cafpleon/filingo-util-utilx" // Nuestro paquete
)

func (s *CityService) Activate(ctx context.Context, cityID uuid.UUID) error {
	// 1. Extraer la identidad de forma segura y tipada
	actorID, err := utilx.GetActorID(ctx)
	if err != nil {
		return fmt.Errorf("operación rechazada, identidad inválida: %w", err)
	}

	tenantID, err := utilx.GetTenantID(ctx)
	if err != nil {
		return fmt.Errorf("operación rechazada, se requiere un Tenant válido: %w", err)
	}

	// 2. Ejecutar la lógica sabiendo que actorID y tenantID son uuid.UUID reales
	return s.repo.UpdateStatus(ctx, tenantID, cityID, actorID)
}
## 🧪 Pruebas Unitarias y Estabilidad
Este paquete cuenta con cobertura de pruebas completa en identity_context_test.go.

Antes de realizar cualquier despliegue o cambio en la infraestructura de seguridad, debes garantizar que la suite de pruebas local corra con éxito:

Bash


go test -v ./...
Si agregas una nueva propiedad global al token de seguridad en el futuro (por ejemplo, el idioma del usuario language o su zona horaria timezone), debes registrar su constante aquí, escribir su respectivo método GetXx y añadir su prueba unitaria antes de pasar a revisión por el Analista Jefe.

