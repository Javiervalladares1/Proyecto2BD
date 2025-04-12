# Proyecto2BD


Proyecto 2 - Simulación de Reservas Concurrentes

Este proyecto simula reservas concurrentes en un evento usando PostgreSQL y Go. Se estudian los niveles de aislamiento (READ COMMITTED, REPEATABLE READ y SERIALIZABLE) para ver cómo afectan el número de reservas exitosas en un escenario donde solo hay 3 asientos disponibles. Además, se realizan pruebas con diferentes cantidades de usuarios (5, 10, 20 y 30).

Requisitos
	•	Docker (para correr contenedores)
	•	Go (para compilar el código de pruebas)
	•	(Opcional) Visual Studio Code (para editar y revisar el código)

Estructura del Proyecto
	•	ddl.sql – Script SQL que crea las tablas de la base de datos.
	•	data.sql – Script SQL que inserta datos de prueba.
En este ejemplo, se asume que hay 6 asientos. Los asientos 1, 3 y 5 quedan libres y los 2, 4 y 6 quedan reservados inicialmente. Además, se insertan al menos 30 usuarios.
	•	run_all_tests_with_init.go – Código en Go que limpia y reinicializa la base de datos para cada prueba, carga los scripts anteriores y ejecuta la simulación mediante la imagen Docker proyecto2bd.
	•	Dockerfile – Instrucciones para construir la imagen Docker de la aplicación (el simulador).

Pasos para Ejecutar el Proyecto

1. Levantar el Contenedor de PostgreSQL

Primero, asegúrate de que no haya un contenedor existente con el mismo nombre:

docker rm -f entradas-postgress

Luego, corre el contenedor de PostgreSQL:

docker run --name entradas-postgress \
  -e POSTGRES_USER=javiervalladares \
  -e POSTGRES_PASSWORD=Database \
  -e POSTGRES_DB=TicketsTracker \
  -p 5436:5432 \
  -d postgres

Nota: Aunque se use TicketsTracker, PostgreSQL crea la base de datos en minúsculas como ticketstracker.

Verifica que el contenedor esté en ejecución con:

docker ps

2. Construir la Imagen de la Aplicación

En la carpeta del proyecto, donde se encuentra el Dockerfile, ejecuta:

docker build -t proyecto2bd .

Esto generará una imagen Docker llamada proyecto2bd que contiene el simulador.

3. Compilar el Programa de Pruebas en Go

Asegúrate de tener en el directorio los archivos ddl.sql, data.sql y run_all_tests_with_init.go. Luego, inicializa el módulo de Go (si aún no lo has hecho) y compila:

go mod init run_all_tests_with_init
go mod tidy
go build -o run_all_tests run_all_tests_with_init.go

Esto creará un ejecutable llamado run_all_tests en tu carpeta.

4. Ejecutar las Pruebas

El programa de pruebas ejecuta 12 combinaciones de pruebas (3 niveles de aislamiento × 4 cantidades de usuarios) y reinicializa la base de datos antes de cada prueba. Para ejecutarlo, simplemente corre:

./run_all_tests

Cada prueba realizará los siguientes pasos:
	•	Recreación de la base de datos: El programa se conecta al contenedor y ejecuta un comando para dropear y crear la base de datos ticketstracker.
	•	Carga de ddl.sql y data.sql: Se crean las tablas y se insertan los datos de prueba. Se aseguran de que los asientos 1, 3 y 5 estén libres.
	•	Ejecución de la simulación: Se invoca la imagen Docker proyecto2bd con las variables de entorno adecuadas (por ejemplo, NUM_USERS=10 y ISOLATION_LEVEL="READ COMMITTED").
	•	Finalmente, el programa muestra la salida de cada prueba, indicando el número de reservas exitosas, fallidas y el tiempo promedio.

5. Verificar el Estado de la Base de Datos (Opcional)

Para asegurarte de que la base de datos se reinicializa correctamente antes de cada prueba, puedes conectarte manualmente y hacer una consulta:

docker exec -it entradas-postgress psql -U javiervalladares -d ticketstracker -c "SELECT id, seat_number, is_reserved FROM seats;"

Deberías ver que los asientos 1, 3 y 5 están libres (is_reserved = false) y los asientos 2, 4 y 6 están marcados como reservados (is_reserved = true).

Resumen
	1.	Levantar PostgreSQL:

docker run --name entradas-postgress -e POSTGRES_USER=javiervalladares -e POSTGRES_PASSWORD=Database -e POSTGRES_DB=TicketsTracker -p 5436:5432 -d postgres


	2.	Construir la imagen Docker:

docker build -t proyecto2bd .


	3.	Compilar el programa de pruebas:

go mod init run_all_tests_with_init
go mod tidy
go build -o run_all_tests run_all_tests_with_init.go


	4.	Ejecutar las pruebas:

./run_all_tests


	5.	Verificar la base de datos (opcional):

docker exec -it entradas-postgress psql -U javiervalladares -d ticketstracker -c "SELECT id, seat_number, is_reserved FROM seats;"



