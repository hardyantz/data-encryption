package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hardyantz/data-encryption/config"
	"github.com/hardyantz/data-encryption/conn"
	memberDelivery "github.com/hardyantz/data-encryption/pkg/member/delivery"
	memberRepository "github.com/hardyantz/data-encryption/pkg/member/repository"
	memberUseCase "github.com/hardyantz/data-encryption/pkg/member/usecase"
	userDelivery "github.com/hardyantz/data-encryption/pkg/user/delivery"
	userRepository "github.com/hardyantz/data-encryption/pkg/user/repository"
	userUseCase "github.com/hardyantz/data-encryption/pkg/user/usecase"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	if err := Load(".env"); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	db := ConnectDB()
	conf := config.NewConfImpl()

	cacheExpire := 1 * time.Hour

	redis := config.NewCacheRedis(conn.RedisConnect(), cacheExpire)

	memberRepo := memberRepository.NewMemberRepoImplementation(db, conf, redis)
	memberUc := memberUseCase.NewMemberImplementation(memberRepo)
	memberHandler := memberDelivery.NewHTTPHandler(memberUc)

	userRepo := userRepository.NewUserImplementation(db, conf)
	userUc := userUseCase.NewUserImplementation(userRepo)
	userHandler := userDelivery.NewHTTPHandler(userUc)

	member := e.Group("/member")
	memberHandler.Mount(member)
	user := e.Group("/user")
	userHandler.Mount(user)

	e.Start(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
