package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Pesan struct {
	Dua                string `json:"2,omitempty"`
	Tiga               string `json:"3,omitempty"`
	Tujuh              string `json:"7,omitempty"`
	Sebelas            string `json:"11,omitempty"`
	Duabelas           string `json:"12,omitempty"`
	Tigabelas          string `json:"13,omitempty"`
	Delapanbelas       string `json:"18,omitempty"`
	Tigapuluhdua       string `json:"32,omitempty"`
	Tigapuluhtiga      string `json:"33,omitempty"`
	Tigapuluhtujuh     string `json:"37,omitempty"`
	Tigapuluhsembilan  string `json:"39,omitempty"`
	Empatpuluhsatu     string `json:"41,omitempty"`
	Empatpuluhtujuh    string `json:"47,omitempty"`
	Empatpuluhsembilan string `json:"49,omitempty"`
	Enampuluhdua       string `json:"62,omitempty"`
	Seratusduapuluh    string `json:"120,omitempty"`
}

func olahRequest(w http.ResponseWriter, r *http.Request) {
	var newPesan Pesan
	err := json.NewDecoder(r.Body).Decode(&newPesan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//-----Cetak data di terminal-----//
	fmt.Println("Data diterima dari App jsonjson / isojson:")
	cetakDataDiTerminal(newPesan)

	switch newPesan.Tiga {
	case "310101":
		newPesan.Tigapuluhsembilan = "00"
	default:
		newPesan.Tigapuluhsembilan = "13"
	}

	//-----Response ke port 8081-----//
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newPesan)
}
