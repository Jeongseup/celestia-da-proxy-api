package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"

	openrpc "github.com/celestiaorg/celestia-openrpc"
	blobtypes "github.com/celestiaorg/celestia-openrpc/types/blob"
	"github.com/celestiaorg/celestia-openrpc/types/header"
	"github.com/celestiaorg/celestia-openrpc/types/share"
)

// must be namespace less than 10bytes
var (
	encodedNamesapce         = base64.StdEncoding.EncodeToString([]byte("CelestiaDragons"))
	celestiaDragonsNamespace = []byte(encodedNamesapce)[:10]
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

	return *resp, nil
}

func GetBlob(ctx context.Context, url string, token string, height uint64, hashStr string) (*blobtypes.Blob, error) {
	client, err := openrpc.NewClient(ctx, url, token)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	// let's post to 0xDEADBEEF namespace
	namespace, err := share.NewBlobNamespaceV0(celestiaDragonsNamespace)
	if err != nil {
		return nil, err
	}

	// share.NewBlobV0(namespace data []byte) (*Blob, error) {
	// hex 값을 string으로 표현
	// hexString := "13BC10A978B617DB8F7837A0B62E7C5FAE843876c916D83EF18D6005466CDF07"

	// string을 []byte로 변환
	data, err := hex.DecodeString(hashStr)
	if err != nil {
		return nil, err
	}

	// fetch the blob back from the network
	blob, err := client.Blob.Get(ctx, height, namespace, data)
	if err != nil {
		return nil, err
	}

	return blob, nil
}

// Blob was included at height 1826272
// Blobs are equal? false
func GetBlobs(ctx context.Context, url string, token string, height uint64) ([]*blobtypes.Blob, error) {
	client, err := openrpc.NewClient(ctx, url, token)
	if err != nil {
		return nil, err
	}

	defer client.Close()

	// let's post to 0xDEADBEEF namespace
	namespace, err := share.NewBlobNamespaceV0(celestiaDragonsNamespace)
	if err != nil {
		return nil, err
	}
	// fetch the blob back from the network
	retrievedBlobs, err := client.Blob.GetAll(ctx, height, []share.Namespace{namespace})
	if err != nil {
		return nil, err
	}

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
	// If it is less than 10 bytes, it will be left padded to size 10 with 0s.
	// namespace, err := share.NewBlobNamespaceV0([]byte{0xDE, 0xAD, 0xBE, 0xEF})
	l.Infof("using namespace bytes: %v", celestiaDragonsNamespace)
	l.Infof("using namespace bytes: %X", celestiaDragonsNamespace)
	l.Infof("using namespace bytes: %s", celestiaDragonsNamespace)
	namespace, err := share.NewBlobNamespaceV0(celestiaDragonsNamespace)
	if err != nil {
		return 0, err
	}

	// create a blob
	createdBlob, err := blobtypes.NewBlobV0(namespace, data)
	if err != nil {
		return 0, err
	}

	// submit the blob to the network
	height, err := client.Blob.Submit(ctx, []*blobtypes.Blob{createdBlob}, openrpc.DefaultGasPrice())
	if err != nil {
		return 0, err
	}

	l.Infof("Blob was included at height %d\n", height)

	// fetch the blob back from the network
	retrievedBlobs, err := client.Blob.GetAll(ctx, height, []share.Namespace{namespace})
	if err != nil {
		return 0, err
	}

	// 그냥 해당 높이에 하나씩만 있다고 가정
	for _, blob := range retrievedBlobs {
		l.Printf("blob commitment: %v \n", blob.Commitment)
		l.Printf("blob Namespace: %v \n", blob.Namespace)
		l.Printf("blob NamespaceVersion: %d \n", blob.NamespaceVersion)
		l.Printf("blob Data: %d \n", len(blob.Data))
		l.Printf("blob index: %d \n", blob.Index)
	}

	// fmt.Printf("Blobs are equal? %v\n", bytes.Equal(helloWorldBlob.Commitment, retrievedBlobs[0].Commitment))
	return height, nil
}
