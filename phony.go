package main

import (
	json2 "encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var port = flag.Int("p", 8042, "the port to listen on")

type mock struct {
	Path       string
	StatusCode int
}

type Config struct {
	Mocks map[string]mock `json:"mocks"`
}

var c Config

func main() {
	GetConfig("mocks.json")

	for _, mock := range c.Mocks {
		http.HandleFunc(mock.Path, phony)
	}

	log.Println("Server started.")
	addr := ":" + strconv.Itoa(*port)
	log.Println("Listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func send(w http.ResponseWriter, file string) {
	content, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func GetConfig(file string) {
	json, err := os.ReadFile(file)
	if err != nil {
		log.Printf("err   #%v ", err)
	}
	err = json2.Unmarshal(json, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func phony(w http.ResponseWriter, req *http.Request) {
	filename := strings.TrimLeft(req.URL.Path, "/")
	defer req.Body.Close()
	send(w, "/app/response/"+filename+".json")
}
