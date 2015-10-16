package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NebulousLabs/Sia/modules"
	"github.com/NebulousLabs/Sia/types"
	"github.com/gorilla/mux"
)

// NewAPIRouter handles any urls that begin with '/api/'
func (es *ExploreServer) NewAPIRouter() {
	r := es.router.PathPrefix("/api").Subrouter()

	r.HandleFunc("/block/hash/{hash}", es.getBlock).
		Methods("GET")
	r.HandleFunc("/block/height/{height}", es.getBlockByHeight).
		Methods("GET")
	r.HandleFunc("/hosts/", es.getHosts).
		Methods("GET")
	r.HandleFunc("/status/", es.getStatus).
		Methods("GET")
}

// writeJson is responsible for sending back headers, and the json response
// as well as the http status code
func writeJson(w http.ResponseWriter, json []byte, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(json)
}

// apiGet does an arbitrary request to the server referenced by link, returns
// as a byte array.
func (es *ExploreServer) apiGet(apiCall string) (response []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", es.url+apiCall, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Sia-Agent")

	// Do a HTTP request to the Sia daemon
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New("Sia Daemon Returned Non-200: " + string(response))
		return
	}
	return
}

// getStatus returns the siae status output
func (es *ExploreServer) getStatus(w http.ResponseWriter, r *http.Request) {
	status, err := es.apiGet("/explorer/status")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		es.logger.Print(err)
		return
	}
	writeJson(w, status, http.StatusOK)
}

// getBlockByHeight takes an integer and returns the corresponding block at
// that height
func (es *ExploreServer) getBlockByHeight(w http.ResponseWriter, r *http.Request) {
	var height types.BlockHeight

	vars := mux.Vars(r)
	fmt.Sscanf(vars["height"], "%d", &height)

	blockSummaries, err := es.getBlockRange(height, height+1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		es.logger.Print(err)
		return
	}

	blockJson, err := es.apiGet("/explorer/gethash?hash=" + blockSummaries[0].ID.String())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		es.logger.Print(err)
		return
	}
	writeJson(w, blockJson, http.StatusOK)
}

// getHosts returns a list of hosts that on the network
func (es *ExploreServer) getHosts(w http.ResponseWriter, r *http.Request) {
	hostsJSON, err := es.apiGet("/hostdb/hosts/active")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		es.logger.Print(err)
		return
	}

	writeJson(w, hostsJSON, http.StatusOK)
}

// getBlock queries siae and returns a block given the block hash
func (es *ExploreServer) getBlock(w http.ResponseWriter, r *http.Request) {
	var hash string
	vars := mux.Vars(r)
	hash = vars["hash"]

	blockJson, err := es.apiGet("/explorer/gethash?hash=" + hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		es.logger.Print(err)
		return
	}

	writeJson(w, blockJson, http.StatusOK)
}
