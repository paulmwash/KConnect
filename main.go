package main

import (
  "log"
  "net/http"
)


type apiConfig struct {
fileserverHits atomic.Int32
}

func 


func handleReadiness(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type","text/plain; charset=utf-8")
w.WriteHeader(http.StatusOK)
w.Write([]byte(http.StatusText(http.StatusOK)))
}


func main() {
const filepathroot = "."
const port = "8080"


mux := http.NewServeMux()
srv := &http.Server{
Addr: ":" + port,
Handler: mux,
}
mux.Handle("/app/",http.StripPrefix("/app/",http.FileServer((http.Dir(filepathroot)))))
mux.HandleFunc("/healthz",handleReadiness)

log.Println("Starting server on  :8080...")

err := srv.ListenAndServe()
if err  != nil {
log.Fatal(err)
}
}


