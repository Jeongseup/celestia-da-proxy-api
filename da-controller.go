package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary Returns Celestia DA node info
// @Description Pings Celestia DA node and returns node info
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /node_info [get]
func NodeInfoController(c *fiber.Ctx) error {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.Println("celestia da rpc node ping...")

	// get rpc da info
	headerInfo, err := NodePing(ctx, celestiaRpcAddress, authToken)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	// jsonBz, _ := headerInfo.MarshalJSON()
	// var payload map[string]interface{}
	// _ = json.Unmarshal(jsonBz, &payload)
	// log.Infof("%s", payload)

	response := Response{
		Success: true,
		Result: fiber.Map{
			"height":     headerInfo.Height(),
			"chain_id":   headerInfo.ChainID(),
			"block_hash": headerInfo.Hash(),
			"timestamp":  headerInfo.Time(),
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func SubmitBlobController(c *fiber.Ctx) error {
	var payload Data

	// 요청 본문을 JSON으로 파싱
	if err := c.BodyParser(&payload); err != nil {
		response := Response{
			Success: false,
			Error:   "cannot parse JSON",
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	l.Infof("received data: %v", payload.Data)

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	svgStr := "data:image/svg+xml;base64,CiAgICAgIDxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB2aWV3Qm94PSIwIDAgMTAwIDEwMCI+CiAgICAgICAgPGNpcmNsZSBjeD0iNTAiIGN5PSI1MCIgcj0iNDAiIGZpbGw9IiNhZDExZjciIC8+CiAgICAgICAgPHRleHQgeD0iNTAiIHk9IjUwIiBmb250LXNpemU9IjEyIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iIGZpbGw9IndoaXRlIj5IYWNrYXRoZW15PC90ZXh0PgogICAgICA8L3N2Zz4KICAgIA=="
	svgBz := []byte(svgStr)

	// submit data
	height, err := SubmitBlob(ctx, celestiaRpcAddress, authToken, svgBz)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	response := Response{
		Success: true,
		Result: fiber.Map{
			"submitted_data": svgBz,
			"height":         height,
		},
	}
	return c.JSON(response)
}

func RetrieveBlobController(c *fiber.Ctx) error {
	var payload Data

	// 요청 본문을 JSON으로 파싱
	if err := c.BodyParser(&payload); err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	l.Infof("received height: %s", payload.RetrieveHeight)

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	retrieveHeight, err := strconv.ParseUint(payload.RetrieveHeight, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// submit data
	retrievedBlobs, err := GetBlobs(ctx, celestiaRpcAddress, authToken, retrieveHeight)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	var jsonBlobs []json.RawMessage

	// 각 블롭을 직렬화하고 리스트에 추가
	for _, blob := range retrievedBlobs {
		l.Printf("blob commitment: %X \n", blob.Commitment)
		l.Printf("blob Namespace: %X \n", blob.Namespace)
		l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
		l.Printf("blob Data: %s \n", blob.Data)
		l.Printf("blob index: %d \n", blob.Index)
		l.Printf("blob ShareVersion: %d \n", blob.ShareVersion)

		jsonBz, err := blob.MarshalJSON()
		if err != nil {
			return err
		}

		jsonBlobs = append(jsonBlobs, jsonBz)
	}

	response := Response{
		Success: true,
		Result:  jsonBlobs,
	}

	return c.JSON(response)
}

func TestBlobController(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	l.Println("=== test case 4 : GetBlobs ===")
	retrievedBlobs, err := GetBlobs(ctx, celestiaRpcAddress, authToken, 1826272)
	if err != nil {
		log.Println(err)
	}
	l.Println("================================")

	for _, blob := range retrievedBlobs {
		jsonBz, err := blob.MarshalJSON()
		if err != nil {
			return err
		}

		l.Infof("json: %s\n", jsonBz) // base64 encoded
		l.Printf("blob commitment: %X \n", blob.Commitment)
		l.Printf("blob Namespace: %X \n", blob.Namespace)
		l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
		l.Printf("blob Data: %s \n", blob.Data)
		l.Printf("blob index: %d \n", blob.Index)
		l.Printf("blob ShareVersion: %d \n", blob.ShareVersion)
	}

	return c.SendStatus(fiber.StatusOK)
}
