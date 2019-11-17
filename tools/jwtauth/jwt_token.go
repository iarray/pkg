package jwtauth

import (
	"errors"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
)

var (
	TokenExpired     error = errors.New("令牌到期")
	TokenNotValidYet error = errors.New("令牌尚未生效")
	TokenMalformed   error = errors.New("无效令牌格式")
	TokenInvalid     error = errors.New("无效的令牌")
	instance         *JwtAuth
	once             sync.Once
)

type JwtAuth struct {
	SigningMethod string
	SigningKey    []byte
	lock          sync.Mutex
}

func Default() *JwtAuth {
	return instance
}

func NewAuth(SigningMethod string, SigningKey string) *JwtAuth {
	return &JwtAuth{SigningMethod: SigningMethod, SigningKey: []byte(SigningKey)}
}

func NewSingleInstance(SigningMethod string, SigningKey string) *JwtAuth {
	if instance == nil {
		once.Do(func() {
			instance = &JwtAuth{SigningMethod: SigningMethod, SigningKey: []byte(SigningKey)}
		})
	}
	return instance
}

func (me *JwtAuth) NewTokenString(clamis jwt.Claims) (string, error) {
	me.lock.Lock()
	defer me.lock.Unlock()
	j := jwt.NewWithClaims(me.getSingingMethod(), clamis)
	token, err := j.SignedString(me.SigningKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (me *JwtAuth) ParseToken(tokenString string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return me.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return TokenNotValidYet
			} else {
				return TokenInvalid
			}
		}
	}

	//claims转换为自定义格式, 并验证是否有效
	//if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
	//	return claims, nil
	//}
	return nil
}

func (me *JwtAuth) getSingingMethod() jwt.SigningMethod {
	switch strings.ToUpper(me.SigningMethod) {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS512":
		return jwt.SigningMethodHS512
	case "NONE":
		return jwt.SigningMethodNone
	case "ES256":
		return jwt.SigningMethodES256
	case "ES512":
		return jwt.SigningMethodES512
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS512":
		return jwt.SigningMethodRS512
	default:
		return jwt.SigningMethodHS256
	}
}
