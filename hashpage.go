package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NebulousLabs/Sia/crypto"
	"github.com/NebulousLabs/Sia/types"
)

// Redefine the structs sent by the blockexplorer here. As they are
// sent through the modules module and the api as an interface{}, they
// are only defined in the blockexplorer module of Sia, which should
// not be imported from an external program
type (
	responseData struct {
		ResponseType string
	}

	blockData struct {
		Block  types.Block
		Height types.BlockHeight
	}

	txData struct {
		Tx       types.Transaction
		ParentID types.BlockID
		TxNum    int
	}

	fcInfo struct {
		Contract  crypto.Hash
		Revisions []crypto.Hash
		Proof     crypto.Hash
	}

	outputTransactions struct {
		OutputTx crypto.Hash
		InputTx  crypto.Hash
	}

	addrData struct {
		Txns []crypto.Hash
	}
)

// hashPage handles the delegation of the pages depending on the hash type
func (es *ExploreServer) hashPageHandler(w http.ResponseWriter, r *http.Request) {
	var d []byte
	_, err := fmt.Sscanf(r.FormValue("h"), "%x", &d)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	itemJSON, err := es.apiGetHash(d)

	// Now decode the json and figure out which display function
	// to dispatch it to.
	var rd responseData
	err = json.Unmarshal(itemJSON, &rd)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	switch rd.ResponseType {
	case "Block":
		es.blockPage(w, itemJSON)
		return
	case "Transaction":
		es.txPage(w, itemJSON)
		return
	case "Output":
		es.outputPage(w, itemJSON, d[:crypto.HashSize])
		return
	case "Address":
		es.addressPage(w, itemJSON, d)
		return
	default:
		http.Error(w, "Siad returned: "+string(itemJSON), 500)
	}
}

// parseTransactions iterates over a list of transactions and puts
// each one into an html template, returning the concatinated list
func (es *ExploreServer) parseTransactions(hashes []crypto.Hash) ([]byte, error) {
	var page []byte

	// Request all the other things to be viewed
	for _, hash := range hashes {
		// Don't attempt to parse zero hashes
		if hash == *new(crypto.Hash) {
			continue
		}
		// Decode into a responseData struct to figure out
		// what type of response it is, then switch on it
		itemJSON, err := es.apiGetHash(hash[:])
		var rd responseData
		err = json.Unmarshal(itemJSON, &rd)
		if err != nil {
			return nil, err
		}

		switch rd.ResponseType {
		case "Block":
			var b blockData
			err := json.Unmarshal(itemJSON, &b)
			if err != nil {
				return nil, err
			}

			// The block page requires additional
			// information contained in the block summary
			blockSummaries, err := es.apiGetBlockData(b.Height, b.Height+1)
			if err != nil {
				return nil, err
			}

			parsed, err := parseTemplate("block.template", blockRoot{
				Block:  b.Block,
				Height: b.Height,
				Target: blockSummaries[0].Target,
				Size:   blockSummaries[0].Size,
			})
			if err != nil {
				return nil, err
			}
			page = append(page, parsed...)

			continue
		case "Transaction":
			var t txData
			err := json.Unmarshal(itemJSON, &t)
			if err != nil {
				return nil, err
			}
			// Parse the main transaction template
			parsed, err := parseTemplate("transaction.template", t)
			if err != nil {
				return nil, err
			}
			page = append(page, parsed...)
			continue
		default:
			continue
		}
	}
	return page, nil
}

// blockPage formats information about a block for viewing
func (es *ExploreServer) blockPage(w http.ResponseWriter, blockJSON []byte) {
	var b blockData
	err := json.Unmarshal(blockJSON, &b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	blockSummaries, err := es.apiGetBlockData(b.Height, b.Height+1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse the header template
	page, err := parseTemplate("header.template", "Block "+fmt.Sprintf("%d", b.Height))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// blockRoot is defined in template.go
	parsed, err := parseTemplate("block.template", blockRoot{
		Block:  b.Block,
		Height: b.Height,
		Target: blockSummaries[0].Target,
		Size:   blockSummaries[0].Size,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	// Parse the footer template
	parsed, err = parseTemplate("footer.template", "Block "+fmt.Sprintf("%d", b.Height))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	w.Write(page)
}

// txPage formats information about a transaction in a format suitable
// for human viewing
func (es *ExploreServer) txPage(w http.ResponseWriter, txJSON []byte) {
	var t txData
	err := json.Unmarshal(txJSON, &t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse the header template
	page, err := parseTemplate("header.template", "Transaction")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse the main transaction template
	parsed, err := parseTemplate("transaction.template", t)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	// Parse the footer template
	parsed, err = parseTemplate("footer.template", "Transaction")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	w.Write(page)
}

// outputPage formats an outputTransactions struct for use in viewing
func (es *ExploreServer) outputPage(w http.ResponseWriter, outJSON []byte, oID []byte) {
	var ot outputTransactions
	err := json.Unmarshal(outJSON, &ot)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse header template
	page, err := parseTemplate("header.template", fmt.Sprintf("Output %x", oID))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	parsed, err := es.parseTransactions([]crypto.Hash{ot.OutputTx, ot.InputTx})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	// Parse the footer template
	parsed, err = parseTemplate("footer.template", fmt.Sprintf("Output ID %x", oID))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	w.Write(page)
}

// addressPage formats the list of transactions that an address
// participated in for human consumption
func (es *ExploreServer) addressPage(w http.ResponseWriter, addrJSON []byte, address []byte) {
	var ad addrData
	err := json.Unmarshal(addrJSON, &ad)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse header template
	page, err := parseTemplate("header.template", fmt.Sprintf("Address %x", address))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Call the function to parse all of the transactions
	parsed, err := es.parseTransactions(ad.Txns)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	// Parse the footer template
	parsed, err = parseTemplate("footer.template", fmt.Sprintf("Address %x", address))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	w.Write(page)
}

func (es *ExploreServer) contractPage(w http.ResponseWriter, fcJSON []byte, fcid []byte) {
	var fi fcInfo
	err := json.Unmarshal(fcJSON, &fi)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse header template
	page, err := parseTemplate("header.template", fmt.Sprintf("File Contract %x", fcid))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	txns := []crypto.Hash{fi.Contract}
	txns = append(txns, fi.Revisions...)
	txns = append(txns, fi.Proof)
	// Call the parseTransactions function to help with all the transactions
	parsed, err := es.parseTransactions(txns)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	// Parse the footer template
	parsed, err = parseTemplate("footer.template", fmt.Sprintf("File Contract %x", fcid))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page = append(page, parsed...)

	w.Write(page)
}
