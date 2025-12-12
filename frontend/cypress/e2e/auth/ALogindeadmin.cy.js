/* eslint-disable no-undef */

describe('Autenticación de Usuarios (Historia A)', () => {

  beforeEach(() => {

    cy.visit('/login');
    cy.url().should('include', '/login');
  });

  it('a. Debería autenticar a un administrador exitosamente y redirigir a /admin', () => {

    cy.get('input[type="text"]').first().type('admin');
    cy.get('input[type="password"]').first().type('admin123');
    cy.get('button[type="submit"]').first().click();


    cy.url({ timeout: 10000 }).should('include', '/admin');
    cy.contains(/portafolio de proyectos/i).should('be.visible');
    cy.contains(/perfiles de usuarios/i).should('be.visible');
  });

  it('a. Debería fallar la autenticación con contraseña incorrecta', () => {
    cy.get('input[type="text"]').first().type('admin');
    cy.get('input[type="password"]').first().type('password-incorrecto');
    cy.get('button[type="submit"]').first().click();


    cy.url().should('include', '/login');


    cy.contains(/credenciales inválidas/i).should('be.visible');
  });

  it('a. Debería fallar la autenticación con usuario incorrecto', () => {
    cy.get('input[type="text"]').first().type('usuario-que-no-existe');
    cy.get('input[type="password"]').first().type('admin123');
    cy.get('button[type="submit"]').first().click();


    cy.url().should('include', '/login');


    cy.contains(/credenciales inválidas/i).should('be.visible');
  });

  it('a. Debería permitir al usuario cerrar sesión', () => {
    // Inicia sesión primero
    cy.get('input[type="text"]').first().type('admin');
    cy.get('input[type="password"]').first().type('admin123');
    cy.get('button[type="submit"]').first().click();
    cy.url({ timeout: 10000 }).should('include', '/admin');


    cy.contains('button', /cerrar sesión/i).click();


    cy.url({ timeout: 10000 }).should('include', '/login');

    cy.get('input[type="password"]').first().should('be.visible');
  });

});