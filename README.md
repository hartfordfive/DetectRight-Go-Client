#DetectRight Go Client
====================

Device Detection and Improved Google Analytics using the DetectRight web service (in Go)

Original code and web service created by Chris Abbott.

Ported to the Go language by Alain Lefebvre



## Usage Example
====================

// Initialize the DetectRigh Go client<br/>
<code>
drc := detectright.InitClient()
</code><br/>

// Store all the headers from the current request in header map<br/>
<code>
drcHeaders := make(map[string]interface{})<br/>
for k, v := range req.Header {<br/>
  drcHeaders[k] = v[0]<br/>
}<br/>
</code><br/>

// Sets the headers of the current rquest<br/>
<code>
drc.SetHeaders(drcHeaders)
</code><br/>

// Fetches the device profile from HQ with the collected headers<br/>
<code>
drc.GetProfileFromHeaders()
</code><br/>

// Get all the profile properties<br/>
<code>
response := drc.GetProperties()
</code><br/>
