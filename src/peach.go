package main

import (
	"./helper"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"bytes"
	"time"
	"io/ioutil"
	"net/url"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/gin-gonic/gin"
)

// ReceivedMessage ...
type ReceivedMessage struct {
	Result []Result `json:"result"`
}

// Result ...
type Result struct {
	ID          string   `json:"id"`
	From        string   `json:"from"`
	FromChannel int      `json:"fromChannel"`
	To          []string `json:"to"`
	ToChannel   int      `json:"toChannel"`
	EventType   string   `json:"eventType"`
	Content     Content  `json:"content"`
}

// Content ...
type Content struct {
	ID          string   `json:"id"`
	ContentType int      `json:"contentType"`
	From        string   `json:"from"`
	CreatedTime int      `json:"createdTime"`
	To          []string `json:"to"`
	ToType      int      `json:"toType"`
	Text        string   `json:"text"`
}

// SendMessage ..
type SendMessage struct {
	To        []string `json:"to"`
	ToChannel int      `json:"toChannel"`
	EventType string   `json:"eventType"`
	Content   Content  `json:"content"`
}

// func main() {
// 	http.HandleFunc("/", helloHandler)
// 	http.HandleFunc("/callback", callbackHandler)
// 	port := os.Getenv("PORT")
// 	addr := fmt.Sprintf(":%s", port)
// 	http.ListenAndServe(addr, nil)
// }

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// ② LINE bot instanceの作成
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	// ③ LINE Messaging API用の Routing設定
	router.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Print(err)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	router.Run(":" + port)
}
func helloHandler(w http.ResponseWriter, r *http.Request)  {

}


// bot機能
func callbackHandler(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var m ReceivedMessage
	err := decoder.Decode(&m)
	if err != nil {
		log.Println(err)
	}
	apiURI := helper.EndPoint + "/v1/events"

	PROJECT := "peach"
	conf,err :=  helper.GetBotData(PROJECT)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err != nil {
		log.Fatalf(err.Error())
	}

	sheetService, d := helper.GetSheet("banana")
	//TODO response受け取る
	updateType := helper.UPDATE_DAILY
	rank := 20
	dbList := d.GetDb(updateType)

	for _,db := range dbList{
		sheet,err := sheetService.Spreadsheets.Values.Get(d.Dbid,db.ColumnRange(rank)).Do()
		if err != nil {
			log.Fatalf(err.Error())
		}
		if sheet != nil{}
		for _, result := range m.Result {
			from := result.Content.From
			text := result.Content.Text
			content := new(Content)
			content.ContentType = result.Content.ContentType
			content.ToType = result.Content.ToType
			content.Text = text
			request(apiURI, "POST", []string{from}, *content,conf.ChannelKey)
		}
	}

}

func request(endpointURL string, method string, to []string, content Content,channel int) {
	m := &SendMessage{}
	m.To = to
	m.ToChannel = channel
	m.EventType = helper.EventType
	m.Content = content
	b, err := json.Marshal(m)
	if err != nil {
		log.Print(err)
	}
	req, err := http.NewRequest(method, endpointURL, bytes.NewBuffer(b))
	if err != nil {
		log.Print(err)
	}
	req = setHeader(req)
	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(getProxyURL())},
		Timeout:   time.Duration(30 * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	var result map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Print(err)
	}
	log.Print(result)
}

func setHeader(req *http.Request) *http.Request {
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("X-Line-ChannelID", os.Getenv("ChannelID"))
	req.Header.Add("X-Line-ChannelSecret", os.Getenv("ChannelSecret"))
	req.Header.Add("X-Line-Trusted-User-With-ACL", os.Getenv("MID"))
	return req
}

func getProxyURL() *url.URL {
	proxyURL, err := url.Parse(os.Getenv("ProxyURL"))
	if err != nil {
		log.Print(err)
	}
	return proxyURL
}