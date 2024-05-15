package main

import (
	"github.com/gofiber/fiber/v2"
)

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}

// @Summary Returns a hello message
// @Description Responds with a simple hello message
// @Produce json
// @Success 200 {object} Response
// @Router /hello [get]
func HelloCheck(c *fiber.Ctx) error {
	response := Response{
		Success: true,
		Result:  "hello",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// @Summary Returns an error message
// @Description Responds with a test error message
// @Produce json
// @Success 400 {object} Response
// @Router /error [get]
func ErrorCheck(c *fiber.Ctx) error {
	response := Response{
		Success: false,
		Error:   "error test",
	}

	return c.Status(fiber.StatusBadRequest).JSON(response)
}
