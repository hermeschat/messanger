package auth

import (
	"encoding/base64"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// ValidateJWT , check given token and secret if matches
// returns aud of token if can decode, error will be returned
// if there is any problems with given token.
func ValidateJWT(token string, secret string) (jwt.MapClaims, error) {
	decodedToken := make([]byte, base64.StdEncoding.DecodedLen(len(secret)))
	_, base64err := base64.StdEncoding.Decode(decodedToken, []byte(secret))
	if base64err != nil {
		return nil, base64err
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	parser := new(jwt.Parser)
	parser.SkipClaimsValidation = true
	prasedToken, parseError := parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return decodedToken, nil
	})

	if parseError != nil {
		fmt.Println("jwt parse error", parseError)
		return nil, errors.New("parse error")
	}

	if _vErr := prasedToken.Claims.Valid(); _vErr != nil {
		if vErr, ok := _vErr.(*jwt.ValidationError); ok {
			if vErr.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				fmt.Println("Token is expired")
				return nil, errors.New("Token is expired")
			} else if vErr.Errors&jwt.ValidationErrorIssuedAt == jwt.ValidationErrorIssuedAt {
				fmt.Println("Token used before issued")
				return nil, errors.New("Token used before issued")
			} else if vErr.Errors&jwt.ValidationErrorNotValidYet == jwt.ValidationErrorNotValidYet {
				fmt.Println("token is not valid yet")
				// return nil, errors.New("token is not valid yet")
			}
		} else {
			fmt.Println("jwt validation error")
			return nil, errors.New("jwt validation error")
		}
	}

	// aud, _ := prasedToken.Claims.(jwt.MapClaims)["aud"].(string)
	claims, ok := prasedToken.Claims.(jwt.MapClaims)
	if ok && prasedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
