package main

import (
	"fmt"
	"tdlib"

	"runtime"

	"github.com/astaxie/beego/logs"

	"tgbot/config"
	"tgbot/tglog"
	"tgbot/utils"
)

var (
	client  *tdlib.Client
	poolOne WorkPool
	logger  *tglog.TGLog
)

func main() {
	logs.SetLogger(logs.AdapterConsole, `{"level":8}`)

	if runtime.GOOS != "windows" {
		logs.SetLogger(logs.AdapterFile, `{"filename":"`+utils.GetAppPath()+`/Logs/tg.log","level":3}`)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"filename":"./Logs/tg.log","level":3}`)
	}
	poolOne.InitPool(50)
	logger = &tglog.TGLog{}

	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./Logs/errors.txt")

	// Create new instance of client
	client = tdlib.NewClient(tdlib.Config{
		APIID:               "793416",
		APIHash:             "021de84fe4f1ac0361c333b0ba6198b6",
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  false,
		UseFileDatabase:     false,
		UseChatInfoDatabase: false,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	// Authorize
	for {
		currentState, _ := client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			go GetProxy()
			fmt.Print("Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				fmt.Printf("Error sending phone number: %v\n", err)
			}
			// p, err := client.CheckAuthenticationBotToken(":AAFomEPDiMQ6hE4dpmDFkKpHrmawsvwA")
			// fmt.Println(p, err)
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				fmt.Printf("Error sending auth code : %v\n", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				fmt.Printf("Error sending auth password: %v\n", err)
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			fmt.Println("登录验证通过!")
			break
		}
	}
	go RunServe()
	go GetMessage()
	go CheckProxy()

	// rawUpdates gets all updates comming from tdlib
	rawUpdates := client.GetRawUpdatesChannel(100)
	for update := range rawUpdates {
		if config.ShowRawMessage {
			logger.Info("raw %s", update.Raw)
		}
	}
}
