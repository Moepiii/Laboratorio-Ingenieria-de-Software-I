/**
 * cypress/e2e/auth/login.cy.js
 * Prueba E2E para el flujo de autenticación (Login y Registro).
 */
describe('Prueba de Autenticación (Login)', () => {

  // --- PRUEBA 1: INICIO DE SESIÓN EXITOSO ---
  it('debe iniciar sesión exitosamente (Happy Path)', () => {
    cy.intercept('POST', '/api/auth/login', {
      statusCode: 200,
      fixture: 'auth-success.json'
    }).as('loginRequest');
    cy.visit('/');
    cy.get('input[id="username"]').type('admin_valido');
    cy.get('input[id="password"]').type('contraseña_valida_123');
    cy.contains('button', 'Iniciar Sesión').click();
    cy.wait('@loginRequest');
    cy.url().should('include', '/admin/proyectos');
    cy.get('input[id="username"]').should('not.exist');
  });


  // --- PRUEBA 2: CREDENCIALES INCORRECTAS ---
  it('debe mostrar un error si las credenciales son incorrectas', () => {
    cy.intercept('POST', '/api/auth/login', {
      statusCode: 401,
      fixture: 'auth-failure-401.json'
    }).as('loginRequestFailed');
    cy.visit('/');
    cy.get('input[id="username"]').type('usuario_invalido');
    cy.get('input[id="password"]').type('contraseña_incorrecta');
    cy.contains('button', 'Iniciar Sesión').click();
    cy.wait('@loginRequestFailed');
    cy.contains('Credenciales inválidas').should('be.visible');
    cy.url().should('not.include', '/admin');
  });

  
  // --- PRUEBA 3: CAMBIO DE MODO (LOGIN <-> REGISTRO) ---
  it('debe cambiar entre el modo Login y Registro', () => {
    cy.visit('/');
    cy.contains('h2', 'Iniciar Sesión').should('be.visible');
    cy.get('input[id="confirmPassword"]').should('not.exist');
    cy.contains('p', '¿No tienes cuenta? Regístrate').click();
    cy.contains('h2', 'Crear Cuenta').should('be.visible');
    cy.get('input[id="confirmPassword"]').should('be.visible');
    cy.contains('p', '¿Ya tienes cuenta? Inicia Sesión').click();
    cy.contains('h2', 'Iniciar Sesión').should('be.visible');
    cy.get('input[id="confirmPassword"]').should('not.exist');
  });


  // --- PRUEBA 4: REGISTRO DE USUARIO EXITOSO ---
  it('debe registrar un nuevo usuario exitosamente (Happy Path)', () => {
    cy.intercept('POST', '/api/auth/register', {
      statusCode: 201,
      fixture: 'auth-register-success.json'
    }).as('registerRequest');
    cy.visit('/');
    cy.contains('p', '¿No tienes cuenta? Regístrate').click();
    cy.get('input[id="username"]').type('usuario_nuevo_test');
    cy.get('input[id="nombre"]').type('Test');
    cy.get('input[id="apellido"]').type('Usuario');
    cy.get('input[id="cedula"]').type('V-123456');
    cy.get('input[id="password"]').type('pass12345');
    cy.get('input[id="confirmPassword"]').type('pass12345');
    cy.contains('button', 'Registrarse').click();
    cy.wait('@registerRequest');
    cy.contains('¡Registro exitoso! Ahora puedes iniciar sesión.').should('be.visible');
    cy.contains('h2', 'Iniciar Sesión').should('be.visible');
  });


  // --- (NUEVA) PRUEBA 5: VALIDACIÓN DE CONTRASEÑAS (REGISTRO) ---
  
  it('debe mostrar un error si las contraseñas no coinciden', () => {
    
    // 1. Interceptamos la llamada a la API y le damos un alias.
    // Usamos 'cy.spy' para verificar que la llamada NUNCA ocurra.
    cy.intercept('POST', '/api/auth/register').as('registerRequest');

    cy.visit('/');

    // 2. Cambiar a modo registro
    cy.contains('p', '¿No tienes cuenta? Regístrate').click();

    // 3. Llenar el formulario con contraseñas que NO coinciden
    cy.get('input[id="username"]').type('usuario_test');
    cy.get('input[id="nombre"]').type('Test');
    cy.get('input[id="apellido"]').type('Usuario');
    cy.get('input[id="cedula"]').type('V-1234567');
    cy.get('input[id="password"]').type('pass123'); // Contraseña 1
    cy.get('input[id="confirmPassword"]').type('pass456'); // Contraseña 2 (diferente)

    // 4. Enviar formulario
    cy.contains('button', 'Registrarse').click();

    // 5. Verificar los resultados
    
    // ⚠️ ASUNCIÓN: Tu AuthForm.js (o AuthContext) 
    // establece el error como "Las contraseñas no coinciden".
    cy.contains('Las contraseñas no coinciden').should('be.visible');

    // 6. (Verificación clave) Asegurarnos de que la llamada a la API NUNCA se hizo.
    // Esto comprueba que la validación fue del lado del cliente.
    cy.get('@registerRequest.all').should('have.length', 0);
  });

});