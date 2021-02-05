package handlers

import (
	"github.com/pusher/pusher-http-go"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/helpers"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// PusherAuth authenticates for pusher
func (repo *DBRepo) PusherAuth(app config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("403 HTTP status code returned!"))
		} else {
			firstName := app.Session.GetString(r.Context(), "userName")
			userID := app.Session.GetInt(r.Context(), "userID")
			params, _ := ioutil.ReadAll(r.Body)

			presenceData := pusher.MemberData{
				UserID: strconv.Itoa(userID),
				UserInfo: map[string]string{
					"name": firstName,
				},
			}

			response, err := app.WsClient.AuthenticatePresenceChannel(params, presenceData)
			if err != nil {
				log.Println(err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(response)
		}
	}
}
