package models

type VentaDetalle struct {
	ID                  string  `json:"id"`
	VentaID             string  `json:"venta_id"`
	ProductoID          string  `json:"producto_id"`
	Cantidad           int     `json:"cantidad"`
	PrecioUnitario     float64 `json:"precio_unitario"`
	EstadoSincronizacion string `json:"estado_sincronizacion"`
}