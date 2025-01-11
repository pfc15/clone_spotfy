package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main(){
	manager := new_manager()
	http.Handle("/", http.FileServer(http.Dir("./front/index")));
	http.HandleFunc("/list", manager.listMusic)
	http.HandleFunc("/uploadMusic", manager.uploadMusic)
	http.HandleFunc("/stream", manager.streamMusic)
	http.HandleFunc("/addPlaylist", manager.addPlaylist)

	port := ":8080"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err) // Log and stop execution if the server fails to start
	}
}	