package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type ClockingDevice struct {
    DeviceID     int    `json:"device_id"`
    DeviceNumber string `json:"device_number"`
    LoftID       int    `json:"loft_id"`
    Type         string `json:"type"`
    Status       string `json:"status"`
}

func GetAllDevices(db *sql.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        rows, err := db.Query("SELECT device_id, device_number, loft_id, type, status FROM ClockingDevices")
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }

        var devices []ClockingDevice
        for rows.Next() {
            var d ClockingDevice
            _ = rows.Scan(&d.DeviceID, &d.DeviceNumber, &d.LoftID, &d.Type, &d.Status)
            devices = append(devices, d)
        }
        return c.JSON(devices)
    }
}
