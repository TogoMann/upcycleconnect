package middlewares

import (
	db "backend/internal/database"
	"backend/internal/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func HasPlan(allowedTiers ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return utils.JwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
				return
			}

			role, _ := claims["role"].(string)
			if role == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			userIdFloat, ok := claims["sub"].(float64)
			if !ok {
				http.Error(w, "Unauthorized: Invalid user ID", http.StatusUnauthorized)
				return
			}
			userId := int64(userIdFloat)

			var tier string
			var until time.Time
			err = db.Pool.QueryRow(r.Context(), "SELECT tier, until FROM subscriptions WHERE subscriber_id = $1 AND until >= CURRENT_DATE ORDER BY until DESC LIMIT 1", userId).Scan(&tier, &until)

			if err != nil {
				tier = "Free"
			}

			isAllowed := false
			for _, allowed := range allowedTiers {
				if strings.EqualFold(tier, allowed) {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, fmt.Sprintf("Forbidden: This feature requires one of the following plans: %v", allowedTiers), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
