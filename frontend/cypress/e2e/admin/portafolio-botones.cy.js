/**
 * cypress/e2e/admin/portafolio-botones.cy.js
 */
describe('Portafolio de Proyectos - Botones', () => {

  beforeEach(() => {
    cy.loginAsAdmin();

    // ⭐️ CORRECCIÓN 1:
    // Interceptamos la carga y usamos tu fixture existente "admin-projects.json".
    // Asumimos que tu componente Portafolio.js espera la clave "projects" (en inglés).
    cy.intercept('POST', '**/api/admin/get-proyectos', {
      statusCode: 200,
      fixture: 'admin-projects.json' 
    }).as('getProjects');
  });

// --- PRUEBA 1 (Corregida para el comportamiento del formulario en la página) ---
  it('debe mostrar el formulario de "Agregar Proyecto" al hacer clic', () => {

    // 1. Visitar la página
    cy.visit('/admin/proyectos');

    // 2. Esperar a que la página esté estable (el título)
    cy.contains('h2', 'Portafolio de Proyectos').should('be.visible');

    // 3. Pre-condición: El formulario NO debe ser visible
    //    (Buscamos el título del formulario que movimos arriba)
    cy.contains('h3', 'Agregar Nuevo Proyecto').should('not.exist');

    // 4. Acción: Clic en el botón
    //    (Uso "Crear" como en tu archivo de prueba)
    cy.contains('Crear').should('be.visible').click();

    // 5. ⭐️ CORRECCIÓN (Esta es la aserción correcta) ⭐️
    //    Verificamos que el formulario ES visible en la MISMA página
    cy.contains('h3', 'Nombre del Proyecto').should('be.visible');

    // 6. Post-condición extra: El botón de agregar ahora está deshabilitado
    cy.contains('Guardar').should('be.disabled');
  });
  // (Aquí pondremos las siguientes pruebas)

});