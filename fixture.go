package speedy

import (
	"net/http"
	"strconv"
)

// debug data
var ResponseHtml = `
<html>
	<head>
		<title>SPDY</title>
		<script type="text/javascript" src="https://localhost:3000/test.js"></script>
	</head>
	<body>
		<h1>Speedy :)</h1>
	</body>
</html>
`

var HeadersFixtureHtml = http.Header{
	":version":       []string{"http/1.1"},
	":status":        []string{"200 OK"},
	":host":          []string{"localhost:3000"},
	":path":          []string{"/"},
	":scheme":        []string{"https"},
	"location":       []string{"https://localhost:3000/"},
	"content-type":   []string{"text/html; charset=utf-8"},
	"content-length": []string{strconv.Itoa(len(ResponseHtml))},
	"server":         []string{"speedy"},
}

var ResponseJS = `
console.log("Speedy");
`

var HeadersFixtureJS = http.Header{
	":version":       []string{"http/1.1"},
	":status":        []string{"200 OK"},
	":host":          []string{"localhost:3000"},
	":path":          []string{"/test.js"},
	":scheme":        []string{"https"},
	"location":       []string{"https://localhost:3000/test.js"},
	"content-type":   []string{"text/javascript; charset=utf-8"},
	"content-length": []string{strconv.Itoa(len(ResponseJS))},
	"server":         []string{"speedy"},
}
