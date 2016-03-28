package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ledyba/gakko-analyzer/nico/client"
)

var user = flag.String("user", "", "user")
var pass = flag.String("pass", "", "pass")
var video = flag.String("video", "", "video")
var when = flag.Int64("when", -1, "unixtime")

func main() {
	log.Printf("Gakkou Gurashi!")
	flag.Parse()
	cl, err := client.LoginWithPassword(*user, *pass)
	if err != nil {
		log.Fatal(err)
	}

	res, err := cl.GetResInfo(*video)
	if err != nil {
		log.Fatal(err)
	}
	wh := *when
	var chats []client.Chat
	for {
		time.Sleep(1 * time.Second)
		_, err = cl.GetComment(res, res.OptionalThreadID, wh)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
		var r []client.Chat
		r, err = cl.GetComment(res, res.ThreadID, wh)
		if err != nil {
			log.Fatal(err)
		}
		if len(r) <= 0 {
			log.Printf("Cooling...")
			time.Sleep(5 * time.Second)
			continue
		}
		chats = append(r, chats...)
		fst := &r[0]
		lst := &r[len(r)-1]
		log.Printf("%d chats: %d(%d) -> %d(%d)", len(r), fst.No, fst.Date, lst.No, lst.Date)
		if fst.No <= 10 {
			break
		}
		wh = fst.Date - 1
	}
	body, err := json.MarshalIndent(chats, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
