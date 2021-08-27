package routes

import (
	"shortener-app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/", controllers.Default)
	api.Get("/l/:id", controllers.Redirect)
	api.Post("/new", controllers.NewShorted)
	api.Post("/signup", controllers.SignUp)
	api.Post("/login", controllers.Login)

	profile := api.Group("/profile", controllers.AuthMiddleware)
	profile.Post("/update", controllers.UpdateProfile)
	profile.Get("/links", controllers.ListLinks)
	api.Delete("/delete", controllers.DeleteUser)
}
