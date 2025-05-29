package tables

import (
	"database/sql"
)

func CrearTablaSucursales(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sucursales (
		id TEXT PRIMARY KEY,
		nombre TEXT NOT NULL,
		direccion TEXT NOT NULL,
		telefono TEXT,
		estado_sincronizacion TEXT DEFAULT 'sincronizado'
	);`
	_, err := db.Exec(query)
	return err
}

func CrearTablaClientes(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS clientes (
		id TEXT PRIMARY KEY,
		nombre TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		telefono TEXT,
		direccion TEXT,
		fecha_registro TEXT DEFAULT CURRENT_TIMESTAMP,
		sucursal_id TEXT,
		estado_sincronizacion TEXT DEFAULT 'sincronizado',
		FOREIGN KEY (sucursal_id) REFERENCES sucursales(id) ON DELETE SET NULL
	);`
	_, err := db.Exec(query)
	return err
}

func CrearTablaProductos(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS productos (
		id TEXT PRIMARY KEY,
		nombre TEXT NOT NULL,
		descripcion TEXT,
		precio REAL NOT NULL,
		activo BOOLEAN DEFAULT TRUE,
		estado_sincronizacion TEXT DEFAULT 'sincronizado'
	);`
	_, err := db.Exec(query)
	return err
}

func CrearTablaProductoSucursal(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS producto_sucursal (
		id TEXT PRIMARY KEY,
		producto_id TEXT,
		sucursal_id TEXT,
		stock INTEGER NOT NULL DEFAULT 0,
		estado_sincronizacion TEXT DEFAULT 'sincronizado',
		FOREIGN KEY (producto_id) REFERENCES productos(id) ON DELETE CASCADE,
		FOREIGN KEY (sucursal_id) REFERENCES sucursales(id) ON DELETE CASCADE
	);`
	_, err := db.Exec(query)
	return err
}

func CrearTablaVentas(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ventas (
		id TEXT PRIMARY KEY,
		cliente_id TEXT,
		sucursal_id TEXT,
		fecha TEXT DEFAULT CURRENT_TIMESTAMP,
		estado_sincronizacion TEXT DEFAULT 'sincronizado',
		FOREIGN KEY (cliente_id) REFERENCES clientes(id) ON DELETE SET NULL,
		FOREIGN KEY (sucursal_id) REFERENCES sucursales(id) ON DELETE SET NULL
	);`
	_, err := db.Exec(query)
	return err
}

func CrearTablaVentaDetalle(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS venta_detalle (
		id TEXT PRIMARY KEY,
		venta_id TEXT,
		producto_id TEXT,
		cantidad INTEGER NOT NULL,
		precio_unitario REAL NOT NULL,
		estado_sincronizacion TEXT DEFAULT 'sincronizado',
		FOREIGN KEY (venta_id) REFERENCES ventas(id) ON DELETE CASCADE,
		FOREIGN KEY (producto_id) REFERENCES productos(id) ON DELETE SET NULL
	);`
	_, err := db.Exec(query)
	return err
}

func CrearTablaPromociones(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS promociones (
		id TEXT PRIMARY KEY,
		nombre TEXT NOT NULL,
		descripcion TEXT,
		descuento_porcentaje REAL CHECK (descuento_porcentaje >= 0 AND descuento_porcentaje <= 100),
		fecha_inicio TEXT,
		fecha_fin TEXT,
		activa BOOLEAN DEFAULT TRUE,
		estado_sincronizacion TEXT DEFAULT 'sincronizado'
	);`
	_, err := db.Exec(query)
	return err
}

func CrearTablaProductoPromocion(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS producto_promocion (
		id TEXT PRIMARY KEY,
		producto_id TEXT,
		promocion_id TEXT,
		estado_sincronizacion TEXT DEFAULT 'sincronizado',
		FOREIGN KEY (producto_id) REFERENCES productos(id) ON DELETE CASCADE,
		FOREIGN KEY (promocion_id) REFERENCES promociones(id) ON DELETE CASCADE
	);`
	_, err := db.Exec(query)
	return err
}
