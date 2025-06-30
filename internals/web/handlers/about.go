package handlers

import "github.com/gofiber/fiber/v2"

func (h *GradeHandler) HandleAbout(c *fiber.Ctx) error {
	return c.Render("about", fiber.Map{})
}
