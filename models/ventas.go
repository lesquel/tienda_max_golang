package models

import "time"

type Venta struct {
	ID                  string    `json:"id"`
	ClienteID           string    `json:"cliente_id"`
	SucursalID          string    `json:"sucursal_id"`
	Fecha               time.Time `json:"fecha"`
	EstadoSincronizacion string   `json:"estado_sincronizacion"`
}