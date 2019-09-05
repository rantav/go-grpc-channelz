package channelz

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	log "google.golang.org/grpc/grpclog"
)

type channelzHandler interface {
	WriteTopChannelsPage(io.Writer)
	WriteChannelPage(io.Writer, int64)
	WriteSubchannelPage(io.Writer, int64)
	WriteServerPage(io.Writer, int64)
}

func createRouter(prefix string, handler channelzHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Route(prefix, func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handler.WriteTopChannelsPage(w)
		})
		r.Get("/channel/{channel}", func(w http.ResponseWriter, r *http.Request) {
			channelStr := chi.URLParam(r, "channel")
			channel, err := strconv.ParseInt(channelStr, 10, 0)
			if err != nil {
				log.Errorf("channelz: Unable to parse int for channel ID. %s", channelStr)
				return
			}
			handler.WriteChannelPage(w, channel)
		})
		r.Get("/subchannel/{channel}", func(w http.ResponseWriter, r *http.Request) {
			channelStr := chi.URLParam(r, "channel")
			channel, err := strconv.ParseInt(channelStr, 10, 0)
			if err != nil {
				log.Errorf("channelz: Unable to parse int for sub-channel ID. %s", channelStr)
				return
			}
			handler.WriteSubchannelPage(w, channel)
		})
		r.Get("/server/{server}", func(w http.ResponseWriter, r *http.Request) {
			serverStr := chi.URLParam(r, "server")
			server, err := strconv.ParseInt(serverStr, 10, 0)
			if err != nil {
				log.Errorf("channelz: Unable to parse int for server ID. %s", serverStr)
				return
			}
			handler.WriteServerPage(w, server)
		})
	})
	return router
}
