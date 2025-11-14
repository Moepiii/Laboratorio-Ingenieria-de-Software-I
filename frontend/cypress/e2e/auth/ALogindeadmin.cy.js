/* eslint-disable no-undef */

describe('Autenticación de Usuarios (Historia A)', () => {

  beforeEach(() => {
    // Visita la página de login antes de cada prueba
    cy.visit('/login');
    cy.url().should('include', '/login');
  });

  it('a. Debería autenticar a un administrador exitosamente y redirigir a /admin', () => {
    // Usamos los mismos selectores que tu otra prueba
    cy.get('input[type="text"]').first().type('admin');
    cy.get('input[type="password"]').first().type('admin123');
    cy.get('button[type="submit"]').first().click();

    // Verifica la redirección y el contenido del dashboard de admin
    cy.url({ timeout: 10000 }).should('include', '/admin');
    cy.contains(/portafolio de proyectos/i).should('be.visible');
    cy.contains(/perfiles de usuarios/i).should('be.visible');
  });

  it('a. Debería fallar la autenticación con contraseña incorrecta', () => {
    cy.get('input[type="text"]').first().type('admin');
    cy.get('input[type="password"]').first().type('password-incorrecto');
    cy.get('button[type="submit"]').first().click();

    // Debería permanecer en la página de login
    cy.url().should('include', '/login');

    // Debería mostrar un mensaje de error
    // (Ajusta "Credenciales inválidas" si tu mensaje de error es diferente)
    cy.contains(/credenciales inválidas/i).should('be.visible');
  });

  it('a. Debería fallar la autenticación con usuario incorrecto', () => {
    cy.get('input[type="text"]').first().type('usuario-que-no-existe');
    cy.get('input[type="password"]').first().type('admin123');
    cy.get('button[type="submit"]').first().click();

    // Debería permanecer en la página de login
    cy.url().should('include', '/login');

    // Debería mostrar un mensaje de error
    cy.contains(/credenciales inválidas/i).should('be.visible');
  });

  it('a. Debería permitir al usuario cerrar sesión', () => {
    // Inicia sesión primero
    cy.get('input[type="text"]').first().type('admin');
    cy.get('input[type="password"]').first().type('admin123');
    cy.get('button[type="submit"]').first().click();
    cy.url({ timeout: 10000 }).should('include', '/admin');

    // Busca y haz clic en el botón de cerrar sesión
    // (Basado en la captura Screenshot_237.png, el botón está en el sidebar)
    cy.contains('button', /cerrar sesión/i).click();

    // Debería redirigir de vuelta al login
    cy.url({ timeout: 10000 }).should('include', '/login');
    // Verifica que el formulario de login esté visible de nuevo
    cy.get('input[type="password"]').first().should('be.visible');
  });

});