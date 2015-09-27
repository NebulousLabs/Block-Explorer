(function () {
	// Storage for variables in page context.
	$.hosts = {};

	function addrSort(a, b) {
		if (a.IPAddress < b.IPAddress) return -1;
		if (a.IPAddress > b.IPAddress) return 1;
		return 0;
	}

	function storageSort(a, b) {
		return a.TotalStorage - b.TotalStorage;
	}

	function priceSort(a, b) {
		return a.Price - b.Price;
	}

	// Converts a count of bytes into a formatted string describing that number
	// in standard disk storage units (default is GB).
	function formatStorage(storageBytes, format) {
		var suffixes = ["B", "KB", "MB", "GB", "TB"]
		if (suffixes.indexOf(format.toUpperCase()) < 0) {
			format = "GB";
		}
		var suffixIndex = suffixes.indexOf(format.toUpperCase())
		if (suffixIndex < 0) {
			suffixIndex = 2;

		}
		var bytesPer = Math.pow(1000, suffixIndex);
		return (storageBytes / bytesPer).toFixed(2) + " " + format;
	}

	// Deterimines the proper sorting function based on the URI query parameter
	// (default is priceSort).
	function getSortFunc() {
		var query = URI.parseQuery(window.location.search);
		var defaultSort = priceSort;
		if (!("sort" in query)) {
			return defaultSort;
		}
		switch (query["sort"]) {
			case "addr":
				return addrSort;
			case "storage":
				return storageSort;
			case "price":
			default:
				return priceSort;
		}
	}

	// Updates the page's URL without reloading the page.
	function setLocation(uri) {
		window.history.pushState({path:uri},'',uri);
	}

	// Sets the sorting mechanism for the table and updates the table with the
	// specified sorting.
	function setTableSorting(sortingValue) {
		var uri = URI(window.location);
		uri.setQuery("sort", sortingValue);
		setLocation(uri.toString());
		populateHostTable();
	}

	// Populates the table of known hosts.
	function populateHostTable() {
		$('#hostsTable > tbody').empty();
		var hosts = $.hosts.hosts;
		sortFunc = getSortFunc();
		hosts.sort(sortFunc);
		hosts.forEach(function (host) {
			var storageFormatted = formatStorage(host.TotalStorage, "GB");
			// 1 SC = 10^24 H
			var hastingsPerSC = Math.pow(10, 24);
			var priceFormatted = (host.Price / hastingsPerSC).toFixed(2) + " SC";

			$('#hostsTable > tbody:last')
				.append($('<tr>')
					.append($('<td>')
						.text(host.IPAddress))
					.append($('<td>', {
						text: storageFormatted,
						class: "storage",
					}))
					.append($('<td>', {
						text: priceFormatted,
						class: "price",
					}))
					);
		});
	}

	function populateHostCount() {
		$('#hostCount').text($.hosts.hosts.length);
	}

	// Initializes all of the dynamic elements of the page.
	function initPage() {
		$('#addrHeader').click(function() { setTableSorting("addr"); });
		$('#storageHeader').click(function() { setTableSorting("storage"); });
		$('#priceHeader').click(function() { setTableSorting("price"); });
		populateHostCount();
		populateHostTable();
	}

	$.getJSON('/hosts.json', function (hosts) {
		$.hosts.hosts = hosts;
		initPage();
	});
}());
