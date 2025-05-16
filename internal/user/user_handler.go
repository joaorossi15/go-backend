package user

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/joaorossi15/gobh/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func GetUserByIdHandler(repo *UserR) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get id
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		// try get user
		usrID, usrName, err := repo.Get(r.Context(), name)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		// return user info
		resp := Response{ID: int64(usrID), Name: usrName}
		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(resp); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b.Bytes())
	}
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

func UserLoginHandler(repo *UserR) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var rq UserInput

		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		// get user info
		_, name, err := repo.Get(r.Context(), rq.Name)
		if err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		hashedPwd, err := repo.GetHashedPassword(r.Context(), name)
		if err != nil {
			http.Error(w, "error getting password", http.StatusInternalServerError)
			return
		}

		// compare password
		if err := bcrypt.CompareHashAndPassword(hashedPwd, []byte(rq.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// generate jwt token
		token, err := middleware.GenerateToken(name)
		if err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var b bytes.Buffer
		response := LoginResponse{Token: token}

		if err := json.NewEncoder(&b).Encode(response); err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b.Bytes())
	}
}
