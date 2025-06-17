package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// =========================== CLUBS ===========================
func CreateClubHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var club struct {
			Name     string `json:"name"`
			Location string `json:"location"`
		}

		if err := c.BodyParser(&club); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
		}

		_, err := db.Exec(`INSERT INTO Clubs (name, location) VALUES ($1, $2)`, club.Name, club.Location)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Club created"})
	}
}

func GetAllClubsHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(`SELECT club_id, name, location, created_at FROM Clubs`)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var clubs []map[string]interface{}
		for rows.Next() {
			var id int
			var name, location, createdAt string
			rows.Scan(&id, &name, &location, &createdAt)
			clubs = append(clubs, fiber.Map{
				"club_id":    id,
				"name":       name,
				"location":   location,
				"created_at": createdAt,
			})
		}

		return c.JSON(clubs)
	}
}

// =========================== DEVICES ===========================
func CreateDeviceHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var d struct {
			UserID       int    `json:"user_id"`
			Name         string `json:"name"`
			SerialNumber string `json:"serial_number"`
		}
		if err := c.BodyParser(&d); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		_, err := db.Exec(`INSERT INTO Devices (user_id, name, serial_number) VALUES ($1, $2, $3)`,
			d.UserID, d.Name, d.SerialNumber)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Device registered"})
	}
}

// =========================== LOFT COORDINATES ===========================
func CreateLoftHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var loft struct {
			UserID    int     `json:"user_id"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}
		if err := c.BodyParser(&loft); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}

		_, err := db.Exec(`INSERT INTO LoftCoordinates (user_id, latitude, longitude) VALUES ($1, $2, $3)`,
			loft.UserID, loft.Latitude, loft.Longitude)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Loft location saved"})
	}
}

// =========================== PIGEONS ===========================
func CreatePigeonHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p struct {
			UserID     int    `json:"user_id"`
			RingNumber string `json:"ring_number"`
			Name       string `json:"name"`
			Color      string `json:"color"`
			Sex        string `json:"sex"`
			Breed      string `json:"breed"`
			BirthDate  string `json:"birth_date"` // Format: YYYY-MM-DD
		}
		if err := c.BodyParser(&p); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		_, err := db.Exec(`
			INSERT INTO Pigeons (user_id, ring_number, name, color, sex, breed, birth_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			p.UserID, p.RingNumber, p.Name, p.Color, p.Sex, p.Breed, p.BirthDate)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Pigeon added"})
	}
}

// =========================== RACES ===========================
func CreateRaceHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var r struct {
			Name        string  `json:"name"`
			Location    string  `json:"location"`
			DistanceKM  float64 `json:"distance_km"`
			ReleaseTime string  `json:"release_time"` // Format: YYYY-MM-DD HH:MM:SS
		}
		if err := c.BodyParser(&r); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		_, err := db.Exec(`INSERT INTO Races (name, location, distance_km, release_time) VALUES ($1, $2, $3, $4)`,
			r.Name, r.Location, r.DistanceKM, r.ReleaseTime)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Race created"})
	}
}

// =========================== PARTICIPANTS ===========================
func RegisterPigeonToRaceHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			RaceID   int `json:"race_id"`
			PigeonID int `json:"pigeon_id"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		_, err := db.Exec(`INSERT INTO RaceParticipants (race_id, pigeon_id) VALUES ($1, $2)`,
			input.RaceID, input.PigeonID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Pigeon registered to race"})
	}
}

// =========================== CLOCKINGS ===========================
func ClockPigeonHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var clk struct {
			PigeonID int     `json:"pigeon_id"`
			RaceID   int     `json:"race_id"`
			UserID   int     `json:"user_id"`
			DeviceID int     `json:"device_id"`
			Arrival  string  `json:"arrival_time"` // YYYY-MM-DD HH:MM:SS
			SpeedKPH float64 `json:"speed_kph"`
		}
		if err := c.BodyParser(&clk); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		_, err := db.Exec(`
			INSERT INTO Clockings (pigeon_id, race_id, user_id, device_id, arrival_time, speed_kph)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			clk.PigeonID, clk.RaceID, clk.UserID, clk.DeviceID, clk.Arrival, clk.SpeedKPH)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Clocking recorded"})
	}
}

// =========================== RACE RESULTS ===========================
func InsertRaceResultHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var res struct {
			RaceID   int     `json:"race_id"`
			PigeonID int     `json:"pigeon_id"`
			SpeedKPH float64 `json:"speed_kph"`
			Arrival  string  `json:"arrival_time"` // YYYY-MM-DD HH:MM:SS
			Rank     int     `json:"rank"`
		}
		if err := c.BodyParser(&res); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		_, err := db.Exec(`
			INSERT INTO RaceResults (race_id, pigeon_id, speed_kph, arrival_time, rank)
			VALUES ($1, $2, $3, $4, $5)`,
			res.RaceID, res.PigeonID, res.SpeedKPH, res.Arrival, res.Rank)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Race result inserted"})
	}
}

// =========================== AUDIT LOGS ===========================
func LogAuditActionHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var logData struct {
			UserID int    `json:"user_id"`
			Action string `json:"action"`
		}
		if err := c.BodyParser(&logData); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		_, err := db.Exec(`INSERT INTO AuditLogs (user_id, action) VALUES ($1, $2)`, logData.UserID, logData.Action)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Audit log recorded"})
	}
}
