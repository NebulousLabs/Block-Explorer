package main

import (
	"encoding/json"
	"net/http"

	"github.com/NebulousLabs/Sia/modules"
)

// Make new struct with more fields than what will be returned
type hostList struct {
	Hosts  []modules.HostSettings
	Length int
}

// hostsHandler handles the /hosts call, which gives a page with
// information about the hosts network and database
func (es *ExploreServer) hostsHandler(w http.ResponseWriter, r *http.Request) {
	// Query the host host database for all hosts
	hlb, err := es.apiGet("/hostdb/hosts/active")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var hl hostList
	err = json.Unmarshal(hlb, &hl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Populate the hostList struct further
	hl.Length = len(hl.Hosts)

	// Put that information through the template
	page, err := es.parseTemplate("hosts.html", hl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(page)
}
