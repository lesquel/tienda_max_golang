package repositories

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"restaurante-crud/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductoPromocionRepository struct {
	db *sql.DB
}

func NewProductoPromocionRepository(db *sql.DB) *ProductoPromocionRepository {
	return &ProductoPromocionRepository{db: db}
}

func (repo *ProductoPromocionRepository) GetPromocionesByProducto(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productoId := params["productoId"]

	rows, err := repo.db.Query(`
		SELECT pr.id, pr.nombre, pr.descripcion, pr.descuento_porcentaje, pr.fecha_inicio, pr.fecha_fin, pr.activa, pr.estado_sincronizacion
		FROM promociones pr
		JOIN producto_promocion pp ON pr.id = pp.promocion_id
		WHERE pp.producto_id = $1
	`, productoId)
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

func (repo *ProductoPromocionRepository) GetProductosByPromocion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	promocionId := params["promocionId"]

	rows, err := repo.db.Query(`
		SELECT p.id, p.nombre, p.descripcion, p.precio, p.activo, p.estado_sincronizacion
		FROM productos p
		JOIN producto_promocion pp ON p.id = pp.producto_id
		WHERE pp.promocion_id = $1
	`, promocionId)
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

func (repo *ProductoPromocionRepository) CreateProductoPromocion(w http.ResponseWriter, r *http.Request) {
	var pp struct {
		ProductoID  string `json:"producto_id"`
		PromocionID string `json:"promocion_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&pp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	_, err = repo.db.Exec("INSERT INTO producto_promocion (id, producto_id, promocion_id) VALUES ($1, $2, $3)",
		id, pp.ProductoID, pp.PromocionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":           id,
		"producto_id":  pp.ProductoID,
		"promocion_id": pp.PromocionID,
	})
}

func (repo *ProductoPromocionRepository) DeleteProductoPromocion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := repo.db.Exec("DELETE FROM producto_promocion WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
