package api

import (
	"mime"
	"net"
	"net/http"
	"path"

	"github.com/euforia/metermaid"
	"github.com/euforia/metermaid/ui"
	"go.uber.org/zap"
)

// API ...
type API struct {
	pricing   *priceAPI
	container *containerAPI
	log       *zap.Logger
}

// New returns a new API instance
func New(mm metermaid.Metermaid, logger *zap.Logger) *API {
	api := &API{
		pricing:   &priceAPI{"/price", mm, logger},
		container: &containerAPI{"/container", mm.Containers()},
		log:       logger,
	}

	http.Handle("/price/", api.pricing)
	http.Handle("/container/", api.container)
	http.HandleFunc("/", handleUI)

	return api
}

// Serve starts serving the API on the listener
func (api *API) Serve(ln net.Listener) error {
	api.log.Info("http server", zap.String("address", ln.Addr().String()))
	err := http.Serve(ln, nil)
	if err != nil {
		api.log.Info("http shutdown unclean", zap.Error(err))
	}
	return err
}

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

func writeErrorReponse(w http.ResponseWriter, e string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(400)
	w.Write([]byte(e))
}

func writeResponse(w http.ResponseWriter, b []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	w.Write(b)
}
