package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
)

type healthResp struct {
    Status    string `json:"status"`
    Timestamp string `json:"timestamp"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    resp := healthResp{Status: "ok", Timestamp: time.Now().Format(time.RFC3339)}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from SRE demo service!\n")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
    delay := os.Getenv("SLOW_MS")
    if delay == "" {
        delay = "0"
    }
    d, err := time.ParseDuration(delay + "ms")
    if err == nil && d > 0 {
        time.Sleep(d)
    }
    fmt.Fprintf(w, "OK after delay %s\n", d)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", rootHandler)
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/slow", slowHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    server := &http.Server{
        Addr:         ":" + port,
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    log.Printf("Starting server on %s\n", server.Addr)
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
