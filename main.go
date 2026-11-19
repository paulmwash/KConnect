package main

import (
  "log"
  "net/http"
  "fmt"
  "sync/atomic"
  "encoding/json"
)

type students struct {
  Name string `json:"name"`
  ID string    `json:"ID"`
}

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
w.Header().Set("Content-Type", "text/html")
w.WriteHeader(http.StatusOK)
html := fmt.Sprintf(`
<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times !</p>
</body>
</html>`,cfg.fileserverHits.Load())
w.Write([]byte(html))

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

func(s *students) getThem(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(s)
}


func main() {

const filepathroot = "."
const port = "8080"

apiCfg := &apiConfig{}
  stud := &students {
    Name : "Paul",
    ID : "KLH545",
  }


mux := http.NewServeMux()
srv := &http.Server{
Addr: ":" + port,
Handler: mux,
}

  //Regest site visits
mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(filepathroot)))))
mux.Handle("/student",apiCfg.middlewareMetricsInc(http.HandlerFunc(stud.getThem)))
//mux.Handle("healthz",apiCfg.middlewareMetricsInc(http.HandlerFunc(apiCfg.handleMetrics))) // Admin should not be recorded
  
mux.HandleFunc("GET /api/healthz",handleReadiness)
mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

log.Println("Starting server on  :8080...")

err := srv.ListenAndServe()
if err  != nil {
log.Fatal(err)
}
}


