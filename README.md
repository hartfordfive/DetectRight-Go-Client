#DetectRight Go Client
====================

Device Detection and Improved Google Analytics using the DetectRight web service (in Go)

Original code and web service created by Chris Abbott.

Ported to the Go language by Alain Lefebvre



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
