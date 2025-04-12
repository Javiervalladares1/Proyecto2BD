-- data.sql
-- 1) Insertar un evento con id=1
INSERT INTO events (name, description, event_date)
VALUES ('Concierto de Rock', 'Concierto con banda reconocida', '2025-05-20 20:00:00');


INSERT INTO seats (event_id, seat_number, is_reserved)
VALUES
  (1, 'A1', false),  -- id=1 => Libre
  (1, 'A2', false),  
  (1, 'A3', false),  -- id=3 => Libre
  (1, 'B1', false),  
  (1, 'B2', false),  -- id=5 => Libre
  (1, 'B3', false);  


INSERT INTO users (username, email)
VALUES
  ('user1', 'user1@example.com'),
  ('user2', 'user2@example.com'),
  ('user3', 'user3@example.com'),
  ('user4', 'user4@example.com'),
  ('user5', 'user5@example.com'),
  ('user6', 'user6@example.com'),
  ('user7', 'user7@example.com'),
  ('user8', 'user8@example.com'),
  ('user9', 'user9@example.com'),
  ('user10', 'user10@example.com'),
  ('user11', 'user11@example.com'),
  ('user12', 'user12@example.com'),
  ('user13', 'user13@example.com'),
  ('user14', 'user14@example.com'),
  ('user15', 'user15@example.com'),
  ('user16', 'user16@example.com'),
  ('user17', 'user17@example.com'),
  ('user18', 'user18@example.com'),
  ('user19', 'user19@example.com'),
  ('user20', 'user20@example.com'),
  ('user21', 'user21@example.com'),
  ('user22', 'user22@example.com'),
  ('user23', 'user23@example.com'),
  ('user24', 'user24@example.com'),
  ('user25', 'user25@example.com'),
  ('user26', 'user26@example.com'),
  ('user27', 'user27@example.com'),
  ('user28', 'user28@example.com'),
  ('user29', 'user29@example.com'),
  ('user30', 'user30@example.com');

-- 4) Reservas iniciales: asientos 2,4,6 (IDs)

INSERT INTO reservations (user_id, seat_id, status)
VALUES
  (1, 2, 'confirmed'),
  (2, 4, 'confirmed'),
  (3, 6, 'confirmed');

-- 5) Marcar is_reserved = true para asientos 2,4,6
UPDATE seats SET is_reserved = true WHERE id IN (2,4,6);