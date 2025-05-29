package repositories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restaurante-crud/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductoRepository struct {
	db *sql.DB
}

func NewProductoRepository(db *sql.DB) *ProductoRepository {
	return &ProductoRepository{db: db}
}

func (repo *ProductoRepository) GetAllProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := repo.db.Query("SELECT id, nombre, descripcion, precio, activo, estado_sincronizacion FROM productos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		err := rows.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Activo, &p.EstadoSincronizacion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		productos = append(productos, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

func (repo *ProductoRepository) GetProducto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var p models.Producto
	err := repo.db.QueryRow("SELECT id, nombre, descripcion, precio, activo, estado_sincronizacion FROM productos WHERE id = $1", id).
		Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.Precio, &p.Activo, &p.EstadoSincronizacion)
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

func (repo *ProductoRepository) CreateProducto(w http.ResponseWriter, r *http.Request) {
	var p models.Producto
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.ID = uuid.New().String()

	_, err = repo.db.Exec("INSERT INTO productos (id, nombre, descripcion, precio, activo) VALUES ($1, $2, $3, $4, $5)",
		p.ID, p.Nombre, p.Descripcion, p.Precio, p.Activo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (repo *ProductoRepository) UpdateProducto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var p models.Producto
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repo.db.Exec("UPDATE productos SET nombre = $1, descripcion = $2, precio = $3, activo = $4 WHERE id = $5",
		p.Nombre, p.Descripcion, p.Precio, p.Activo, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (repo *ProductoRepository) DeleteProducto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := repo.db.Exec("DELETE FROM productos WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}