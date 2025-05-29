package repositories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restaurante-crud/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type VentaDetalleRepository struct {
	db *sql.DB
}

func NewVentaDetalleRepository(db *sql.DB) *VentaDetalleRepository {
	return &VentaDetalleRepository{db: db}
}

func (repo *VentaDetalleRepository) GetDetallesByVenta(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ventaId := params["ventaId"]

	rows, err := repo.db.Query("SELECT id, venta_id, producto_id, cantidad, precio_unitario, estado_sincronizacion FROM venta_detalle WHERE venta_id = $1", ventaId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var detalles []models.VentaDetalle
	for rows.Next() {
		var d models.VentaDetalle
		err := rows.Scan(&d.ID, &d.VentaID, &d.ProductoID, &d.Cantidad, &d.PrecioUnitario, &d.EstadoSincronizacion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		detalles = append(detalles, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detalles)
}

func (repo *VentaDetalleRepository) GetVentaDetalle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var d models.VentaDetalle
	err := repo.db.QueryRow("SELECT id, venta_id, producto_id, cantidad, precio_unitario, estado_sincronizacion FROM venta_detalle WHERE id = $1", id).
		Scan(&d.ID, &d.VentaID, &d.ProductoID, &d.Cantidad, &d.PrecioUnitario, &d.EstadoSincronizacion)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (repo *VentaDetalleRepository) CreateVentaDetalle(w http.ResponseWriter, r *http.Request) {
	var d models.VentaDetalle
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.ID = uuid.New().String()

	// Obtener el precio del producto
	var precio float64
	err = repo.db.QueryRow("SELECT precio FROM productos WHERE id = $1", d.ProductoID).Scan(&precio)
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusBadRequest)
		return
	}

	// Usar el precio del producto si no se especifica otro
	if d.PrecioUnitario == 0 {
		d.PrecioUnitario = precio
	}

	_, err = repo.db.Exec("INSERT INTO venta_detalle (id, venta_id, producto_id, cantidad, precio_unitario) VALUES ($1, $2, $3, $4, $5)",
		d.ID, d.VentaID, d.ProductoID, d.Cantidad, d.PrecioUnitario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

func (repo *VentaDetalleRepository) UpdateVentaDetalle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var d models.VentaDetalle
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repo.db.Exec("UPDATE venta_detalle SET producto_id = $1, cantidad = $2, precio_unitario = $3 WHERE id = $4",
		d.ProductoID, d.Cantidad, d.PrecioUnitario, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (repo *VentaDetalleRepository) DeleteVentaDetalle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := repo.db.Exec("DELETE FROM venta_detalle WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}