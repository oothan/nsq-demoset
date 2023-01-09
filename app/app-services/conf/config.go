package conf

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/joho/godotenv"
	logger "nsq-demoset/app/_applib"
	"os"
)

var (
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
	AppHost       string

	NsqAddr string
)

func init() {
	//getting working directory
	dir, err := os.Getwd()
	if err != nil {
		logger.Sugar.Error("Error on getting directory : ", err.Error())
	}

	/*// logs create directory
	if err = os.MkdirAll(fmt.Sprintf("%s/logs", dir), 0755); err != nil {
		logger.Sugar.Error("Error on creating directory logs: ", err)
	}*/

	err = godotenv.Load(dir + "/conf/.env")
	if err != nil {
		logger.Sugar.Error("error on getting directory : ", err)
	}

	// Load rsa [private]
	privateBytes, err := os.ReadFile(os.Getenv("RSA_PRIVATE"))
	if err != nil {
		logger.Sugar.Error("Error on loading private key: ", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		logger.Sugar.Error("Error on parsing private key: ", err)
	}
	PrivateKey = privateKey

	// Load rsa [public]
	publicBytes, err := os.ReadFile(os.Getenv("RSA_PUBLIC"))
	if err != nil {
		logger.Sugar.Error("Error on loading public key: ", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		logger.Sugar.Error("Error on parsing key: ", err)
	}
	PublicKey = publicKey

	// Load rsa [secret]
	RefreshSecret = os.Getenv("RSA_SECRET")

	// App host
	AppHost = os.Getenv("APP_DOMAIN")

	NsqAddr = os.Getenv("NSQ_ADDR")

}
