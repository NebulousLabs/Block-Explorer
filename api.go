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
)

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
	v := url.Values{}
	v.Set("start", strconv.Itoa(int(start)))
	v.Add("finish", strconv.Itoa(int(end)))
	blockSumJson, err := es.apiGet("/explorer/blockdata?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var blockSummaries []modules.ExplorerBlockData
	err = json.Unmarshal(blockSumJson, &blockSummaries)
	return blockSummaries, err
}

// apiGetHash queries siad and returns the raw data. The JSON data can
// be decoded based on the ResponseType field
func (es *ExploreServer) apiGetHash(hash []byte) ([]byte, error) {
	return es.apiGet(fmt.Sprintf("/explorer/gethash?hash=%x", hash))
}
