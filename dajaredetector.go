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
	"github.com/kurehajime/dajarep"
	"github.com/mattn/go-mastodon"
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
			if t.Status.Reblogged == false && len(t.Status.Mentions) == 0 {
				if t.Status.Visibility == "public" || t.Status.Visibility == "unlisted" {
					fmt.Println(t.Status.Account.Acct + ": " + html.UnescapeString(removeTag(t.Status.Content)))
					snt, key := dajarep.Dajarep(html.UnescapeString(removeTag(t.Status.Content)))
					if snt != nil {
						s, err := c.PostStatus(context.Background(), &mastodon.Toot{
							Status:      "@" + t.Status.Account.Acct + " ダジャレを検出しました（検出ワード: " + strings.Join(key, ", ") + "）",
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
		}
	}
}
