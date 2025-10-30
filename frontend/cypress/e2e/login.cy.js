/* eslint-disable no-undef */

const testUsername = `testuser_${Date.now()}`;

describe('Gestión de Usuarios - Admin', () => {

  beforeEach(() => {
    cy.visit('/login');
    cy.get('input[type="text"]').first().type('admin');
    cy.get('input[type="password"]').first().type('admin123');
    cy.get('button[type="submit"]').first().click();

    // Espera redirección
    cy.url({ timeout: 10000 }).should('include', '/admin');

    // Ir a la sección de usuarios
    cy.contains(/perfiles de usuarios/i, { timeout: 10000 }).click();
    cy.url().should('include', '/admin/usuarios');
  });


  // --- PRUEBA 1: CREAR UN NUEVO USUARIO ---
  it('1. Debería crear un nuevo usuario', () => {

    console.log('=== INICIANDO CREACIÓN DE USUARIO ===');
    console.log('Username a crear:', testUsername);

    // Esperar que cargue el formulario
    cy.contains(/crear nuevo usuario/i, { timeout: 10000 }).should('be.visible');

    // 🟢 LLENAR CAMPOS DEL FORMULARIO 🟢
    cy.get('input[placeholder="Nombre de Usuario"]').type(testUsername);
    cy.get('input[placeholder="Contraseña (mín. 6 caracteres)"]').type('password123');
    cy.get('input[placeholder="Nombre"]').type('Test');
    cy.get('input[placeholder="Apellido"]').type('User');

    // 🟣 Captura segura (sin romper si falla)
    cy.then(() => {
      try {
        cy.screenshot('antes-de-crear-usuario', { capture: 'viewport' });
      } catch (err) {
        cy.log('⚠️ Error ignorado al tomar screenshot:', err.message);
      }
    });

    // 🔵 AHORA SÍ: hacer clic en "Crear Usuario"
    cy.contains('button', /crear\s+usuario/i)
      .should('be.visible', { timeout: 10000 })
      .click({ force: true });

    // 🟢 Verificar mensaje de éxito o fallback si no aparece texto visible
    cy.wait(2000);
    cy.get('body').then(($body) => {
      const texto = $body.text().toLowerCase();

      if (
        texto.includes('exitosamente') ||
        texto.includes('éxito') ||
        texto.includes('correctamente') ||
        texto.includes('creado') ||
        texto.includes('agregado') ||
        texto.includes('registrado')
      ) {
        cy.log('✅ Usuario creado exitosamente');
      } else {
        cy.log('⚠️ No se encontró mensaje visible de éxito; puede ser un toast o mensaje oculto.');
      }
    });

    // 🧩 Confirmar que el usuario aparece en la tabla (resultado real)
    cy.contains('td', testUsername, { timeout: 10000 }).should('exist');
  });


  // --- PRUEBA 2: VERIFICAR QUE EL USUARIO EXISTE EN LA TABLA ---
  it('2. Debería verificar que el usuario se creó en la lista', () => {
    cy.reload();

    // Espera a que la tabla cargue
    cy.get('table tbody tr', { timeout: 10000 }).should('have.length.greaterThan', 0);

    // Busca el usuario recién creado
    cy.contains('td', testUsername)
      .should('be.visible')
      .parent('tr')
      .within(() => {
        cy.contains('td', 'Test').should('be.visible');
        cy.contains('td', 'User').should('be.visible');
        cy.contains('td', 'user').should('be.visible');
      });
  });

});
