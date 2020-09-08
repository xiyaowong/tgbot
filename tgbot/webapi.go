// Package main provides ...
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"

	"tdlib"
	"tgbot/config"
	"tgbot/utils"
)

func BuildResponse(c *gin.Context, err error, data interface{}) {
	if err == nil {
		c.JSON(http.StatusOK, &struct {
			Code int         `json:"code"`
			Err  string      `json:"err"`
			Data interface{} `json:"data"`
		}{
			Code: 0,
			Data: data,
		})
	} else {
		c.JSON(http.StatusOK, &struct {
			Code int         `json:"code"`
			Err  string      `json:"err"`
			Data interface{} `json:"data"`
		}{
			Code: 1,
			Err:  err.Error(),
		})
	}
}

// WebSendText 发送文字
func WebSendText(c *gin.Context) {
	logger.Notice("Web SendText")
	var text struct {
		ChatID int    `json:"chat_id" binding:"required"`
		Text   string `json:"text" binding:"required"`
	}
	var err error
	if err = c.ShouldBindJSON(&text); err == nil {
		inputMsgTxt := tdlib.NewInputMessageText(tdlib.NewFormattedText(text.Text, nil), true, false)
		_, err = client.SendMessage(int64(text.ChatID), 0, nil, nil, inputMsgTxt)
	}
	BuildResponse(c, err, nil)
}

// WebSendPhoto 发送图片
func WebSendPhoto(c *gin.Context) {
	logger.Notice("Web SendPhoto")
	// decode json
	var photo struct {
		ChatID int    `json:"chat_id" binding:"required"`
		Text   string `json:"text"`
		Base64 string `json:"base64"`
		URL    string `json:"url"`
	}
	var err error
	if err = c.ShouldBindJSON(&photo); err == nil {
		// get image data
		var imgData []byte
		if photo.Base64 != "" {
			// decode base64
			imgData, err = base64.StdEncoding.DecodeString(photo.Base64)
		} else if photo.URL != "" {
			// request images data
			imgData, err, _ = utils.HttpRequest("GET", photo.URL, nil, nil, nil)
		} else {
			imgData, err = nil, fmt.Errorf("Base64 和 URL 必选一项")
		}
		if err == nil {
			// convert []byte to local path
			tempDir, _ := ioutil.TempDir("", "tgbot")
			filePath := path.Join(tempDir, "photo")
			defer os.RemoveAll(tempDir)
			if err = ioutil.WriteFile(filePath, imgData, os.ModePerm); err == nil {
				inputMsg := tdlib.NewInputMessagePhoto(tdlib.NewInputFileLocal(filePath), nil, nil, 400, 400,
					tdlib.NewFormattedText(photo.Text, nil), 0)
				_, err = client.SendMessage(int64(photo.ChatID), 0, nil, nil, inputMsg)
			}
		}
	}
	BuildResponse(c, err, nil)
}

// WebDeleteMessages 删除消息
func WebDeleteMessages(c *gin.Context) {
	logger.Notice("Web DeleteMessages")
	var del struct {
		ChatID    int  `json:"chat_id" binding:"required"`
		MessageID int  `json:"message_id" binding:"required"`
		Revoke    bool `json:"revoke"`
	}
	var err error
	if err = c.ShouldBindJSON(&del); err == nil {
		_, err = client.DeleteMessages(int64(del.ChatID), []int64{int64(del.MessageID)}, del.Revoke)
	}
	BuildResponse(c, err, nil)
}

func WebGetUser(c *gin.Context) {
	logger.Notice("Web GetUser")
	var user struct {
		UserID int `json:"user_id" binding:"required"`
	}
	var err error
	var data = make(map[string]interface{})
	if err = c.ShouldBindJSON(&user); err == nil {
		if info, err := client.GetUser(int32(user.UserID)); err == nil {
			data["id"] = info.ID
			data["first_name"] = info.FirstName
			data["last_name"] = info.LastName
			data["username"] = info.Username
			data["phone_number"] = info.PhoneNumber
			data["is_contact"] = info.IsContact
			data["status"] = info.Status
			data["is_scam"] = info.IsScam
			data["is_verified"] = info.IsVerified
		}
	}
	BuildResponse(c, err, data)
}

func WebGetChat(c *gin.Context) {
	logger.Notice("Web GetChat")
	var chat struct {
		ChatID int `json:"chat_id" binding:"required"`
	}
	var err error
	var data = make(map[string]interface{})
	if err = c.ShouldBindJSON(&chat); err == nil {
		if info, err := client.GetChat(int64(chat.ChatID)); err == nil {
			data["id"] = info.ID
			data["chat_title"] = info.Title
		}
	}
	BuildResponse(c, err, data)
}

func RunServe() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.Recovery())

	router.POST("/sendText", WebSendText)
	router.POST("/sendPhoto", WebSendPhoto)
	router.POST("/deleteMessages", WebDeleteMessages)
	router.POST("/getUser", WebGetUser)
	router.POST("/getChat", WebGetChat)

	fmt.Printf("Listening on port %d\n", config.ServerPort)
	router.Run(fmt.Sprintf(":%d", config.ServerPort))
}
