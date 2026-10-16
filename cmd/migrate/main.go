// Tool migration sendiri (bukan CLI golang-migrate eksternal) supaya jalan
// persis sama di Windows, macOS, maupun Linux — cukup "go run ./cmd/migrate ...",
// tidak perlu install CLI terpisah ataupun export DATABASE_URL manual.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL belum diatur di berkas .env")
	}

	if len(os.Args) < 2 {
		log.Fatal("penggunaan: go run ./cmd/migrate <up|down|version|force> [angka versi]")
	}

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("gagal membuka migrations: %v", err)
	}

	switch os.Args[1] {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	case "version":
		version, dirty, verErr := m.Version()
		if verErr != nil {
			log.Fatalf("gagal membaca versi: %v", verErr)
		}
		fmt.Printf("versi migration saat ini: %d (dirty: %v)\n", version, dirty)
		return
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("penggunaan: go run ./cmd/migrate force <angka versi>")
		}
		version, convErr := strconv.Atoi(os.Args[2])
		if convErr != nil {
			log.Fatalf("angka versi tidak valid: %v", convErr)
		}
		err = m.Force(version)
	default:
		log.Fatalf("perintah tidak dikenal: %q (pakai up/down/version/force)", os.Args[1])
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migration gagal: %v", err)
	}
	fmt.Println("selesai")
}
