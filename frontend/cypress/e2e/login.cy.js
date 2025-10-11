describe('Login Flow - React + Go', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000'); 
  });

  it('✅ Login exitoso con intercept y token', () => {
    cy.intercept('POST', '/api/login', {
      statusCode: 200,
      body: { token: 'fake-jwt-token' },
    }).as('loginRequest');

    cy.get('input#username').type('David');
    cy.get('input#password').type('123456');
    cy.get('button[type="submit"]').click();

    cy.wait('@loginRequest');
    cy.contains('Bienvenido').should('exist');
  });

  it('❌ Login fallido con credenciales incorrectas', () => {
    cy.intercept('POST', '/api/login', {
      statusCode: 401,
      body: { error: 'Credenciales inválidas' },
    }).as('loginFail');

    cy.get('input#username').type('usuario@ejemplo.com');
    cy.get('input#password').type('claveIncorrecta');
    cy.get('button[type="submit"]').click();

    cy.wait('@loginFail');
    cy.contains('Usuario o contraseña incorrectos.').should('exist');
  });

  // --- PRUEBA PARA ADMINISTRADOR (Historia 1.b) ---
  it('✅ Admin: Creación de un nuevo perfil de usuario', () => {
    // 1. Mockear el login como administrador (importante: role: 'admin')
    cy.intercept('POST', '/api/login', {
      statusCode: 200,
      body: { 
        token: 'fake-admin-token',
        usuario: 'AdminUser',
        role: 'admin'
      },
    }).as('adminLogin');

    // 2. Mockear la llamada para añadir un nuevo usuario
    cy.intercept('POST', '/api/admin/add-user', (req) => {
      // Puedes verificar que el cuerpo de la solicitud sea el correcto
      expect(req.body.username).to.eq('NuevoUsuario');
      expect(req.body.password).to.eq('password123');
      req.reply({
        statusCode: 200,
        body: { mensaje: 'Usuario NuevoUsuario creado exitosamente.' },
      });
    }).as('addUserRequest');

    // 3. Simular el login del administrador (para renderizar AdminDashboard)
    cy.get('input#username').type('admin');
    cy.get('input#password').type('admin123');
    cy.get('button[type="submit"]').click();
    cy.wait('@adminLogin');
    
    // Asersión para verificar que se cargó el Dashboard de Admin
    cy.contains('Panel de Administración').should('exist');
    
    // 4. Ingresar datos del nuevo usuario
    cy.get('input[placeholder="Nombre de Usuario"]').type('NuevoUsuario');
    cy.get('input[placeholder^="Contraseña (mín. 6 caracteres"]').type('password123'); // Selector que busca el inicio del placeholder
    
    // 5. Hacer clic en el botón de Crear Usuario
    cy.contains('Crear Usuario').click();

    // 6. Esperar el mock de la API de creación de usuario
    cy.wait('@addUserRequest');
    
    // 7. Verificar el mensaje de éxito que se muestra en la UI
    cy.contains('Usuario NuevoUsuario creado exitosamente.').should('be.visible');
  });
  // ----------------------------------------------------

  it('🧼 Limpieza de sesión y cookies', () => {
    cy.clearCookies();
    cy.clearLocalStorage();
    cy.visit('http://localhost:3000');
    cy.get('input#username').should('exist');
  });
});