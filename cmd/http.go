package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/euforia/gossip"
	"github.com/euforia/metermaid/storage"
)

type containerAPI struct {
	prefix string
	store  storage.Containers
}

func (api *containerAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, api.prefix)
	if p == "/" {
		list, _ := api.store.List()
		b, _ := json.Marshal(list)
		w.WriteHeader(200)
		w.Write(b)
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

type nodeAPI struct {
	prefix string
	pool   *gossip.Pool
}

func (api *nodeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	list := api.pool.Members()
	b, err := json.Marshal(list)
	if err == nil {
		w.WriteHeader(200)
		w.Write(b)
		return
	}
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}
