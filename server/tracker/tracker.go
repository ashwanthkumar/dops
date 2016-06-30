package tracker

import (
	"log"
	"net/http"

	"github.com/ashwanthkumar/dops/config"
	"github.com/chihaya/chihaya/tracker"
	"github.com/julienschmidt/httprouter"
)

// Tracker wraps Chihaya's Tracker related implementation
type Tracker struct {
	tkr *tracker.Tracker
	cfg *config.Config
}

// New creates new Tracker instance
func New(config *config.Config) (*Tracker, error) {
	tkr, err := tracker.NewTracker(config.ToChihayaTrackerConfig())
	if err != nil {
		return nil, err
	}
	return &Tracker{
		tkr: tkr,
		cfg: config,
	}, nil
}

// ServeAnnounce is the handler for /announce URI
func (t *Tracker) ServeAnnounce(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req, err := announceRequest(r, &t.cfg.TrackerConfig)
	if err != nil {
		writeError(w, err)
		return
	}

	resp, err := t.tkr.HandleAnnounce(req)
	if err != nil {
		writeError(w, err)
		return
	}

	err = writeAnnounceResponse(w, resp)
	if err != nil {
		log.Println("error serializing response", err)
	}
}

// ServeScrape is the handler for /scrape URI
func (t *Tracker) ServeScrape(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	req, err := scrapeRequest(r, &t.cfg.TrackerConfig)
	if err != nil {
		writeError(w, err)
		return
	}

	resp, err := t.tkr.HandleScrape(req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeScrapeResponse(w, resp)
}
