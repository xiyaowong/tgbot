module tgbot

go 1.15

replace tdlib => ../tdlib

require (
	github.com/astaxie/beego v1.12.2
	github.com/cjoudrey/gluahttp v0.0.0-20200626084403-ae897a63b78b
	github.com/gin-gonic/gin v1.6.3
	github.com/junhsieh/goexamples v0.0.0-20190721045834-1c67ae74caa6 // indirect
	github.com/tengattack/gluasql v0.0.0-20181229041402-2e5ed630c4cf
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
	gopkg.in/ini.v1 v1.60.2
	layeh.com/gopher-luar v1.0.8
	tdlib v0.0.0-00010101000000-000000000000
)
