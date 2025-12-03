/* eslint-disable no-undef */

describe('Módulo: Recurso Humano', () => {

    // 1. Datos simulados para la TABLA
    // Total esperado: 500 + 300 = 800
    const mockRecursos = [
        {
            id: 1,
            actividad: 'Cosecha',
            accion: 'Manual',
            nombre: 'Juan Perez',
            cedula: '12345',
            tiempo: 8,
            cantidad: 5,
            costo_unitario: 12.5,
            monto: 500
        },
        {
            id: 2,
            actividad: 'Riego',
            accion: 'Supervisión',
            nombre: 'Maria Lopez',
            cedula: '67890',
            tiempo: 4,
            cantidad: 1,
            costo_unitario: 75,
            monto: 300
        },
    ];

    // 2. Datos para los Selects (Configuración)
    const mockConfig = {
        actividades: [{ id: 1, actividad: 'Cosecha' }, { id: 2, actividad: 'Riego' }],
        labores: [{ id: 1, descripcion: 'Manual' }, { id: 2, descripcion: 'Supervisión' }],
        encargados: [{ id: 1, nombre: 'Juan', apellido: 'Perez', cedula: '12345' }, { id: 2, nombre: 'Maria', apellido: 'Lopez', cedula: '67890' }]
    };

    beforeEach(() => {
        // Interceptores
        cy.intercept('POST', '**/api/admin/get-recursos', { body: { recursos: mockRecursos } }).as('getRecursos');
        cy.intercept('POST', '**/api/admin/get-datos-proyecto', { body: mockConfig }).as('getConfig');
        cy.intercept('POST', '**/api/auth/login').as('loginRequest');

        // Login
        cy.visit('/login');
        cy.get('input[type="text"]').first().type('admin', { force: true });
        cy.get('input[type="password"]').first().type('admin123', { force: true });
        cy.get('button[type="submit"]').first().click({ force: true });

        // Esperar login
        cy.wait('@loginRequest').its('response.statusCode').should('eq', 200);
        cy.url({ timeout: 10000 }).should('include', '/admin');

        // Navegar a la ruta de Recursos Humanos (Proyecto ID 15 simulado)
        cy.visit('/admin/planes-accion/proyecto/15/recursos');

        // Esperar cargas
        cy.wait('@getConfig');
        cy.wait('@getRecursos');
    });

    // --- PRUEBA 1: Verificar Tabla y Total ---
    it('1. Muestra la tabla y verifica la suma total del Talento Humano', () => {
        // Verificar filas
        cy.get('table tbody tr').should('have.length', 2);

        // Verificar datos visuales
        cy.get('table tbody tr').first().should('contain', 'Cosecha');
        cy.get('table tbody tr').first().should('contain', 'Juan Perez');

        // ⭐️ VERIFICAR TOTAL EN EL PIE DE PÁGINA (500 + 300 = 800)
        // Buscamos la celda que contiene "800.00"
        cy.contains('td', '800.00').should('be.visible');

        // Verificar etiqueta del total
        cy.contains('Monto Total Talento Humano ($):').should('be.visible');
    });

    // --- PRUEBA 2: Verificar la FÓRMULA ---
    it('2. Calcula automáticamente el Monto con la fórmula (Tiempo/Cant * Costo * Cant)', () => {
        cy.contains('button', '+ Añadir').click();

        // 1. Llenar Tiempo = 10
        cy.contains('label', 'Tiempo (Días/Horas)').parent().find('input').type('10');

        // 2. Llenar Cantidad = 2
        cy.contains('label', 'Cantidad (Personas)').parent().find('input').type('2');

        // 3. Llenar Costo ($) = 50
        // Cálculo esperado: (10 / 2) * 50 * 2 
        // Paso 1: 5 * 50 = 250
        // Paso 2: 250 * 2 = 500
        cy.contains('label', 'Costo ($)').parent().find('input').type('50');

        // 4. Verificar Monto
        cy.contains('label', 'Monto ($)').parent().find('input')
            .should('have.value', '500.00');
    });

    // --- PRUEBA 3: Crear Nuevo Recurso ---
    it('3. Crea un nuevo recurso humano exitosamente', () => {
        cy.intercept('POST', '**/api/admin/create-recurso', {
            statusCode: 201,
            body: { mensaje: "Recurso creado" }
        }).as('createRecurso');

        // Simular que la tabla crece
        const nuevoRecurso = { ...mockRecursos[0], id: 3, actividad: 'Riego', monto: 100 };
        cy.intercept('POST', '**/api/admin/get-recursos', {
            body: { recursos: [...mockRecursos, nuevoRecurso] }
        }).as('getRecursosUpdate');

        cy.contains('button', '+ Añadir').click();

        // Llenar formulario
        cy.contains('label', 'Actividad').next('select').select('Riego');
        cy.contains('label', 'Acción').next('select').select('Supervisión');

        // Seleccionar responsable (debe autocompletar la cédula, aunque esté oculta)
        cy.contains('label', 'Responsable').next('select').select('Maria Lopez');

        // Llenar números
        cy.contains('label', 'Tiempo').parent().find('input').type('5');
        cy.contains('label', 'Cantidad').parent().find('input').type('1');
        cy.contains('label', 'Costo ($)').parent().find('input').type('100');

        // Guardar
        cy.contains('button', 'Guardar').click();

        // Validar envío al backend
        cy.wait('@createRecurso').then((interception) => {
            const body = interception.request.body;
            expect(body.tiempo).to.eq(5);
            expect(body.cantidad).to.eq(1);
            expect(body.costo_unitario).to.eq(100);
            expect(body.monto).to.eq(500); // (5/1)*100*1 = 500
        });

        // Validar tabla actualizada
        cy.wait('@getRecursosUpdate');
        cy.get('table tbody tr').should('have.length', 3);
    });

    // --- PRUEBA 4: Eliminar Recurso ---
    it('4. Elimina un recurso existente', () => {
        cy.intercept('POST', '**/api/admin/delete-recurso', { body: { mensaje: "OK" } }).as('deleteRecurso');
        // Simulamos que queda solo 1
        cy.intercept('POST', '**/api/admin/get-recursos', { body: { recursos: [mockRecursos[0]] } }).as('getRecursosDelete');

        // Clic en Borrar (segunda fila)
        cy.get('table tbody tr').eq(1).contains('button', 'Borrar').click();

        // (Cypress acepta confirm automáticamente)

        cy.wait('@deleteRecurso');
        cy.wait('@getRecursosDelete');

        cy.get('table tbody tr').should('have.length', 1);
    });

});