/* eslint-disable no-undef */

// ⭐️ CAMBIO: Título actualizado para reflejar el rol de Gerente
describe('Gestión de Proyectos - Gerente', () => {

    beforeEach(() => {
        cy.visit('/login');

        // ⭐️ CAMBIO: Inicia sesión como 'gerente'
        // (¡Asegúrate de que 'gerente' y 'gerente123' sean credenciales válidas en tu sistema!)
        cy.get('input[type="text"]').first().type('fff');
        cy.get('input[type="password"]').first().type('fff123');
        cy.get('button[type="submit"]').first().click();

        // Esperar redirección (un gerente también debe ir a /admin)
        cy.url({ timeout: 10000 }).should('include', '/admin');

        // 🟢 NAVEGACIÓN FLEXIBLE A PROYECTOS (Sin cambios) 🟢
        cy.get('body').then(($body) => {
            const texto = $body.text().toLowerCase();
            const textosProyectos = [
                'portafolio de proyectos', // Texto de tu app
                'gestión de proyectos',
                'proyectos',
            ];

            let encontrado = false;

            for (const textoBuscar of textosProyectos) {
                if (texto.includes(textoBuscar)) {
                    console.log(`✅ Encontrado: ${textoBuscar}`);
                    cy.contains(new RegExp(textoBuscar, 'i'), { timeout: 10000 })
                        .should('be.visible')
                        .click({ force: true });
                    encontrado = true;
                    break;
                }
            }

            if (!encontrado) {
                cy.log('⚠️ No se encontró texto específico, buscando por link "proyecto"');
                cy.get('a, button').contains(/proyecto|project/i, { timeout: 10000 })
                    .should('be.visible')
                    .click({ force: true });
            }
        });

        // Verificar que estamos en la sección correcta
        cy.url({ timeout: 10000 }).should(($url) => {
            // Un gerente debería ser redirigido a /admin/proyectos
            expect($url).to.match(/\/admin\/proyectos/i);
        });
    });

    // --- PRUEBA 1: LISTAR PROYECTOS (Sin cambios) ---
    it('1. Debería listar todos los proyectos', () => {
        cy.get('table, .table, tbody', { timeout: 10000 })
            .should('exist');
        cy.get('body').then(($body) => {
            const texto = $body.text().toLowerCase();
            if (texto.includes('no hay') || texto.includes('sin proyectos') || texto.includes('empty')) {
                cy.log('ℹ️ No hay proyectos existentes, continuando...');
            } else {
                cy.get('tbody tr').should('have.length.greaterThan', 0);
                cy.log('✅ Proyectos encontrados en la tabla');
            }
        });
    });

    // --- PRUEBA 2: CREAR NUEVO PROYECTO (Sin cambios) ---
    it('2. Debería crear un nuevo proyecto', () => {
        const projectName = `Proyecto Gerente ${Date.now()}`;

        // 🟢 BUSCAR BOTÓN DE CREAR PROYECTO (Tu lógica flexible) 🟢
        cy.get('body').then(($body) => {
            const textosCrear = [
                'crear proyecto', // Este es el texto en tu app
                'crear nuevo proyecto',
                'nuevo proyecto',
                'crear' // Botón "Crear" de la toolbar
            ];
            let encontrado = false;
            for (const textoBuscar of textosCrear) {
                if ($body.text().toLowerCase().includes(textoBuscar)) {
                    cy.contains(new RegExp(textoBuscar, 'i'), { timeout: 10000 })
                        .should('be.visible')
                        .click({ force: true });
                    encontrado = true;
                    break;
                }
            }
            if (!encontrado) {
                cy.get('button').contains(/crear|nuevo|add|new/i, { timeout: 10000 })
                    .should('be.visible').first().click({ force: true });
            }
        });

        // Esperar que cargue el formulario
        cy.get('form, input[type="text"]', { timeout: 10000 }).should('be.visible');

        // 🟢 LLENAR CAMPOS DEL FORMULARIO (Tu lógica flexible) 🟢
        cy.get('body').then(($body) => {
            // Campo nombre (basado en el placeholder de tu app: "Ej: Proyecto Titán")
            // O podemos usar el ID que establecimos: #nombre
            cy.get('#nombre').should('be.visible').type(projectName);

            // Campo fecha inicio (basado en ID: #fecha_inicio)
            cy.get('#fecha_inicio').should('be.visible').type('2024-01-20');

            // Campo fecha cierre (basado en ID: #fecha_cierre)
            cy.get('#fecha_cierre').should('be.visible').type('2024-12-31');
        });

        // 🔵 Hacer clic en botón de guardar (Busca "Crear Proyecto" dentro del form)
        cy.get('form').contains('button', /crear proyecto/i, { timeout: 10000 })
            .should('be.visible')
            .click({ force: true });

        // 🟢 Verificar éxito (Tu lógica flexible) 🟢
        cy.wait(2000);
        cy.get('body', { timeout: 10000 }).then(($body) => {
            const texto = $body.text().toLowerCase();
            if (texto.includes('éxito') || texto.includes('creado')) {
                cy.log('✅ Proyecto creado exitosamente');
            } else {
                cy.log('⚠️ No se encontró mensaje visible, verificando en tabla...');
            }
        });

        // 🧩 Confirmar que el proyecto aparece en la tabla
        cy.contains('td', projectName, { timeout: 10000 }).should('exist');
    });

});