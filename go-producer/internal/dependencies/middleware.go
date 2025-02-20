package dependencies

import "github.com/gofiber/fiber/v2"

// InjectDependencies adds dependencies to Fiber context
func InjectDependencies(di *AppDependencies) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("di", di)
		return c.Next()
	}
}

// ExtractDependencies retrieves dependencies from Fiber context
func ExtractDependencies(c *fiber.Ctx) *AppDependencies {
	if di, ok := c.Locals("di").(*AppDependencies); ok {
		return di
	}
	return nil
}
