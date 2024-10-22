package keepalive

import (
	"io"
	"net/http"
	"time"

	"github.com/helloacai/spindle/pkg/log"
)

func Start() {
	// TODO: this is also terrible. workaround to make render not shut down our instances
	// also should add better shutdown
	go func() {
		for {
			resp, err := http.Get("https://yelpagent.onrender.com")
			if err != nil {
				log.Err(err).Msg("YelpAgent keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("YelpAgent keepalive body error")
				} else {
					log.Debug().Msg("YelpAgent keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			resp, err = http.Get("https://coordinatingagent.onrender.com")
			if err != nil {
				log.Err(err).Msg("CoordinatingAgent keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("CoordinatingAgent keepalive body error")
				} else {
					log.Debug().Msg("CoordinatingAgent keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			resp, err = http.Get("https://gcalagent.onrender.com")
			if err != nil {
				log.Err(err).Msg("GCalAgent keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("GCalAgent keepalive body error")
				} else {
					log.Debug().Msg("GCalAgent keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			resp, err = http.Get("https://spindle.onrender.com/healthz")
			if err != nil {
				log.Err(err).Msg("spindle keepalive error")
			} else {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Err(err).Msg("spindle keepalive body error")
				} else {
					log.Debug().Msg("spindle keepalive: " + string(b))
				}
			}
			resp.Body.Close()

			time.Sleep(10 * time.Second)
		}
	}()
}
