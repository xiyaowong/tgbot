module tgbot

go 1.15

replace tdlib => ../tdlib

require (
	github.com/astaxie/beego v1.12.2
	github.com/cjoudrey/gluahttp v0.0.0-20200626084403-ae897a63b78b
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.2.0
	github.com/junhsieh/goexamples v0.0.0-20190721045834-1c67ae74caa6 // indirect
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/tengattack/gluasql v0.0.0-20181229041402-2e5ed630c4cf
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
	google.golang.org/grpc v1.31.1 // indirect
	gopkg.in/ini.v1 v1.60.2
	gopkg.in/src-d/go-vitess.v0 v0.1.0 // indirect
	layeh.com/gopher-luar v1.0.8
	tdlib v0.0.0-00010101000000-000000000000
)
