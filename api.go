package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/NebulousLabs/Sia/modules"
	"github.com/NebulousLabs/Sia/types"
	"github.com/gorilla/mux"
)

// NewAPIRouter handles any urls that begin with '/api/'
func (es *ExploreServer) NewAPIRouter() {
	r := es.router.PathPrefix("/api").Subrouter()

	r.HandleFunc("/block/hash/{hash}", es.getBlock).
		Methods("GET")
	r.HandleFunc("/block/hash/height/{height}", es.getHashByHeight).
		Methods("GET")
	r.HandleFunc("/hosts/", es.getHosts).
		Methods("GET")
	r.HandleFunc("/status/", es.getStatus).
		Methods("GET")
}

func writeJson(w http.ResponseWriter, obj interface{}) {
	if json.NewEncoder(w).Encode(obj) != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Does an arbitrary request to the server referenced by link, returns as a byte array.
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

// ExplorerState queries the locally running status.
func (es *ExploreServer) apiExplorerState() (explorerStatus modules.ExplorerStatus, err error) {
	stateJSON, err := es.apiGet("/explorer/status")
	if err != nil {
		return
	}

	err = json.Unmarshal(stateJSON, &explorerStatus)
	return
}

// GetBlockData queries a range of blocks from the server, and returns that list
func (es *ExploreServer) apiGetBlockData(start types.BlockHeight, end types.BlockHeight) ([]modules.ExplorerBlockData, error) {
	return nil, nil
}

// apiGetHash queries siad and returns the raw data. The JSON data can
// be decoded based on the ResponseType field
func (es *ExploreServer) apiGetHash(hash []byte) ([]byte, error) {
	return es.apiGet(fmt.Sprintf("/explorer/gethash?hash=%x", hash))
}

func (es *ExploreServer) getStatus(w http.ResponseWriter, r *http.Request) {
	status, err := es.apiGet("/explorer/status")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(status)
}

func (es *ExploreServer) getHashByHeight(w http.ResponseWriter, r *http.Request) {
	var height types.BlockHeight

	vars := mux.Vars(r)
	fmt.Sscanf(vars["height"], "%d", &height)
	//_, err := fmt.Sscanf(vars["height"
	v := url.Values{}
	v.Set("start", strconv.Itoa(int(height)))
	v.Add("finish", strconv.Itoa(int(height+1)))
	blockSumJson, err := es.apiGet("/blockexplorer/blockdata?" + v.Encode())
	fmt.Println(blockSumJson)
	if err != nil {
		fmt.Println(err)
		//http.Error(w, "Failed to get blockdata", http.StatusInternalServerError)
	}

	//var blockSummaries []modules.ExplorerBlockData
	//err = json.Unmarshal(blockSumJson, &blockSummaries)

	w.Write(blockSumJson)
	//writeJson(w, blockSumJson)
}

func (es *ExploreServer) getHash(w http.ResponseWriter, r *http.Request) {
}

func (es *ExploreServer) getHosts(w http.ResponseWriter, r *http.Request) {
	hostsJSON, err := es.apiGet("/hostdb/hosts/active")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(hostsJSON)
}

func (es *ExploreServer) getBlock(w http.ResponseWriter, r *http.Request) {
	var hash string
	vars := mux.Vars(r)
	hash = vars["hash"]

	fmt.Println(hash)
	blockJson, err := es.apiGet("/explorer/gethash?hash=" + hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(blockJson)
}
