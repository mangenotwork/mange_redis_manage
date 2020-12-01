module github.com/mangenotwork/mange_redis_manage

go 1.13

replace github.com/mangenotwork/mange_redis_manage => ./

require (
	github.com/astaxie/beego v1.12.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/garyburd/redigo v1.6.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ini/ini v1.57.0
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/patrickmn/go-cache v2.1.0+incompatible
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
)
