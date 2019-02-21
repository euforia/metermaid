package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/euforia/metermaid/fl"
	"github.com/euforia/metermaid/storage"
	"github.com/euforia/metermaid/types"
)

type containerAPI struct {
	prefix string
	store  storage.Containers
}

func (api *containerAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, api.prefix)
	if p == "/" {
		api.handleQuery(w, r)
		return
	}

	switch r.Method {
	case "GET":
		resp, err := api.store.Get(r.URL.Path[1:])
		switch err {
		case nil:
			b, _ := json.Marshal(resp)
			w.WriteHeader(200)
			w.Write(b)
		case storage.ErrNotFound:
			w.WriteHeader(404)
		default:
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
		}
	default:
		w.WriteHeader(405)
	}
}

func (api *containerAPI) handleQuery(w http.ResponseWriter, r *http.Request) {
	query := fl.ParseQuery(r.URL.Query())
	out := make([]types.Container, 0)
	api.store.Iter(func(c types.Container) error {
		if c.Match(query) {
			out = append(out, c)
		}
		return nil
	})

	b, _ := json.Marshal(out)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	w.Write(b)
}
