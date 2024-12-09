package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// initTable crea la tabla "messages" si no existe
func initTable(w http.ResponseWriter, r *http.Request) {
	// SQL para crear la tabla si no existe
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		from VARCHAR(255) NOT NULL,
		to VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		constraint from_to check (from <> to),
		constraint from_check CHECK (position('@' IN from) > 0),
		constraint to_check CHECK (position('@' IN to) > 0)
	);`

	// Ejecuta la creación de la tabla
	_, err := db.Exec(createTableSQL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al crear la tabla: %v", err), http.StatusInternalServerError)
		return
	}

	// Responde indicando que la tabla se creó o ya existe
	w.Write([]byte("Tabla 'messages' creada o ya existe"))
}

// checkTable verifica si la tabla "messages" existe
func checkTable(w http.ResponseWriter, r *http.Request) {
	// SQL para verificar si la tabla existe
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'messages');`
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al verificar la tabla: %v", err), http.StatusInternalServerError)
		return
	}

	// Responde según si la tabla existe
	if exists {
		w.Write([]byte("La tabla 'messages' existe"))
	} else {
		w.Write([]byte("La tabla 'messages' no existe"))
	}
}

func main() {
	// Obtener las variables de entorno
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	const dbService = "postgresql"

	// Si alguna variable de entorno no está definida, el programa falla
	if dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Las variables de entorno POSTGRES_USER, POSTGRES_PASSWORD y POSTGRES_DB deben estar definidas")
	}

	// Construir la cadena de conexión
	// connStr := "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbService, dbName)

	// Conexión a PostgreSQL
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	defer db.Close()

	// Verifica que la base de datos se pueda acceder
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}

	// Manejadores de las rutas
	http.HandleFunc("/init", withLogging(initTable))
	http.HandleFunc("/status", withLogging(checkTable))
	//manejador por defecto 404
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Ruta no encontrada: %s %s", r.Method, r.URL.Path)
		http.Error(w, "Ruta no encontrada", http.StatusNotFound)
	})

	// Inicia el servidor en el puerto 8080
	fmt.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Middleware para registrar solicitudes HTTP
func withLogging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Registrar información de la solicitud
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Ejecutar el manejador original
		handler(w, r)

		// Registrar información adicional (tiempo de respuesta)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}
