package models

import "time"

type Cliente struct {
	ID                   string    `json:"id"`
	Nombre               string    `json:"nombre"`
	Email                string    `json:"email"`
	Telefono             string    `json:"telefono"`
	Direccion            string    `json:"direccion"`
	FechaRegistro        time.Time `json:"fecha_registro"`
	SucursalID           *string    `json:"sucursal_id"`
	EstadoSincronizacion string    `json:"estado_sincronizacion"`
}
