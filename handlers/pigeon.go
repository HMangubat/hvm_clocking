package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Pigeon struct {
    PigeonID   int    `json:"pigeon_id"`
    RingNumber string `json:"ring_number"`
    Name       string `json:"name"`
    Breed      string `json:"breed"`
    Gender     string `json:"gender"`
    LoftID     int    `json:"loft_id"`
    OwnerID    int    `json:"owner_id"`
}

func GetAllPigeons(db *sql.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        rows, err := db.Query(`SELECT pigeon_id, ring_number, name, breed, gender, loft_id, owner_id FROM Pigeons`)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        var pigeons []Pigeon
        for rows.Next() {
            var p Pigeon
            _ = rows.Scan(&p.PigeonID, &p.RingNumber, &p.Name, &p.Breed, &p.Gender, &p.LoftID, &p.OwnerID)
            pigeons = append(pigeons, p)
        }

        return c.JSON(pigeons)
    }
}
