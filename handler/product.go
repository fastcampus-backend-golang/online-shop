package handler

import (
	"database/sql"
	"errors"

	"github.com/fastcampus-backend-golang/online-shop/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil data produk dari database
		products, err := model.SelectProduct(db)
		if err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// tampilkan data produk
		c.JSON(200, products)
	}
}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil id produk dari URL
		id := c.Param("id")

		// ambil data produk dari database
		product, err := model.SelectProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
				return
			}

			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// tampilkan data produk
		c.JSON(200, product)
	}
}

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil data produk dari request body
		var product model.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(400, gin.H{"error": "Data produk tidak valid"})
			return
		}

		// atur id dari UUID
		product.ID = uuid.New().String()

		// simpan data produk ke database
		if err := model.InsertProduct(db, product); err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// tampilkan data produk yang disimpan
		c.JSON(201, product)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil id produk dari URL
		id := c.Param("id")

		// ambil data produk dari request body
		var productReq model.Product
		if err := c.BindJSON(&productReq); err != nil {
			c.JSON(400, gin.H{"error": "Data produk tidak valid"})
			return
		}

		// ambil data produk dari database
		product, err := model.SelectProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
				return
			}
		}

		// update nama produk jika tidak kosong
		if productReq.Name != "" {
			product.Name = productReq.Name
		}

		// update harga produk jika tidak kosong
		if productReq.Price != 0 {
			product.Price = productReq.Price
		}

		// update data produk ke database
		if err := model.UpdateProduct(db, product); err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// tampilkan data produk yang diupdate
		c.JSON(200, product)
	}
}

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil id produk dari URL
		id := c.Param("id")

		// hapus data produk dari database
		if err := model.DeleteProduct(db, id); err != nil {
			c.JSON(500, gin.H{"error": "Terjadi kesalahan pada server"})
			return
		}

		// tampilkan data produk yang dihapus
		c.JSON(204, nil)
	}
}
