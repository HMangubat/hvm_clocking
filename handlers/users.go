package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
}

// func RegisterHandler(db *sql.DB) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		log.Println("📥 Received register request")

// 		username := c.FormValue("username")
// 		password := c.FormValue("password")
// 		fullName := c.FormValue("full_name")
// 		email := c.FormValue("email")
// 		phone := c.FormValue("phone_number")

// 		log.Printf("Parsed data: username=%s fullName=%s phone=%s", username, fullName, phone)

// 		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 		if err != nil {
// 			log.Println("❌ Failed to hash password:", err)
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "hashing error"})
// 		}

// 		_, err = db.Exec(`
// 			INSERT INTO Users (username, password_hash, full_name, email, phone_number)
// 			VALUES ($1, $2, $3, $4, $5)
// 		`, username, string(hash), fullName, email, phone)
// 		if err != nil {
// 			log.Println("❌ Failed to insert user:", err)
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		log.Printf("✅ Registered user: %s\n", username)
// 		return c.JSON(fiber.Map{"message": "registered"})
// 	}
// }

func RegisterHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("📥 Received register request")

		var input struct {
			Username    string `json:"username"`
			Password    string `json:"password"`
			FullName    string `json:"full_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
		}

		if err := c.BodyParser(&input); err != nil {
			log.Println("❌ Failed to parse input:", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		log.Printf("Parsed data: username=%s fullName=%s phone=%s\n", input.Username, input.FullName, input.PhoneNumber)

		if input.Username == "" || input.Password == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Username and password are required"})
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("❌ Failed to hash password:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Hashing error"})
		}

		_, err = db.Exec(`
			INSERT INTO Users (username, password_hash, full_name, email, phone_number)
			VALUES ($1, $2, $3, $4, $5)
		`, input.Username, string(hash), input.FullName, input.Email, input.PhoneNumber)
		if err != nil {
			log.Println("❌ Failed to insert user:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		log.Printf("✅ Registered user: %s\n", input.Username)
		return c.JSON(fiber.Map{"message": "User registered successfully"})
	}
}

func LoginHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}

		var input LoginInput
		log.Println("📥 Received login request")

		if err := c.BodyParser(&input); err != nil {
			log.Println("❌ Failed to parse login input:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		log.Printf("🔎 Looking up user: %s\n", input.Username)
		var hashedPassword string
		query := "SELECT password_hash FROM Users WHERE username=$1"
		err := db.QueryRow(query, input.Username).Scan(&hashedPassword)
		if err != nil {
			log.Println("❌ User not found or DB error:", err)
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
		}

		log.Println("🔐 Verifying password")
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
			log.Println("❌ Password verification failed")
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid username or password")
		}

		log.Printf("✅ User %s logged in successfully\n", input.Username)
		return c.JSON(fiber.Map{
			"status":   "success",
			"message":  "Login successful",
			"redirect": "/dashboard",
		})

	}
}

func GetAllUsers(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("📤 Fetching all users from database")

		rows, err := db.Query(`SELECT user_id, username, full_name, email, phone_number, role, created_at FROM Users`)
		if err != nil {
			log.Println("❌ Failed to fetch users:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var users []User
		for rows.Next() {
			var u User
			err := rows.Scan(&u.UserID, &u.Username, &u.FullName, &u.Email, &u.PhoneNumber, &u.Role, &u.CreatedAt)
			if err != nil {
				log.Println("⚠️ Failed to scan user row:", err)
				continue
			}
			users = append(users, u)
		}

		log.Printf("✅ Fetched %d users\n", len(users))
		return c.JSON(users)
	}
}

func UpdateUser(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var u struct {
			FullName    string `json:"full_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
		}

		if err := c.BodyParser(&u); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		_, err := db.Exec(`UPDATE Users SET full_name=$1, email=$2, phone_number=$3 WHERE user_id=$4`,
			u.FullName, u.Email, u.PhoneNumber, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "User updated"})
	}
}
