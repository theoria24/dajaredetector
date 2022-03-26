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

func removeTag(str string) string {
	rep := regexp.MustCompile(`<("[^"]*"|'[^']*'|[^'">])*>`)
	str = rep.ReplaceAllString(str, "")
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
	}
	for e := range q {
		if t, ok := e.(*mastodon.UpdateEvent); ok {
			if t.Status.Reblog == nil && len(t.Status.Mentions) == 0 {
				if t.Status.Visibility == "public" || t.Status.Visibility == "unlisted" {
					mainText := html.UnescapeString(removeTag(t.Status.SpoilerText))
					if t.Status.SpoilerText != "" {
						mainText += " "
					}
					mainText += html.UnescapeString(removeTag(t.Status.Content))
					fmt.Println(t.Status.Account.Acct + ": " + mainText)
					snt, key := dajarep.Dajarep(mainText, 3, true)
					if snt != nil {
						cont := "@" + t.Status.Account.Acct + " ダジャレを検出しました（検出ワード: "
						for i := 0; i < len(key); i++ {
							if i != 0 {
								cont += ", "
							}
							cont += strings.Join(key[i], ", ")
						}
						cont += "）"
						s, err := c.PostStatus(context.Background(), &mastodon.Toot{
							Status:      cont,
							InReplyToID: t.Status.ID,
							Visibility:  t.Status.Visibility,
						})
						if err == nil {
							fmt.Println(s)
						} else {
							fmt.Println(err)
						}
					} else {
						fmt.Println("ダジャレじゃない")
					}
				} else {
					fmt.Print("Private Toot\n")
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
