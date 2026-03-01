package main

import (
	"LibrosElectronicosGolang/internal/service"
	"LibrosElectronicosGolang/internal/store"
	"LibrosElectronicosGolang/internal/transport"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 1. Conectar a SQLite: Se establece la conexión con la base de datos local
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Esquema de Base de Datos: Crear tabla si no existe al arrancar
	q := `
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL
		)
	`
	if _, err := db.Exec(q); err != nil {
		log.Fatal("Error creando tabla: ", err)
	}

	// 3. Inyección de Dependencias: Se inicializan las capas de adentro hacia afuera
	bookStore := store.New(db)
	bookService := service.New(bookStore)
	bookHandler := transport.New(bookService)

	// Configurar rutas
	// Primero las rutas fijas y específicas (estadísticas y estado del sistema)
	http.HandleFunc("/health", bookHandler.HealthCheck)
	http.HandleFunc("/books/stats", bookHandler.GetStats)

	// Luego las rutas generales (CRUD)
	http.HandleFunc("/books", bookHandler.HandleBooks)
	http.HandleFunc("/books/", bookHandler.HandleBookByID)

	port := ":8080"

	// 5. Interfaz de Consola para el Desarrollador
	fmt.Println(" ")
	fmt.Printf(" Servidor ejecutándose en http://localhost%s\n\n", port)
	fmt.Println(" API Endpoints disponibles:")
	fmt.Println("-------------------------------------------------------")
	fmt.Println(" MÉTODO   RUTA              DESCRIPCIÓN")
	fmt.Println("-------------------------------------------------------")
	fmt.Println(" GET      /books            Obtener todos los libros")
	fmt.Println(" POST     /books            Crear un nuevo libro")
	fmt.Println(" GET      /books/{id}       Obtener un libro específico")
	fmt.Println(" PUT      /books/{id}       Actualizar un libro")
	fmt.Println(" DELETE   /books/{id}       Eliminar un libro")
	fmt.Println(" GET      /books?author=x   Buscar libros por autor")
	fmt.Println(" GET      /books/stats      Estadisticas del catalogo")
	fmt.Println(" GET      /health           Estado del sistema")
	fmt.Println("--------------------------------------------------------")

	// 6. Servir archivos estáticos (HTML, CSS, JS) desde la carpeta /static
	// Esto permite que al entrar a localhost:8080 se cargue automáticamente la interfaz visual.
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs) // El index.html será la página de inicio

	// 7. Iniciar Servidor Web
	log.Fatal(http.ListenAndServe(port, nil))
}
