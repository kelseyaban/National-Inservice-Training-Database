// Filename: cmd/examples/cors/basic/main.go
package main

import (
	"flag"
	"log"
	"net/http"
)

// create a simple HTML page with some JS added.
// in a professional setting, we would have the JS code in a script file
const html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
</head>
<body>
    <h1>Quotes CORS</h1>
    <pre id="output"></pre>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            fetch("http://localhost:4000/v1/quotes")
                .then(res => res.json())
                .then(data => {
                    document.getElementById("output").textContent = JSON.stringify(data, null, 2);
                })
                .catch(err => {
                    document.getElementById("output").textContent = "Error: " + err;
                });
        });
    </script>
</body>
</html>`

// Simple HTTP server
func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()

	log.Printf("starting server on %s", *addr)
	err := http.ListenAndServe(*addr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(html))
		}))
	log.Fatal(err)
}
