package main

import (
	"encoding/base64"
	"io"

	"github.com/gofiber/fiber/v2"
)

func ReceiveFormData(c *fiber.Ctx) error {
	// 파일 업로드 처리
	file, err := c.FormFile("image")
	if err != nil {
		l.Errorln(err)

		response := Response{
			Success: false,
			Error:   "cannot parse form file",
		}

		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	l.Infoln(file.Filename)

	// 파일 읽기
	f, err := file.Open()
	if err != nil {
		response := Response{
			Success: false,
			Error:   "cannot open uploaded file",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		response := Response{
			Success: false,
			Error:   "cannot read uploaded file",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// 기본 정보를 로그로 출력
	l.Infof("received file: %s", file.Filename)
	l.Infof("file size: %d bytes", len(fileBytes))

	response := Response{
		Success: true,
		Result:  base64.StdEncoding.EncodeToString(fileBytes),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func ReceiveJSON(c *fiber.Ctx) error {
	var payload Data

	if err := c.BodyParser(&payload); err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	l.Infof("received data: %v", payload.Data)

	response := Response{
		Success: true,
		Result:  payload.Data,
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
