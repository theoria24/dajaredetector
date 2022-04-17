package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mattn/go-mastodon"
	"github.com/theoria24/dajarep"
)

func normalizeText(str string) string {
	rep := regexp.MustCompile(`<("[^"]*"|'[^']*'|[^'">])*>`)
	str = rep.ReplaceAllString(str, "")
	rep = regexp.MustCompile("[˗֊‐‑‒–⁃⁻₋−]+")
	str = rep.ReplaceAllString(str, "-")
	rep = regexp.MustCompile("[﹣－ｰ—―─━ー]+")
	str = rep.ReplaceAllString(str, "ー")
	return str
}

func main() {
	err := godotenv.Load()

	c := mastodon.NewClient(&mastodon.Config{
		Server:       os.Getenv("MSTDN_SERVER"),
		ClientID:     os.Getenv("MSTDN_CLIENT_ID"),
		ClientSecret: os.Getenv("MSTDN_CLIENT_SECRET"),
		AccessToken:  os.Getenv("MSTDN_ACCESS_TOKEN"),
	})

	wsc := c.NewWSClient()
	q, err := wsc.StreamingWSUser(context.Background())
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Start Watching")
	}
	for e := range q {
		if t, ok := e.(*mastodon.UpdateEvent); ok {
			if t.Status.Reblog == nil && len(t.Status.Mentions) == 0 {
				if t.Status.Visibility == "public" || t.Status.Visibility == "unlisted" {
					mainText := html.UnescapeString(normalizeText(t.Status.SpoilerText))
					if t.Status.SpoilerText != "" {
						mainText += " "
					}
					mainText += html.UnescapeString(normalizeText(t.Status.Content))
					fmt.Println(t.Status.Account.Acct + ": " + mainText)
					snt, key := dajarep.Dajarep(mainText, 2, true)
					if snt != nil {
						hitKey := []string{}
						for i := 0; i < len(snt); i++ {
							if len(key[i]) > 1 {
								hitKey = append(hitKey, key[i]...)
							} else {
								length := 0
								for j := 0; j < len(key[i]); j++ {
									if len([]rune(key[i][j])) > length {
										length = len([]rune(key[i][j]))
									}
								}
								if length > 2 {
									hitKey = append(hitKey, key[i]...)
								}
							}
						}
						fmt.Println(hitKey)
						if len(hitKey) > 0 {
							s, err := c.PostStatus(context.Background(), &mastodon.Toot{
								Status:      "@" + t.Status.Account.Acct + " ダジャレを検出しました（検出ワード: " + strings.Join(hitKey, ", ") + "）",
								InReplyToID: t.Status.ID,
								Visibility:  t.Status.Visibility,
							})
							if err == nil {
								fmt.Println(s)
							} else {
								fmt.Println(err)
							}
						} else {
							fmt.Println("ダジャレだけど条件を満たさない")
						}
					} else {
						fmt.Println("ダジャレじゃない")
					}
				} else {
					fmt.Println("Private Toot")
				}
			} else {
				fmt.Print("BT or Reply\n")
			}
		} else if t, ok := e.(*mastodon.NotificationEvent); ok {
			if t.Notification.Type == "follow" {
				fmt.Printf("followed by %v (ID: %v)\n", t.Notification.Account.Acct, t.Notification.Account.ID)
				_, err = c.AccountFollow(context.Background(), t.Notification.Account.ID)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("follow!")
				}
			}
		}
	}
}
