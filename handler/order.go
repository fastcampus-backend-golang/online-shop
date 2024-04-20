package handler

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/fastcampus-backend-golang/online-shop/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil data pesanan dari request body
		var checkoutOrder model.Checkout
		if err := c.BindJSON(&checkoutOrder); err != nil {
			c.JSON(400, gin.H{"error": "Data pesanan tidak valid"})
			return
		}

		// daftar ID produk dan jumlah pesanan
		ids := make([]string, len(checkoutOrder.Products))
		orderQuantity := make(map[string]int32)
		for i, p := range checkoutOrder.Products {
			ids[i] = p.ID
			orderQuantity[p.ID] = p.Quantity
		}

		// ambil data produk dari database
		product, err := model.SelectProductIn(db, ids)
		if err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// pastikan semua produk ada
		if len(product) != len(checkoutOrder.Products) {
			c.JSON(400, gin.H{"error": "Produk tidak ditemukan"})
			return
		}

		// siapkan passcode
		passcode := generatePasscode(5)

		// hash passcode untuk disimpan di database
		hashPasscode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)
		if err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// ubah menjadi string
		hashedPasscodeStr := string(hashPasscode)

		// siapkan data order dan detail order
		order := model.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			Passcode:   &hashedPasscodeStr,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}

		// buat detail dan hitung total harga
		for _, p := range product {
			// hitung total untuk produk ini
			total := p.Price * int64(orderQuantity[p.ID])

			// tambahkan detail pesanan
			detail := model.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQuantity[p.ID],
				Price:     p.Price,
				Total:     total,
			}
			details = append(details, detail)

			// tambahkan total harga pesanan
			order.GrandTotal += total
		}

		// simpan data order dan detail order ke database
		if err := model.CreateOrder(db, order, details); err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// tampilkan passcode yang tidak dihash
		// agar user bisa menyimpannya untuk mengakses pesanan
		order.Passcode = &passcode

		// buat response
		response := model.OrderWithDetail{
			Order:  order,
			Detail: details,
		}

		// tampilkan data order yang disimpan
		c.JSON(201, response)
	}
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil id order dari URL
		id := c.Param("id")

		// baca request body
		var confirm model.Confirm
		if err := c.BindJSON(&confirm); err != nil {
			c.JSON(400, gin.H{"error": "Data konfirmasi tidak valid"})
			return
		}

		// ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Pesanan tidak ditemukan"})
				return
			}

			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// pastikan passcode tidak kosong agar tidak terjadi panic
		if order.Passcode == nil {
			c.JSON(500, gin.H{"error": "Data pesanan tidak valid"})
			return
		}

		// cocokkan passcode
		if err := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(confirm.Passcode)); err != nil {
			c.JSON(401, gin.H{"error": "Passcode tidak valid"})
			return
		}

		// izinkan hanya untuk pesanan yang belum dibayar
		if order.PaidAt != nil {
			c.JSON(400, gin.H{"error": "Pesanan sudah dibayar"})
			return
		}

		// cocokkan jumlah pembayaran
		if order.GrandTotal != confirm.Amount {
			c.JSON(400, gin.H{"error": "Jumlah pembayaran tidak sesuai"})
			return
		}

		// update status pesanan
		currentTime := time.Now()
		if err := model.UpdateOrderStatus(db, id, confirm, currentTime); err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// ambil detail order dari database
		details, err := model.SelectOrderDetailByOrderID(db, id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// jangan tampilkan passcode
		// cukup ditampilkan ketika pesanan dibuat
		order.Passcode = nil

		// update response dengan data konfirmasi pesanan
		order.PaidAt = &currentTime
		order.PaidBank = &confirm.Bank
		order.PaidAccountNumber = &confirm.AccountNumber

		// buat response
		response := model.OrderWithDetail{
			Order:  order,
			Detail: details,
		}

		// tampilkan data order yang sudah dikonfirmasi
		c.JSON(200, response)
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil id order dari URL
		id := c.Param("id")

		// ambil passcode dari query URL
		passcode := c.Query("passcode")

		// ambil data order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"error": "Pesanan tidak ditemukan"})
				return
			}

			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// pastikan passcode tidak kosong agar tidak terjadi panic
		if order.Passcode == nil {
			c.JSON(500, gin.H{"error": "Data pesanan tidak valid"})
			return
		}

		// cocokkan passcode
		if err := bcrypt.CompareHashAndPassword([]byte(*order.Passcode), []byte(passcode)); err != nil {
			c.JSON(401, gin.H{"error": "Passcode tidak valid"})
			return
		}

		// ambil detail order dari database
		details, err := model.SelectOrderDetailByOrderID(db, id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// jangan tampilkan passcode
		// cukup ditampilkan ketika pesanan dibuat
		order.Passcode = nil

		// buat response
		response := model.OrderWithDetail{
			Order:  order,
			Detail: details,
		}

		// tampilkan data order
		c.JSON(200, response)
	}
}

func generatePasscode(length int) string {
	charSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	random := make([]byte, length)
	for i := range random {
		random[i] = charSet[randomGen.Intn(len(charSet))]
	}

	return string(random)
}
