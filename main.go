package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	_ "github.com/Jeongseup/celestia-da-proxy-api/docs" // yourproject 경로를 실제 프로젝트 경로로 변경하세요
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	celestiaRpcAddress string
	authToken          string
	l                  *logrus.Logger
	db                 *sql.DB
)

// @title Fiber Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
func main() {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// get celestia da rpc infos from envs
	celestiaRpcAddress = os.Getenv("CELESTIA_DA_RPC_ADDRESSS")
	authToken = os.Getenv("RPC_AUTH_TOKEN")

	// logrus 설정
	l = logrus.New()
	l.SetOutput(os.Stdout)

	// 환경 변수로부터 로그 레벨 설정
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // 기본 로그 레벨
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	l.SetLevel(level)

	// Initialize SQLite database
	db, err = InitDB("./data.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app := fiber.New()

	// Fiber 로거 미들웨어를 logrus와 통합
	app.Use(logger.New(logger.Config{
		Output: log.Writer(),
	}))

	// Swagger 문서 라우트 설정
	app.Get("/swagger/*", swagger.HandlerDefault)

	// routes for test
	app.Post("/test_receive_jsondata", ReceiveJSON)
	app.Post("/test_receive_formdata", ReceiveFormData)
	app.Get("/test_blob", TestBlobController)

	// routes for default
	app.Get("/hello", HelloCheck)
	app.Get("/error", ErrorCheck)

	// routes for da
	app.Get("/node_info", NodeInfoController)
	app.Post("/submit_metadata", SubmitJSONDataController)
	app.Post("/submit_formdata", SubmitFormDataController)
	app.Post("/retrieve_blob", RetrieveBlobController2)
	app.Get("/retrieve_blob", RetrieveBlobController)

	// start server...
	port := os.Getenv("PORT")
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}

/*


	// test 2
	log.Println("=== test case 2 : SubmitBlob ===")
	err = SubmitBlob(ctx, rpcAddress, authToken)
	if err != nil {
		log.Println(err)
	}
	log.Println("================================")

	// test 3
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("=== test case 3 : Balance ===")
	err = Balance(ctx, rpcAddress, authToken)
	if err != nil {
		log.Println(err)
	}
	log.Println("================================")
*/
