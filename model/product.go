package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Product adalah representasi dari data produk di database dan API
type Product struct {
	ID        string `json:"id" binding:"len=0"` // mencegah ID diisi oleh user
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	IsDeleted *bool  `json:"is_deleted,omitempty"`
}

// SelectProduct adalah fungsi untuk mengambil data produk dari database
func SelectProduct(db *sql.DB) ([]Product, error) {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return nil, errors.New("tidak ada koneksi ke database")
	}

	// query untuk mengambil data produk
	query := `SELECT id, name, price FROM products WHERE is_deleted = FALSE`

	// eksekusi query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// ubah data hasil query ke bentuk slice
	products := []Product{}
	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// SelectProductByID adalah fungsi untuk mengambil data produk berdasarkan ID dari database
func SelectProductByID(db *sql.DB, id string) (Product, error) {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return Product{}, errors.New("gagal melakukan koneksi ke database")
	}

	// query untuk mengambil data produk berdasarkan ID
	query := `SELECT id, name, price FROM products WHERE is_deleted = FALSE AND id = $1`

	// eksekusi query
	product := Product{}
	err := db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// SelectProductIn adalah fungsi untuk mengambil data produk berdasarkan ID-ID dari database
func SelectProductIn(db *sql.DB, ids []string) ([]Product, error) {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return nil, errors.New("tidak ada koneksi ke database")
	}

	// buat placeholder & args untuk query
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// buat query dengan placeholder
	query := fmt.Sprintf(`SELECT id, name, price FROM products WHERE is_deleted = FALSE AND id IN (%s)`, strings.Join(placeholders, ","))

	// eksekusi query dengan args berisi id-id produk
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// ubah data hasil query ke bentuk slice
	products := []Product{}
	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// InsertProduct adalah fungsi untuk menyimpan data produk ke database
func InsertProduct(db *sql.DB, product Product) error {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return errors.New("gagal melakukan koneksi ke database")
	}

	// query untuk insert data produk
	query := `INSERT INTO products (id, name, price) VALUES ($1, $2, $3)`

	// eksekusi query
	_, err := db.Exec(query, product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct adalah fungsi untuk mengubah data produk di database
func UpdateProduct(db *sql.DB, product Product) error {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return errors.New("gagal melakukan koneksi ke database")
	}

	// query untuk update data produk
	query := `UPDATE products SET name = $1, price = $2 WHERE id = $3`

	// eksekusi query
	_, err := db.Exec(query, product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct adalah fungsi untuk menghapus data produk dari database
func DeleteProduct(db *sql.DB, id string) error {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return errors.New("gagal melakukan koneksi ke database")
	}

	// query untuk menghapus data produk
	query := `UPDATE products SET is_deleted = TRUE WHERE id = $1`

	// eksekusi query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
