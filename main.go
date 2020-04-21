package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

func envVar() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BOT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func run() {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("token"))
	if err != nil {
		log.Panic("new bot ", err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func server() {
	engine := gin.New()

	engine.POST("/webhook", func(c *gin.Context) {
		raw, err := c.GetRawData()
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		fmt.Println(string(raw))

		c.JSON(http.StatusOK, "")
	})

	engine.Run()
}

func main() {
	fmt.Println("bot start")
	envVar()
	run()

	server()
}
