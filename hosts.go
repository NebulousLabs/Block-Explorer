package main

import (
	"encoding/json"
	"net/http"

	"github.com/NebulousLabs/Sia/api"
	"github.com/NebulousLabs/Sia/modules"
	"github.com/NebulousLabs/Sia/types"
)

// hostDisplayInfo contains the data about each host that will be displayed on
// the hosts page.
type hostDisplayInfo struct {
	IPAddress    modules.NetAddress
	TotalStorage int64
	Price        types.Currency
}

func NewHostDisplayInfo(hostSettings modules.HostSettings) (hdi hostDisplayInfo) {
	return hostDisplayInfo{hostSettings.IPAddress, hostSettings.TotalStorage, hostSettings.Price}
}

func (es *ExploreServer) hostsHandler(w http.ResponseWriter, r *http.Request) {
	// Query the host host database for all hosts
	hostMessage, err := es.apiGet("/hostdb/hosts/active")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var hl api.ActiveHosts
	err = json.Unmarshal(hostMessage, &hl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert each host into display format.
	displayHosts := make([]hostDisplayInfo, len(hl.Hosts))
	for i, host := range hl.Hosts {
		displayHosts[i] = NewHostDisplayInfo(host)
	}

	hostsJson, err := json.Marshal(displayHosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(hostsJson)
}
