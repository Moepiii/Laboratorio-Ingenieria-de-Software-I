import React from 'react';

// Este componente recibirá `apiCall` para hacer peticiones al backend
const PortafolioProyectos = ({ apiCall }) => {
  
  // Aquí iría la lógica para (con useEffect y useState):
  // 1. Cargar los proyectos (apiCall('GET', '/api/admin/proyectos', null))
  // 2. Manejar el estado de los inputs de búsqueda, agregar, etc.

  return (
    <div>
      <h2>Portafolio de Proyectos</h2>
      
      {/* Aquí puedes añadir los botones de tu maqueta */}
      <div style={{ margin: '1rem 0' }}>
        <button>Agregar Nuevo Proyecto</button>
        <button style={{ marginLeft: '10px' }}>Modificar Proyecto</button>
        <button style={{ marginLeft: '10px' }}>Eliminar Proyecto</button>
        <input type="text" placeholder="Buscar Proyecto..." style={{ marginLeft: '20px' }} />
      </div>

      {/* Aquí iría la tabla de proyectos */}
      <p>Aquí se mostrará la tabla con los proyectos...</p>
      {/* Ejemplo de cómo se vería la tabla (necesitas datos reales) */}
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr style={{ backgroundColor: '#eee' }}>
            <th style={{ padding: '8px', border: '1px solid #ddd' }}>ID</th>
            <th style={{ padding: '8px', border: '1px solid #ddd' }}>Descripción</th>
            <th style={{ padding: '8px', border: '1px solid #ddd' }}>Inicio</th>
            <th style={{ padding: '8px', border: '1px solid #ddd' }}>Cierre</th>
          </tr>
        </thead>
        <tbody>
          {/* Aquí mapearías los datos de los proyectos */}
          <tr>
            <td style={{ padding: '8px', border: '1px solid #ddd' }}>1</td>
            <td style={{ padding: '8px', border: '1px solid #ddd' }}>Proyecto 1</td>
            <td style={{ padding: '8px', border: '1px solid #ddd' }}>01/11/2022</td>
            <td style={{ padding: '8px', border: '1px solid #ddd' }}>01/01/2023</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

export default PortafolioProyectos;