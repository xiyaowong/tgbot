package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"

	"tdlib"
	"tgbot/config"
	"tgbot/utils"
)

func HookPost(instances []interface{}) {
	if len(instances) == 0 {
		return
	}
	if config.HookPostURL != "" {
		if data, err := json.Marshal(instances[0]); err == nil {
			utils.HttpRequest("POST", config.HookPostURL, data, nil, map[string]string{"content-type": "application/json"})
		}
	}
}

func GetMessage() {
	eventFilter := func(msg *tdlib.TdMessage) bool {
		return true
	}
	receiver := client.AddEventReceiver(&tdlib.UpdateNewMessage{}, eventFilter, 5)
	for newMsg := range receiver.Chan {
		updateMsg := (newMsg).(*tdlib.UpdateNewMessage)
		// Text Message
		if msgText, ok := updateMsg.Message.Content.(*tdlib.MessageText); ok {
			m := make(map[string]interface{})

			m["ChatID"] = updateMsg.Message.ChatID
			chat, err := client.GetChat(updateMsg.Message.ChatID)
			if err == nil {
				m["ChatTitle"] = chat.Title
			} else {
				m["ChatTitle"] = ""
			}
			m["SenderUserID"] = updateMsg.Message.SenderUserID
			user, err := client.GetUser(updateMsg.Message.SenderUserID)
			if err == nil {
				m["SenderUserFirstName"] = user.FirstName
				m["SenderUserLastName"] = user.LastName
				m["SenderUserName"] = user.Username
				m["SenderUserPhoneNumber"] = user.PhoneNumber
			} else {
				m["SenderUserFirstName"] = ""
				m["SenderUserLastName"] = ""
				m["SenderUserName"] = ""
				m["SenderUserPhoneNumber"] = ""

			}
			m["MsgType"] = "MessageText"
			m["MessageID"] = updateMsg.Message.ID
			m["Content"] = msgText.Text.Text
			//协程池执行lua插件
			poolOne.Run(TGLuaVMRun, m)
			// 上报数据
			poolOne.Run(HookPost, m)
			continue
		}
		// Photo Message
		if msgPhoto, ok := updateMsg.Message.Content.(*tdlib.MessagePhoto); ok {
			m := make(map[string]interface{})

			m["ChatID"] = updateMsg.Message.ChatID
			chat, err := client.GetChat(updateMsg.Message.ChatID)
			if err == nil {
				m["ChatTitle"] = chat.Title
			} else {
				m["ChatTitle"] = ""
			}
			m["SenderUserID"] = updateMsg.Message.SenderUserID
			user, err := client.GetUser(updateMsg.Message.SenderUserID)
			if err == nil {
				m["SenderUserFirstName"] = user.FirstName
				m["SenderUserLastName"] = user.LastName
				m["SenderUserName"] = user.Username
				m["SenderUserPhoneNumber"] = user.PhoneNumber
			} else {
				m["SenderUserFirstName"] = ""
				m["SenderUserLastName"] = ""
				m["SenderUserName"] = ""
				m["SenderUserPhoneNumber"] = ""

			}
			m["MessageID"] = updateMsg.Message.ID
			m["MsgType"] = "MessagePhoto"

			// 如果文件下载成功，取base64设置为Content字段
			photo := msgPhoto.Photo.Sizes[len(msgPhoto.Photo.Sizes)-1]
			photoFile, err := client.DownloadFile(photo.Photo.ID, 1, photo.Photo.Local.DownloadOffset, 0, true)
			if err != nil {
				logger.Trace("图片文件下载失败")
				continue
			}
			if !photoFile.Local.IsDownloadingCompleted {
				continue
			}
			photoPath := photoFile.Local.Path
			if photoPath == "" {
				continue
			}
			photoBytes, err := ioutil.ReadFile(photoPath)
			os.RemoveAll(photoFile.Local.Path)
			if err != nil {
				logger.Trace("读取图片文件错误")
				continue
			}
			photoBase64 := base64.StdEncoding.EncodeToString(photoBytes)
			m["PhotoBase64"] = photoBase64
			m["Content"] = msgPhoto.Caption.Text

			poolOne.Run(TGLuaVMRun, m)
			poolOne.Run(HookPost, m)
			continue
		}
	}
}
