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
		// NUEVO: Revisamos si la URL contiene el parámetro 'author' (Query Param)
		authorParam := r.URL.Query().Get("author")

		// Si el parámetro existe, ejecutamos el Servicio 6 (Búsqueda por Autor)
		if authorParam != "" {
			libros, err := h.service.BuscarPorAutor(authorParam)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(libros)
			return // Retornamos para salir de la función y no ejecutar la búsqueda general
		}

		// Si no hay parámetro 'author', ejecutamos el Servicio 2 normal (Listar Todos)
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

// GetStats maneja la ruta GET /books/stats
// Se comunica con el servicio para obtener el total de libros y lo devuelve en formato JSON.
func (h *BookHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	total, err := h.service.GetTotalBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_libros": total,
		"mensaje":      "Estadisticas recuperadas con exito",
	})
}

// HealthCheck maneja la ruta GET /health
// Servicio vital para monitoreo: verifica que el servidor esté encendido y respondiendo.
// Devuelve un JSON con el estado operativo del sistema y la versión.
func (h *BookHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "Servidor Operativo ",
		"database": "SQLite Conectada ",
		"version":  "1.0.0",
	})
}
