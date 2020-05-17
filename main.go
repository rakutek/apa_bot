package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/line/line-bot-sdk-go/linebot"

	"encoding/json"
	// "reflect"
	"os"
	"time"
)

type Data struct {
	Price int
	Area  []string
	Taipu string
	Room  string
	Link  string
	Name  string
	// Names []string
	Date   string
	Way    string
	Search []string
	Image  string
}

func main() {

	bot, err := linebot.New(
		"5da7dd0eed121538aed757beb6046423",
		"qyGwKFJqzdiGMZe5Nc6vcyRJtaIRRi35ke5L6g8hFe81hdKEjLx6kGXB8TupRFPdAFAy0sm4TTGFAqVYAu6NI8XmmjMYbdImGw3fX6jOo9NQT7UiulJJtQDcrFwgoVPJOJ04xb/Ixk7yR8HhDuOHCQdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}
	dic := tokenizer.SysDicSimple()
	t := tokenizer.NewWithDic(dic)

	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "apa-bot-276111")

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
		}

		for _, event := range events {

			if event.Type == linebot.EventTypeMessage {

				switch m := event.Message.(type) {

				case *linebot.TextMessage:
					var pages []Data
					var num int
					tokens := t.Tokenize(m.Text)
					for _, token := range tokens {
						features := strings.Join(token.Features(), ",")

						if strings.Contains(features, "固有名詞") {
							fmt.Println(token.Surface)

							t := time.Now().UTC()
							df := t.Format("2006-01-02")

							log.Printf("duo_" + df[8:])
							query := datastore.NewQuery("duo_"+df[8:]).
								Filter("Search =", token.Surface)

							it := client.Run(ctx, query)

							for {
								var data Data
								_, errr := it.Next(&data)

								if errr != nil {
									break
								} else {
									pages = append(pages, data)
									fmt.Println(data.Name)
								}
							}
							num = len(pages)
							fmt.Println(num)

						}

					}
					sura := []map[string]interface{}{}

					if num == 0 {

						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ヒット0件")).Do(); err != nil {
							log.Print(err)
						}
					} else {

						for i := 0; i < num; i++ {

							val := map[string]interface{}{
								"type":   "text",
								"text":   "料金",
								"size":   "lg",
								"color":  "#999999",
								"margin": "md",
								"flex":   0,
							}

							val2 := map[string]interface{}{
								"type":   "text",
								"text":   strconv.Itoa(pages[i].Price),
								"weight": "bold",
								"size":   "xl",
								"margin": "lg",
							}

							con := []map[string]interface{}{}

							con = append(con, val)

							con = append(con, val2)

							ac := map[string]interface{}{
								"type":  "uri",
								"label": "予約へ進む",
								"uri":   pages[i].Link,
							}

							m4 := map[string]interface{}{
								"type": "bubble",
								"size": "kilo",

								"hero": map[string]interface{}{
									"type":        "image",
									"url":         pages[i].Image,
									"size":        "full",
									"aspectRatio": "20:13",
									"aspectMode":  "cover",
									"action": map[string]interface{}{
										"type": "uri",
										"uri":  "http://linecorp.com/",
									},
								},

								"body": map[string]interface{}{
									"type":   "box",
									"layout": "vertical",
									"contents": []map[string]interface{}{

										map[string]interface{}{
											"type":   "text",
											"text":   pages[i].Name,
											"weight": "bold",
											"size":   "lg",
										},

										map[string]interface{}{
											"type":   "box",
											"layout": "baseline",
											"margin": "md",

											"contents": con,
										},

										map[string]interface{}{
											"type":    "box",
											"layout":  "vertical",
											"margin":  "lg",
											"spacing": "sm",
											"contents": []map[string]interface{}{
												map[string]interface{}{
													"type":    "box",
													"layout":  "baseline",
													"spacing": "sm",

													"contents": []map[string]interface{}{
														map[string]interface{}{
															"type":  "text",
															"text":  "行き方",
															"color": "#aaaaaa",
															"size":  "sm",
															"flex":  1,
														},

														map[string]interface{}{
															"type":  "text",
															"text":  pages[i].Way,
															"wrap":  true,
															"color": "#666666",
															"size":  "sm",
															"flex":  5,
														},
													},
												},

												map[string]interface{}{
													"type":    "box",
													"layout":  "baseline",
													"spacing": "sm",

													"contents": []map[string]interface{}{
														map[string]interface{}{
															"type":  "text",
															"text":  "空室",
															"color": "#aaaaaa",
															"size":  "sm",
															"flex":  1,
														},

														map[string]interface{}{
															"type":  "text",
															"text":  pages[i].Room,
															"wrap":  true,
															"color": "#666666",
															"size":  "sm",
															"flex":  5,
														},
													},
												},
											},
										},
									},
								},

								"footer": map[string]interface{}{
									"type":    "box",
									"layout":  "vertical",
									"spacing": "sm",
									"flex":    0,

									"contents": []map[string]interface{}{
										map[string]interface{}{
											"type":   "button",
											"style":  "link",
											"height": "sm",

											"action": ac,
										},
									},
								},
							}

							sura = append(sura, m4)
						}

						m1 := map[string]interface{}{
							"type":     "carousel",
							"contents": sura,
						}
						i, _ := json.MarshalIndent(m1, "", "   ")

						container, err := linebot.UnmarshalFlexMessageJSON(i)

						if err != nil {
							fmt.Println(err)
						}
						for _, event := range events {
							if event.Type == linebot.EventTypeMessage {
								switch event.Message.(type) {

								case *linebot.TextMessage:

									if _, err := bot.ReplyMessage(
										event.ReplyToken,
										linebot.NewFlexMessage("alt text", container),
									).Do(); err != nil {
									}
								}
							}
						}
					}
				}
			}
		}
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
