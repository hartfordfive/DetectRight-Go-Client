#DetectRight Go Client
====================

Device Detection and Improved Google Analytics using the DetectRight web service (in Go)

Original code and web service created by Chris Abbott.

Ported to the Go language by Alain Lefebvre

## Installation

Simply intsall the package with the "go get" tool:

go get github.com/DetectRight/DetectRight-Go-Client/detectright


## Usage Example
====================

<pre><code>
// Initialize the DetectRigh Go client
drc := detectright.InitClient()

// Store all the headers from the current request in header map
drcHeaders := make(map[string]interface{})
for k, v := range req.Header {
  drcHeaders[k] = v[0]
}

// Sets the headers of the current rquest
drc.SetHeaders(drcHeaders)

// Fetches the device profile from HQ with the collected headers
drc.GetProfileFromHeaders()

// Get all the profile properties
response := drc.GetProperties()
</code></pre>

## Version Highlights
=======================

###V 0.3.0-beta:
- Added in process device profile caching
- Added cache_ttl parameter to config
- Added analytics reporting for cached device profiles
