package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mangenotwork/mange_redis_manage/structs"
)

var jwtSecret []byte

func init() {
	//取jwtSecret,先取缓存，缓存没有取数据库
	jwtSecret = []byte("mangeredismanage20200808")
}

type UserClaims struct {
	*structs.UserParameter
	jwt.StandardClaims
}

func UserToken(userinfo *structs.UserParameter) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := UserClaims{
		userinfo,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "man-ge",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func UserParseToken(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
