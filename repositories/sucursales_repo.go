package repositories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restaurante-crud/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SucursalRepository struct {
	db *sql.DB
}

func NewSucursalRepository(db *sql.DB) *SucursalRepository {
	return &SucursalRepository{db: db}
}

func (repo *SucursalRepository) GetAllSucursales(w http.ResponseWriter, r *http.Request) {
	rows, err := repo.db.Query("SELECT id, nombre, direccion, telefono, estado_sincronizacion FROM sucursales")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sucursales []models.Sucursal
	for rows.Next() {
		var s models.Sucursal
		err := rows.Scan(&s.ID, &s.Nombre, &s.Direccion, &s.Telefono, &s.EstadoSincronizacion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sucursales = append(sucursales, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sucursales)
}

func (repo *SucursalRepository) GetSucursal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var s models.Sucursal
	err := repo.db.QueryRow("SELECT id, nombre, direccion, telefono, estado_sincronizacion FROM sucursales WHERE id = $1", id).
		Scan(&s.ID, &s.Nombre, &s.Direccion, &s.Telefono, &s.EstadoSincronizacion)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (repo *SucursalRepository) CreateSucursal(w http.ResponseWriter, r *http.Request) {
	var s models.Sucursal
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.ID = uuid.New().String()

	_, err = repo.db.Exec("INSERT INTO sucursales (id, nombre, direccion, telefono) VALUES ($1, $2, $3, $4)",
		s.ID, s.Nombre, s.Direccion, s.Telefono)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (repo *SucursalRepository) UpdateSucursal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var s models.Sucursal
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repo.db.Exec("UPDATE sucursales SET nombre = $1, direccion = $2, telefono = $3 WHERE id = $4",
		s.Nombre, s.Direccion, s.Telefono, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (repo *SucursalRepository) DeleteSucursal(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := repo.db.Exec("DELETE FROM sucursales WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
