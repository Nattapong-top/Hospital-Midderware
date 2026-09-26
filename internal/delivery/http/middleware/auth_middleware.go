package middleware

import (
	"net/http"
	"strings"

	"Hospital-Midderware/internal/domain"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenProvider domain.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 1. เช็กว่ามี Header Authorization และขึ้นต้นด้วย "Bearer " หรือไม่
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		// 2. ตัดคำว่า "Bearer " ออก
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. ตรวจสอบความถูกต้องของ Token
		// (หมายเหตุ: ถ้าใน domain.TokenProvider ของป๋ามี Method ValidateToken ให้เรียกตรงนี้ได้เลย)
		claims, err := tokenProvider.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// 4. ฝากข้อมูลลงใน Context
		c.Set("username", claims.Username)
		c.Set("hospital", claims.Hospital)

		c.Next()
	}
}
