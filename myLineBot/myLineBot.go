package myLineBot

import (
	"fmt"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	mysecretmanager "github.com/ozaki-physics/gcp-training/mySecretManager"
)

const projectId = "smart-ruler-277318"

// lineChannel callback が呼び出されるたびに Secret Manager から呼び出すのが嫌だったから
// 先に struct を作るようにしてみる(Secret Manager への呼び出しが減るか未検証)
type lineChannel struct {
	secret string
	token  string
}

func Main() {
	// 1個目の LINE bot のクレデンシャルを取得
	channelSecret, err := mysecretmanager.GetGCPSecretValue(projectId, "LINE_CHANNEL_SECRET", 1)
	if err != nil {
		log.Fatal(err)
	}
	channelToken, err := mysecretmanager.GetGCPSecretValue(projectId, "LINE_CHANNEL_TOKEN", 1)
	if err != nil {
		log.Fatal(err)
	}
	lineCredential := lineChannel{channelSecret, channelToken}

	// 2個目の LINE bot のクレデンシャルを取得
	channelSecret02, err := mysecretmanager.GetGCPSecretValue(projectId, "LINE_CHANNEL_SECRET_02", 1)
	if err != nil {
		log.Fatal(err)
	}
	channelToken02, err := mysecretmanager.GetGCPSecretValue(projectId, "LINE_CHANNEL_TOKEN_02", 1)
	if err != nil {
		log.Fatal(err)
	}
	lineCredential02 := lineChannel{channelSecret02, channelToken02}

	// URL ハンドリング
	http.HandleFunc("/callback-go", lineCredential.linebotHandler)
	http.HandleFunc("/callback-python", lineCredential02.linebotHandler)
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(":"+"8080", nil); err != nil {
		log.Fatal(err)
	}
}

// indexHandler ブラウザの表示
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	name := "test"
	value, err := mysecretmanager.GetGCPSecretValue(projectId, name, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, value+" Hello, World!")
}

// linebotHandler LINE bot でオウム返しをする
func (l *lineChannel) linebotHandler(w http.ResponseWriter, req *http.Request) {
	bot, err := linebot.New(l.secret, l.token)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(req)
	log.Printf("line request: %v\n", events)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Printf("line request message: %s\n", message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ozaki-physics: "+message.Text)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf("sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
