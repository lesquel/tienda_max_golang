# Usa una imagen oficial de Go para compilar
FROM golang:1.20 AS builder

WORKDIR /app

# Copia los archivos necesarios
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Compila la app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Imagen final mínima
FROM alpine:latest

WORKDIR /root/

# Copia el binario desde la etapa anterior
COPY --from=builder /app/main .

# El puerto que usará Cloud Run
EXPOSE 8080

# Comando para ejecutar
CMD ["./main"]
