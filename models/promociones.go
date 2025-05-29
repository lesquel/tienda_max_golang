package models

import "time"

type Promocion struct {
	ID                  string    `json:"id"`
	Nombre              string    `json:"nombre"`
	Descripcion         string    `json:"descripcion"`
	DescuentoPorcentaje float64   `json:"descuento_porcentaje"`
	FechaInicio         time.Time `json:"fecha_inicio"`
	FechaFin            time.Time `json:"fecha_fin"`
	Activa              bool      `json:"activa"`
	EstadoSincronizacion string   `json:"estado_sincronizacion"`
}