package main

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/NebulousLabs/Sia/modules"
	"github.com/NebulousLabs/Sia/types"
)

type (
	// This will be the root struct given to the template parser
	overviewRoot struct {
		Explorer       modules.ExplorerStatus
		BlockSummaries []modules.ExplorerBlockData
	}

	blockRoot struct {
		Block  types.Block
		Height types.BlockHeight
		Target types.Target
		Size   uint64
	}
)

// funcMap is passed to the template engine so that templates may have
// access to these funcitons
var funcMap = template.FuncMap{
	"siacoinString": siacoinString,
	"byteString":    byteString,
	"hashAvgString": hashAvgString,
	"timeString":    timeString,
	"increment":     func(x types.BlockHeight) types.BlockHeight { return x + 1 },
}

// Used in siacoinString and byteString
var coinPostfixes []string = []string{"SC", "KS", "MS", "GS", "TS", "PS"}
var bytePostfixes []string = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func siacoinString(siacoins types.Currency) string {
	coins := float64(siacoins.Div(types.SiacoinPrecision).Big().Int64())

	var i int = 0
	for coins > 1000 {
		coins = coins / 1000
		i++
	}

	return fmt.Sprintf("%f %s", coins, coinPostfixes[i])
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
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	t := time.Unix(int64(epoch), 0)
	return t.Format(layout)
}

// hashAvgString is a wrapper for the hashrate function, found in hashrate.go
func hashAvgString(numBlocks types.BlockHeight, o overviewRoot) (s string) {
	if int(numBlocks) >= len(o.BlockSummaries) {
		return rateString(hashrate(o.BlockSummaries))
	}
	s = rateString(hashrate(o.BlockSummaries[o.Explorer.Height-numBlocks : o.Explorer.Height]))
	return
}

// parseTemplate is a more generic funciton to parse a template given
// just a template filename and an object to be put into the template,
func parseTemplate(templateName string, data interface{}) ([]byte, error) {
	t, err := template.New(templateName).Funcs(funcMap).ParseFiles("templates/" + templateName)
	if err != nil {
		s := fmt.Sprintf("Error parsing template %s : %s", templateName, err.Error())
		fmt.Println(s)
		return nil, errors.New(s)
	}

	// Use a bytes buffer to avoid requiring the writing
	// object. This way we can return a byte slice
	var doc bytes.Buffer
	err = t.Execute(&doc, data)
	if err != nil {
		s := fmt.Sprintf("Error executing template %s: %s", templateName, err.Error())
		fmt.Println(s)
		return nil, errors.New(s)
	}

	s := doc.Bytes()
	return s, nil
}
