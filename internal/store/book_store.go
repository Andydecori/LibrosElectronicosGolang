package store

import (
	"LibrosElectronicosGolang/internal/model"
	"database/sql"
	"fmt"
)

// Store define el contrato (interfaz para las operaciones de la base de datos.
// Esto permite la abstracción y facilita el testeo o cambio de base de datos a futuro.
type Store interface {
	GetALL() ([]*model.Libro, error)
	GetByID(id int) (*model.Libro, error)
	Create(libro *model.Libro) (*model.Libro, error)
	Update(id int, libro *model.Libro) (*model.Libro, error)
	Delete(id int) error
	CountLibros() (int, error)
	GetByAuthor(author string) ([]*model.Libro, error)
}

// store es la implementación concreta de la interfaz Store.
// Nota: Se usa minúscula para aplicar el principio de ENCAPSULACIÓN,
// ocultando la conexión directa a la BD desde otros paquetes.
type store struct {
	db *sql.DB
}

// New crea y retorna una nueva instancia de la capa Store.
func New(db *sql.DB) Store {
	return &store{db: db}
}

// GetALL ejecuta una consulta para traer todos los registros de libros de SQLite.
func (s *store) GetALL() ([]*model.Libro, error) {
	q := `SELECT id, title, author FROM books`
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Asegura que se liberen los recursos de la BD al terminar

	var libros []*model.Libro
	for rows.Next() {
		b := model.Libro{}
		// Scan mapea las columnas de la BD a los campos del struct Libro
		if err := rows.Scan(&b.ID, &b.Titulo, &b.Autor); err != nil {
			return nil, err
		}
		libros = append(libros, &b)
	}
	return libros, nil
}

// GetByID busca un libro específico utilizando su identificador único.
func (s *store) GetByID(id int) (*model.Libro, error) {
	q := `SELECT id, title, author FROM books WHERE id = ?`
	b := model.Libro{}
	// QueryRow se usa porque esperamos exactamente un solo resultado
	err := s.db.QueryRow(q, id).Scan(&b.ID, &b.Titulo, &b.Autor)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// Create inserta un nuevo registro de libro y le asigna el ID autogenerado.
func (s *store) Create(libro *model.Libro) (*model.Libro, error) {
	q := `INSERT INTO books (title, author) VALUES (?, ?)`
	resp, err := s.db.Exec(q, libro.Titulo, libro.Autor)
	if err != nil {
		return nil, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}
	libro.ID = int(id) // Actualizamos el modelo con el nuevo ID de la BD
	return libro, nil
}

// Update modifica los datos de un libro existente basándose en su ID.
func (s *store) Update(id int, libro *model.Libro) (*model.Libro, error) {
	q := `UPDATE books SET title = ?, author = ? WHERE id = ?`
	result, err := s.db.Exec(q, libro.Titulo, libro.Autor, id)
	if err != nil {
		return nil, err
	}

	// Validación: Comprobamos si el ID realmente existía en la base de datos
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("el libro con ID %d no existe en la base de datos", id)
	}
	libro.ID = id
	return libro, nil
}

// Delete elimina de forma permanente el registro de un libro mediante su ID.
func (s *store) Delete(id int) error {
	q := `DELETE from books WHERE id = ?`
	result, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	// Validación sobre si se borro algo o el id no existia.
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("imposible eliminar: el libro con ID %d no existe", id)
	}
	return nil
}

// CountLibros devuelve la cantidad total de registros en la tabla
// Ejecuta una consulta SQL de agregación (COUNT) para obtener estadísticas rápidas.
func (s *store) CountLibros() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetByAuthor filtra los libros por autor utilizando una coincidencia parcial (LIKE).
// Permite encontrar coincidencias como "Gabriel" para "Gabriel Garcia Marquez" Por ejemplo.
func (s *store) GetByAuthor(author string) ([]*model.Libro, error) {
	// Usamos LIKE en SQLite y concatenamos '%' para buscar subcadenas
	q := `SELECT id, title, author FROM books WHERE author LIKE ?`
	rows, err := s.db.Query(q, "%"+author+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []*model.Libro
	for rows.Next() {
		b := model.Libro{}
		if err := rows.Scan(&b.ID, &b.Titulo, &b.Autor); err != nil {
			return nil, err
		}
		libros = append(libros, &b)
	}
	return libros, nil
}
