package line

import (
	"cloud.google.com/go/datastore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"strconv"
	"strings"
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



func Replay(bot *linebot.Client,events []*linebot.Event) {

	dic := tokenizer.SysDicSimple()
	t := tokenizer.NewWithDic(dic)

	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "apa-bot-276111")


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
				contents := []map[string]interface{}{}

				if num == 0 {

					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ヒット0件")).Do(); err != nil {
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

						contents = append(contents, m4)
					}

					m1 := map[string]interface{}{
						"type":     "carousel",
						"contents": contents,
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


}