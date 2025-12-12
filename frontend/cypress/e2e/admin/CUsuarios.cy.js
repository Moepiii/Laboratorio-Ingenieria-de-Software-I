/* eslint-disable no-undef */

describe('Gestión de Usuarios (Historia C)', () => {

  // Generamos datos aleatorios que cumplan con tus requisitos
  const randomNum = Math.floor(Math.random() * 90000000) + 10000000; // Genera un número de 8 dígitos
  const nuevoUsuario = {
    username: `testuser_${Math.floor(Math.random() * 1000)}`,
    password: 'passwordSegura123', // Mínimo 6 caracteres
    nombre: 'Cypress',
    apellido: 'Test',
    cedula: `V-${randomNum}` // Formato V-XXXXXXXX
  };

  // --- PASO PREVIO: Iniciar sesión como Admin ---
  beforeEach(() => {
    cy.visit('/login');
    // Esperamos a que el formulario de login sea visible
    cy.get('form').should('be.visible');

    // Login
    cy.get('input[id="username"]').type('admin');
    cy.get('input[id="password"]').type('admin123');
    cy.get('button[type="submit"]').click();

    // Esperamos redirección
    cy.url({ timeout: 10000 }).should('include', '/admin');

    // Navegamos a Usuarios y esperamos carga
    cy.contains('Perfiles de usuarios').click();
    cy.url().should('include', '/admin/usuarios');


    cy.contains('button', 'Agregar Usuario').should('be.visible');
  });

  // --- PRUEBA 1: Crear un Nuevo Usuario ---
  it('C.1 Debería crear un nuevo usuario exitosamente', () => {


    // Username
    cy.get('input[name="username"]')
      .should('be.visible')
      .click()
      .type(nuevoUsuario.username);

    // Password (mínimo 6 caracteres)
    cy.get('input[name="password"]')
      .click()
      .type(nuevoUsuario.password);

    // Nombre
    cy.get('input[name="nombre"]')
      .click()
      .type(nuevoUsuario.nombre);

    // Apellido
    cy.get('input[name="apellido"]')
      .click()
      .type(nuevoUsuario.apellido);

    // Cédula (Formato V-XXXXXXX)
    cy.get('input[name="cedula"]')
      .click()
      .type(nuevoUsuario.cedula);

    // 2. Guardar (Click en "Agregar Usuario")
    cy.contains('button', 'Agregar Usuario').click();

    // 3. Verificaciones
    // Esperamos un momento para que la tabla se actualice
    cy.wait(1000);

    // Verificar que el usuario aparece en la tabla
    cy.contains('tr', nuevoUsuario.username).should('exist');
    cy.contains('tr', nuevoUsuario.cedula).should('exist');
  });

  // --- PRUEBA 2: Cambiar Rol y Asignar Proyecto ---
  it('C.2 Debería actualizar el rol y asignar proyecto', () => {
    // Localizamos la fila del usuario
    cy.contains('tr', nuevoUsuario.username).within(() => {

      // 1. CAMBIAR ROL
      // Seleccionamos 'gerente' y guardamos
      cy.get('select').first().select('gerente');
      cy.contains('button', 'Guardar Rol').click();
    });

    // Verificamos
    cy.wait(500);
    cy.contains('tr', nuevoUsuario.username).find('select').first().should('have.value', 'gerente');

    // 2. ASIGNAR PROYECTO (Si existen proyectos)
    cy.contains('tr', nuevoUsuario.username).within(() => {
      cy.get('select').last().then($select => {
        const options = $select.find('option');
        if (options.length > 1) {
          // Selecciona el segundo proyecto de la lista
          cy.get('select').last().select(options[1].value);
          // Esperamos un momento para asegurar que la petición se envió (el select tiene onChange automático)
          cy.wait(500);
        } else {
          cy.log('⚠️ No hay proyectos para asignar, saltando este paso.');
        }
      });
    });
  });

  // --- PRUEBA 3: Eliminar el Usuario ---
  it('C.3 Debería eliminar el usuario creado', () => {
    // 1. Localizar y borrar
    cy.contains('tr', nuevoUsuario.username).within(() => {
      cy.contains('button', 'Borrar').click();
    });

    // 2. Aceptar confirmación
    cy.on('window:confirm', () => true);

    // 3. Verificar que desapareció
    cy.wait(1000);
    cy.contains('tr', nuevoUsuario.username).should('not.exist');
  });
});