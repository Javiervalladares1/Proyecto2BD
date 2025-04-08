package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

// reserveSeat intenta reservar el asiento indicado para el usuario actual
// usando un transaction con el nivel de aislamiento especificado.
func reserveSeat(db *sql.DB, isolationLevel string, seatID int, userID int, resultChan chan<- string) {
	// Inicia la transacción
	tx, err := db.Begin()
	if err != nil {
		resultChan <- fmt.Sprintf("Usuario %d: error al iniciar transacción: %v", userID, err)
		return
	}

	// Configura el nivel de aislamiento de la transacción
	query := fmt.Sprintf("SET TRANSACTION ISOLATION LEVEL %s", isolationLevel)
	if _, err := tx.Exec(query); err != nil {
		tx.Rollback()
		resultChan <- fmt.Sprintf("Usuario %d: error al establecer aislamiento (%s): %v", userID, isolationLevel, err)
		return
	}

	// Bloquea el registro del asiento para evitar conflictos (consulta FOR UPDATE)
	var isReserved bool
	err = tx.QueryRow("SELECT is_reserved FROM seats WHERE id = $1 FOR UPDATE", seatID).Scan(&isReserved)
	if err != nil {
		tx.Rollback()
		resultChan <- fmt.Sprintf("Usuario %d: error al consultar el asiento %d: %v", userID, seatID, err)
		return
	}

	// Si el asiento ya está reservado, se cancela la transacción
	if isReserved {
		tx.Rollback()
		resultChan <- fmt.Sprintf("Usuario %d: el asiento %d ya está reservado", userID, seatID)
		return
	}

	// Inserta el registro en la tabla reservations
	_, err = tx.Exec("INSERT INTO reservations(user_id, seat_id, status) VALUES ($1, $2, 'confirmed')", userID, seatID)
	if err != nil {
		tx.Rollback()
		resultChan <- fmt.Sprintf("Usuario %d: error al insertar reserva: %v", userID, err)
		return
	}

	// Actualiza el estado del asiento a reservado
	_, err = tx.Exec("UPDATE seats SET is_reserved = TRUE WHERE id = $1", seatID)
	if err != nil {
		tx.Rollback()
		resultChan <- fmt.Sprintf("Usuario %d: error al actualizar el asiento: %v", userID, err)
		return
	}

	// Intenta confirmar la transacción
	if err = tx.Commit(); err != nil {
		resultChan <- fmt.Sprintf("Usuario %d: error al confirmar transacción: %v", userID, err)
		return
	}

	resultChan <- fmt.Sprintf("Usuario %d: reserva confirmada para el asiento %d", userID, seatID)
}

func main() {
	// Parámetros de línea de comandos
	isolationPtr := flag.String("isolation", "READ COMMITTED", "Nivel de aislamiento: READ COMMITTED, REPEATABLE READ, SERIALIZABLE")
	numUsersPtr := flag.Int("users", 5, "Número de usuarios concurrentes")
	seatIDPtr := flag.Int("seat", 1, "ID del asiento a reservar")
	flag.Parse()

	// Cadena de conexión a PostgreSQL (ajusta host, usuario, contraseña y nombre de base de datos según tu configuración)
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=reservasdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error conectando a la base de datos: ", err)
	}
	defer db.Close()

	// Prueba la conexión a la base de datos
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la base de datos: ", err)
	}

	// Canal para recibir resultados de cada goroutine
	resultChan := make(chan string, *numUsersPtr)
	var wg sync.WaitGroup

	fmt.Printf("Simulación de %d usuarios intentando reservar el asiento %d con nivel de aislamiento %s\n", *numUsersPtr, *seatIDPtr, *isolationPtr)

	// Lanza las goroutines para simular usuarios concurrentes
	for i := 1; i <= *numUsersPtr; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			reserveSeat(db, *isolationPtr, *seatIDPtr, userID, resultChan)
		}(i)
	}

	wg.Wait()
	close(resultChan)

	// Muestra los resultados obtenidos
	for res := range resultChan {
		fmt.Println(res)
	}
}
