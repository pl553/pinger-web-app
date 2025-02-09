package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func ReportPing(backendUrl, ip string, pingTimeMs int32, status string) {
	data := map[string]any{
		"container_ip": ip,
		"ping_time_ms": pingTimeMs,
		"status":       status,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	r, err := http.Post(backendUrl+"/containers/ping", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf(err.Error())
		return
	}
	if err = r.Body.Close(); err != nil {
		log.Printf(err.Error())
	}
}
