package main

import (
	"apa_bot/line"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)


func main() {

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {


		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}else{
			w.WriteHeader(200)
		}

		line.Replay(bot,events)


	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

// gcloud beta emulators datastore start --host-port localhost:8059 --project apa-bot-274108  --data-dir='/Users/raku/Library/Mobile Documents/com~apple~CloudDocs/dev/apabot/data/'
// env DATASTORE_EMULATOR_HOST=localhost:8059 DATASTORE_PROJECT_ID=apa-bot-274108 go run line.go
// google-cloud-gui

// dev_appserver.py app.yaml --enable_host_checking=false
