package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Race struct {
    RaceID       int     `json:"race_id"`
    Name         string  `json:"name"`
    ReleasePoint string  `json:"release_point"`
    ReleaseLat   float64 `json:"release_lat"`
    ReleaseLng   float64 `json:"release_lng"`
    ScheduledAt  string  `json:"scheduled_at"`
    DistanceKm   float64 `json:"distance_km"`
    Status       string  `json:"status"`
}

func GetAllRaces(db *sql.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        rows, err := db.Query(`SELECT race_id, name, release_point, release_lat, release_lng, scheduled_at, distance_km, status FROM Races`)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        var races []Race
        for rows.Next() {
            var r Race
            _ = rows.Scan(&r.RaceID, &r.Name, &r.ReleasePoint, &r.ReleaseLat, &r.ReleaseLng, &r.ScheduledAt, &r.DistanceKm, &r.Status)
            races = append(races, r)
        }

        return c.JSON(races)
    }
}
