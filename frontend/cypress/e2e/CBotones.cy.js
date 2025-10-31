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
        // --- ⭐️ ARREGLO 1: Define TODAS las interceptaciones PRIMERO ---

        // 1. Intercepta el Login del Admin
        cy.intercept('POST', '**/api/login', {
            statusCode: 200,
            body: {
                token: "jwt.token.simulado.admin",
                user: { username: "admin", nombre: "Admin", apellido: "User" },
                userId: 1,
                role: "admin"
            },
        }).as('loginAdmin');

        // 2. Intercepta la carga INICIAL de Proyectos (la que se llama al cargar /admin)
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: {
                proyectos: [proyectoHabilitado, proyectoCerrado]
            }
        }).as('getProjects');

        // --- ⭐️ ARREGLO 2: Ahora visita la página e interactúa ---
        cy.visit('/login');

        // 3. Realiza el Login
        cy.get('input[type="text"]').first().type('admin');
        cy.get('input[type="password"]').first().type('admin123');
        cy.get('button[type="submit"]').first().click();

        // 4. Espera a que el login (y la redirección implícita) ocurra
        cy.wait('@loginAdmin');
        cy.url({ timeout: 10000 }).should('include', '/admin');

        // 5. Navega a Proyectos
        // (Esta llamada ya está interceptada por '@getProjects')
        cy.contains('a', /Portafolio de Proyectos/i).click();
        cy.url().should('include', '/admin/proyectos');

        // 6. Espera a que la tabla cargue
        // (Esto ahora espera de forma fiable la llamada API que se disparó en el paso 5)
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
        // ⭐️ ARREGLO: Define las interceptaciones específicas de ESTA prueba
        cy.intercept('POST', '**/api/admin/set-proyecto-estado', (req) => {
            expect(req.body.id).to.equal(proyectoCerrado.id);
            expect(req.body.estado).to.equal('habilitado');
            req.reply({ statusCode: 200, body: { mensaje: "Estado actualizado a 'habilitado'" } });
        }).as('setEstado');

        // Sobrescribe el 'get-proyectos' solo para la RECARGA
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: { proyectos: [proyectoHabilitado, { ...proyectoCerrado, estado: 'habilitado' }] }
        }).as('reloadProjects');

        // Acción
        cy.contains('td', proyectoCerrado.nombre).click();
        cy.get('main.main-content').contains('button', 'Habilitar').click();
        cy.on('window:confirm', () => true);

        // Verificación
        cy.wait('@setEstado');
        cy.wait('@reloadProjects'); // Espera la recarga
        cy.contains("Estado actualizado a 'habilitado'").should('be.visible');
        cy.contains('td', proyectoCerrado.nombre).parent('tr').contains('Habilitado');
    });

    // --- Prueba 4: Probar funcionalidad "Cerrar" ---
    it("debería permitir CERRAR un proyecto habilitado", () => {
        // ⭐️ ARREGLO: Define las interceptaciones específicas de ESTA prueba
        cy.intercept('POST', '**/api/admin/set-proyecto-estado', (req) => {
            expect(req.body.id).to.equal(proyectoHabilitado.id);
            expect(req.body.estado).to.equal('cerrado');
            req.reply({ statusCode: 200, body: { mensaje: "Proyecto cerrado." } });
        }).as('setEstado');

        // Sobrescribe el 'get-proyectos' solo para la RECARGA
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: { proyectos: [{ ...proyectoHabilitado, estado: 'cerrado' }, proyectoCerrado] }
        }).as('reloadProjects');

        // Acción
        cy.contains('td', proyectoHabilitado.nombre).click();
        cy.get('main.main-content').contains('button', 'Cerrar').click();
        cy.on('window:confirm', () => true);

        // Verificación
        cy.wait('@setEstado');
        cy.wait('@reloadProjects');
        cy.contains("Proyecto cerrado.").should('be.visible');
        cy.contains('td', proyectoHabilitado.nombre).parent('tr').contains('Cerrado');
    });

    // --- Prueba 5: Probar funcionalidad "Borrar" ---
    it("debería permitir BORRAR un proyecto", () => {
        // ⭐️ ARREGLO: Define las interceptaciones específicas de ESTA prueba
        cy.intercept('POST', '**/api/admin/delete-proyecto', (req) => {
            expect(req.body.id).to.equal(proyectoHabilitado.id);
            expect(req.body.admin_username).to.equal('admin');
            req.reply({ statusCode: 200, body: { mensaje: "Proyecto eliminado" } });
        }).as('deleteProject');

        // Sobrescribe el 'get-proyectos' solo para la RECARGA
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: { proyectos: [proyectoCerrado] } // Solo queda el proyecto 2
        }).as('reloadProjectsAfterDelete');

        // Acción
        cy.contains('td', proyectoHabilitado.nombre).click();
        cy.get('main.main-content').contains('button', 'Borrar').click();
        cy.on('window:confirm', () => true);

        // Verificación
        cy.wait('@deleteProject');
        cy.wait('@reloadProjectsAfterDelete');
        cy.contains("Proyecto eliminado").should('be.visible');
        cy.contains('td', proyectoHabilitado.nombre).should('not.exist');
        cy.contains('td', proyectoCerrado.nombre).should('be.visible');
    });

    // --- Prueba 6: Probar funcionalidad "Modificar" ---
    it("debería permitir MODIFICAR un proyecto habilitado", () => {
        const nuevoNombre = 'Proyecto Activo (Modificado)';
        const nuevaFecha = '2025-02-02';

        // ⭐️ ARREGLO: Define las interceptaciones específicas de ESTA prueba
        cy.intercept('POST', '**/api/admin/update-proyecto', (req) => {
            expect(req.body.id).to.equal(proyectoHabilitado.id);
            expect(req.body.nombre).to.equal(nuevoNombre);
            expect(req.body.fecha_inicio).to.equal(nuevaFecha);
            expect(req.body.admin_username).to.equal('admin');
            req.reply({ statusCode: 200, body: { mensaje: "Proyecto actualizado con éxito" } });
        }).as('updateProject');

        // Sobrescribe el 'get-proyectos' solo para la RECARGA
        cy.intercept('POST', '**/api/admin/get-proyectos', {
            statusCode: 200,
            body: {
                proyectos: [
                    { ...proyectoHabilitado, nombre: nuevoNombre, fecha_inicio: `${nuevaFecha}T00:00:00Z` },
                    proyectoCerrado
                ]
            }
        }).as('reloadProjectsAfterUpdate');

        // Acción
        cy.contains('td', proyectoHabilitado.nombre).click();
        cy.get('main.main-content').contains('button', 'Modificar').click();

        // Llenar formulario
        cy.get('form').should('be.visible');
        cy.get('input#nombre').should('have.value', proyectoHabilitado.nombre);
        cy.get('input#fecha_inicio').should('have.value', '2025-01-01');
        cy.get('input#nombre').clear().type(nuevoNombre);
        cy.get('input#fecha_inicio').clear().type(nuevaFecha);
        cy.get('form').contains('button', /Guardar Cambios/i).click();

        // Verificación
        cy.wait('@updateProject');
        cy.wait('@reloadProjectsAfterUpdate');
        cy.contains("Proyecto actualizado con éxito").should('be.visible');
        cy.get('form').should('not.exist');
        cy.contains('td', nuevoNombre).should('be.visible');
        cy.contains('td', proyectoHabilitado.nombre).should('not.exist');
    });

});