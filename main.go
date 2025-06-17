// package main

// import (
// 	"hvm_clocking/config"
// 	"hvm_clocking/handlers"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	app := fiber.New()
// 	db := config.InitDB()
// 	defer db.Close()

// 	// Static files
// 	app.Static("/static", "./static")

// 	// Templates (rendering if needed)
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendFile("./templates/login.html")
// 	})
// 	app.Get("/dashboard", func(c *fiber.Ctx) error {
// 		return c.SendFile("./templates/dashboard.html")
// 	})
// 	app.Get("/register", func(c *fiber.Ctx) error {
// 		return c.SendFile("./templates/register.html")
// 	})

// 	// Routes
// 	app.Post("/register", handlers.RegisterHandler(db))
// 	app.Post("/login", handlers.LoginHandler(db))

// 	app.Get("/users", handlers.GetAllUsers(db))
// 	app.Get("/pigeons", handlers.GetAllPigeons(db))
// 	app.Get("/lofts", handlers.GetAllLofts(db))
// 	app.Get("/races", handlers.GetAllRaces(db))

// 	app.Listen(":2000")
// }

package main

import (
	"hvm_clocking/config"
	"hvm_clocking/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Load HTML templates from ./templates with .html extension
	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	db := config.InitDB()
	defer db.Close()

	// Static files (CSS, JS, images)
	app.Static("/static", "./static")

	// View-rendered pages
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{}) // login.html
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.Render("dashboard", fiber.Map{
			"ActivePage": "dashboard",
		})
	})
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.Render("users", fiber.Map{
			"ActivePage": "users",
		})
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{}) // register.html
	})

	// Handlers
	app.Post("/register", handlers.RegisterHandler(db))
	app.Post("/login", handlers.LoginHandler(db))

	//app.Get("/users", handlers.GetAllUsers(db))
	app.Get("/api/users", handlers.GetAllUsers(db))
	app.Put("/api/users/:id", handlers.UpdateUser(db))

	app.Get("/pigeons", handlers.GetAllPigeons(db))
	app.Get("/lofts", handlers.GetAllLofts(db))
	app.Get("/races", handlers.GetAllRaces(db))

	app.Post("/api/clubs", handlers.CreateClubHandler(db))
	app.Get("/api/clubs", handlers.GetAllClubsHandler(db))

	app.Post("/api/devices", handlers.CreateDeviceHandler(db))
	app.Post("/api/lofts", handlers.CreateLoftHandler(db))
	app.Post("/api/pigeons", handlers.CreatePigeonHandler(db))
	app.Post("/api/races", handlers.CreateRaceHandler(db))
	app.Post("/api/race-participants", handlers.RegisterPigeonToRaceHandler(db))
	app.Post("/api/clockings", handlers.ClockPigeonHandler(db))
	app.Post("/api/race-results", handlers.InsertRaceResultHandler(db))
	app.Post("/api/audit-logs", handlers.LogAuditActionHandler(db))

	app.Listen(":2000")
}
