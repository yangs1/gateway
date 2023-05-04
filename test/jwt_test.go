package main

import (
	"fmt"
	"gateway/config"
	"gateway/initialize"
	"gateway/utils"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	initialize.Viper(&config.GVA_CONFIG, "../config/app.debug.yaml")

	fmt.Println(config.GVA_CONFIG)

	config.GVA_LOG = initialize.Zap()
	config.GVA_DB = initialize.GormMysql()
	rand.Seed(time.Now().UnixNano())
	id := strconv.Itoa(rand.Intn(20))
	expiredAt := time.Now().Add(time.Hour)
	claim := jwt.StandardClaims{
		Id:        id,
		ExpiresAt: expiredAt.Unix(),
	}

	token, _ := utils.JwtEncode(claim)
	fmt.Println(id, token)
}
