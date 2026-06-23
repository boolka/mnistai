package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/boolka/goai/pkg/network"
)

func NewServer(port int, net *network.Network) (*http.Server, error) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting CWD:", err)
	}
	fmt.Println("Current working directory:", dir)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /prediction", func(rw http.ResponseWriter, r *http.Request) {
		var netInput PredictionRequest

		if err := json.NewDecoder(r.Body).Decode(&netInput); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		out := net.Activate(netInput.Inputs)

		netOutput := PredictionResponse{
			Outputs: out,
		}

		if err := json.NewEncoder(rw).Encode(&netOutput); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}
	})

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", port),
	}

	if err := srv.ListenAndServe(); err != nil {
		return nil, err
	}

	return srv, nil
}
