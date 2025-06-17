package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Loft struct {
    LoftID    int     `json:"loft_id"`
    Name      string  `json:"name"`
    Location  string  `json:"location"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    UserID    int     `json:"user_id"`
}

func GetAllLofts(db *sql.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        rows, err := db.Query("SELECT loft_id, name, location, latitude, longitude, user_id FROM Lofts")
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        var lofts []Loft
        for rows.Next() {
            var l Loft
            _ = rows.Scan(&l.LoftID, &l.Name, &l.Location, &l.Latitude, &l.Longitude, &l.UserID)
            lofts = append(lofts, l)
        }
        return c.JSON(lofts)
    }
}
