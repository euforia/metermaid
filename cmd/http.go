package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/euforia/metermaid/fl"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/storage"
)

type nodeAPI struct {
	prefix string
	store  storage.Nodes
}

func (api *nodeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, api.prefix)
	if p == "/" {
		api.handleQuery(w, r)
		return
	}
	w.WriteHeader(404)
}

func (api *nodeAPI) handleQuery(w http.ResponseWriter, r *http.Request) {
	query := fl.ParseQuery(r.URL.Query())
	gb, ok := query["groupBy"]
	if ok {
		delete(query, "groupBy")
	}

	// Filter
	nodes := make(node.Nodes, 0)
	api.store.Iter(func(c node.Node) error {
		if c.Match(query) {
			nodes = append(nodes, c)
		}
		return nil
	})

	// Group
	var out interface{}
	if ok {
		out = nodes.GroupBy(gb[0].Values[0])
	} else {
		out = nodes
	}

	b, _ := json.Marshal(out)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	w.Write(b)
}
