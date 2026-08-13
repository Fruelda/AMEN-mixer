package main

import (
	"log"
	"net/http"
	"time"
)

func StartRealtimeServer(
	server *RealtimeServer,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/ws",
		server.HandleWebSocket,
	)

	mux.HandleFunc(
		"/",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			_, _ = w.Write(
				[]byte(
					"AMEN MIXER REALTIME SERVER OK",
				),
			)
		},
	)

	httpServer := &http.Server{
		Addr:              "0.0.0.0:8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Println("[WS] =================================")
		log.Println("[WS] AMEN REALTIME SERVER")
		log.Println("[WS] Listening :8081")
		log.Println("[WS] WebSocket: ws://0.0.0.0:8081/ws")

		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("[WS]", err)
		}
	}()
}
