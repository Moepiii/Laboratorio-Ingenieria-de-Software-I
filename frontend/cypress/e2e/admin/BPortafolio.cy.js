
describe('Portafolio de Proyectos (Historia B)', () => {

    // --- PASO 1: Iniciar sesión ---
    beforeEach(() => {
        cy.visit('/login');
        cy.get('input[type="text"]').first().type('admin');
        cy.get('input[type="password"]').first().type('admin123');
        cy.get('button[type="submit"]').first().click();

        cy.url({ timeout: 10000 }).should('include', '/admin');
        cy.contains('Portafolio de Proyectos').click();
        cy.url().should('include', '/admin/proyectos');
    });

    // --- PRUEBA 2: Estado inicial de los botones ---
    it('B.1 Debería mostrar la barra de herramientas y los botones (Datos, Borrar, etc) deshabilitados', () => {
        cy.get('h2').contains('Portafolio de Proyectos').should('be.visible');
        cy.contains('button', /^Crear$/).should('be.visible').and('not.be.disabled');
        cy.contains('a', 'Datos').should('have.css', 'cursor', 'not-allowed');
        cy.contains('button', /^Borrar$/).should('be.disabled');
        cy.contains('button', /^Habilitar$/).should('be.disabled');
        cy.contains('button', /^Cerrar$/).should('be.disabled');
    });

    // --- PRUEBA 3: Funcionalidad del botón "Crear" ---
    it('B.2 Debería mostrar el formulario al hacer clic en "Crear"', () => {
        cy.contains('button', /^Crear$/).click();
        cy.contains('Crear Nuevo Proyecto').should('be.visible');
        cy.contains('button', /^Guardar$/).should('be.visible');
        cy.contains('button', /^Cancelar$/).should('be.visible');
    });

    // --- PRueba 4: Habilitar botones al seleccionar ---
    it('B.3 Debería habilitar los botones de acción al seleccionar un proyecto', () => {
        cy.get('table tbody tr').first().find('input[type="radio"]').click();
        cy.contains('a', 'Datos').should('not.have.css', 'cursor', 'not-allowed');
        cy.contains('button', /^Borrar$/).should('not.be.disabled');
    });

    // --- PRUEBA 5: (MEJORADA) CICLO DE VIDA COMPLETO --- 
    it('B.4 Debería Crear, Cerrar y Habilitar un proyecto', () => {
        const nombreProyecto = 'Proyecto Ciclo-Vida ' + Date.now();

        // --- 1. CREAR ---
        cy.contains('button', /^Crear$/).click();
        cy.contains('label', 'Nombre del Proyecto').next('input').type(nombreProyecto);
        cy.contains('label', 'Fecha Inicio').next('input').type('2025-01-01');
        cy.contains('label', 'Fecha Cierre').next('input').type('2025-12-31');
        cy.contains('button', /^Guardar$/).click();


        cy.contains('tr', nombreProyecto)
            .find('span')
            .should('contain.text', 'Activo'); // 

        // --- 3. PROBAR BOTÓN "CERRAR" ---
        cy.contains('tr', nombreProyecto).find('input[type="radio"]').click();


        cy.contains('button', /^Habilitar$/).should('not.be.disabled');
        cy.contains('button', /^Cerrar$/).should('not.be.disabled');

        // Haz clic en "Cerrar"
        cy.contains('button', /^Cerrar$/).click();

        // --- 4. VERIFICAR ESTADO "cerrado" ---
        cy.contains('tr', nombreProyecto)
            .find('span')
            .should('contain.text', 'cerrado');


        cy.contains('button', /^Habilitar$/).should('not.be.disabled');
        cy.contains('button', /^Cerrar$/).should('be.disabled');

        // Haz clic en "Habilitar"
        cy.contains('button', /^Habilitar$/).click();

        // --- 6. VERIFICAR ESTADO "habilitado" ---
        // El botón "Habilitar" setea el estado a "habilitado" (minúscula)
        cy.contains('tr', nombreProyecto)
            .find('span')
            .should('contain.text', 'habilitado'); 
    });

});