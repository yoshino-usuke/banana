package gae

import (
	"net/http"
	"github.com/dongri/line-bot-sdk-go/linebot"
	"os"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

var botClient *linebot.Client

type Data struct {
	Secretchannel string
	Channelkey string
}

// BotEventHandler ...
type BotEventHandler struct{
	linebot.EventHandler
}

func init() {

	err := godotenv.Load("./line.env")
	if err != nil {
		panic(err)
	}
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESSTOKEN")
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")

	botClient = linebot.NewClient(channelAccessToken)
	botClient.SetChannelSecret(channelSecret)

	// EventHandler
	var myEvent linebot.EventHandler = NewEventHandler()
	botClient.SetEventHandler(myEvent)
	http.HandleFunc("/callback", callbackHandler)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	fmt.Fprint(w, "OK")
	fmt.Fprintln(w,botClient)
}


// NewEventHandler ...
func NewEventHandler() *BotEventHandler {
	return new(BotEventHandler)
}

func (be *BotEventHandler) OnTextMessage(source linebot.EventSource, replyToken, text string) {
	message := linebot.NewTextMessage(text + "じゃねぇよ！")
	result, err := botClient.ReplyMessage(replyToken, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
