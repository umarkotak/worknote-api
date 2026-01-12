package auth_handler

import (
	"github.com/gofiber/fiber/v2"

	"worknote-api/contract"
	"worknote-api/services/auth_service"
	"worknote-api/utils/render"
)

// GoogleAuth handles POST /auth/google
func GoogleAuth(c *fiber.Ctx) error {
	var req contract.GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return render.BadRequest(c, "invalid request body")
	}

	if req.IDToken == "" {
		return render.BadRequest(c, "id_token is required")
	}

	// Authenticate with Google
	authResp, err := auth_service.AuthenticateWithGoogle(c.Context(), req.IDToken)
	if err != nil {
		return render.Unauthorized(c, err.Error())
	}

	return render.JSON(c, fiber.StatusOK, authResp)
}
