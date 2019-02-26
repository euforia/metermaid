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

// ServeHTTP serves the pricing api.  All information is from the perspective of
// the node
func (api *priceAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	start, end, err := parseDateRange(r.URL.Query())
	if err != nil {
		writeErrorReponse(w, err.Error()+"\n"+
			"must be RFC3339 https://tools.ietf.org/html/rfc3339\n")
		return
	}

	priceHistory, err := api.mm.PriceReport(start, end)
	if err == nil {
		var b []byte
		if b, err = json.Marshal(priceHistory); err == nil {
			writeResponse(w, b)
			return
		}
	}

	writeErrorReponse(w, err.Error())
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
