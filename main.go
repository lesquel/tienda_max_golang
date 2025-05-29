package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"restaurante-crud/repositories"
	"restaurante-crud/tables"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

func main() {
	// Conexión a la base de datos
	db, err := sql.Open("sqlite", "tiendamax.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables.CrearTablaSucursales(db)
	tables.CrearTablaClientes(db)
	tables.CrearTablaProductos(db)
	tables.CrearTablaProductoSucursal(db)
	tables.CrearTablaVentas(db)
	tables.CrearTablaVentaDetalle(db)
	tables.CrearTablaPromociones(db)
	tables.CrearTablaProductoPromocion(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conexión exitosa a PostgreSQL")

	// Inicializar repositorios
	sucursalRepo := repositories.NewSucursalRepository(db)
	productoRepo := repositories.NewProductoRepository(db)
	clienteRepo := repositories.NewClienteRepository(db)
	ventaRepo := repositories.NewVentaRepository(db)
	ventaDetalleRepo := repositories.NewVentaDetalleRepository(db)
	promocionRepo := repositories.NewPromocionRepository(db)
	productoPromocionRepo := repositories.NewProductoPromocionRepository(db)

	// Configurar rutas
	r := mux.NewRouter()

	// Rutas para Sucursales
	r.HandleFunc("/sucursales", sucursalRepo.GetAllSucursales).Methods("GET")
	r.HandleFunc("/sucursales/{id}", sucursalRepo.GetSucursal).Methods("GET")
	r.HandleFunc("/sucursales", sucursalRepo.CreateSucursal).Methods("POST")
	r.HandleFunc("/sucursales/{id}", sucursalRepo.UpdateSucursal).Methods("PUT")
	r.HandleFunc("/sucursales/{id}", sucursalRepo.DeleteSucursal).Methods("DELETE")

	// Rutas para Productos
	r.HandleFunc("/productos", productoRepo.GetAllProductos).Methods("GET")
	r.HandleFunc("/productos/{id}", productoRepo.GetProducto).Methods("GET")
	r.HandleFunc("/productos", productoRepo.CreateProducto).Methods("POST")
	r.HandleFunc("/productos/{id}", productoRepo.UpdateProducto).Methods("PUT")
	r.HandleFunc("/productos/{id}", productoRepo.DeleteProducto).Methods("DELETE")

	// Rutas para Clientes
	r.HandleFunc("/clientes", clienteRepo.GetAllClientes).Methods("GET")
	r.HandleFunc("/clientes/{id}", clienteRepo.GetCliente).Methods("GET")
	r.HandleFunc("/clientes", clienteRepo.CreateCliente).Methods("POST")
	r.HandleFunc("/clientes/{id}", clienteRepo.UpdateCliente).Methods("PUT")
	r.HandleFunc("/clientes/{id}", clienteRepo.DeleteCliente).Methods("DELETE")

	// Rutas para Ventas
	r.HandleFunc("/ventas", ventaRepo.GetAllVentas).Methods("GET")
	r.HandleFunc("/ventas/{id}", ventaRepo.GetVenta).Methods("GET")
	r.HandleFunc("/ventas", ventaRepo.CreateVenta).Methods("POST")
	r.HandleFunc("/ventas/{id}", ventaRepo.UpdateVenta).Methods("PUT")
	r.HandleFunc("/ventas/{id}", ventaRepo.DeleteVenta).Methods("DELETE")

	// Rutas para VentaDetalle
	r.HandleFunc("/ventas/{ventaId}/detalles", ventaDetalleRepo.GetDetallesByVenta).Methods("GET")
	r.HandleFunc("/ventas/detalles/{id}", ventaDetalleRepo.GetVentaDetalle).Methods("GET")
	r.HandleFunc("/ventas/detalles", ventaDetalleRepo.CreateVentaDetalle).Methods("POST")
	r.HandleFunc("/ventas/detalles/{id}", ventaDetalleRepo.UpdateVentaDetalle).Methods("PUT")
	r.HandleFunc("/ventas/detalles/{id}", ventaDetalleRepo.DeleteVentaDetalle).Methods("DELETE")

	// Rutas para Promociones
	r.HandleFunc("/promociones", promocionRepo.GetAllPromociones).Methods("GET")
	r.HandleFunc("/promociones/{id}", promocionRepo.GetPromocion).Methods("GET")
	r.HandleFunc("/promociones", promocionRepo.CreatePromocion).Methods("POST")
	r.HandleFunc("/promociones/{id}", promocionRepo.UpdatePromocion).Methods("PUT")
	r.HandleFunc("/promociones/{id}", promocionRepo.DeletePromocion).Methods("DELETE")

	// Rutas para ProductoPromocion
	r.HandleFunc("/productos/{productoId}/promociones", productoPromocionRepo.GetPromocionesByProducto).Methods("GET")
	r.HandleFunc("/promociones/{promocionId}/productos", productoPromocionRepo.GetProductosByPromocion).Methods("GET")
	r.HandleFunc("/producto-promocion", productoPromocionRepo.CreateProductoPromocion).Methods("POST")
	r.HandleFunc("/producto-promocion/{id}", productoPromocionRepo.DeleteProductoPromocion).Methods("DELETE")

	fmt.Println("Servidor escuchando en el puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
