package root

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/poeticmetric/poeticmetric/backend/pkg/frontend"
)

func Add(app *fiber.App) {
	app.Get("/", index)
}

func index(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("%s%s", c.Protocol(), frontend.GenerateUrl("/")))
}
