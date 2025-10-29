/**
 * Prueba E2E para el flujo de creación de un proyecto por Admin/Gerente.
 * Asume que el usuario ya está autenticado.
 */
describe('Flujo Admin - Creación de Proyecto', () => {

    // Antes de cada prueba ('it'), simula estar logueado y visita la página
    beforeEach(() => {
        // 1. Simula la sesión guardando datos en localStorage
        const adminUser = { username: "adminTest", nombre: "Admin", apellido: "Prueba" };
        localStorage.setItem('token', 'token-admin-simulado');
        localStorage.setItem('user', JSON.stringify(adminUser));
        localStorage.setItem('role', 'admin'); // o 'gerente'
        localStorage.setItem('userId', '1'); // localStorage guarda strings

        // 2. Intercepta la carga INICIAL de proyectos (vacía)
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: { proyectos: [] }
        }).as('loadInitialProjects');

        // 3. Visita directamente la página de proyectos del admin
        cy.visit('/admin/proyectos');

        // 4. Espera a que la carga inicial se complete
        cy.wait('@loadInitialProjects');
    });

    // --- Prueba Única: Crear Proyecto ---
    it('debería permitir crear un proyecto y mostrarlo en la lista', () => {
        // 1. Verifica estado inicial (tabla vacía)
        cy.contains('Portafolio de Proyectos').should('be.visible');
        cy.contains('td', 'No hay proyectos creados.').should('be.visible');

        // 2. Define las interceptaciones para ESTA prueba específica
        // Intercepta la llamada de CREACIÓN
        cy.intercept('POST', '**/api/admin/create-proyecto', (req) => {
            // Verifica que el frontend envíe los datos correctos
            expect(req.body.nombre).to.equal('Proyecto Test Cypress');
            expect(req.body.fecha_inicio).to.equal('2024-12-01');
            expect(req.body.admin_username).to.equal('adminTest'); // Verifica el username enviado
            // Simula la respuesta exitosa del backend
            req.reply({
                statusCode: 201,
                body: { mensaje: "Proyecto 'Proyecto Test Cypress' creado." }
            });
        }).as('createProjectApi');

        // Intercepta la llamada para RECARGAR la lista DESPUÉS de crear
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: {
                proyectos: [ // Ahora la lista incluye el nuevo proyecto
                    { id: 5, nombre: 'Proyecto Test Cypress', fecha_inicio: '2024-12-01T00:00:00Z', fecha_cierre: null, estado: 'habilitado' }
                ]
            }
        }).as('reloadProjectsApi');

        // 3. Interactúa con la UI para crear
        cy.contains('button', 'Crear').click(); // Clic en el botón "Crear" de la toolbar

        // Llena el formulario (asegúrate que los IDs existen)
        cy.get('#nombre').should('be.visible').type('Proyecto Test Cypress');
        cy.get('#fecha_inicio').should('be.visible').type('2024-12-01');

        // Envía el formulario (busca el botón DENTRO del form)
        cy.contains('form button[type="submit"]', 'Crear Proyecto').click();

        // 4. Espera a que las llamadas API específicas terminen
        cy.wait('@createProjectApi');
        cy.wait('@reloadProjectsApi'); // Espera la recarga

        // 5. Verifica el resultado final
        cy.contains('Proyecto creado con éxito').should('be.visible'); // Mensaje de éxito
        cy.get('#nombre').should('not.exist'); // Formulario se cierra
        cy.contains('td', 'Proyecto Test Cypress').should('be.visible'); // Proyecto en la tabla
        cy.contains('td', '01/12/2024').should('be.visible'); // Fecha formateada
        cy.contains('span', 'Habilitado').should('be.visible'); // Estado
    });

});