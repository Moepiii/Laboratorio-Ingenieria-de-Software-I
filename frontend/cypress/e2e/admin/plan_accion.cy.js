/* eslint-disable no-undef */

describe('Módulo: Plan de Acción', () => {

    const mockPlanes = [
        { id: 1, actividad: 'Fertilización', accion: 'Manual', fecha_inicio: '2025-01-01', fecha_cierre: '2025-01-05', horas: 10, responsable: 'Juan', costo_unitario: 50, monto: 500 },
        { id: 2, actividad: 'Poda', accion: 'Mecánica', fecha_inicio: '2025-02-01', fecha_cierre: '2025-02-02', horas: 5, responsable: 'Pedro', costo_unitario: 40, monto: 200 },
    ];

    // Datos simulados para que los Selects tengan opciones
    const mockConfig = {
        actividades: [{ id: 1, actividad: 'Fertilización' }, { id: 2, actividad: 'Poda' }],
        labores: [{ id: 1, descripcion: 'Manual' }, { id: 2, descripcion: 'Mecánica' }],
        encargados: [{ id: 1, nombre: 'Juan', apellido: 'Perez' }, { id: 2, nombre: 'Pedro', apellido: 'Gomez' }]
    };

    beforeEach(() => {
        cy.intercept('POST', '**/api/admin/get-planes', { body: { planes: mockPlanes } }).as('getPlanes');
        cy.intercept('POST', '**/api/admin/get-datos-proyecto', { body: mockConfig }).as('getConfig');
        cy.intercept('POST', '**/api/auth/login').as('loginRequest');

        cy.visit('/login');
        cy.get('input[type="text"]').first().type('admin', { force: true });
        cy.get('input[type="password"]').first().type('admin123', { force: true });
        cy.get('button[type="submit"]').first().click({ force: true });

        cy.wait('@loginRequest').its('response.statusCode').should('eq', 200);
        cy.url({ timeout: 10000 }).should('include', '/admin');

        cy.visit('/admin/planes-accion/proyecto/15/general');

        cy.wait('@getConfig');
        cy.wait('@getPlanes');
    });

    // --- PRUEBA 1: Verificar Tabla y SUMA TOTAL ---
    it('1. Muestra la tabla correctamente y verifica la Suma Total', () => {
        // Verificar filas
        cy.get('table tbody tr').should('have.length', 2);
        cy.get('table tbody tr').first().should('contain', 'Fertilización');

        // Verificar Total en el pie de página (500 + 200 = 700)
        // Buscamos la celda que contiene el texto "700.00"
        cy.contains('td', '700.00').should('be.visible');
    });

    // --- PRUEBA 2: Verificar FÓRMULA (Cálculo Automático) ---
    it('2. Calcula automáticamente el Monto al ingresar Horas y Dinero', () => {
        cy.contains('button', '+ Añadir').click();

        // Seleccionar Actividad para habilitar flujo (aunque no es estrictamente necesario para el cálculo, es más real)
        cy.contains('label', 'Actividad').next('select').select('Poda');

        // ⭐️ AQUÍ ESTABA EL ERROR: Corregido a "Cantidad Horas"
        cy.contains('label', 'Cantidad Horas').parent().find('input').type('5');

        // Ingresar Dinero ($)
        cy.contains('label', 'Dinero ($)').parent().find('input').type('20');

        // Verificar cálculo (5 * 20 = 100)
        cy.contains('label', 'Monto ($)').parent().find('input')
            .should('have.value', '100.00');
    });

    // --- PRUEBA 3: Crear un Nuevo Plan ---
    it('3. Envía los datos correctos al crear un nuevo plan', () => {
        cy.intercept('POST', '**/api/admin/create-plan', {
            statusCode: 201,
            body: { mensaje: "Plan creado" }
        }).as('createPlan');

        const nuevoPlan = { ...mockPlanes[0], id: 3, actividad: 'Poda', monto: 100 };
        cy.intercept('POST', '**/api/admin/get-planes', {
            body: { planes: [...mockPlanes, nuevoPlan] }
        }).as('getPlanesAfterCreate');

        cy.contains('button', '+ Añadir').click();

        // Llenar formulario
        cy.contains('label', 'Actividad').next('select').select('Poda');

        // Corregido a "Acción" (antes era Acción Específica)
        cy.contains('label', 'Acción').next('select').select('Manual');

        cy.contains('label', 'Fecha Inicio').parent().find('input').type('2025-05-01');
        cy.contains('label', 'Fecha Cierre').parent().find('input').type('2025-05-05');

        // Corregido a "Cantidad Horas"
        cy.contains('label', 'Cantidad Horas').parent().find('input').type('2');

        cy.contains('label', 'Responsable').parent().find('select').select('Juan Perez');
        cy.contains('label', 'Dinero ($)').parent().find('input').type('50');

        cy.contains('button', 'Guardar').click();

        // Validar envío
        cy.wait('@createPlan').then((interception) => {
            const body = interception.request.body;
            expect(body.horas).to.eq(2);
            expect(body.costo_unitario).to.eq(50);
            expect(body.monto).to.eq(100);
        });

        cy.wait('@getPlanesAfterCreate');
        cy.get('table tbody tr').should('have.length', 3);
    });

    // --- PRUEBA 4: Eliminar un Plan ---
    it('4. Elimina un plan existente', () => {
        cy.intercept('POST', '**/api/admin/delete-plan', { body: { mensaje: "OK" } }).as('deletePlan');
        cy.intercept('POST', '**/api/admin/get-planes', { body: { planes: [mockPlanes[0]] } }).as('reloadTable');

        // Click Borrar en la segunda fila
        cy.get('table tbody tr').eq(1).contains('button', 'Borrar').click();

        cy.wait('@deletePlan');
        cy.wait('@reloadTable');

        cy.get('table tbody tr').should('have.length', 1);
    });

});