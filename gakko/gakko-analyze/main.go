package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ledyba/gakko-analyzer/nico/client"
)

var logfile = flag.String("file", "", "log.json")
var words = flag.String("words", "", "target words")

func main() {
	log.Printf("Gakkou Gurashi!")
	flag.Parse()

	var logs []*client.Chat
	var grps []*Group

	{
		f, err := os.Open(*logfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		dec.Decode(&logs)
		grps = makeGraph(logs, strings.Split(*words, ","))
	}

	render := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		drawGraph(w, logs, grps)
	}

	http.Handle("/", http.HandlerFunc(render))
	log.Printf("Listen: http://localhost:2003/")
	err := http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
