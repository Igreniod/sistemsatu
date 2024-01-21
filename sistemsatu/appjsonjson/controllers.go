package main

import (
	"bytes"
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

type ResponsePesan struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func olahRequest(w http.ResponseWriter, r *http.Request) {
	var newPesan Pesan
	err := json.NewDecoder(r.Body).Decode(&newPesan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//-----Cetak data di terminal-----//
	fmt.Println("Data diterima dari HTTP Req:")
	cetakDataDiTerminal(newPesan)

	//-----Kirim/POST data ke app responder-----//
	pesanInJSON, err := json.MarshalIndent(newPesan, "", "  ")
	if err != nil {
		fmt.Println("Error mencetak JSON:", err)
		return
	}
	URLAppResponder := "http://localhost:8081/request"
	respon, err := http.Post(URLAppResponder, "application/json", bytes.NewBuffer(pesanInJSON))
	if err != nil {
		fmt.Println("Error saat POST ke app responder", err)
		return
	}

	//-----decode reponse dari app responder-----//
	json.NewDecoder(respon.Body).Decode(&newPesan)
	fmt.Println("Response dari app responder :")
	cetakDataDiTerminal(newPesan)
	defer respon.Body.Close()

	responseKeRequester := ResponsePesan{
		Status: "",
		Data:   newPesan,
	}

	//-----Mengatur apakah pesan sukses atau gagal-----//
	switch newPesan.Tigapuluhsembilan {
	case "00":
		responseKeRequester.Status = "Success"
	case "13":
		responseKeRequester.Status = "Failed"
	}

	// //-----Response ke web channel-----//
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseKeRequester)

}
