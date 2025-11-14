/**
 * cypress/e2e/admin/usuarios.cy.js
 * Pruebas para la página de gestión de perfiles de usuario.
 */
describe('Gestión de Usuarios (Página Admin)', () => {

  beforeEach(() => {
    // --- PASO 1: Iniciar sesión como Admin ---
    cy.loginAsAdmin();

    // --- PASO 2: Interceptar las APIs de CARGA INICIAL ---
    // ⭐️ CORRECCIÓN: Usamos '**' para un "matching" de ruta más robusto
    
    cy.intercept('POST', '**/api/admin/users', {
      statusCode: 200,
      fixture: 'admin-users.json'
    }).as('getUsers');

    cy.intercept('POST', '**/api/admin/get-proyectos', {
      statusCode: 200,
      fixture: 'admin-projects.json' // (Este fixture DEBE tener el formato {"projects": [...]})
    }).as('getProjects');
  });


  // --- PRUEBA 1: Carga de la página ---
  it('debe mostrar la lista de usuarios al cargar (Happy Path)', () => {
    
    cy.visit('/admin/usuarios');
    cy.wait(['@getUsers', '@getProjects']);

    cy.contains('h2', 'Perfiles de Usuarios').should('be.visible');
    cy.contains('h3', 'Agregar Nuevo Usuario').should('be.visible');
    cy.contains('td', 'Usuario').should('be.visible');
    
    cy.contains('td', 'admin_test')
      .parent('tr') 
      .contains('(Tú)') 
      .scrollIntoView() 
      .should('be.visible');
  });


  // --- PRUEBA 2: Agregar un nuevo usuario ---
  
  it('debe permitir a un admin agregar un nuevo usuario (Happy Path)', () => {
    
    // ⭐️ CORRECCIÓN: Usamos '**' para un "matching" de ruta más robusto
    cy.intercept('POST', '**/api/admin/add-user', {
      statusCode: 201,
      fixture: 'admin-users-updated.json'
    }).as('addUser');

    // Visitar la página y esperar la carga
    cy.visit('/admin/usuarios');
    cy.wait(['@getUsers', '@getProjects']);

    // Llenar el formulario
    cy.get('input[name="username"]').type('nuevo_usuario_test');
    cy.get('input[name="password"]').type('pass123');
    cy.get('input[name="nombre"]').type('Test');
    cy.get('input[name="apellido"]').type('Creado');
    cy.get('input[name="cedula"]').type('V-333');
    
    cy.get('select[name="role"]').select('user');

    // (Esta es la línea que fallaba)
    // Ahora esperará a que el fixture con formato correcto
    // sea renderizado por React.
    cy.get('select[name="proyecto_id"]')
      .find('option[value="2"]')
      .should('exist');

    // Ahora que sabemos que la opción existe, la seleccionamos
    cy.get('select[name="proyecto_id"]').select('2'); 

    // Enviar el formulario
    cy.get('form').contains('button', 'Agregar Usuario').click();

    // Verificar los resultados
    cy.wait('@addUser');
    cy.contains('Usuario agregado con éxito').should('be.visible');
    cy.get('input[name="username"]').should('have.value', '');
    
    // Verificar que el NUEVO usuario está AHORA en la tabla
    cy.contains('td', 'usuario_prueba').should('be.visible');
  });

});