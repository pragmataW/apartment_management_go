package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware(jwtKey string, expectedRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Token'ı al
        token := c.Cookies("Authentication")

        // Token olup olmadığını kontrol et
        if token == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Unauthorized: Missing JWT token",
            })
        }

        // Token'ı doğrula
        claims := jwt.MapClaims{}
        tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(jwtKey), nil
        })

        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "message": "Unauthorized: Invalid JWT signature",
                })
            }
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Internal Server Error: JWT parsing error",
            })
        }

        if !tkn.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Unauthorized: Invalid JWT",
            })
        }

        // Email doğrulaması yap
        email, ok := claims["email"].(string)
        if !ok {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Internal Server Error: Could not extract email from JWT claims",
            })
        }

        // Role doğrulaması yap
        role, ok := claims["role"].(string)
        if !ok {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Internal Server Error: Could not extract role from JWT claims",
            })
        }

        // Beklenen roller arasında mı kontrol et
        roleAllowed := false
        for _, expectedRole := range expectedRoles {
            if role == expectedRole {
                roleAllowed = true
                break
            }
        }

        if !roleAllowed {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "message": "Unauthorized: Insufficient permissions",
            })
        }

        c.Locals("email", email)

        // Middleware'i geç
        return c.Next()
    }
}