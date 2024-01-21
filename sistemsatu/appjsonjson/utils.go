package main

import (
	"encoding/json"
	"fmt"
)

func cetakDataDiTerminal(dataDiterima Pesan) {
	dataFormatJSON, err := json.MarshalIndent(dataDiterima, "", "  ")
	if err != nil {
		fmt.Println("Error mencetak JSON:", err)
		return
	}

	fmt.Println(string(dataFormatJSON))
}
