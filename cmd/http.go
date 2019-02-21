package main

import (
	"encoding/json"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/euforia/metermaid/tsdb"

	"github.com/euforia/metermaid/fl"
	"github.com/euforia/metermaid/types"

	"go.uber.org/zap"

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
	out := make([]node.Node, 0)
	api.store.Iter(func(c node.Node) error {
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

type price struct {
	Total   float64
	History tsdb.DataPoints
}

func newPrice(data tsdb.DataPoints, per time.Duration) *price {
	out := &price{History: data}
	list := data.Per(per)
	out.Total = list.Sum()
	return out
}

type priceAPI struct {
	prefix string
	mm     *meterMaid
	log    *zap.Logger
}

func (api *priceAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start, end, err := parseDateRange(r.URL.Query())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error() + "\n" +
			"must be RFC3339 https://tools.ietf.org/html/rfc3339\n"))
		return
	}

	list, err := api.mm.BurnHistory(start, end)
	if err == nil {
		// Per hour as that is what aws provides
		d, _ := time.ParseDuration("1h")
		out := newPrice(list, d)

		var b []byte
		if b, err = json.Marshal(out); err == nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(200)
			w.Write(b)
			return
		}
	}

	w.WriteHeader(400)
	w.Write([]byte(err.Error()))
}

func parseDateRange(params url.Values) (start, end time.Time, err error) {
	startStr := params["start"]
	if len(startStr) > 0 && startStr[0] != "" {
		start, err = time.Parse(time.RFC3339, startStr[0])
	}

	endStr := params["end"]
	if len(endStr) > 0 && endStr[0] != "" {
		end, err = time.Parse(time.RFC3339, endStr[0])
	} else {
		end = time.Now()
	}
	return
}
