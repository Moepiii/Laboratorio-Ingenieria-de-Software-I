# Sistema de Gestión de Proyectos Agrícolas

Sistema web full-stack para la gestión integral de proyectos agrícolas, incluyendo administración de usuarios, proyectos, labores agronómicas, equipos, actividades, planes de acción, recursos humanos, materiales e insumos.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Tecnologías](#tecnologías)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Requisitos Previos](#requisitos-previos)
- [Instalación](#instalación)
- [Configuración](#configuración)
- [Uso](#uso)
- [Estructura de la Base de Datos](#estructura-de-la-base-de-datos)
- [API Endpoints](#api-endpoints)
- [Roles y Permisos](#roles-y-permisos)
- [Testing](#testing)
  - [Pruebas Unitarias e Integración en Go](#pruebas-unitarias-e-integración-en-go)
  - [Pruebas E2E con Cypress](#pruebas-e2e-con-cypress)
- [Desarrollo](#desarrollo)

## ✨ Características

- **Autenticación y Autorización**: Sistema de login con JWT y roles (admin, gerente, user)
- **Gestión de Proyectos**: Creación, edición y administración de proyectos agrícolas
- **Labores Agronómicas**: Registro y seguimiento de labores por proyecto
- **Equipos e Implementos**: Control de inventario de equipos e implementos
- **Actividades**: Registro detallado de actividades con asignación de responsables
- **Planes de Acción**: Planificación y seguimiento de acciones por proyecto
- **Recursos Humanos**: Gestión de personal asignado a proyectos
- **Materiales e Insumos**: Control de inventario de materiales
- **Unidades de Medida**: Configuración de unidades de medida personalizadas
- **Sistema de Auditoría**: Logger de eventos para seguimiento de acciones
- **Dashboard Diferenciado**: Interfaces distintas para administradores y usuarios regulares

## 🛠 Tecnologías

### Backend
- **Go 1.25.1**: Lenguaje de programación
- **SQLite**: Base de datos embebida
- **Gorilla Handlers**: Middleware para CORS
- **JWT**: Autenticación basada en tokens
- **bcrypt**: Encriptación de contraseñas

### Frontend
- **React 19.2.0**: Biblioteca de UI
- **React Router DOM 7.9.4**: Enrutamiento
- **Lucide React**: Iconos
- **Cypress**: Testing end-to-end

## 📁 Estructura del Proyecto

```
proyecto/
├── backend/
│   ├── internal/
│   │   ├── actividades/      # Servicio de actividades
│   │   ├── auth/             # Autenticación y autorización
│   │   ├── database/         # Configuración y queries de BD
│   │   ├── equipos/          # Servicio de equipos
│   │   ├── handlers/         # Controladores HTTP
│   │   ├── labores/          # Servicio de labores
│   │   ├── logger/           # Servicio de auditoría
│   │   ├── models/           # Modelos de datos
│   │   ├── proyectos/        # Servicio de proyectos
│   │   ├── unidades/         # Servicio de unidades
│   │   └── users/            # Servicio de usuarios
│   ├── main.go               # Punto de entrada del servidor
│   ├── main_test.go          # Tests del servidor
│   ├── go.mod                # Dependencias de Go
│   └── users.db              # Base de datos SQLite
│
└── frontend/
    ├── src/
    │   ├── components/       # Componentes reutilizables
    │   │   └── auth/         # Componentes de autenticación
    │   ├── context/          # Context API (AuthContext)
    │   ├── pages/            # Páginas principales
    │   ├── services/         # Servicios de API
    │   ├── App.js            # Componente principal
    │   └── index.js          # Punto de entrada
    ├── cypress/              # Tests E2E
    ├── public/               # Archivos estáticos
    └── package.json          # Dependencias de Node.js
```

## 📦 Requisitos Previos

- **Go**: Versión 1.25.1 o superior
- **Node.js**: Versión 14 o superior
- **npm**: Versión 6 o superior
- **Git**: Para clonar el repositorio

## 🚀 Instalación

### 1. Clonar el repositorio

```bash
git clone <url-del-repositorio>
cd proyecto
```

### 2. Configurar el Backend

```bash
cd backend
go mod download
```

### 3. Configurar el Frontend

```bash
cd frontend
npm install
```

## ⚙️ Configuración

### Backend

El servidor backend se ejecuta por defecto en el puerto `8080`. La base de datos SQLite se crea automáticamente en `backend/users.db` al iniciar el servidor.

**Usuario por defecto:**
- Username: `admin`
- Password: `admin123`
- Rol: `admin`

### Frontend

El frontend se ejecuta por defecto en el puerto `3000`. Asegúrate de que el backend esté corriendo antes de iniciar el frontend.

## 🎯 Uso

### Iniciar el Backend

```bash
cd backend
go run main.go
```

El servidor estará disponible en `http://localhost:8080`

### Iniciar el Frontend

```bash
cd frontend
npm start
```

La aplicación estará disponible en `http://localhost:3000`

### Acceso a la Aplicación

1. Abre tu navegador en `http://localhost:3000`
2. Inicia sesión con las credenciales de administrador:
   - Username: `admin`
   - Password: `admin123`

## 🗄️ Estructura de la Base de Datos

El sistema utiliza las siguientes tablas principales:

- **users**: Usuarios del sistema
- **proyectos**: Proyectos agrícolas
- **labores_agronomicas**: Labores agronómicas por proyecto
- **equipos_implementos**: Equipos e implementos
- **actividades**: Actividades del proyecto
- **unidades_medida**: Unidades de medida personalizadas
- **planes_accion**: Planes de acción
- **recursos_humanos**: Recursos humanos asignados
- **materiales_insumos**: Materiales e insumos
- **event_logs**: Logs de auditoría

## 🔌 API Endpoints

### Autenticación
- `POST /api/auth/register` - Registro de usuarios
- `POST /api/auth/login` - Inicio de sesión

### Usuarios (Admin)
- `GET /api/admin/users` - Listar usuarios
- `POST /api/admin/add-user` - Crear usuario
- `POST /api/admin/delete-user` - Eliminar usuario
- `POST /api/admin/update-user` - Actualizar rol de usuario
- `POST /api/admin/assign-project` - Asignar proyecto a usuario

### Proyectos (Admin)
- `GET /api/admin/get-proyectos` - Listar proyectos
- `POST /api/admin/create-proyecto` - Crear proyecto
- `POST /api/admin/update-proyecto` - Actualizar proyecto
- `POST /api/admin/delete-proyecto` - Eliminar proyecto
- `POST /api/admin/set-proyecto-estado` - Cambiar estado del proyecto

### Labores Agronómicas (Admin)
- `GET /api/admin/get-labores` - Listar labores
- `POST /api/admin/create-labor` - Crear labor
- `POST /api/admin/update-labor` - Actualizar labor
- `POST /api/admin/delete-labor` - Eliminar labor

### Equipos e Implementos (Admin)
- `GET /api/admin/get-equipos` - Listar equipos
- `POST /api/admin/create-equipo` - Crear equipo
- `POST /api/admin/update-equipo` - Actualizar equipo
- `POST /api/admin/delete-equipo` - Eliminar equipo

### Unidades de Medida (Admin)
- `GET /api/admin/get-unidades` - Listar unidades
- `POST /api/admin/create-unidad` - Crear unidad
- `POST /api/admin/update-unidad` - Actualizar unidad
- `POST /api/admin/delete-unidad` - Eliminar unidad

### Actividades (Admin)
- `GET /api/admin/get-datos-proyecto` - Obtener datos del proyecto
- `POST /api/admin/create-actividad` - Crear actividad
- `POST /api/admin/update-actividad` - Actualizar actividad
- `POST /api/admin/delete-actividad` - Eliminar actividad

### Planes de Acción (Admin)
- `GET /api/admin/get-planes` - Listar planes
- `POST /api/admin/create-plan` - Crear plan
- `POST /api/admin/update-plan` - Actualizar plan
- `POST /api/admin/delete-plan` - Eliminar plan

### Recursos Humanos (Admin)
- `GET /api/admin/get-recursos` - Listar recursos
- `POST /api/admin/create-recurso` - Crear recurso
- `POST /api/admin/update-recurso` - Actualizar recurso
- `POST /api/admin/delete-recurso` - Eliminar recurso

### Materiales e Insumos (Admin)
- `GET /api/admin/get-materiales` - Listar materiales
- `POST /api/admin/create-material` - Crear material
- `POST /api/admin/update-material` - Actualizar material
- `POST /api/admin/delete-material` - Eliminar material

### Logger/Auditoría (Admin)
- `GET /api/admin/get-logs` - Obtener logs
- `POST /api/admin/delete-logs` - Eliminar logs
- `POST /api/admin/delete-logs-range` - Eliminar logs por rango

### Usuario Regular
- `GET /api/user/project-details` - Detalles del proyecto asignado

## 👥 Roles y Permisos

### Admin
- Acceso completo a todas las funcionalidades
- Gestión de usuarios y proyectos
- Acceso al sistema de auditoría/logs

### Gerente
- Gestión de proyectos asignados
- Acceso a configuraciones y planes de acción
- Sin acceso al sistema de logs

### User
- Acceso de solo lectura a su proyecto asignado
- Visualización de detalles del proyecto
- Sin permisos de edición

## 🧪 Testing

El proyecto incluye dos tipos de pruebas: pruebas unitarias/integración en Go para el backend y pruebas end-to-end (E2E) con Cypress para el frontend.

### Pruebas Unitarias e Integración en Go

Las pruebas del backend están ubicadas en `backend/main_test.go` y cubren el flujo completo de la aplicación, incluyendo pruebas de integración y seguridad.

#### Estructura de las Pruebas

Las pruebas en Go utilizan el paquete estándar `testing` y siguen un enfoque de integración que prueba el flujo completo de la aplicación:

1. **Setup y Teardown**: Cada ejecución de pruebas crea una base de datos temporal (`test_integration.db`) que se elimina al finalizar
2. **Flujo Completo**: Las pruebas verifican el flujo completo desde registro hasta operaciones CRUD
3. **Pruebas de Seguridad**: Incluyen validación de permisos y acceso no autorizado

#### Casos de Prueba Implementados

- ✅ **Registro de Usuario**: Verifica la creación de nuevos usuarios
- ✅ **Autenticación**: Prueba el login y obtención de tokens JWT
- ✅ **Gestión de Proyectos**: Creación de proyectos agrícolas
- ✅ **Unidades de Medida**: Creación y gestión de unidades
- ✅ **Equipos e Implementos**: CRUD de equipos
- ✅ **Labores Agronómicas**: Gestión de labores por proyecto
- ✅ **Materiales e Insumos**: Registro de materiales
- ✅ **Seguridad**: Validación de acceso no autorizado (usuarios sin permisos no pueden acceder a rutas protegidas)

#### Ejecutar las Pruebas del Backend

```bash
cd backend
go test ./...
```

Para ejecutar con más detalles:

```bash
go test -v ./...
```

Para ejecutar un test específico:

```bash
go test -v -run TestFlujoCompleto
```

#### Ejemplo de Prueba

```go
// Las pruebas verifican el flujo completo:
// 1. Registro → 2. Login → 3. Crear Proyecto → 4. Operaciones CRUD
// 5. Validación de seguridad (acceso no autorizado)
```

### Pruebas E2E con Cypress

Cypress se utiliza para realizar pruebas end-to-end que simulan el comportamiento real del usuario en la aplicación web.

#### Configuración

El archivo `frontend/cypress.config.js` configura Cypress con:
- **Base URL**: `http://localhost:3000`
- **Viewport**: 1280x720
- **Patrón de specs**: `cypress/e2e/**/*.cy.js`

#### Estructura de Pruebas Cypress

Las pruebas están organizadas en las siguientes categorías:

```
cypress/
├── e2e/
│   ├── auth/                    # Pruebas de autenticación
│   │   └── ALogindeadmin.cy.js  # Login de administrador
│   ├── admin/                   # Pruebas del dashboard de admin
│   │   ├── BPortafolio.cy.js    # Gestión de portafolio
│   │   ├── CUsuarios.cy.js      # Gestión de usuarios
│   │   ├── DConfiguraciones.cy.js # Configuraciones
│   │   ├── logger.cy.js         # Sistema de logs
│   │   ├── materiales.cy.js      # Materiales e insumos
│   │   ├── plan_accion.cy.js    # Planes de acción
│   │   └── RecursoHumano.cy.js  # Recursos humanos
│   └── 2-advanced-examples/     # Ejemplos avanzados de Cypress
├── fixtures/                    # Datos de prueba (JSON)
├── support/
│   └── commands.js              # Comandos personalizados
└── screenshots/                 # Capturas de pantalla de pruebas
```

#### Casos de Prueba E2E Implementados

**Autenticación (Historia A)**
- ✅ Login exitoso de administrador
- ✅ Validación de credenciales incorrectas
- ✅ Cierre de sesión

**Gestión de Usuarios (Historia C)**
- ✅ Crear nuevo usuario
- ✅ Editar usuario existente
- ✅ Eliminar usuario
- ✅ Asignar proyecto a usuario

**Portafolio de Proyectos (Historia B)**
- ✅ Crear proyecto
- ✅ Editar proyecto
- ✅ Eliminar proyecto
- ✅ Cambiar estado de proyecto

**Configuraciones (Historia D)**
- ✅ Gestión de labores agronómicas
- ✅ Gestión de equipos e implementos
- ✅ Gestión de unidades de medida

**Otros Módulos**
- ✅ Planes de acción
- ✅ Recursos humanos
- ✅ Materiales e insumos
- ✅ Sistema de logger/auditoría

#### Comandos Personalizados

Cypress incluye comandos personalizados en `cypress/support/commands.js`:

- `cy.loginAsAdmin()`: Simula el login de un administrador para pruebas rápidas

#### Ejecutar las Pruebas E2E

**Modo Interactivo (Recomendado para desarrollo):**

```bash
cd frontend
npx cypress open
```

Esto abre la interfaz gráfica de Cypress donde puedes seleccionar qué pruebas ejecutar.

**Modo Headless (Para CI/CD):**

```bash
cd frontend
npx cypress run
```

**Ejecutar una suite específica:**

```bash
npx cypress run --spec "cypress/e2e/auth/ALogindeadmin.cy.js"
```

**Ejecutar todas las pruebas de admin:**

```bash
npx cypress run --spec "cypress/e2e/admin/**/*.cy.js"
```

#### Requisitos para Ejecutar Pruebas E2E

1. **Backend corriendo**: El servidor debe estar ejecutándose en `http://localhost:8080`
2. **Frontend corriendo**: La aplicación React debe estar en `http://localhost:3000`
3. **Base de datos**: Asegúrate de que la base de datos tenga el usuario `admin` con contraseña `admin123`

#### Fixtures (Datos de Prueba)

Cypress utiliza archivos JSON en `cypress/fixtures/` para datos de prueba:
- `auth-success.json`: Respuesta exitosa de autenticación
- `auth-failure-401.json`: Respuesta de error de autenticación
- `admin-projects.json`: Datos de proyectos
- `admin-users.json`: Datos de usuarios

### Cobertura de Pruebas

#### Backend (Go)
- ✅ Flujo completo de registro y autenticación
- ✅ Operaciones CRUD de todas las entidades principales
- ✅ Validación de seguridad y permisos
- ✅ Manejo de errores

#### Frontend (Cypress)
- ✅ Flujos de usuario completos
- ✅ Interfaz de administrador
- ✅ Validación de formularios
- ✅ Navegación entre páginas
- ✅ Interacciones con la UI

### Mejores Prácticas

1. **Ejecutar pruebas antes de commit**: Siempre ejecuta las pruebas antes de hacer commit
2. **Pruebas aisladas**: Cada prueba debe ser independiente y poder ejecutarse sola
3. **Datos de prueba**: Usa fixtures y datos aleatorios para evitar conflictos
4. **Limpieza**: Las pruebas deben limpiar después de ejecutarse (las pruebas de Go lo hacen automáticamente)
5. **Nombres descriptivos**: Usa nombres claros que describan qué prueba cada test

### Troubleshooting

**Problema**: Las pruebas de Cypress fallan con errores de conexión
- **Solución**: Asegúrate de que tanto el backend como el frontend estén corriendo

**Problema**: Las pruebas de Go fallan por base de datos bloqueada
- **Solución**: Cierra cualquier conexión activa a la base de datos de prueba

**Problema**: Las pruebas E2E son inconsistentes
- **Solución**: Aumenta los timeouts en Cypress o verifica que la aplicación responda correctamente

## 💻 Desarrollo

### Estructura de Código

El proyecto sigue una arquitectura en capas:

1. **Handlers**: Manejan las peticiones HTTP
2. **Services**: Contienen la lógica de negocio
3. **Database**: Maneja las consultas a la base de datos
4. **Models**: Define las estructuras de datos

### Convenciones

- Los handlers validan la autenticación y autorización
- Los servicios contienen la lógica de negocio
- Las queries de base de datos están separadas en archivos específicos
- El logger registra todas las acciones administrativas

### CORS

El backend está configurado para aceptar peticiones desde `http://localhost:3000`. Para producción, actualiza la configuración CORS en `backend/main.go`.

## 📝 Notas Adicionales

- La base de datos se crea automáticamente al iniciar el servidor
- El usuario administrador se crea automáticamente si no existe
- Todos los endpoints administrativos requieren autenticación JWT
- El sistema utiliza Write-Ahead Logging (WAL) para SQLite

## 📄 Licencia

Este proyecto es parte de un trabajo académico para la materia CI3715.

## 👤 Autor

Este proyecto es parte de un trabajo académico para la materia **CI3715**, desarrollado como parte del **Laboratorio de CI3715**.

| Nombre               | Correo electrónico         | Rol              |
|----------------------|----------------------------|------------------|
| Jean Carlos Guzmán   | jguzman106@gmail.com       | Agile Coach      |
| David Pereira        | 18-10245@usb.ve            | Miembro del Equipo |
| Rafael Valera        | 16-11202@usb.ve            | Miembro del Equipo |

📬 Para más información o soporte, contacta al equipo de desarrollo.

