package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// ReservationResult almacena el resultado de cada usuario al intentar reservar
type ReservationResult struct {
	userID   int
	success  bool
	duration time.Duration
	message  string
}

// reserveSeat maneja la lógica de reservar un asiento con transacciones y reintentos en caso de error de serialización
func reserveSeat(db *sql.DB, isolationLevel string, seatID int, userID int) ReservationResult {
	const maxRetries = 5
	start := time.Now()
	var finalMsg string
	var success bool

	for attempt := 1; attempt <= maxRetries; attempt++ {
		tx, err := db.Begin()
		if err != nil {
			finalMsg = fmt.Sprintf("error starting transaction: %v", err)
			break
		}

		// Ajusta el nivel de aislamiento
		setQuery := fmt.Sprintf("SET TRANSACTION ISOLATION LEVEL %s", isolationLevel)
		if _, err := tx.Exec(setQuery); err != nil {
			tx.Rollback()
			finalMsg = fmt.Sprintf("error setting isolation level: %v", err)
			break
		}

		// Bloquea el registro del asiento con FOR UPDATE
		var isReserved bool
		err = tx.QueryRow("SELECT is_reserved FROM seats WHERE id = $1 FOR UPDATE", seatID).Scan(&isReserved)
		if err != nil {
			tx.Rollback()
			finalMsg = fmt.Sprintf("error selecting seat %d: %v", seatID, err)
			break
		}

		// Si el asiento ya está reservado, no reintentar
		if isReserved {
			tx.Rollback()
			finalMsg = fmt.Sprintf("seat %d already reserved", seatID)
			break
		}

		// Inserta la reserva en la tabla
		_, err = tx.Exec("INSERT INTO reservations(user_id, seat_id, status) VALUES ($1, $2, 'confirmed')", userID, seatID)
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				finalMsg = fmt.Sprintf("seat %d already reserved", seatID)
				break
			} else if strings.Contains(err.Error(), "could not serialize access") {
				// Error de serialización => reintento
				if attempt == maxRetries {
					finalMsg = fmt.Sprintf("failed to insert reservation for seat %d after %d attempts: %v", seatID, attempt, err)
					break
				}
				time.Sleep(time.Duration(200*attempt) * time.Millisecond)
				continue
			} else {
				finalMsg = fmt.Sprintf("error inserting reservation: %v", err)
				break
			}
		}

		// Actualiza el estado del asiento a reservado
		_, err = tx.Exec("UPDATE seats SET is_reserved = TRUE WHERE id = $1", seatID)
		if err != nil {
			tx.Rollback()
			finalMsg = fmt.Sprintf("error updating seat: %v", err)
			break
		}

		// Intenta confirmar la transacción
		err = tx.Commit()
		if err != nil {
			if strings.Contains(err.Error(), "could not serialize access") {
				if attempt == maxRetries {
					finalMsg = fmt.Sprintf("failed to commit reservation for seat %d after %d attempts: %v", seatID, attempt, err)
					break
				}
				time.Sleep(time.Duration(200*attempt) * time.Millisecond)
				continue
			} else {
				finalMsg = fmt.Sprintf("error committing transaction: %v", err)
				break
			}
		}

		// Si llega aquí, la reserva fue exitosa.
		finalMsg = fmt.Sprintf("reservation confirmed for seat %d", seatID)
		success = true
		break
	}

	duration := time.Since(start)
	return ReservationResult{
		userID:   userID,
		success:  success,
		duration: duration,
		message:  finalMsg,
	}
}

func main() {
	// Lee variables de entorno o usa valores por defecto
	isolationLevel := os.Getenv("ISOLATION_LEVEL")
	if isolationLevel == "" {
		isolationLevel = "READ COMMITTED"
	}

	numUsers := 5
	if env := os.Getenv("NUM_USERS"); env != "" {
		if parsed, err := strconv.Atoi(env); err == nil {
			numUsers = parsed
		}
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5436"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "ticketstracker"
	}

	// Construir la cadena de conexión
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	// Verifica la conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}
	fmt.Println("Connected to PostgreSQL!")
	fmt.Printf("Simulating %d users attempting to reserve seats with isolation %s\n", numUsers, isolationLevel)

	// Lista de asientos libres, ajústalo a tu data.sql

	freeSeats := []int{1, 3, 5}

	// Canal y WaitGroup para los resultados
	var wg sync.WaitGroup
	resultChan := make(chan ReservationResult, numUsers)

	// Lógica de round-robin: cada usuario coge un asiento del array freeSeats
	for i := 1; i <= numUsers; i++ {
		wg.Add(1)
		seatID := freeSeats[(i-1)%len(freeSeats)]
		go func(userID, seat int) {
			defer wg.Done()
			res := reserveSeat(db, isolationLevel, seat, userID)
			resultChan <- res
		}(i, seatID)
	}

	wg.Wait()
	close(resultChan)

	// Recolectar resultados
	var successCount, failureCount int
	var totalTime time.Duration
	for res := range resultChan {
		fmt.Printf("User %d: %s (in %v)\n", res.userID, res.message, res.duration)
		totalTime += res.duration
		if res.success {
			successCount++
		} else {
			failureCount++
		}
	}

	avgTime := totalTime / time.Duration(numUsers)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Usuarios Concurrentes: %d\tAislamiento: %s\tExitosas: %d\tFallidas: %d\tPromedio: %d ms\n",
		numUsers, isolationLevel, successCount, failureCount, avgTime.Milliseconds())

}
