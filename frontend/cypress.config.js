const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    // ---- Configuración Principal para E2E ----
    baseUrl: "http://localhost:3000",
    specPattern: "cypress/e2e/**/*.cy.js",
    viewportWidth: 1280,
    viewportHeight: 720,
    video: false,
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
  },

  // (El bloque 'component' se puede borrar)
});