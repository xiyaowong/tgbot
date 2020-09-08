package main

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"tdlib"
	"tgbot/utils"
)

func GetProxy() {
	body, err, _ := utils.HttpRequest("GET", "http://129.204.103.68:58897/v1/TGProxys", nil, nil, nil)

	if err == nil && body != nil {
		var _jArray []interface{}
		json.Unmarshal(body, &_jArray)

		logger.Critical("%v", _jArray)
		for _, v := range _jArray {
			link := v.(string)
			link = strings.ReplaceAll(link, "?", "&")
			if val, err := url.ParseQuery(link); err == nil {
				port, _ := strconv.Atoi(val.Get("port"))
				_, err := client.AddProxy(val.Get("server"), int32(port), true, tdlib.NewProxyTypeMtproto(val.Get("secret")))
				logger.Trace("AddProxy %v", err)
			}
		}
	}
}

func CheckProxy() {
	heartbeat := time.NewTicker(30 * time.Second)
	pullproxy := time.NewTicker(1 * time.Hour) // 1个小时拉取一次代理列表
	proxyFlag := false

	pmap := make(map[int32]int)

	for {
		select {
		case <-heartbeat.C:
			proxys, err := client.GetProxies()
			if err != nil {
				break
			}
			for _, v := range proxys.Proxies {
				t, err := client.PingProxy(v.ID)
				coust := 0.0
				if err != nil || t.Seconds == 0 {
					count := 0
					if _, ok := pmap[v.ID]; ok {
						pmap[v.ID] += 1
						count = pmap[v.ID]
					} else {
						pmap[v.ID] = 1
					}
					logger.Error("Proxy %d Err %v Try %d", v.ID, err, count)
					if count == 15 {
						if v.IsEnabled {
							proxyFlag = true
						}
						delete(pmap, v.ID)
						client.RemoveProxy(v.ID)
					}
					continue
				} else {
					if _, ok := pmap[v.ID]; ok {
						pmap[v.ID] = 0
					}
					coust = t.Seconds
				}
				if l, err := client.GetProxyLink(v.ID); err == nil {
					if proxyFlag {
						client.EnableProxy(v.ID)
						proxyFlag = false
					}
					logger.Emergency("Check Ok Proxy %d Ping %fs Link %s IsEnabled %v", v.ID, coust, l.Text, v.IsEnabled)
				}
			}
		case <-pullproxy.C:
			GetProxy()
		}
	}
}
