package main

import (
	"encoding/json"
	"fmt"
	"github.com/Bejdenn/url-shortener/functions/url-shortening/url"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func main() {
	addresses := []string{
		"http://www.facebook.com",
		"http://www.baidu.com",
		"http://www.yahoo.com",
		"http://www.amazon.com",
		"http://www.wikipedia.org",
		"http://www.taobao.com",
		"http://www.instagram.com",
		"http://www.sina.com.cn",
		"http://www.yahoo.co.jp",
		"http://www.yandex.ru",
		"http://www.hao123.com",
		"http://www.ebay.com",
		"http://www.tmall.com",
		"http://www.amazon.co.jp",
		"http://www.paypal.com",
		"http://www.stackoverflow.com",
		"http://www.aliexpress.com",
		"http://www.naver.com",
		"http://www.apple.com",
		"http://www.chinadaily.com.cn",
		"http://www.google.ca",
		"http://www.whatsapp.com",
		"http://www.amazon.in",
		"http://www.tianya.cn",
		"http://www.rakuten.co.jp",
		"http://www.craigslist.org",
		"http://www.amazon.de",
		"http://www.xinhuanet.com",
		"http://www.outbrain.com",
		"http://www.alibaba.com",
		"http://www.alipay.com",
		"http://www.google.com.au",
		"http://www.popads.net",
		"http://www.amazon.co.uk",
		"http://www.wikia.com",
		"http://www.googleadservices.com",
		"http://www.accuweather.com",
		"http://www.answers.yahoo.com",
	}

	relations := make([]*url.Relation, 0)

	for _, addr := range addresses {
		for i := 0; i < rand.Intn(25)+1; i++ {
			api := "https://api-72ey6bex.nw.gateway.dev/url-shortening"
			method := "POST"

			payload := strings.NewReader("longUrl=" + addr)

			client := &http.Client{}
			req, err := http.NewRequest(method, api, payload)

			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}

			log.Printf("Request short URL for: %s\n", addr)

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			rel := &url.Relation{}
			err = json.Unmarshal(body, rel)
			if err != nil {
				log.Printf("could not unmarshall response for addr %s: %v", addr, err)
				return
			}

			relations = append(relations, rel)
		}
	}
}
