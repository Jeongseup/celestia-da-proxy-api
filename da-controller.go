package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

func SubmitJSONDataController(c *fiber.Ctx) error {
	var payload Data

	// 요청 본문을 JSON으로 파싱
	if err := c.BodyParser(&payload); err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var sendData []byte
	if len(payload.MetaData) > 0 {
		l.Infoln("received metadata!")
		sendData = payload.MetaData
	}

	if len(payload.Data) > 0 {
		l.Infoln("received data!")
		sendData = payload.Data
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// svgStr := "data:image/svg+xml;base64,CiAgICAgIDxzdmcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiB2aWV3Qm94PSIwIDAgMTAwIDEwMCI+CiAgICAgICAgPGNpcmNsZSBjeD0iNTAiIGN5PSI1MCIgcj0iNDAiIGZpbGw9IiNhZDExZjciIC8+CiAgICAgICAgPHRleHQgeD0iNTAiIHk9IjUwIiBmb250LXNpemU9IjEyIiBmb250LWZhbWlseT0ic2Fucy1zZXJpZiIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iIGZpbGw9IndoaXRlIj5IYWNrYXRoZW15PC90ZXh0PgogICAgICA8L3N2Zz4KICAgIA=="
	// svgBz := []byte(svgStr)

	// submit data
	height, err := SubmitBlob(ctx, celestiaRpcAddress, authToken, sendData)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	l.Infof("Blob was included at height %d\n", height)

	// fetch the blob back from the network
	// retrievedBlobs, err := GetAll(ctx, height, []share.Namespace{namespace})
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("Blobs are equal? %v\n", bytes.Equal(helloWorldBlob.Commitment, retrievedBlobs[0].Commitment))

	response := Response{
		Success: true,
		Result: fiber.Map{
			// "submitted_data": payload.Data,
			"submitted_height": height,
		},
	}
	return c.JSON(response)
}

func SubmitFormDataController(c *fiber.Ctx) error {
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

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// submit data
	height, err := SubmitBlob(ctx, celestiaRpcAddress, authToken, fileBytes)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	l.Infof("Blob was included at height %d\n", height)

	// retrieve blobs
	retrievedBlobs, err := GetBlobs(ctx, celestiaRpcAddress, authToken, height)
	if err != nil {
		// l.Errorf("Successfully, submitted formdata to celestia da but, failed to retrieve blobs in the height. unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   fmt.Sprintf("Successfully, submitted formdata to celestia da but, failed to retrieve blobs in the height. unexpected err: %s", err.Error()),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	var hashStr string
	for _, blob := range retrievedBlobs {
		if bytes.Equal(blob.Data, fileBytes) {
			// l.Infoln("equal!")
			l.Printf("blob commitment hash: %X \n", blob.Commitment)
			// l.Printf("blob commitment: %s \n", blob.Commitment)
			hashStr = fmt.Sprintf("%X", blob.Commitment)
			break
		}
	}

	// insert db hash for saving height
	err = InsertBlob(db, hashStr, int(height))
	if err != nil {
		resp := Response{
			Success: false,
			Error:   fmt.Sprintf("Successfully, submitted formdata to celestia da but, failed to save commitment hash & height in db. unexpected err: %s", err.Error()),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	response := Response{
		Success: true,
		Result: fiber.Map{
			"hash":             hashStr,
			"submitted_height": height,
		},
	}

	return c.JSON(response)
}

func RetrieveBlobController(c *fiber.Ctx) error {
	retrieveHeightStr := c.Query("height")

	l.Infof("received height: %s", retrieveHeightStr)

	retrieveHeight, err := strconv.ParseUint(retrieveHeightStr, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

	var data []byte

	// 각 블롭을 직렬화하고 리스트에 추가 (그냥 하나만 있다고 가정)
	for _, blob := range retrievedBlobs {
		l.Printf("blob commitment: %X \n", blob.Commitment)
		l.Printf("blob Namespace: %X \n", blob.Namespace)
		l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
		l.Printf("blob Data: %d \n", len(blob.Data))
		l.Printf("blob index: %d \n", blob.Index)
		l.Printf("blob ShareVersion: %d \n", blob.ShareVersion)

		data = blob.Data
		break
	}

	return c.Send(data)
}

func RetrieveBlobByCommitment(c *fiber.Ctx) error {

	hashStr := c.Params("hash") // URL 파라미터에서 해시 값 가져오기
	l.Infof("Received hash: %s\n", hashStr)

	retrieveHeight, err := GetBlobHeight(db, hashStr)
	if err != nil {
		return err
	}

	l.Infof("Found height by hash: %d\n", retrieveHeight)

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// submit data
	blob, err := GetBlob(ctx, celestiaRpcAddress, authToken, uint64(retrieveHeight), hashStr)
	if err != nil {
		l.Errorf("unexpected err: %s", err)
		resp := Response{
			Success: false,
			Error:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	l.Printf("blob commitment: %X \n", blob.Commitment)
	l.Printf("blob Namespace: %X \n", blob.Namespace)
	l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
	l.Printf("blob Data: %d \n", len(blob.Data))
	l.Printf("blob index: %d \n", blob.Index)

	return c.Send(blob.Data)
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
