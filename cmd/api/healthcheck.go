// Filename: cmd/api/healthcheck.go

package main

import (
	"net/http"
)

// healthcheckHandler gives us the health of the system
func (a *application) healthcheckHandler(w http.ResponseWriter,
	r *http.Request) {

	// panic("Apples & Oranges") // deliberate panic

	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.env,
			"version":     a.config.version,
		},
	}

	err := a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
