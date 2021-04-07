package hooker

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var hookAddr = "http://www.servicechain.ru/hook/?token=%s&TicketID=%v"

type IChecker interface {
	IsAllow(clientID int) bool
}

func Hook(id int) {
	token := os.Getenv("TOKEN")
	if os.Getenv("HOOK_URL") != "" {
		hookAddr = os.Getenv("HOOK_URL")
	}
	url := fmt.Sprintf(hookAddr, token, id)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithField("id", id).Error(err)
	}
	if resp.StatusCode != 200 {
		logrus.WithField("id", id).Error("StatusCode=" + resp.Status)
	}
}

func Handle(client IChecker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		issue := &Issue{}
		err := json.NewDecoder(r.Body).Decode(issue)
		if err != nil {
			logrus.WithField("issue", issue).Error(err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		w.WriteHeader(200)
		w.Write([]byte("ok"))

		if client.IsAllow(issue.CompanyID) {
			go Hook(issue.ID)
		}
		return
	}
}

type Issue struct {
	ID        int `json:"id"`
	CompanyID int `json:"company_id"`
}
