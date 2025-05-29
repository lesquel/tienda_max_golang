package repositories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restaurante-crud/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type VentaRepository struct {
	db *sql.DB
}

func NewVentaRepository(db *sql.DB) *VentaRepository {
	return &VentaRepository{db: db}
}

func (repo *VentaRepository) GetAllVentas(w http.ResponseWriter, r *http.Request) {
	rows, err := repo.db.Query("SELECT id, cliente_id, sucursal_id, fecha, estado_sincronizacion FROM ventas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ventas []models.Venta
	for rows.Next() {
		var v models.Venta
		err := rows.Scan(&v.ID, &v.ClienteID, &v.SucursalID, &v.Fecha, &v.EstadoSincronizacion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ventas = append(ventas, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ventas)
}

func (repo *VentaRepository) GetVenta(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var v models.Venta
	err := repo.db.QueryRow("SELECT id, cliente_id, sucursal_id, fecha, estado_sincronizacion FROM ventas WHERE id = $1", id).
		Scan(&v.ID, &v.ClienteID, &v.SucursalID, &v.Fecha, &v.EstadoSincronizacion)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (repo *VentaRepository) CreateVenta(w http.ResponseWriter, r *http.Request) {
	var v models.Venta
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v.ID = uuid.New().String()
	v.Fecha = time.Now()

	_, err = repo.db.Exec("INSERT INTO ventas (id, cliente_id, sucursal_id, fecha) VALUES ($1, $2, $3, $4)",
		v.ID, v.ClienteID, v.SucursalID, v.Fecha)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(v)
}

func (repo *VentaRepository) UpdateVenta(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var v models.Venta
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repo.db.Exec("UPDATE ventas SET cliente_id = $1, sucursal_id = $2 WHERE id = $3",
		v.ClienteID, v.SucursalID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (repo *VentaRepository) DeleteVenta(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := repo.db.Exec("DELETE FROM ventas WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}