package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestBlobController(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	l.Println("=== test case 4 : GetBlobs ===")
	retrievedBlobs, err := GetBlobs(ctx, celestiaRpcAddress, authToken, 1826272, celestiaDragonsNamespace)
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
