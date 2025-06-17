package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
	Password     string `json:"password,omitempty"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Role         string `json:"role"`
	CreatedAt    string `json:"created_at"`
	LatitudeDMS  string `json:"latitude_dms"`
	LongitudeDMS string `json:"longitude_dms"`
}

func DMSStringToDecimal(dms string) (float64, error) {
	// Example input: "14:09:12.42 N" or "121:15:58.30 E"
	parts := strings.Fields(dms)
	if len(parts) != 2 {
		return 0, errors.New("invalid DMS format")
	}
	dmsPart := parts[0]
	direction := parts[1]

	dmsSplit := strings.Split(dmsPart, ":")
	if len(dmsSplit) != 3 {
		return 0, errors.New("invalid DMS values")
	}

	degrees, err := strconv.ParseFloat(dmsSplit[0], 64)
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.ParseFloat(dmsSplit[1], 64)
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.ParseFloat(dmsSplit[2], 64)
	if err != nil {
		return 0, err
	}

	decimal := degrees + (minutes / 60) + (seconds / 3600)
	if direction == "S" || direction == "W" {
		decimal = -decimal
	}

	return decimal, nil
}

func RegisterHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("📥 Received register request")

		var input struct {
			Username     string `json:"username"`
			Password     string `json:"password"`
			FullName     string `json:"full_name"`
			Email        string `json:"email"`
			PhoneNumber  string `json:"phone_number"`
			LatitudeDMS  string `json:"latitude_dms"`
			LongitudeDMS string `json:"longitude_dms"`
		}

		if err := c.BodyParser(&input); err != nil {
			log.Println("❌ Failed to parse input:", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		log.Printf("Parsed: %s (%s, %s)", input.Username, input.LatitudeDMS, input.LongitudeDMS)

		if input.Username == "" || input.Password == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Username and password are required"})
		}

		// Convert DMS to decimal
		latitudeDecimal, err := DMSStringToDecimal(input.LatitudeDMS)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude format"})
		}
		longitudeDecimal, err := DMSStringToDecimal(input.LongitudeDMS)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude format"})
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("❌ Password hash error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Hashing error"})
		}

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "DB transaction failed"})
		}

		var userID int
		err = tx.QueryRow(`
			INSERT INTO Users (username, password_hash, full_name, email, phone_number)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING user_id
		`, input.Username, string(hash), input.FullName, input.Email, input.PhoneNumber).Scan(&userID)

		if err != nil {
			tx.Rollback()
			log.Println("❌ Insert user error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = tx.Exec(`
			INSERT INTO LoftCoordinates (user_id, latitude_dms, longitude_dms, latitude, longitude)
			VALUES ($1, $2, $3, $4, $5)
		`, userID, input.LatitudeDMS, input.LongitudeDMS, latitudeDecimal, longitudeDecimal)

		if err != nil {
			tx.Rollback()
			log.Println("❌ Insert loft coordinates error:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err := tx.Commit(); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Commit failed"})
		}

		log.Printf("✅ Registered user %s with loft coordinates\n", input.Username)
		return c.JSON(fiber.Map{"message": "User registered with loft location"})
	}
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

// func RegisterHandler(db *sql.DB) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		log.Println("📥 Received register request")

// 		var input struct {
// 			Username    string `json:"username"`
// 			Password    string `json:"password"`
// 			FullName    string `json:"full_name"`
// 			Email       string `json:"email"`
// 			PhoneNumber string `json:"phone_number"`
// 		}

// 		if err := c.BodyParser(&input); err != nil {
// 			log.Println("❌ Failed to parse input:", err)
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
// 		}

// 		log.Printf("Parsed data: username=%s fullName=%s phone=%s\n", input.Username, input.FullName, input.PhoneNumber)

// 		if input.Username == "" || input.Password == "" {
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Username and password are required"})
// 		}

// 		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// 		if err != nil {
// 			log.Println("❌ Failed to hash password:", err)
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Hashing error"})
// 		}

// 		_, err = db.Exec(`
// 			INSERT INTO Users (username, password_hash, full_name, email, phone_number)
// 			VALUES ($1, $2, $3, $4, $5)
// 		`, input.Username, string(hash), input.FullName, input.Email, input.PhoneNumber)
// 		if err != nil {
// 			log.Println("❌ Failed to insert user:", err)
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		log.Printf("✅ Registered user: %s\n", input.Username)
// 		return c.JSON(fiber.Map{"message": "User registered successfully"})
// 	}
// }

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

		rows, err := db.Query(`
			SELECT u.user_id, u.username, u.full_name, u.email, u.phone_number, u.role, u.created_at,
				l.latitude_dms, l.longitude_dms
			FROM Users u
			LEFT JOIN LoftCoordinates l ON u.user_id = l.user_id
		`)

		if err != nil {
			log.Println("❌ Failed to fetch users:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var users []User
		for rows.Next() {
			var u User
			err := rows.Scan(&u.UserID, &u.Username, &u.FullName, &u.Email, &u.PhoneNumber, &u.Role, &u.CreatedAt, &u.LatitudeDMS, &u.LongitudeDMS)
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
