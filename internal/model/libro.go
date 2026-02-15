package model

// Libro representa la entidad principal del dominio.
// Se utiliza para transferir datos entre la base de datos (Store) y la respuesta HTTP (JSON).
type Libro struct {
	ID     int    `json:"id"`     // Identificador único autoincremental
	Titulo string `json:"title"`  // Título de la obra (Obligatorio)
	Autor  string `json:"author"` // Creador de la obra (Obligatorio)
}
