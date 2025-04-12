package main

import (
	"fmt"
	"os/exec"
	"time"
)

func cleanSchema() error {
	fmt.Println("Cleaning schema 'public' in ticketstracker...")
	cmd := exec.Command("sh", "-c", `echo "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" | docker exec -i entradas-postgress psql -U javiervalladares -d ticketstracker`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error cleaning schema: %s\n", string(output))
		return err
	}
	return nil
}

func recreateDB() error {
	fmt.Println("Recreating database 'ticketstracker'...")
	cmd := exec.Command("sh", "-c", `echo "DROP DATABASE IF EXISTS ticketstracker; CREATE DATABASE ticketstracker;" | docker exec -i entradas-postgress psql -U javiervalladares -d postgres`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error recreating database: %s\n", string(output))
		return err
	}
	return nil
}

func cleanDB() error {
	if err := cleanSchema(); err != nil {
		// Si falla por que la base no existe, recrea la BD
		fmt.Println("Cleaning schema failed, attempting to recreate database...")
		return recreateDB()
	}
	return nil
}

func initDB() error {
	if err := cleanDB(); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Loading ddl.sql...")
	cmdDDL := exec.Command("sh", "-c", "cat ddl.sql | docker exec -i entradas-postgress psql -U javiervalladares -d ticketstracker")
	if output, err := cmdDDL.CombinedOutput(); err != nil {
		fmt.Printf("Error running ddl.sql: %s\n", string(output))
		return err
	}

	fmt.Println("Loading data.sql...")
	cmdData := exec.Command("sh", "-c", "cat data.sql | docker exec -i entradas-postgress psql -U javiervalladares -d ticketstracker")
	if output, err := cmdData.CombinedOutput(); err != nil {
		fmt.Printf("Error running data.sql: %s\n", string(output))
		return err
	}
	return nil
}

func runTest(isolationLevel string, numUsers int) (string, error) {
	args := []string{
		"run", "--rm", "-i",
		"-e", "DB_HOST=host.docker.internal",
		"-e", "DB_PORT=5436",
		"-e", "DB_USER=javiervalladares",
		"-e", "DB_PASSWORD=Database",
		"-e", "DB_NAME=ticketstracker",
		"-e", fmt.Sprintf("ISOLATION_LEVEL=%s", isolationLevel),
		"-e", fmt.Sprintf("NUM_USERS=%d", numUsers),
		"proyecto2bd",
	}
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {

	isolationLevels := []string{"READ COMMITTED", "REPEATABLE READ", "SERIALIZABLE"}
	userCounts := []int{5, 10, 20, 30}

	for _, iso := range isolationLevels {
		for _, users := range userCounts {
			fmt.Println("============================================")
			fmt.Printf("Test: %d users, Isolation: %s\n", users, iso)
			fmt.Println("============================================")

			if err := initDB(); err != nil {
				fmt.Printf("Error during DB initialization, skipping test for %d users, isolation %s\n", users, iso)
				continue
			}

			time.Sleep(1 * time.Second)

			output, err := runTest(iso, users)
			if err != nil {
				fmt.Printf("Error running simulation: %v\n", err)
			}
			fmt.Println(output)
			fmt.Println("============================================\n")

			time.Sleep(2 * time.Second)
		}
	}
}
