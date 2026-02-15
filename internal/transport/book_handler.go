package transport

import (
	"LibrosElectronicosGolang/internal/model"
	"LibrosElectronicosGolang/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// BookHandler actúa como la capa de transporte/comunicación.
// Se encarga de recibir peticiones HTTP, parsear JSON y devolver respuestas al cliente.
type BookHandler struct {
	service *service.Service
}

// New crea el handler inyectando el servicio que contiene la lógica de negocio.
func New(s *service.Service) *BookHandler {
	return &BookHandler{
		service: s,
	}
}

// HandleBooks es el controlador (multiplexor) para la ruta general '/books'.
// Gestiona peticiones GET (listar) y POST (crear).
func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		libros, err := h.service.ObtenTodosLosLibros()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Formatea y envía la respuesta como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(libros)

	case http.MethodPost:
		var libro model.Libro
		// Decodifica el cuerpo de la petición (JSON) hacia el struct Libro
		if err := json.NewDecoder(r.Body).Decode(&libro); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		created, err := h.service.CrearLibro(libro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated) // HTTP 201
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)

	default:
		// Si usan un método no soportado (ej. PATCH)
		http.Error(w, "Metodo no disponible", http.StatusMethodNotAllowed)
	}
}

// HandleBookByID gestiona rutas con parámetros dinámicos, ej: '/books/1'.
// Soporta GET (buscar uno), PUT (actualizar) y DELETE (borrar).
func (h *BookHandler) HandleBookByID(w http.ResponseWriter, r *http.Request) {
	// Limpieza de URL para extraer únicamente el número de ID
	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		libro, err := h.service.ObtenLibroPorID(id)
		if err != nil {
			http.Error(w, "No lo encontramos", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(libro)

	case http.MethodPut:
		var libro model.Libro
		if err := json.NewDecoder(r.Body).Decode(&libro); err != nil {
			http.Error(w, "input invalido", http.StatusBadRequest)
			return
		}

		update, err := h.service.UpdateAlLibro(id, libro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(update)

	case http.MethodDelete:
		err := h.service.RemoverLibro(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent) // HTTP 204: Borrado exitoso sin contenido a devolver

	default:
		http.Error(w, "metodo no disponible", http.StatusMethodNotAllowed)
	}
}
