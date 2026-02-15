package service

import (
	"LibrosElectronicosGolang/internal/model"
	"LibrosElectronicosGolang/internal/store"
	"errors"
)

// Service actúa como la capa de lógica de negocio del sistema.
// Encapsula la interfaz Store para interactuar con los datos sin saber qué BD se usa.
type Service struct {
	store store.Store
}

// New inicializa el servicio inyectando la dependencia de Store.
func New(s store.Store) *Service {
	return &Service{
		store: s,
	}
}

// ObtenTodosLosLibros solicita a la capa store el catálogo completo de libros.
func (s *Service) ObtenTodosLosLibros() ([]*model.Libro, error) {
	libros, err := s.store.GetALL()
	if err != nil {
		return nil, err
	}
	return libros, nil
}

// ObtenLibroPorID pasa la solicitud de búsqueda específica a la capa inferior.
func (s *Service) ObtenLibroPorID(id int) (*model.Libro, error) {
	return s.store.GetByID(id)
}

// CrearLibro valida los datos de entrada antes de intentar guardarlos.
// Retorna un error explícito si falla una regla de negocio (ej. título vacío).
func (s *Service) CrearLibro(libro model.Libro) (*model.Libro, error) {
	if libro.Titulo == "" {
		return nil, errors.New("necesitamos el titulo")
	}
	return s.store.Create(&libro)
}

// UpdateAlLibro asegura que los datos a modificar cumplan con las reglas de negocio.
func (s *Service) UpdateAlLibro(id int, libro model.Libro) (*model.Libro, error) {
	if libro.Titulo == "" {
		return nil, errors.New("necesitamos el titulo")
	}
	return s.store.Update(id, &libro)
}

// RemoverLibro solicita la eliminación física de un registro en la base de datos.
func (s *Service) RemoverLibro(id int) error {
	return s.store.Delete(id)
}
