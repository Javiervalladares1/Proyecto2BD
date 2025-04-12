-- Elimina las tablas existentes en orden de dependencia
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS users;

-- Crear tabla de eventos
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    event_date TIMESTAMP NOT NULL
);

-- Crear tabla de asientos
CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    seat_number VARCHAR(10) NOT NULL,
    is_reserved BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT unique_seat UNIQUE (event_id, seat_number)
);

-- Crear tabla de usuarios
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE
);

-- Crear tabla de reservas
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    seat_id INTEGER NOT NULL,
    reserved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_seat FOREIGN KEY (seat_id) REFERENCES seats(id) ON DELETE CASCADE,
    CONSTRAINT unique_reservation UNIQUE (seat_id)
);