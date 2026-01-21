package apiserver

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/pkg/api"
)

// RequireAPIKey enforces CLI authentication for HTTP management APIs.
//
// Expected header:
//   Authorization: Bearer <api_key>
func RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := strings.TrimSpace(c.GetHeader("Authorization"))
		if authz == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{Success: false, Message: "missing Authorization header"})
			return
		}

		parts := strings.Fields(authz)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{Success: false, Message: "invalid Authorization header"})
			return
		}

		expected, err := configmanager.GetSystemConfig("api_key")
		if err != nil || expected == "" {
			// Server misconfigured: serve already checks existence, but keep defense-in-depth.
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.Response{Success: false, Message: "server api_key not configured"})
			return
		}

		if subtle.ConstantTimeCompare([]byte(parts[1]), []byte(expected)) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, api.Response{Success: false, Message: "invalid api key"})
			return
		}

		c.Next()
	}
}

