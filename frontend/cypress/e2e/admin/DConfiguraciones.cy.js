

describe('Módulo de Configuraciones (Historia D)', () => {

    const randomId = Math.floor(Math.random() * 1000);
    const nuevaUnidad = {
        nombre: `Saco Test ${randomId}`,
        abreviatura: `sc${randomId}`,
        dimension: '50.5',
        tipo: 'Peso'
    };

    // --- PASO PREVIO: Login y Navegación ---
    beforeEach(() => {
        cy.visit('/login');
        cy.get('input[id="username"]').should('be.visible').type('admin');
        cy.get('input[id="password"]').should('be.visible').type('admin123');
        cy.get('button[type="submit"]').click();

        // Esperar Dashboard
        cy.url({ timeout: 10000 }).should('include', '/admin');

        // Ir a Configuraciones
        cy.get('a[href="/admin/configuraciones"]').should('be.visible').click();
        cy.url().should('include', '/admin/configuraciones');

        // Esperar tabla y seleccionar el primer proyecto para desplegar el menú
        cy.get('tbody tr', { timeout: 10000 }).should('have.length.greaterThan', 0);
        cy.get('tbody tr').first().click();

        // Verificar que el menú lateral se expandió
        cy.contains('a', 'Unidades de Medida', { timeout: 5000 }).should('be.visible');
    });

    // --- PRUEBA 1: CREACIÓN (Ya verificada) ---
    it('D.1 Debería crear una nueva Unidad de Medida exitosamente', () => {
        cy.contains('a', 'Unidades de Medida').click();
        cy.contains('h2', 'Unidades de Medida').should('be.visible');

        // Abrir Modal (Botón al lado del título)
        cy.contains('h2', 'Unidades de Medida').parent().find('button').click();

        // Llenar formulario
        cy.get('input[name="nombre"]').should('be.visible').type(nuevaUnidad.nombre);
        cy.get('input[name="abreviatura"]').type(nuevaUnidad.abreviatura);
        cy.get('input[name="dimension"]').type(nuevaUnidad.dimension);
        cy.get('select[name="tipo"]').select(nuevaUnidad.tipo);

        // Guardar
        cy.get('form').contains('button', 'Guardar').click();

        // Verificar
        cy.get('input[name="nombre"]').should('not.exist');
        cy.contains('tr', nuevaUnidad.nombre).should('be.visible');
    });


    it('D.2 Debería navegar correctamente entre los submódulos del menú', () => {

        // 1. Ir a Labores Agronómicas
        cy.log('➡️ Navegando a Labores...');
        cy.contains('a', 'Labores Agronómicas').click();

        // Verificaciones
        cy.url().should('include', '/labores'); // La URL debe cambiar
        cy.get('h2').should('contain', 'Labores Agronómicas'); // El título debe coincidir

        // 2. Ir a Equipos e Implementos
        cy.log('➡️ Navegando a Equipos...');
        cy.contains('a', 'Equipos e Implementos').click();

        // Verificaciones
        cy.url().should('include', '/equipos');
        cy.get('h2').should('contain', 'Equipos e Implementos');

        // 3. Ir a Unidades de Medida
        cy.log('➡️ Navegando a Unidades...');
        cy.contains('a', 'Unidades de Medida').click();

        // Verificaciones
        cy.url().should('include', '/unidades');
        cy.get('h2').should('contain', 'Unidades de Medida');
    });

});