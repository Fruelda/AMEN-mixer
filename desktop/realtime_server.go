package main

import (
	"log"
	"net/http"
	"time"
)

// ============================================================
// START REALTIME SERVER
// ============================================================

func StartRealtimeServer(
	server *RealtimeServer,
) {

	mux :=
		http.NewServeMux()

	// ========================================================
	// WEBSOCKET ENDPOINT
	// ========================================================

	mux.HandleFunc(
		"/ws",
		server.HandleWebSocket,
	)

	// ========================================================
	// HEALTH CHECK
	// ========================================================

	mux.HandleFunc(
		"/",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			_, _ =
				w.Write(
					[]byte(
						"AMEN MIXER REALTIME SERVER OK",
					),
				)
		},
	)

	// ========================================================
	// START HTTP SERVER
	// ========================================================

	go func() {

		log.Println(
			"[WS] =================================",
		)

		log.Println(
			"[WS] AMEN REALTIME SERVER",
		)

		log.Println(
			"[WS] Listening :8081",
		)

		log.Println(
			"[WS] WebSocket: ws://0.0.0.0:8081/ws",
		)

		// ====================================================
		// SERVER
		// ====================================================

		httpServer :=
			&http.Server{
				Addr: "0.0.0.0:8081",

				Handler: mux,

				ReadHeaderTimeout: 5 * time.Second,
			}

		// ====================================================
		// LISTEN
		// ====================================================

		err :=
			httpServer.ListenAndServe()

		if err != nil {

			log.Println(
				"[WS]",
				err,
			)
		}
	}()
}
