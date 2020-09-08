// Package main provides ...
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

	"tdlib"
	"tgbot/config"
	"tgbot/utils"
)

var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()

		uni := ut.New(enT, zhT, enT)

		var ok bool
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

type Response struct {
	Code int         `json:"code"`
	Err  string      `json:"err"`
	Data interface{} `json:"data"`
}

func BuildResponse(c *gin.Context, err error, data interface{}) {
	if err == nil {
		c.JSON(http.StatusOK, &Response{
			Code: 0,
			Data: data,
		})
	} else {
		if errs, ok := err.(validator.ValidationErrors); ok {
			var errStrs []string
			for _, errStr := range errs.Translate(trans) {
				errStrs = append(errStrs, errStr)
			}
			c.JSON(http.StatusOK, &Response{
				Code: 1,
				Err:  strings.Join(errStrs, ","),
			})
		} else {
			c.JSON(http.StatusOK, &Response{
				Code: 1,
				Err:  err.Error(),
			})
		}
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
	if err := InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
	}

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
