package auth

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func strMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	strSig := m.Sum(nil)
	return hex.EncodeToString(strSig)
}

func GetUserPassword(password, salt string) string {
	saltMd5 := strMd5(salt)
	var bt bytes.Buffer
	bt.WriteString(saltMd5)
	bt.WriteString(password)
	return strMd5(bt.String())
}

// Encrypt encrypts the plain text with bcrypt.
func Encrypt(source string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare compares the encrypted text with the plain text if it's the same.
func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Sign issue a jwt token based on secretID, secretKey, iss and aud.
func Sign(secretID, secretKey, iss, aud string) string {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Add(0).Unix(),
		"aud": aud,
		"iss": iss,
	}

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = secretID

	// Sign the token with the specified secret.
	tokenString, _ := token.SignedString([]byte(secretKey))

	return tokenString
}
