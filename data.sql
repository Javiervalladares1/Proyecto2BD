-- data.sql: Inserción de datos de prueba

-- Insertar un evento
INSERT INTO events (name, description, event_date)
VALUES
('Concierto de Rock', 'Concierto en vivo de una banda reconocida', '2025-05-20 20:00:00');

-- Insertar asientos para el evento (asumiendo que el evento recién insertado tiene id = 1)
INSERT INTO seats (event_id, seat_number)
VALUES
(1, 'A1'),
(1, 'A2'),
(1, 'A3'),
(1, 'B1'),
(1, 'B2'),
(1, 'B3');

-- Insertar usuarios de prueba
INSERT INTO users (username, email)
VALUES
('juan', 'juan@example.com'),
('maria', 'maria@example.com'),
('carlos', 'carlos@example.com'),
('ana', 'ana@example.com');

-- Insertar reservas iniciales
-- Se asume que los asientos insertados tienen los siguientes IDs: A1 -> id 1, A2 -> id 2, B1 -> id 4, A3 -> id 3
INSERT INTO reservations (user_id, seat_id)
VALUES
(1, 1),   -- Juan reserva el asiento A1
(2, 2),   -- Maria reserva el asiento A2
(3, 4),   -- Carlos reserva el asiento B1
(4, 3);   -- Ana reserva el asiento A3
