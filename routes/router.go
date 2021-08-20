package routes

import (
	"shortener-app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/", controllers.Default)
	app.Get("/l/:id", controllers.Redirect)
	app.Post("/new", controllers.NewShorted)
	app.Post("/signup", controllers.SignUp)
	app.Post("/login", controllers.Login)

	profile := app.Group("/profile", controllers.AuthMiddleware)
	profile.Post("/update", controllers.UpdateProfile)
}
