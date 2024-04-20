package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/fastcampus-backend-golang/online-shop/handler"
	"github.com/fastcampus-backend-golang/online-shop/middleware"

	"github.com/gin-gonic/gin"
)

// routes digunakan untuk inisiasi endpoint-endpoint API
func routes(db *sql.DB) (http.Handler, error) {
	// pastikan koneksi ke database tidak nil
	if db == nil {
		return nil, errors.New("tidak ada koneksi ke database")
	}

	// init router
	r := gin.Default()

	// endpoint publik
	r.GET("/api/v1/products", handler.ListProducts(db))
	r.GET("/api/v1/products/:id", handler.GetProduct(db))
	r.POST("/api/v1/checkout", handler.CheckoutOrder(db))

	// endpoint pelanggan dengan passcode
	r.POST("/api/v1/orders/:id/confirm", handler.ConfirmOrder(db))
	r.GET("/api/v1/orders/:id", handler.GetOrder(db))

	// endpoint admin (dengan verifikasi header)
	r.POST("/admin/products", middleware.AdminOnly(), handler.CreateProduct(db))
	r.PUT("/admin/products/:id", middleware.AdminOnly(), handler.UpdateProduct(db))
	r.DELETE("/admin/products/:id", middleware.AdminOnly(), handler.DeleteProduct(db))

	return r, nil
}
