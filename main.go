//
//
// webhook api document:
// https://docs.gitlab.com/ce/user/project/integrations/webhooks.html
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

type project struct {
	Path string `json:"path_with_namespace"`
	URL  string `json:"url"`
}

type attributes struct {
	ID     int    `json:"id"`
	IID    int    `json:"iid"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
	Action string `json:"action"`
}

type webhook struct {
	ObjectKind string `json:"object_kind"`

	User             user
	Project          project
	ObjectAttributes *attributes `json:"object_attributes"`
}

func (wh *webhook) Notification() string {
	switch wh.ObjectKind {
	case "merge_request":
		return wh.mrNotification()
	case "issue":
		return wh.issueNotification()
	case "note":
		return wh.commentNotification()
	default:
		fmt.Println("webhook", wh.ObjectKind)
	}
	return ""
}

func (wh *webhook) mrNotification() string {
	switch wh.ObjectAttributes.Action {
	case "open", "merge", "close":
		return fmt.Sprintf("%s %s MR [\\!%d](%s) %s at %s",
			wh.User.Username,
			wh.ObjectAttributes.Action,
			wh.ObjectAttributes.IID,
			wh.ObjectAttributes.URL,
			wh.ObjectAttributes.Title,
			wh.Project.Path,
		)
	default:
		return ""
	}
}

func (wh *webhook) issueNotification() string {
	switch wh.ObjectAttributes.Action {
	case "open", "merge", "close":
		return fmt.Sprintf("%s %s issue [\\#%d](%s) %s at %s",
			wh.User.Username,
			wh.ObjectAttributes.Action,
			wh.ObjectAttributes.IID,
			wh.ObjectAttributes.URL,
			wh.ObjectAttributes.Title,
			wh.Project.Path,
		)
	default:
		return ""
	}
}

func (wh *webhook) commentNotification() string {
	return ""
}

func markdownV2Escape(s string) string {
	return strings.ReplaceAll(s, "-", "\\-")
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
			c.JSON(http.StatusOK, gin.H{"message": "no chat id"})
			return
		}

		if message := wh.Notification(); message != "" {
			msg := tgbotapi.NewMessage(
				viper.GetInt64("chat.id"),
				markdownV2Escape(message),
			)
			msg.ParseMode = "MarkdownV2"
			_, err := bot.Send(msg)
			if err != nil {
				fmt.Println("send message error", err)
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
