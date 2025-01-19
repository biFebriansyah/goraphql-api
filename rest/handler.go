package rest

import (
	"fmt"
	"strings"

	"github.com/biFebriansyah/goraphql/graph/service"
	"github.com/biFebriansyah/goraphql/utils"
	"github.com/gofiber/fiber/v2"
)

type RestHandler struct {
	UserService *service.UserService
}

func (rest *RestHandler) SignIn(ctx *fiber.Ctx) error {
	body := new(SigninInput)
	if err := ctx.BodyParser(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userData, err := rest.UserService.GetByEmail(body.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if !utils.CheckPasswordHash(body.Password, userData.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "worng password")
	}

	token, err := utils.GenerateJwt(userData.ID, *userData.Admin)
	if err != nil {
		return fmt.Errorf("fail when generate token: %w", err)
	}

	return ctx.JSON(fiber.Map{"token": token})
}

func (rest *RestHandler) AuthMiddleware(ctx *fiber.Ctx) error {
	headers := ctx.GetReqHeaders()["Authorization"]
	if len(headers) <= 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "need login")
	}

	tokens := strings.Replace(headers[0], "Bearer ", "", 1)
	claims, err := utils.ParseJwt(tokens)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	ctx.Locals("userId", claims.ID)
	return ctx.Next()
}
