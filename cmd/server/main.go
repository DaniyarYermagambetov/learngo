package main

import (
	"fmt"
	"net/http"
)

type MemStorage struct {
	storage map[string][]byte
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Post)

	var gauge float64
	var counter int64
	fmt.Println("Listening on port 8080", gauge, counter)
	http.ListenAndServe(":8080", mux)

}

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unresolved method", http.StatusMethodNotAllowed)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

}
