/**
 * Pruebas E2E para los flujos de Login.
 * Verifica redirecciones y contenido inicial de dashboards.
 */
describe('Flujo de Autenticación - Login', () => {

  // Antes de cada prueba ('it' block), limpia el estado y visita la página de login
  beforeEach(() => {
    cy.clearLocalStorage(); // Asegura que no haya sesión previa
    cy.visit('/login');    // Ve a la página de login
    cy.url().should('include', '/login'); // Confirma que estamos en /login
  });

  // --- Prueba 1: Login Admin/Gerente Exitoso ---
  it('debería loguear a un admin/gerente y redirigir a /admin/proyectos', () => {
    // 1. Simula la respuesta EXITOSA del API de login para un admin/gerente
    cy.intercept('POST', '**/api/login', {
      statusCode: 200,
      body: { // Debe coincidir con lo que AuthContext espera
        token: 'token-admin-simulado',
        user: { username: 'adminTest', nombre: 'Admin', apellido: 'Prueba' },
        userId: 1,
        role: 'admin' // O 'gerente'
      }
    }).as('loginAdmin'); // Alias para esperar esta llamada

    // 2. Simula la respuesta del API que carga los proyectos INICIALMENTE (vacía)
    cy.intercept('POST', '**/api/admin/get-proyectos', {
      statusCode: 200,
      body: { proyectos: [] } // Empieza sin proyectos
    }).as('getAdminProjects');

    // 3. Interactúa con el formulario de login
    cy.get('#username').should('be.visible').type('adminTest');
    cy.get('#password').should('be.visible').type('password123');
    cy.contains('button', 'Iniciar Sesión').click();

    // 4. Espera a que las llamadas API terminen
    cy.wait('@loginAdmin'); // Espera el login
    cy.wait('@getAdminProjects'); // Espera la carga inicial de proyectos

    // 5. Verifica la redirección y el contenido del dashboard admin
    cy.url().should('include', '/admin/proyectos'); // Verifica URL final
    cy.contains('Portafolio de Proyectos').should('be.visible'); // Verifica título
    cy.contains('No hay proyectos creados.').should('be.visible'); // Verifica tabla vacía
  });

  // --- Prueba 2: Login Usuario Normal Exitoso (Sin Proyecto) ---
  it('debería loguear a un usuario normal y redirigir a /dashboard (mostrando "sin proyecto")', () => {
    // 1. Simula respuesta EXITOSA del API de login para un usuario 'user'
    cy.intercept('POST', '**/api/login', {
      statusCode: 200,
      body: {
        token: 'token-user-simulado',
        user: { username: 'userTest', nombre: 'Usuario', apellido: 'Prueba' },
        userId: 15,
        role: 'user'
      }
    }).as('loginUser');

    // 2. Simula respuesta del API de detalles (404 manejado -> sin proyecto)
    cy.intercept('POST', '**/api/user/project-details', {
      statusCode: 404, // Simula el caso "no encontrado"
      body: null      // Respuesta vacía
    }).as('getUserDashboard');

    // 3. Interactúa con el formulario
    cy.get('#username').should('be.visible').type('userTest');
    cy.get('#password').should('be.visible').type('password123');
    cy.contains('button', 'Iniciar Sesión').click();

    // 4. Espera llamadas API
    cy.wait('@loginUser');
    cy.wait('@getUserDashboard');

    // 5. Verifica redirección y contenido del dashboard de usuario
    cy.url().should('include', '/dashboard'); // Verifica URL final
    // Verifica elementos clave del estado "sin proyecto"
    cy.get('h1').should('contain', 'Bienvenido, Usuario'); // Verifica saludo
    cy.contains('p', 'Actualmente no estás asignado a ningún proyecto.').should('be.visible');
    cy.contains('button', 'Cerrar Sesión').should('be.visible'); // Verifica botón logout
  });

  // --- Prueba 3: Login Fallido ---
  it('debería mostrar un error si las credenciales son incorrectas', () => {
    // 1. Simula respuesta de ERROR 401 del API de login
    cy.intercept('POST', '**/api/login', {
      statusCode: 401,
      body: { error: 'Credenciales inválidas' } // El mensaje exacto de tu backend
    }).as('loginFail');

    // 2. Interactúa con el formulario (datos incorrectos)
    cy.get('#username').should('be.visible').type('usuarioErroneo');
    cy.get('#password').should('be.visible').type('passErronea');
    cy.contains('button', 'Iniciar Sesión').click();

    // 3. Espera la llamada fallida
    cy.wait('@loginFail');

    // 4. Verifica que NO redirige y muestra el error
    cy.url().should('include', '/login'); // Sigue en la misma página
    // Busca el texto de error (ajusta si se muestra diferente)
    cy.contains('Credenciales inválidas').should('be.visible');
  });

});