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
	// Do a http request to the sia daemon
	resp, err := http.Get(es.url + apiCall)
	if err != nil {
		// err is already set
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		err = errors.New("Sia Daemon Returned Non-200: " + string(response))
		return
	}

	return
}

// ExplorerState queries the locally running statu
func (es *ExploreServer) apiExplorerState() (explorerStatus modules.ExplorerStatus, err error) {
	stateJSON, err := es.apiGet("/blockexplorer/status")
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
	blockSumJson, err := es.apiGet("/blockexplorer/blockdata?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var blockSummaries []modules.ExplorerBlockData

	// Attepmt to interpret as a block
	err = json.Unmarshal(blockSumJson, &blockSummaries)

	return blockSummaries, err
}

// apiGetHash queries the running instance of siad and returns the raw
// data to the calling function. The json data could decode into
// multiple types of structures, all of which *should* have the
// ResponseType field as a string for identifying the type of data
func (es *ExploreServer) apiGetHash(hash []byte) ([]byte, error) {
	v := url.Values{}
	v.Set("hash", fmt.Sprintf("%x", hash))
	resultJSON, err := es.apiGet("/blockexplorer/gethash?" + v.Encode())
	if err != nil {
		return nil, err
	}

	return resultJSON, nil
}
