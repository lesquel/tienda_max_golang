package models

type Producto struct {
	ID                  string  `json:"id"`
	Nombre              string  `json:"nombre"`
	Descripcion         string  `json:"descripcion"`
	Precio             float64 `json:"precio"`
	Activo             bool    `json:"activo"`
	EstadoSincronizacion string `json:"estado_sincronizacion"`
}