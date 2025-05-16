package message

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joaorossi15/gobh/internal/middleware"
	"github.com/joaorossi15/gobh/internal/sqlc"
	"github.com/joaorossi15/gobh/internal/user"
)

func CreateMessageHandler(repo *MessageRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "type not allowed", http.StatusUnsupportedMediaType)
			return
		}

		// get creation params
		var params sqlc.CreateMessageParams
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "incorrect json format", http.StatusBadRequest)
			return
		}

		id, err := repo.CreateMessage(r.Context(), params.SenderID, params.RecID, params.Body)
		if err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(id); err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(b.Bytes())
	}
}

func GetConversationMessagesHandler(repo *MessageRepo, userRepo *user.UserR) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// add logic to only get
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get creation params
		v := r.Context().Value(middleware.UserIDKey)
		userName, ok := v.(string)
		if !ok {
			http.Error(w, "error in getting user: "+userName, http.StatusBadRequest)
			return
		}

		userID, _, err := userRepo.Get(r.Context(), userName)

		recID, err := strconv.Atoi(r.PathValue("recID"))
		if err != nil {
			http.Error(w, "error in URI: "+err.Error(), http.StatusBadRequest)
			return
		}

		data, err := repo.GetConversationMessages(r.Context(), int64(userID), int64(recID))
		if err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(data); err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b.Bytes())
	}
}
