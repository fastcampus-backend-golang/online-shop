package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil kunci rahasia admin
		key := os.Getenv("ADMIN_SECRET")
		if key == "" {
			c.JSON(500, gin.H{"error": "Gagal memproses permintaan"})
			c.Abort()
			return
		}

		// ambil header Authorization
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "Akses tidak diizinkan"})
			c.Abort()
			return
		}

		// validasi kunci admin dengan header
		if auth != key {
			c.JSON(401, gin.H{"error": "Akses tidak diizinkan"})
			c.Abort()
			return
		}

		// melanjutkan ke handler selanjutnya
		c.Next()
	}
}
