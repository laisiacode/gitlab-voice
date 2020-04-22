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

type user struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type attributes struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
	Action string `json:"action"`
}

type webhook struct {
	ObjectKind string `json:"object_kind"`

	User             user
	ObjectAttributes *attributes `json:"object_attributes"`
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
		wh := new(webhook)
		err := c.BindJSON(wh)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}

		if viper.GetInt64("chat.id") == 0 {
			c.JSON(http.StatusOK, "")
			return
		}

		switch wh.ObjectKind {
		case "merge_request":
			bot.Send(tgbotapi.NewMessage(
				viper.GetInt64("chat.id"),
				fmt.Sprintf("%s %s MR %s", wh.User.Username, wh.ObjectAttributes.State, wh.ObjectAttributes.Title),
			))
		}

		c.JSON(http.StatusOK, "")
	})

	engine.Run()
}

func main() {
	fmt.Println("bot start")
	envVar()
	go runBot()
	server()
}
