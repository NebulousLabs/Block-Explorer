package main

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/NebulousLabs/Sia/crypto"
	"github.com/NebulousLabs/Sia/modules"
	"github.com/NebulousLabs/Sia/types"
)

type (
	// This will be the root struct given to the template parser
	overviewRoot struct {
		Explorer       modules.ExplorerStatus
		BlockSummaries []modules.ExplorerBlockData
		NodeVersion    string
	}

	blockRoot struct {
		Block  types.Block
		Height types.BlockHeight
		Target types.Target
		Hashes uint64
		Size   uint64
	}

	outputRoot struct {
		OutputTx types.TransactionID
		InputTx  types.TransactionID
		OutputID crypto.Hash
	}

	filecontractRoot struct {
		Contract  types.TransactionID
		Revisions []types.TransactionID
		Proof     types.TransactionID
		FcID      types.FileContractID
	}

	addressRoot struct {
		Txns []types.TransactionID
		Addr []byte
	}
)

// Used in siacoinString and byteString
var coinPostfixes []string = []string{"H", "zS", "aS", "fS", "pS", "nS", "µS", "mS", "SC", "KS", "MS", "GS", "TS", "PS"}
var bytePostfixes []string = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// siacoinString, byteString, and timeString, are used for formatting numbers
// in a human readable way inside the template
func siacoinString(siacoins types.Currency) string {
	if siacoins.Cmp(types.NewCurrency64(0)) == 0 {
		return "0 SC"
	}

	coinStr := siacoins.String()

	unitIndex := (len(coinStr) - 1) / 3

	if unitIndex >= len(coinPostfixes) {
		unitIndex = len(coinPostfixes) - 1
	}

	decPoint := len(coinStr) - (unitIndex * 3)

	postStr := coinStr[decPoint:]
	if len(postStr) > 6 {
		postStr = postStr[:6]
	}

	return fmt.Sprintf("%s.%s %s", coinStr[:decPoint], postStr, coinPostfixes[unitIndex])
}

func byteString(numBytes uint64) string {
	numBytesF := float64(numBytes)

	var i int = 0
	for numBytesF > 1024 {
		numBytesF = numBytesF / 1024
		i++
	}

	return fmt.Sprintf("%f %s", numBytesF, bytePostfixes[i])
}

func timeString(epoch types.Timestamp) string {
	// layout shows by example how the reference time should be
	// represented. #Golang Magic
	return time.Unix(int64(epoch), 0).Format("Jan 2, 2006 at 3:04pm (MST)")
}

// hashAvgString is a wrapper for the hashrate function, found in hashrate.go
func hashAvgString(numBlocks types.BlockHeight, o overviewRoot) string {
	if int(numBlocks) >= len(o.BlockSummaries) {
		return rateString(hashrate(o.BlockSummaries))
	}
	return rateString(hashrate(o.BlockSummaries[o.Explorer.Height-numBlocks : o.Explorer.Height]))
}

// parseTemplate is a more generic function to parse a template given
// just a template filename and an object to be put into the template,
func (es *ExploreServer) parseTemplate(templateName string, data interface{}) ([]byte, error) {

	// funcMap is passed to the template engine so that templates may have
	// access to these functions. Defined here to give access to functions of es.
	var funcMap = template.FuncMap{
		"siacoinString":    siacoinString,
		"byteString":       byteString,
		"hashAvgString":    hashAvgString,
		"timeString":       timeString,
		"parseTemplate":    es.parseTemplate,
		"parseTransaction": es.parseTransaction,
		"findOutput":       es.findOutput,
		"increment":        func(x types.BlockHeight) types.BlockHeight { return x + 1 },
		"uint64":           func(x int) uint64 { return uint64(x) },
		"int64_uint64":     func(x int64) uint64 { return uint64(x) },
		"txidToHash":       func(x types.TransactionID) crypto.Hash { return crypto.Hash(x) },
	}

	t, err := template.New(templateName).Funcs(funcMap).ParseFiles("templates/" + templateName)
	if err != nil {
		s := fmt.Sprintf("Error parsing template %s: %v", templateName, err)
		fmt.Println(s)
		return nil, errors.New(s)
	}

	// Use a bytes buffer to avoid requiring the writing
	// object. This way we can return a byte slice
	var doc bytes.Buffer
	err = t.Execute(&doc, data)
	if err != nil {
		s := fmt.Sprintf("Error executing template %s: %v", templateName, err)
		fmt.Println(s)
		return nil, errors.New(s)
	}

	return doc.Bytes(), nil
}
