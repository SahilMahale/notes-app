package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt/v5"
)

/* func getAuthToken(c *fiber.Ctx) string {
	return strings.Split(c.Get("Authorization"), "Bearer ")[1]
} */

func (B *notesService) initMiddleware() {
	// Adding logger to the app
	B.app.Use(requestid.New())
	B.app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	B.app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	B.app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://localhost:4200,http://localhost:8080",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}

func (B *notesService) initAuth() {
	secretsFolderPath := os.Getenv("APP_AUTH")
	if secretsFolderPath == "no-auth" || secretsFolderPath == "" {
		// run app without jwt auth
		return
	}
	privateKeyPath := fmt.Sprintf("%s/private_key.pem", secretsFolderPath)
	publicKeyPath := fmt.Sprintf("%s/public_key.pem.pub", secretsFolderPath)
	err := readPrivateKeyFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	err = readPublicKeyFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	B.app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    publicKey,
		},
		ContextKey: "acces-key-token",
	}))
}

/*
	 func checkIfAdmin(role string) bool {
		return role == string(constants.Admin)
	}
*/
func readPrivateKeyFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	decodedKey, _ := pem.Decode(buffer)

	pk, errP := x509.ParsePKCS8PrivateKey(decodedKey.Bytes)
	if errP != nil {
		return errP
	}
	pKey, ok := pk.(*rsa.PrivateKey)
	if !ok {
		panic("PrivateKey type error")
	}
	privateKey = pKey
	return nil
}

func readPublicKeyFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return err
	}
	decodedKey, _ := pem.Decode(buffer)

	pk, errP := x509.ParsePKIXPublicKey(decodedKey.Bytes)
	if errP != nil {
		return errP
	}
	pubKey, ok := pk.(*rsa.PublicKey)
	if !ok {
		panic("PrivateKey type error")
	}
	publicKey = pubKey
	return nil
}

func makeTokenWithClaims(username string) (token string, err error) {
	// create claims for user
	claims := MyCustomClaims{
		Name: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	/* if checkIfAdmin {
		claims.Type = string(constants.Admin)
	} else {
		claims.Type = string(constants.User)
	} */
	// create token for the claims
	ctoken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	token, errp := ctoken.SignedString(privateKey)
	if errp != nil {
		log.Error(errp)
		return "", errp
	}
	return token, nil
}
