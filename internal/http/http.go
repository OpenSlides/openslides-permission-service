// Package http provides HTTP handler to give other services access to to
// permission service.
package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const prefix = "/internal/permission"

// IsAlloweder provides the IsAllowed method.
type IsAlloweder interface {
	IsAllowed(ctx context.Context, name string, userID int, dataList [](map[string]json.RawMessage)) (bool, error)
}

// IsAllowed registers a handler, to connect to the IsAllowed method.
func IsAllowed(mux *http.ServeMux, provider IsAlloweder) {
	url := prefix + "/is_allowed"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Can't read request body: "+err.Error(), 500)
			return
		}

		var requestData struct {
			Name     string                         `json:"name"`
			UserID   int                            `json:"user_id"`
			DataList [](map[string]json.RawMessage) `json:"data"`
		}
		if err := json.Unmarshal(b, &requestData); err != nil {
			http.Error(w, fmt.Sprintf("Can not decode request body '%s': %v", b, err), 500)
			return
		}

		allowed, err := provider.IsAllowed(r.Context(), requestData.Name, requestData.UserID, requestData.DataList)

		if err != nil {
			http.Error(w, "Internal Error. Norman, Do not sent it to client: "+err.Error(), 500)
			return

		}

		value := "false"
		if allowed {
			value = "true"
		}
		fmt.Fprintln(w, value)
	})

	mux.Handle(url, handler)
}

type allrouter interface {
	AllRoutes() ([]string, []string)
}

// Health registers a handler, that tells, if the service is running.
func Health(mux *http.ServeMux, router allrouter) {
	url := prefix + "/health"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var rData struct {
			Info struct {
				Routes struct {
					Read  []string `json:"read"`
					Write []string `json:"write"`
				} `json:"routes"`
			} `json:"healthinfo"`
		}
		rData.Info.Routes.Read, rData.Info.Routes.Write = router.AllRoutes()
		if err := json.NewEncoder(w).Encode(rData); err != nil {
			http.Error(w, "Something went wrong", 500)
			return
		}
	})

	mux.Handle(url, handler)
}
