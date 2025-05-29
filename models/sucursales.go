package models

type Sucursal struct {
	ID                   string `json:"id"`
	Nombre               string `json:"nombre"`
	Direccion            string `json:"direccion"`
	Telefono             string `json:"telefono"`
	EstadoSincronizacion string `json:"estado_sincronizacion"`
}
