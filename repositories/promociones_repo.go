package repositories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restaurante-crud/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PromocionRepository struct {
	db *sql.DB
}

func NewPromocionRepository(db *sql.DB) *PromocionRepository {
	return &PromocionRepository{db: db}
}

func (repo *PromocionRepository) GetAllPromociones(w http.ResponseWriter, r *http.Request) {
	rows, err := repo.db.Query("SELECT id, nombre, descripcion, descuento_porcentaje, fecha_inicio, fecha_fin, activa, estado_sincronizacion FROM promociones")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var promociones []models.Promocion
	for rows.Next() {
		var p models.Promocion
		err := rows.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.DescuentoPorcentaje, &p.FechaInicio, &p.FechaFin, &p.Activa, &p.EstadoSincronizacion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		promociones = append(promociones, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(promociones)
}

func (repo *PromocionRepository) GetPromocion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var p models.Promocion
	err := repo.db.QueryRow("SELECT id, nombre, descripcion, descuento_porcentaje, fecha_inicio, fecha_fin, activa, estado_sincronizacion FROM promociones WHERE id = $1", id).
		Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.DescuentoPorcentaje, &p.FechaInicio, &p.FechaFin, &p.Activa, &p.EstadoSincronizacion)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (repo *PromocionRepository) CreatePromocion(w http.ResponseWriter, r *http.Request) {
	var p models.Promocion
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.ID = uuid.New().String()

	_, err = repo.db.Exec("INSERT INTO promociones (id, nombre, descripcion, descuento_porcentaje, fecha_inicio, fecha_fin, activa) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		p.ID, p.Nombre, p.Descripcion, p.DescuentoPorcentaje, p.FechaInicio, p.FechaFin, p.Activa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (repo *PromocionRepository) UpdatePromocion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var p models.Promocion
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repo.db.Exec("UPDATE promociones SET nombre = $1, descripcion = $2, descuento_porcentaje = $3, fecha_inicio = $4, fecha_fin = $5, activa = $6 WHERE id = $7",
		p.Nombre, p.Descripcion, p.DescuentoPorcentaje, p.FechaInicio, p.FechaFin, p.Activa, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (repo *PromocionRepository) DeletePromocion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := repo.db.Exec("DELETE FROM promociones WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
