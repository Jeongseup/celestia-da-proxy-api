package main

import (
	"context"
	"log"

	openrpc "github.com/celestiaorg/celestia-openrpc"
	"github.com/celestiaorg/celestia-openrpc/types/blob"
	"github.com/celestiaorg/celestia-openrpc/types/header"
	"github.com/celestiaorg/celestia-openrpc/types/share"
)

func NodePing(ctx context.Context, url string, token string) (header.ExtendedHeader, error) {
	client, err := openrpc.NewClient(context.Background(), url, token)
	if err != nil {
		return header.ExtendedHeader{}, err
	}
	defer client.Close()

	resp, err := client.Header.NetworkHead(ctx)
	if err != nil {
		return header.ExtendedHeader{}, err
	}

	// log.Printf("chain id: %s", resp.ChainID())
	// log.Printf("height: %d", resp.Height())

	// info, err := client.Node.Info(ctx)
	// if err != nil {
	// 	return header.ExtendedHeader{}, err
	// }

	// log.Printf("api version: %s", info.APIVersion)

	// Type defines the Node type (e.g. `light`, `bridge`) for identity purposes.
	// The zero value for Type is invalid.
	// log.Printf("node type: %v", info.Type)

	return *resp, nil
}

// Blob was included at height 1826272
// Blobs are equal? false
func GetBlobs(ctx context.Context, url string, token string, height uint64) ([]*blob.Blob, error) {
	client, err := openrpc.NewClient(ctx, url, token)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	// let's post to 0xDEADBEEF namespace
	namespace, err := share.NewBlobNamespaceV0([]byte{0xDE, 0xAD, 0xBE, 0xEF})
	if err != nil {
		return nil, err
	}
	// fetch the blob back from the network
	retrievedBlobs, err := client.Blob.GetAll(ctx, height, []share.Namespace{namespace})
	if err != nil {
		return nil, err
	}

	// for _, blob := range retrievedBlobs {
	// 	jsonBz, err := blob.MarshalJSON()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	fmt.Printf("json: %s\n", jsonBz) // base64 encoded

	// 	fmt.Printf("blob commitment: %X \n", blob.Commitment)
	// 	fmt.Printf("blob Namespace: %X \n", blob.Namespace)
	// 	fmt.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
	// 	fmt.Printf("blob commitment: %s \n", blob.Data)
	// 	fmt.Printf("blob index: %d \n", blob.Index)
	// 	fmt.Printf("blob ShareVersion: %d \n", blob.ShareVersion)
	// 	// fmt.Printf("%v blob in %d height\n", blob, height)
	// 	// fmt.Printf("%v blob in %d height\n", blob, height)
	// 	// fmt.Printf("%v blob in %d height\n", blob, height)
	// 	// fmt.Printf("%v blob in %d height\n", blob, height)
	// }

	// fmt.Printf("Blobs are equal? %v ", retrievedBlobs[0].Commitment)
	return retrievedBlobs, nil
}

// SubmitBlob submits a blob containing "Hello, World!" to the 0xDEADBEEF namespace. It uses the default signer on the running node.
func SubmitBlob(ctx context.Context, url string, token string, data []byte) (uint64, error) {
	client, err := openrpc.NewClient(ctx, url, token)
	if err != nil {
		return 0, err
	}

	defer client.Close()

	// let's post to 0xDEADBEEF namespace
	namespace, err := share.NewBlobNamespaceV0([]byte{0xDE, 0xAD, 0xBE, 0xEF})
	if err != nil {
		return 0, err
	}

	// create a blob
	createdBlob, err := blob.NewBlobV0(namespace, data)
	if err != nil {
		return 0, err
	}

	// submit the blob to the network
	height, err := client.Blob.Submit(ctx, []*blob.Blob{createdBlob}, openrpc.DefaultGasPrice())
	if err != nil {
		return 0, err
	}

	l.Infof("Blob was included at height %d\n", height)

	// fetch the blob back from the network
	// retrievedBlobs, err := client.Blob.GetAll(ctx, height, []share.Namespace{namespace})
	// if err != nil {
	// 	return err
	// }

	// for
	// l.Infoln(retrievedBlobs)
	// fmt.Printf("Blobs are equal? %v\n", bytes.Equal(helloWorldBlob.Commitment, retrievedBlobs[0].Commitment))
	return height, nil
}

func Balance(ctx context.Context, url string, token string) error {
	client, err := openrpc.NewClient(ctx, url, token)
	if err != nil {
		return err
	}

	defer client.Close()

	// account, err := client.State.AccountAddress(ctx)
	// if err != nil {
	// 	fmt.Println("here")
	// 	return err

	// }

	log.Println("account address : celestia1ynrht4f0jnltat30vqlzcw62z3jra49fc0xdx5")

	balance, err := client.State.Balance(ctx)
	if err != nil {
		return err
	}

	log.Println(balance.Denom, balance.Amount)
	return nil
}
