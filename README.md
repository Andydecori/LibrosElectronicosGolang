# Sistema de Gestión de Biblioteca Digital

**Nombre del proyecto:** Sistema de gestión de libros electrónicos
**Integrantes:** Andrés Guerrero
**Materia / Docente:** Programación Orientada A Objetos - Milton Palacios

## Objetivos del Sistema
El objetivo central es desplegar una API REST que exponga un ciclo CRUD completo (Crear, Leer, Actualizar y Eliminar). Como valor agregado, el sistema:
* Integra una interfaz web alojada en una capa estática, permitiendo una interacción fluida para usuarios finales.
* Funciona como una herramienta técnica capaz de manejar errores de forma inteligente (como la validación preventiva de registros inexistentes).
* Ofrece vistas HTML limpias para navegación y expone servicios web estructurados en formato JSON para su consumo externo.

Este proyecto es una aplicación web completa desarrollada para la asignatura de Programación, integrando un backend robusto en **Go**, una base de datos **SQLite** y una interfaz dinámica en **HTML/JS**.

## Funcionalidades Implementadas (8 Servicios Web)
El sistema expone 8 servicios web clave que cubren el ciclo de vida de la gestión de información:

1. **Crear Libro (`POST /books`):** Registra un nuevo ejemplar en la base de datos.
2. **Listar Catálogo (`GET /books`):** Obtiene todos los libros disponibles.
3. **Buscar por ID (`GET /books/{id}`):** Recupera información detallada de un libro específico.
4. **Actualizar Datos (`PUT /books/{id}`):** Modifica la información de un registro existente.
5. **Eliminar Registro (`DELETE /books/{id}`):** Borra permanentemente un libro del sistema.
6. **Búsqueda por Autor (`GET /books?author={nombre}`):** Filtra el catálogo dinámicamente según el nombre del autor mediante Query Params.
7. **Servicio Analítico (`GET /books/stats`):** Generación de estadísticas en tiempo real (conteo total) sobre el catálogo.
8. **Monitor de Salud (`GET /health`):** Endpoint de diagnóstico para verificar el estado operativo del servidor y la base de datos.

## Tecnologías Utilizadas
* **Lenguaje:** Go (Golang) 
* **Base de Datos:** SQLite3 para persistencia ligera y eficiente.
* **Arquitectura:** Clean Architecture (Capas: Store, Service, Transport, Model).
* **Frontend:** HTML5, CSS (Bootstrap) y JavaScript (Fetch API).

## Estructura del Proyecto
```text
LibrosElectronicosGolang/
├── internal/
│   ├── model/      # Definición de estructuras (Entidad Libro)
│   ├── service/    # Lógica de negocio y validaciones preventivas
│   ├── store/      # Consultas SQL, persistencia y encapsulación
│   └── transport/  # Manejadores HTTP (Handlers) y enrutamiento
├── static/         # Interfaz web (Frontend)
├── books.db        # Base de datos local autogenerada
└── main.go         # Punto de entrada e inyección de dependencias
```
## Cómo Ejecutar

1. Clona el repositorio.
2. Abre la terminal en la carpeta raíz.
3. Ejecuta el comando:
   ```bash
   go run main.go
   ```
4. Abre tu navegador en: http://localhost:8080

## Visión del Futuro
Actualmente, este proyecto representa una base monolítica sólida. A futuro, visualizo la evolución de esta herramienta hacia una **arquitectura de microservicios en la nube**, lo que permitirá gestionar múltiples bibliotecas de forma simultánea y altamente escalable. 

Además, planeo integrar herramientas de análisis de datos e **Inteligencia Artificial**. Esto transformará el sistema actual en una plataforma educativa inteligente, capaz de analizar patrones de lectura y utilizar algoritmos avanzados para recomendar libros personalizados según el historial de cada usuario.
