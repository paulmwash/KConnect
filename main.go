package main

import (
  "log"
  "net/http"
  "fmt"
  "sync/atomic"
)


type apiConfig struct {
fileserverHits atomic.Int32
}

func(cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
return http.HandlerFunc(func ( w http.ResponseWriter, r *http.Request) {
cfg.fileserverHits.Add(1)
next.ServeHTTP(w,r)
})
}

func(cfg *apiConfig) handleMetrics( w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "text/plain; charset=utf-8")
w.WriteHeader(http.StatusOK)
w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func(cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
cfg.fileserverHits.Store(0)
w.Header().Set("Content-Type", "text/plain; charset=utf-8")
w.WriteHeader(http.StatusOK)
w.Write([]byte("Hits reset to 0"))
}


func handleReadiness(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type","text/plain; charset=utf-8")
w.WriteHeader(http.StatusOK)
w.Write([]byte(http.StatusText(http.StatusOK)))
}


func main() {

const filepathroot = "."
const port = "8080"

apiCfg := &apiConfig{}


mux := http.NewServeMux()
srv := &http.Server{
Addr: ":" + port,
Handler: mux,
}
mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(filepathroot)))))

mux.HandleFunc("/healthz",handleReadiness)
mux.HandleFunc("/metrics", apiCfg.handleMetrics)
mux.HandleFunc("/reset", apiCfg.handlerReset)

log.Println("Starting server on  :8080...")

err := srv.ListenAndServe()
if err  != nil {
log.Fatal(err)
}
}


