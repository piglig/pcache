package main

import (
	"fmt"
	"log"
	"net/http"
	"pcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	pcache.NewGroup("scores", 2<<10, pcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[Slow DB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not found", key)
		}))

	addr := "localhost:9999"
	peers := pcache.NewHTTPPool(addr)
	log.Println("pcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))

}
