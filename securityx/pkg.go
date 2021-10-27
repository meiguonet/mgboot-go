package securityx

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-errors/errors"
	"github.com/meiguonet/mgboot-go"
	"io/ioutil"
)

func ParseJsonWebToken(token string) (*jwt.Token, error) {
	keyBytes := loadJwtPublicKeyPem(mgboot.JwtPublicKeyPemFile())

	if len(keyBytes) < 1 {
		return nil, errors.New("fail to load public key from pem file")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)

	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
		}

		return publicKey, nil
	})
}

func loadJwtPublicKeyPem(arg0 interface{}) []byte {
	var fpath string

	if s1, ok := arg0.(string); ok && s1 != "" {
		fpath = s1
	} else if s1, ok := arg0.(*JwtSettings); ok && s1 != nil {
		fpath = s1.PublicKeyPemFile()
	}

	if fpath == "" {
		return make([]byte, 0)
	}

	buf, err := ioutil.ReadFile(fpath)

	if err != nil {
		return make([]byte, 0)
	}

	return buf
}
