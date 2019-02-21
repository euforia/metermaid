package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/euforia/metermaid"
	"go.uber.org/zap"
)

type priceAPI struct {
	prefix string
	mm     metermaid.Metermaid
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

	priceHistory, err := api.mm.BurnHistory(start, end)
	if err == nil {
		var b []byte
		if b, err = json.Marshal(priceHistory); err == nil {
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
