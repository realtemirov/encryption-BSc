package main

import (
	"github.com/gin-gonic/gin"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	handler "github.com/realtemirov/encryption/handler"
	"github.com/realtemirov/encryption/repo"
	"github.com/realtemirov/encryption/service"
)

func main() {

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", nil)
	})
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Home",
		})
	})

	bot, err := tg.NewBotAPI("YOUR_TOKEN_HERE")
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	db := repo.NewDB()
	s := service.NewService(db)
	h := handler.NewHandler(db, s, bot)

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			h.Messages(update.Message)
		}
	}()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
