const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    // ---- Configuración Principal para E2E ----

    // 1. URL base de tu aplicación React
    //    Asegúrate de que tu app esté corriendo en este puerto (npm start)
    baseUrl: 'http://localhost:3000',

    // 2. (Opcional) Define dónde están tus archivos de prueba
    //    Por defecto, Cypress busca en 'cypress/e2e/**/*.cy.{js,jsx,ts,tsx}'
    //    Si tus archivos están ahí, esta línea no es estrictamente necesaria.
    specPattern: 'cypress/e2e/**/*.cy.js',

    // 3. (Opcional) Configura el tamaño de la ventana del navegador
    viewportWidth: 1280,
    viewportHeight: 720,

    // 4. (Opcional) Deshabilita videos si no los necesitas (acelera un poco)
    video: false,

    // 5. Función setupNodeEvents (generalmente vacía al principio)
    //    Se usa para plugins o tareas más avanzadas.
    setupNodeEvents(on, config) {
      // implement node event listeners here
      // (Puedes dejarla vacía por ahora)
    },
  },

  // (Opcional) Otras configuraciones globales de Cypress
  // defaultCommandTimeout: 5000, // Aumenta si tus comandos fallan por tiempo
});