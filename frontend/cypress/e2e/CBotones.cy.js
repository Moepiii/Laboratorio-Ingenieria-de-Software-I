/* eslint-disable no-undef */

describe('Gestión de Estado de Proyectos (Habilitar/Cerrar/Borrar/Modificar)', () => {

    // --- Datos Simulados ---
    const proyectoHabilitado = {
        id: 1,
        nombre: 'Proyecto Activo (Test)',
        fecha_inicio: '2025-01-01T00:00:00Z',
        fecha_cierre: null,
        estado: 'habilitado'
    };
    const proyectoCerrado = {
        id: 2,
        nombre: 'Proyecto Inactivo (Test)',
        fecha_inicio: '2024-01-01T00:00:00Z',
        fecha_cierre: '2024-12-31T00:00:00Z',
        estado: 'cerrado'
    };

    beforeEach(() => {
        // --- Login como Admin ---
        cy.visit('/login');
        cy.get('input[type="text"]').first().type('admin');
        cy.get('input[type="password"]').first().type('admin123');
        cy.get('button[type="submit"]').first().click();
        cy.url({ timeout: 10000 }).should('include', '/admin');

        // --- Interceptar carga INICIAL de Proyectos ---
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: {
                proyectos: [proyectoHabilitado, proyectoCerrado]
            }
        }).as('getProjects');

        // --- Navegar a Proyectos ---
        cy.contains('a', /Portafolio de Proyectos/i).click();
        cy.url().should('include', '/admin/proyectos');

        // Espera a que la tabla cargue
        cy.wait('@getProjects');
        cy.contains('td', proyectoHabilitado.nombre).should('be.visible');
        cy.contains('td', proyectoCerrado.nombre).should('be.visible');
    });

    // --- Prueba 1: Verificar botones (Proyecto Cerrado) ---
    it('debería deshabilitar Modificar y Cerrar si el proyecto está Cerrado', () => {
        cy.contains('td', proyectoCerrado.nombre).click();
        cy.get('main.main-content').within(() => {
            cy.contains('button', 'Modificar').should('be.disabled');
            cy.contains('button', 'Cerrar').should('be.disabled');
            cy.contains('button', 'Habilitar').should('be.enabled');
            cy.contains('button', 'Borrar').should('be.enabled');
        });
    });

    // --- Prueba 2: Verificar botones (Proyecto Habilitado) ---
    it('debería deshabilitar Habilitar si el proyecto está Habilitado', () => {
        cy.contains('td', proyectoHabilitado.nombre).click();
        cy.get('main.main-content').within(() => {
            cy.contains('button', 'Modificar').should('be.enabled');
            cy.contains('button', 'Cerrar').should('be.enabled');
            cy.contains('button', 'Habilitar').should('be.disabled');
            cy.contains('button', 'Borrar').should('be.enabled');
        });
    });

    // --- Prueba 3: Probar funcionalidad "Habilitar" ---
    it("debería permitir HABILITAR un proyecto cerrado", () => {
        cy.intercept('POST', '**/api/admin/set-proyecto-estado', (req) => {
            expect(req.body.id).to.equal(proyectoCerrado.id);
            expect(req.body.estado).to.equal('habilitado');
            req.reply({ statusCode: 200, body: { mensaje: "Estado actualizado a 'habilitado'" } });
        }).as('setEstado');
        cy.intercept('POST', '**/api/admin/get-proyectos', { statusCode: 200, body: { proyectos: [proyectoHabilitado, { ...proyectoCerrado, estado: 'habilitado' }] } }).as('reloadProjects');

        cy.contains('td', proyectoCerrado.nombre).click();
        cy.get('main.main-content').contains('button', 'Habilitar').click();
        cy.on('window:confirm', () => true);

        cy.wait('@setEstado');
        cy.wait('@reloadProjects');
        cy.contains("Estado actualizado a 'habilitado'").should('be.visible');
        cy.contains('td', proyectoCerrado.nombre).parent('tr').contains('Habilitado');
    });

    // --- Prueba 4: Probar funcionalidad "Cerrar" ---
    it("debería permitir CERRAR un proyecto habilitado", () => {
        cy.intercept('POST', '**/api/admin/set-proyecto-estado', (req) => {
            expect(req.body.id).to.equal(proyectoHabilitado.id);
            expect(req.body.estado).to.equal('cerrado');
            req.reply({ statusCode: 200, body: { mensaje: "Proyecto cerrado." } });
        }).as('setEstado');
        cy.intercept('POST', '**/api/admin/get-proyectos', { statusCode: 200, body: { proyectos: [{ ...proyectoHabilitado, estado: 'cerrado' }, proyectoCerrado] } }).as('reloadProjects');

        cy.contains('td', proyectoHabilitado.nombre).click();
        cy.get('main.main-content').contains('button', 'Cerrar').click();
        cy.on('window:confirm', () => true);

        cy.wait('@setEstado');
        cy.wait('@reloadProjects');
        cy.contains("Proyecto cerrado.").should('be.visible');
        cy.contains('td', proyectoHabilitado.nombre).parent('tr').contains('Cerrado');
    });

    // --- Prueba 5: Probar funcionalidad "Borrar" ---
    it("debería permitir BORRAR un proyecto", () => {
        cy.intercept('POST', '**/api/admin/delete-proyecto', (req) => {
            expect(req.body.id).to.equal(proyectoHabilitado.id);
            expect(req.body.admin_username).to.equal('admin');
            req.reply({ statusCode: 200, body: { mensaje: "Proyecto eliminado" } });
        }).as('deleteProject');
        cy.intercept('POST', '**/api/admin/get-proyectos', { statusCode: 200, body: { proyectos: [proyectoCerrado] } }).as('reloadProjectsAfterDelete');

        cy.contains('td', proyectoHabilitado.nombre).click();
        cy.get('main.main-content').contains('button', 'Borrar').click();
        cy.on('window:confirm', () => true);

        cy.wait('@deleteProject');
        cy.wait('@reloadProjectsAfterDelete');
        cy.contains("Proyecto eliminado").should('be.visible');
        cy.contains('td', proyectoHabilitado.nombre).should('not.exist');
        cy.contains('td', proyectoCerrado.nombre).should('be.visible');
    });

    // --- ⭐️ PRUEBA 6 AÑADIDA: Probar funcionalidad "Modificar" ⭐️ ---
    it("debería permitir MODIFICAR un proyecto habilitado", () => {
        // 1. Define los nuevos datos
        const nuevoNombre = 'Proyecto Activo (Modificado)';
        const nuevaFecha = '2025-02-02'; // Formato YYYY-MM-DD para el input[type=date]

        // 2. Interceptar la API de ACTUALIZACIÓN (update-proyecto)
        cy.intercept('POST', '**/api/admin/update-proyecto', (req) => {
            // Verifica que el frontend envía los datos nuevos
            expect(req.body.id).to.equal(proyectoHabilitado.id); // ID 1
            expect(req.body.nombre).to.equal(nuevoNombre);
            expect(req.body.fecha_inicio).to.equal(nuevaFecha);
            expect(req.body.admin_username).to.equal('admin');
            req.reply({
                statusCode: 200,
                body: { mensaje: "Proyecto actualizado con éxito" } // Mensaje de tu app
            });
        }).as('updateProject');

        // 3. Interceptar la RECARGA de proyectos (get-proyectos)
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: {
                proyectos: [
                    // El proyecto 1 ahora tiene los datos nuevos
                    { ...proyectoHabilitado, nombre: nuevoNombre, fecha_inicio: `${nuevaFecha}T00:00:00Z` },
                    proyectoCerrado // El proyecto 2 sigue igual
                ]
            }
        }).as('reloadProjectsAfterUpdate');

        // 4. Acción: Seleccionar el proyecto HABILITADO
        cy.contains('td', proyectoHabilitado.nombre).click(); // Clic en "Proyecto Activo (Test)"

        // 5. Acción: Clic en Modificar (dentro del main)
        cy.get('main.main-content').contains('button', 'Modificar').click();

        // 6. Verificar que el formulario aparece y está pre-llenado
        cy.get('form').should('be.visible');
        cy.get('input#nombre').should('have.value', proyectoHabilitado.nombre); // Verifica nombre antiguo
        // Verifica fecha antigua (formateada como YYYY-MM-DD)
        cy.get('input#fecha_inicio').should('have.value', '2025-01-01');

        // 7. Acción: Modificar los campos del formulario
        cy.get('input#nombre').clear().type(nuevoNombre);
        cy.get('input#fecha_inicio').clear().type(nuevaFecha);
        // (Opcional: modificar fecha_cierre si es necesario)

        // 8. Acción: Guardar cambios
        // (El botón de submit en el form se llama "Guardar Cambios")
        cy.get('form').contains('button', /Guardar Cambios/i).click();

        // 9. Verificar
        cy.wait('@updateProject');
        cy.wait('@reloadProjectsAfterUpdate');

        // Verifica que el mensaje de éxito aparece
        cy.contains("Proyecto actualizado con éxito").should('be.visible');

        // Verifica que el formulario desapareció (el input #nombre ya no existe)
        cy.get('form').should('not.exist');
        cy.get('input#nombre').should('not.exist');

        // Verifica que el proyecto con el NUEVO nombre está en la tabla
        cy.contains('td', nuevoNombre).should('be.visible');

        // Verifica que el proyecto con el VIEJO nombre ya NO está
        cy.contains('td', proyectoHabilitado.nombre).should('not.exist');
    });

});