package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// buat koneksi database
	db, err := sql.Open("pgx", os.Getenv("DB_URI"))
	if err != nil {
		fmt.Printf("Gagal membuat koneksi ke database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// lakukan verifikasi koneksi database
	if err = db.Ping(); err != nil {
		fmt.Printf("Gagal memverifikasi koneksi database: %v\n", err)
		os.Exit(1)
	}

	// lakukan migrasi tabel database
	if err = migrate(db); err != nil {
		fmt.Printf("Gagal melakukan migrasi database: %v\n", err)
		os.Exit(1)
	}

	// inisiasi router
	r, err := routes(db)
	if err != nil {
		fmt.Printf("Gagal membuat router: %v\n", err)
		os.Exit(1)
	}

	// buat server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// jalankan server
	if err = server.ListenAndServe(); err != nil {
		fmt.Printf("Gagal menjalankan server: %v\n", err)
		os.Exit(1)
	}
}
