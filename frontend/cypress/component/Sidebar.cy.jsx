
import React from 'react';
import { MemoryRouter } from 'react-router-dom';


import { AuthContext } from '../../src/context/AuthContext';
import Sidebar from '../../src/Sidebar';


const mountWithProviders = (
    component,
    authValue,
    initialRoute = ['/admin']
) => {
    return cy.mount(
        <AuthContext.Provider value={authValue}>
            <MemoryRouter initialEntries={initialRoute}>
                {component}
            </MemoryRouter>
        </AuthContext.Provider>
    );
};


// --- INICIO DE LAS PRUEBAS ---
describe('<Sidebar /> Component Test (en cypress/component)', () => {

    it('1. Debe mostrar TODOS los links para un rol "admin"', () => {
        const adminUser = {
            userRole: 'admin',
            logout: cy.stub().as('logoutStub')
        };
        mountWithProviders(<Sidebar />, adminUser);

        cy.contains('Portafolio de Proyectos').should('be.visible');
        cy.contains('Perfiles de usuarios').should('be.visible');
        cy.contains('Logger de eventos').should('be.visible');
    });


    it('2. Debe OCULTAR "Logger de eventos" para un rol "gerente"', () => {
        const gerenteUser = {
            userRole: 'gerente',
            logout: cy.stub().as('logoutStub')
        };
        mountWithProviders(<Sidebar />, gerenteUser);

        cy.contains('Portafolio de Proyectos').should('be.visible');
        cy.contains('Perfiles de usuarios').should('be.visible');
        cy.contains('Logger de eventos').should('not.exist');
    });


    it('3. Debe resaltar el link activo basado en la URL', () => {
        const adminUser = {
            userRole: 'admin',
            logout: cy.stub().as('logoutStub')
        };

        // Simulamos que la URL actual es '/admin/usuarios'
        mountWithProviders(<Sidebar />, adminUser, ['/admin/usuarios']);

        cy.contains('Perfiles de usuarios').should('have.class', 'active');
        cy.contains('Portafolio de Proyectos').should('not.have.class', 'active');
    });


    it('4. Debe llamar a la función logout al hacer clic en "Cerrar Sesión"', () => {
        const adminUser = {
            userRole: 'admin',
            logout: cy.stub().as('logoutStub')
        };
        mountWithProviders(<Sidebar />, adminUser);

        cy.contains('Cerrar Sesión').click();
        cy.get('@logoutStub').should('have.been.calledOnce');
    });

});