
Cypress.Commands.add('loginAsAdmin', () => {

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