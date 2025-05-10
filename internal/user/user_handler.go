package user

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func PostUserHandler(repo *UserR) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "type not allowed", http.StatusUnsupportedMediaType)
			return
		}

		// decode json
		var body UserInput
		r.Body = http.MaxBytesReader(w, r.Body, 1<<10)
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// check if user is already registered
		if id, _, _ := repo.Get(r.Context(), body.Name); id != 0 {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		// register if not
		id, err := repo.Create(r.Context(), body.Name, []byte(body.Password))
		if err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// return response
		resp := Response{ID: id, Name: body.Name}
		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(resp); err != nil {
			log.Printf("response: %v", err)
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(b.Bytes())
	}
}
