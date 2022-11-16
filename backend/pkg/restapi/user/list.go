package user

import (
	"github.com/gofiber/fiber/v2"
	am "github.com/poeticmetric/poeticmetric/backend/pkg/restapi/middleware/authentication"
	dm "github.com/poeticmetric/poeticmetric/backend/pkg/restapi/middleware/depot"
	"github.com/poeticmetric/poeticmetric/backend/pkg/service/user"
)

func list(c *fiber.Ctx) error {
	dp := dm.Get(c)
	auth := am.Get(c)

	users, err := user.List(dp, auth.Organization.Id)
	if err != nil {
		return err
	}

	return c.JSON(users)
}
