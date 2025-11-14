// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })

// cypress/support/commands.js

/**
 * Comando personalizado para iniciar sesión como Admin.
 * Esto es mucho más rápido que llenar el formulario de UI.
 * Simplemente "setea" el localStorage como lo haría un login real.
 */
Cypress.Commands.add('loginAsAdmin', () => {
  // Datos simulados de un admin logueado
  // (Basado en nuestro 'auth-success.json' y 'AuthContext.js')
  const user = {
    id: 1,
    username: 'admin_test',
    nombre: 'Admin',
    apellido: 'Tester'
  };

  const token = 'jwt.token.simulado.para.pruebas';
  const role = 'admin';
  const userId = 1;

  // Seteamos el localStorage
  cy.window().then((win) => {
    win.localStorage.setItem('token', token);
    win.localStorage.setItem('user', JSON.stringify(user));
    win.localStorage.setItem('role', role);
    win.localStorage.setItem('userId', String(userId));
  });
});