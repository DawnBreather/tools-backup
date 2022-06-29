package main

import(
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-resty/resty"
	"log"
)


type botMessage struct{
	Channel string `json:"channel"`
	Secret string	`json:"secret"`
	Text string		`json:"text"`
	Title string	`json:"title"`
	Subtitle string	`json:"subtitle"`
	Type string		`json:"type"`
}


var _channel string
var _secret string
var _text string
var _title string
var _subtitle string
var _type string

var _url string


func main(){
	client := resty.New()
	flag.StringVar(&_channel, "channel", "", "Bot channel")
	flag.StringVar(&_secret, "secret", "", "Bot secret")
	flag.StringVar(&_text, "text", "", "Message text")
	flag.StringVar(&_title, "title", "", "Message title")
	flag.StringVar(&_subtitle, "subtitle", "", "Message subtitle")
	flag.StringVar(&_type, "type", "", "Message type")
	flag.StringVar(&_url, "botUrl", "https://touch.epm-esp.projects.epam.com/bot-esp/message", "Bot API URL endpoint")
	flag.Parse()

	msg := botMessage{
		Channel: _channel,
		Secret: _secret,
		Text: _text,
		Title: _title,
		Subtitle: _subtitle,
		Type: _type,
	}

	msgJson, _ := json.Marshal(msg)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("accept-version", "1.0.0").
		SetBody(string(msgJson)).
		Post(_url)

	if err != nil{
		log.Fatalln(fmt.Sprintf("ERROR posting request: %v", err))
	} else {
		fmt.Println(string(resp.Body()))
	}

}