package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/th0th/poeticmetric/backend/pkg/restapi/helpers"
	am "github.com/th0th/poeticmetric/backend/pkg/restapi/middleware/authentication"
	dm "github.com/th0th/poeticmetric/backend/pkg/restapi/middleware/depot"
	user2 "github.com/th0th/poeticmetric/backend/pkg/service/user"
	"github.com/th0th/poeticmetric/backend/pkg/validator"
)

func read(c *fiber.Ctx) error {
	dp := dm.Get(c)
	auth := am.Get(c)

	id, err := helpers.IdParam(c, "id")
	if err != nil {
		return err
	}

	if !validator.OrganizationUserId(dp, auth.Organization.Id, id) {
		return fiber.ErrNotFound
	}

	user, err := user2.Read(dp, id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}
