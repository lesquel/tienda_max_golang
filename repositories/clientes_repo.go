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

type ClienteRepository struct {
	db *sql.DB
}

func NewClienteRepository(db *sql.DB) *ClienteRepository {
	return &ClienteRepository{db: db}
}

func (repo *ClienteRepository) GetAllClientes(w http.ResponseWriter, r *http.Request) {
	rows, err := repo.db.Query("SELECT id, nombre, email, telefono, direccion, fecha_registro, sucursal_id, estado_sincronizacion FROM clientes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var c models.Cliente
		err := rows.Scan(&c.ID, &c.Nombre, &c.Email, &c.Telefono, &c.Direccion, &c.FechaRegistro, &c.SucursalID, &c.EstadoSincronizacion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clientes = append(clientes, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientes)
}

func (repo *ClienteRepository) GetCliente(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var c models.Cliente
	err := repo.db.QueryRow("SELECT id, nombre, email, telefono, direccion, fecha_registro, sucursal_id, estado_sincronizacion FROM clientes WHERE id = $1", id).
		Scan(&c.ID, &c.Nombre, &c.Email, &c.Telefono, &c.Direccion, &c.FechaRegistro, &c.SucursalID, &c.EstadoSincronizacion)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (repo *ClienteRepository) CreateCliente(w http.ResponseWriter, r *http.Request) {
	var c models.Cliente
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.ID = uuid.New().String()
	c.FechaRegistro = time.Now()

	_, err = repo.db.Exec("INSERT INTO clientes (id, nombre, email, telefono, direccion, fecha_registro, sucursal_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		c.ID, c.Nombre, c.Email, c.Telefono, c.Direccion, c.FechaRegistro, c.SucursalID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (repo *ClienteRepository) UpdateCliente(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var c models.Cliente
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repo.db.Exec("UPDATE clientes SET nombre = $1, email = $2, telefono = $3, direccion = $4, sucursal_id = $5 WHERE id = $6",
		c.Nombre, c.Email, c.Telefono, c.Direccion, c.SucursalID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func (repo *ClienteRepository) DeleteCliente(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := repo.db.Exec("DELETE FROM clientes WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}