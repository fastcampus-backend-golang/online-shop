package model

import (
	"database/sql"
	"errors"
	"time"
)

// ProductQuantity adalah representasi dari data produk dan kuantitas di API
type ProductQuantity struct {
	ID       string `json:"id" binding:"required"`
	Quantity int32  `json:"quantity" binding:"required"`
}

// Checkout adalah representasi dari data checkout di API
type Checkout struct {
	Email    string            `json:"email" binding:"required,email"`
	Address  string            `json:"address" binding:"required"`
	Products []ProductQuantity `json:"products" binding:"min=1"` // minimal 1 produk
}

// Confirm adalah representasi dari data konfirmasi pembayaran di API
type Confirm struct {
	Amount        int64  `json:"amount" binding:"required"`
	Bank          string `json:"bank" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	Passcode      string `json:"passcode"`
}

// Order adalah representasi dari data pesanan di database
type Order struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	Address           string     `json:"address"`
	GrandTotal        int64      `json:"grandTotal"`
	Passcode          *string    `json:"passcode,omitempty"`
	PaidAt            *time.Time `json:"paidAt,omitempty"`
	PaidBank          *string    `json:"paidBank,omitempty"`
	PaidAccountNumber *string    `json:"paidAccountNumber,omitempty"`
}

// OrderDetail adalah representasi dari detail data pesanan di database dan API
type OrderDetail struct {
	ID        string `json:"id"`
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	Quantity  int32  `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}

// OrderWithDetail adalah representasi dari data pesanan dengan detail untuk API (tidak menampilkan passcode)
type OrderWithDetail struct {
	Order
	Detail []OrderDetail `json:"detail"`
}

// CreateOrder adalah fungsi untuk menyimpan data pesanan ke database
func CreateOrder(db *sql.DB, order Order, details []OrderDetail) error {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return errors.New("tidak ada koneksi ke database")
	}

	// buat transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// query untuk simpan data order
	queryOrder := `INSERT INTO orders (id, email, address, passcode, grand_total) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(queryOrder, order.ID, order.Email, order.Address, order.Passcode, order.GrandTotal)
	if err != nil {
		tx.Rollback()
		return err
	}

	// query untuk simpan data detail order
	queryDetail := `INSERT INTO order_details (id, order_id, product_id, quantity, price, total) VALUES ($1, $2, $3, $4, $5, $6)`
	for _, detail := range details {
		_, err = tx.Exec(queryDetail, detail.ID, detail.OrderID, detail.ProductID, detail.Quantity, detail.Price, detail.Total)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// UpdateOrderStatus adalah fungsi untuk mengubah status pesanan menjadi sudah dibayar
func UpdateOrderStatus(db *sql.DB, id string, confirmation Confirm, paidAt time.Time) error {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return errors.New("tidak ada koneksi ke database")
	}

	// query untuk update status pesanan
	query := `UPDATE orders SET paid_at = $1, paid_bank = $2, paid_account_number = $3 WHERE id = $4`
	_, err := db.Exec(query, paidAt, confirmation.Bank, confirmation.AccountNumber, id)
	if err != nil {
		return err
	}

	return nil
}

// SelectOrderByID adalah fungsi untuk mengambil data pesanan berdasarkan ID
func SelectOrderByID(db *sql.DB, id string) (Order, error) {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return Order{}, errors.New("tidak ada koneksi ke database")
	}

	// query untuk mengambil data order
	queryOrder := `SELECT id, email, address, passcode, grand_total, paid_at, paid_bank, paid_account_number FROM orders WHERE id = $1`
	row := db.QueryRow(queryOrder, id)

	// siapkah variabel untuk menampung data order
	order := Order{}

	// ambil data dari row
	err := row.Scan(&order.ID, &order.Email, &order.Address, &order.Passcode, &order.GrandTotal, &order.PaidAt, &order.PaidBank, &order.PaidAccountNumber)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

// SelectOrderDetailByOrderID adalah fungsi untuk mengambil data detail pesanan berdasarkan ID pesanan
func SelectOrderDetailByOrderID(db *sql.DB, orderID string) ([]OrderDetail, error) {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return nil, errors.New("tidak ada koneksi ke database")
	}

	// query untuk mengambil data detail order
	queryDetail := `SELECT id, order_id, product_id, quantity, price, total FROM order_details WHERE order_id = $1`
	rows, err := db.Query(queryDetail, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// siapkan variabel untuk menampung data detail order
	var details []OrderDetail

	// ambil data dari rows
	for rows.Next() {
		detail := OrderDetail{}
		err := rows.Scan(&detail.ID, &detail.OrderID, &detail.ProductID, &detail.Quantity, &detail.Price, &detail.Total)
		if err != nil {
			return nil, err
		}

		details = append(details, detail)
	}

	return details, nil
}
