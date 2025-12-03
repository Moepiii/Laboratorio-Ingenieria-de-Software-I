/* eslint-disable no-undef */

describe('Módulo: Logger de Eventos', () => {

    const mockLogs = [
        { id: 101, timestamp: '2025-01-01', usuario_username: 'admin', accion: 'CREACIÓN', entidad: 'Proyectos', entidad_id: 15 },
        { id: 102, timestamp: '2025-01-02', usuario_username: 'pepe', accion: 'LOGIN', entidad: 'Auth', entidad_id: 0 }
    ];

    beforeEach(() => {
        cy.intercept('POST', '**/api/admin/get-logs', { body: mockLogs }).as('getLogsInit');

        cy.visit('/login');
        cy.get('input[type="text"]').first().type('admin', { force: true });
        cy.get('input[type="password"]').first().type('admin123', { force: true });
        cy.get('button[type="submit"]').first().click({ force: true });

        cy.url({ timeout: 10000 }).should('include', '/admin');
        cy.contains('a', 'Logger de eventos').click({ force: true });
        cy.wait('@getLogsInit');
    });

    // --- PRUEBA 1: Eliminar UNO ---
    it('1. Elimina un evento y verifica que la tabla se reduce', () => {
        cy.intercept('POST', '**/api/admin/delete-logs', { body: { mensaje: "OK" } }).as('deleteReq');
        cy.intercept('POST', '**/api/admin/get-logs', { body: [mockLogs[1]] }).as('reloadTable');

        // 1. Seleccionar el primero
        cy.get('table tbody tr').first().find('input[type="checkbox"]').check({ force: true });

        // 2. Click en borrar (Cypress acepta la alerta automáticamente)
        cy.contains('button', 'Eliminar Seleccionados').click({ force: true });

        // 3. Esperar
        cy.wait('@deleteReq');
        cy.wait('@reloadTable');

        // 4. Verificar SOLO dentro de la tabla
        cy.get('table tbody').should('not.contain', 'Proyectos');
    });

    // --- PRUEBA 2: Eliminar TODOS ---
    it('2. Elimina todos y verifica que la tabla queda vacía', () => {
        cy.intercept('POST', '**/api/admin/delete-logs', { body: { mensaje: "OK" } }).as('deleteReq');
        // Simulamos respuesta vacía
        cy.intercept('POST', '**/api/admin/get-logs', { body: [] }).as('reloadTableEmpty');

        // 1. Seleccionar todo
        cy.get('table thead tr th input[type="checkbox"]').check({ force: true });

        // 2. Click en borrar
        cy.contains('button', 'Eliminar Seleccionados').click({ force: true });

        // 3. Esperar
        cy.wait('@deleteReq');
        cy.wait('@reloadTableEmpty');

        // ⭐️ CORRECCIÓN AQUÍ ⭐️
        // Verificamos que NO existan filas de datos (checkboxes)
        cy.get('table tbody tr input[type="checkbox"]').should('not.exist');

        // Y si verificamos texto, lo hacemos SOLO dentro de la tabla (tbody)
        // (Porque la palabra "Proyectos" existe en el menú lateral)
        cy.get('table tbody').should('not.contain', 'Proyectos');
    });

});