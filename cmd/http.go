package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/euforia/gossip"
	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/storage"
	"github.com/euforia/metermaid/ui"
)

func handleUI(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path[1:]
	if upath == "" {
		upath = "index.html"
	}

	data, err := ui.Asset(upath)
	if err == nil {
		contentType := mime.TypeByExtension(path.Ext(upath))
		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(200)
		w.Write(data)
		return
	}

	w.WriteHeader(404)
}

type containerAPI struct {
	prefix string
	node   *node.Node
	store  storage.Containers
}

func (api *containerAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, api.prefix)
	if p == "/" {
		list, _ := api.store.List()
		b, _ := json.Marshal(list)
		w.Header().Set("Node-Name", api.node.Name)
		w.Header().Set("Node-Addr", api.node.Address)
		w.Header().Set("Node-CPU", fmt.Sprintf("%d", api.node.CPUShares))
		w.Header().Set("Node-Memory", fmt.Sprintf("%d", api.node.Memory))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Node-Name,Node-Addr,Node-CPU,Node-Memory")
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

	nodes := make([]node.Node, len(list))
	for i, item := range list {
		nodes[i] = *node.NewFromMemberlistNode(item)
	}

	b, err := json.Marshal(nodes)
	if err == nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		w.Write(b)
		return
	}
	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}
