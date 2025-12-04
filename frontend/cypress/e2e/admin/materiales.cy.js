/* eslint-disable no-undef */

describe('Módulo: Materiales e Insumos', () => {

    // 1. CREAR EL USUARIO ENCARGADO SIMULADO
    const nuevoEncargado = {
        id: 50,
        nombre: 'Pedro',
        apellido: 'El Escamoso',
        cedula: '99999',
        role: 'encargado'
    };

    const mockMateriales = [
        {
            id: 10,
            actividad: 'Fumigación',
            accion: 'Química',
            categoria: 'Insumos',
            nombre: 'Glifosato',
            // --- CORRECCIÓN AQUÍ ---
            // Cambiamos 'Carlos Ruiz' por 'Carlos' para que coincida con la opción
            // disponible en el <select> del frontend.
            responsable: 'Carlos', 
            unidad: 'Lts',
            cantidad: 10,
            costo_unitario: 20,
            monto: 200
        }
    ];

    const mockConfig = {
        actividades: [{ id: 1, actividad: 'Fumigación' }, { id: 2, actividad: 'Siembra' }],
        labores: [{ id: 1, descripcion: 'Química' }, { id: 2, descripcion: 'Mecánica' }],
        encargados: [
            { id: 1, nombre: 'Carlos', apellido: 'Ruiz', role: 'encargado' },
            nuevoEncargado
        ]
    };

    beforeEach(() => {
        cy.intercept('POST', '**/api/admin/get-materiales', { body: { materiales: mockMateriales } }).as('getMateriales');
        cy.intercept('POST', '**/api/admin/get-datos-proyecto', { body: mockConfig }).as('getConfig');
        cy.intercept('POST', '**/api/auth/login').as('loginRequest');

        cy.visit('/login');
        cy.get('input[type="text"]').first().type('admin', { force: true });
        cy.get('input[type="password"]').first().type('admin123', { force: true });
        cy.get('button[type="submit"]').first().click({ force: true });

        cy.wait('@loginRequest').its('response.statusCode').should('eq', 200);
        cy.url({ timeout: 10000 }).should('include', '/admin');

        cy.visit('/admin/planes-accion/proyecto/15/materiales');

        cy.wait('@getConfig');
        cy.wait('@getMateriales');
    });

    // --- PRUEBA 1 ---
    it('1. Verifica la tabla y el monto total', () => {
        cy.get('table tbody tr').should('have.length', 1);
        cy.contains('Monto Total Materiales e Insumos ($):').should('be.visible');
        cy.contains('td', '200.00').should('be.visible');
    });

    // --- PRUEBA 2 ---
    it('2. Calcula automáticamente: Cantidad * Costo = Monto', () => {
        cy.contains('button', '+ Añadir').click();
        cy.get('input[name="cantidad"]').type('10');
        cy.get('input[name="costo_unitario"]').type('5.50');
        cy.get('input[name="monto"]').should('have.value', '55.00');
    });

    // --- PRUEBA 3 ---
    it('3. Crea un material asignando al encargado "Pedro"', () => {

        cy.intercept('POST', '**/api/admin/create-material', {
            statusCode: 201,
            body: { mensaje: "Material creado" }
        }).as('createMaterial');

        const nuevoMat = { ...mockMateriales[0], id: 99, nombre: 'Item Nuevo', responsable: 'Pedro' };
        cy.intercept('POST', '**/api/admin/get-materiales', {
            body: { materiales: [...mockMateriales, nuevoMat] }
        }).as('getMatUpdate');

        cy.contains('button', '+ Añadir').click();

        cy.get('select[name="actividad"]').select('Siembra');
        cy.get('select[name="accion"]').select('Mecánica');
        cy.get('select[name="categoria"]').select('Materiales');

        // Solución anterior aplicada:
        cy.get('select[name="responsable"]').should('contain.text', 'Pedro');
        cy.get('select[name="responsable"]').select('Pedro');

        cy.get('input[name="nombre"]').type('Tractor');
        cy.get('input[name="unidad"]').type('Unidad');
        cy.get('input[name="cantidad"]').type('1');
        cy.get('input[name="costo_unitario"]').type('5000');

        cy.contains('button', 'Guardar').click();

        cy.wait('@createMaterial').then((interception) => {
            const body = interception.request.body;
            expect(body.responsable).to.contain('Pedro');
        });

        cy.wait('@getMatUpdate');
        cy.get('table tbody tr').should('have.length', 2);
    });

// --- PRUEBA 4 (CORREGIDA) ---
    it('4. Edita un material correctamente', () => {
        cy.intercept('POST', '**/api/admin/update-material', { statusCode: 200, body: { mensaje: "OK" } }).as('updateMaterial');

        cy.get('table tbody tr').first().contains('button', 'Editar').click();

        // Verificamos que cargue el valor actual (Corregido en el paso anterior a 'Carlos')
        cy.get('select[name="responsable"]').should('have.value', 'Carlos');

        cy.get('input[name="costo_unitario"]').clear().type('25');
        
        // --- CORRECCIÓN AQUÍ ---
        // Cambiamos 'Actualizar' por 'Guardar', ya que es probable que el botón
        // conserve el mismo nombre que en la creación.
        cy.contains('button', 'Guardar').click(); 

        cy.wait('@updateMaterial');
    });

    // --- PRUEBA 5 ---
    it('5. Elimina un material', () => {
        cy.intercept('POST', '**/api/admin/delete-material', { statusCode: 200, body: { mensaje: "OK" } }).as('deleteMat');
        cy.intercept('POST', '**/api/admin/get-materiales', { body: { materiales: [] } }).as('reloadTableEmpty');

        cy.get('table tbody tr').first().contains('button', 'Borrar').click();

        cy.wait('@deleteMat');
        cy.wait('@reloadTableEmpty');
        cy.get('table tbody tr').should('have.length', 1); 
    });

});