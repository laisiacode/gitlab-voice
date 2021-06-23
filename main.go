//
//
// webhook api document:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html
//
// telegram api document:
// https://core.telegram.org/bots/api
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nebulosa-studio/gitlab-voice/voice"
	"github.com/spf13/viper"
)

var bot *tgbotapi.BotAPI

func envVar() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BOT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func runBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(viper.GetString("token"))
	if err != nil {
		log.Panic("new bot ", err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		fmt.Println(
			update.Message.Chat.ID,
			update.Message.From.String(),
			update.Message.Text,
		)
	}
}

func server() {
	engine := gin.New()

	engine.POST("/webhook/raw", func(c *gin.Context) {
		raw, err := c.GetRawData()
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		fmt.Println(string(raw))
		c.JSON(http.StatusOK, "")
	})

	engine.POST("/webhook", func(c *gin.Context) {
		wh := new(voice.Webhook)
		err := c.BindJSON(wh)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		if viper.GetInt64("chat.id") == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "no chat id"})
			return
		}

		if message := wh.Notification(); message != "" {
			msg := tgbotapi.NewMessage(
				viper.GetInt64("chat.id"),
				message,
			)
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println("Send message error", err, "|msg|", msg)
			}
		}

		c.JSON(http.StatusOK, "")
	})

	if err := engine.Run(); err != nil {
		fmt.Println("engine error", err)
	}
}

func main() {
	fmt.Println("bot start")
	envVar()
	go runBot()
	server()
}
