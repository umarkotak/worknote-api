package render

import (
	"github.com/gofiber/fiber/v2"
)

// JSON writes a JSON response with the given status code
func JSON(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}

// Error writes a JSON error response
func Error(c *fiber.Ctx, status int, message string) error {
	return JSON(c, status, fiber.Map{"error": message})
}

// Unauthorized writes a 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// BadRequest writes a 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// Itoa converts an int64 to string without using strconv
func Itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte(n%10) + '0'
		n /= 10
	}
	return string(b[i:])
}
