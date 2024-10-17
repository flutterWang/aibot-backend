package api

import (
	"aibot-backend/io"
	"aibot-backend/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Start(configPath string) {
	srv := newServer()

	_, err := toml.DecodeFile(configPath, srv.conf)
	if err != nil {
		log.Fatal(err)
	}

	srv.io = bot_io.NewIo()
	if srv.io == nil {
		panic("new io error")
	}

	if err := srv.io.RegisterMongoClient(srv.conf.Mongo.URI, srv.conf.Mongo.Database, 10*time.Second); err != nil {
		panic(err)
	}

	err = srv.io.CreateIndex("kb", "id", true)
	if err != nil {
		fmt.Printf("kb create index error: %v\n", err)
	}

	err = srv.io.CreateIndex("category", "id", true)
	if err != nil {
		fmt.Printf("kb create index error: %v\n", err)
	}

	err = srv.io.CreateIndex("labels", "id", true)
	if err != nil {
		fmt.Printf("kb create index error: %v\n", err)
	}

	err = srv.io.CreateIndex("assistant", "id", true)
	if err != nil {
		fmt.Printf("kb create index error: %v\n", err)
	}

	err = srv.io.CreateIndex("device", "id", true)
	if err != nil {
		fmt.Printf("device create index error: %v\n", err)
	}

	srv.startHttpApiServer(fmt.Sprintf("0.0.0.0:%s", srv.conf.Port))
}

func (s *server) startHttpApiServer(addr string) {
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	kbGroup := engine.Group("/api/v1")
	{
		kbGroup.StaticFS("/static", http.Dir("./static"))

		kbGroup.GET("/kbs", s.handlerKbList)
		kbGroup.POST("/kb/create", s.handlerKbCreate)
		kbGroup.GET("/kb/detail", s.handlerKbDetail)
		kbGroup.POST("/kb/update", s.handlerKbUpdate)
		kbGroup.GET("/delete", s.handlerKbDelete)
		kbGroup.POST("/kb/upload/doc", s.handlerKbUploadDoc)

		kbGroup.GET("/labels", s.handlerKbLabelsList)
		kbGroup.POST("/labels/create", s.handlerLabelCreate)
		kbGroup.GET("/labels/delete", s.handlerLabelDelete)
		kbGroup.GET("/labels/kbs", s.handlerLabelKbs)

		kbGroup.GET("/categories", s.handlerKbCategoriesList)
		kbGroup.POST("/categories/create", s.handlerCategoryCreate)
		kbGroup.GET("/categories/delete", s.handlerCategoryDelete)
		kbGroup.GET("/categories/kbs", s.handlerCategoryKbs)

		kbGroup.POST("/search", s.handlerKbSearch)

		kbGroup.GET("/assistants", s.handlerAssistantList)
		kbGroup.POST("/assistants/create", s.handlerAssistantCreate)
		kbGroup.POST("/assistants/check/doc", s.handlerAssistantCheckDocs)
		kbGroup.GET("/assistants/async_check/doc", s.handlerAssistantAsyncCheckDocs)
		kbGroup.GET("/assistants/assistant/docs", s.handlerAssistantDocs)
		kbGroup.POST("/assistants/upload/doc", s.handlerAssistantUploadDoc)
		kbGroup.POST("/assistants/doc/ok", s.handlerAssistantUpdateDocsOK)
		kbGroup.POST("/assistants/check/doc/result", s.handlerAssistantCheckDocResult)
		kbGroup.POST("/assistants/check/doc/clear", s.handlerAssistantCheckDocClear)

		kbGroup.POST("/assistants/chat/list", s.handlerAssistantChatList)
		kbGroup.POST("/assistants/chat/create", s.handlerAssistantChatCreate)
		kbGroup.POST("/assistants/chat/update", s.handlerAssistantChatUpdate)
		kbGroup.POST("/assistants/chat/delete", s.handlerAssistantChatDelete)
		kbGroup.GET("/assistants/chat/send", s.handlerAssistantChatSend)
		kbGroup.GET("/assistants/chat/start", s.handlerAssistantChatStart)
		kbGroup.POST("/assistants/chat/conversation/save", s.handlerAssistantChatSaveConversation)

		kbGroup.GET("/devices", s.handlerDeviceList)
		kbGroup.GET("/devices/info", s.handlerGetDeviceInfo)
		kbGroup.POST("/devices/collect", s.handlerDeviceCollect)
		kbGroup.POST("/devices/create", s.handlerCreateDevices)
		kbGroup.POST("/devices/update", s.handlerUpdateDevice)

		kbGroup.GET("/chats", s.handlerChatList)
		kbGroup.POST("/chats/create", s.handlerCreateChat)
		kbGroup.GET("/chats/conversations", s.handlerChatConversations)
		kbGroup.POST("/chats/conversations/save", s.handlerSaveChatConversations)
		kbGroup.POST("/chats/delete", s.handlerChatDelete)

		kbGroup.POST("/user/login", s.handlerUserLogin)
		kbGroup.GET("/user", JWTMiddleware(s.conf.JWTSecret), s.handlerGetUserInfo)
		kbGroup.POST("/user/upload", JWTMiddleware(s.conf.JWTSecret), s.handlerUserUpload)
	}

	err := engine.Run(addr)

	go func() {
		if err := engine.Run(addr); err != nil {
			log.Println("http inner api server run error", zap.Error(err))
			panic(err)
		} else {
			log.Println("http inner api server running")
		}

	}()
	log.Printf("http api server run error: %v", err)
}

func Stop() {
	log.Println("stop")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")

		slice := strings.Split(bearerToken, " ")
		if len(slice) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := slice[1]
		claims, err := utils.ParseToken(token, jwtSecret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		key := "open_id"
		c.Set(key, claims[key])

		c.Next()
	}
}
