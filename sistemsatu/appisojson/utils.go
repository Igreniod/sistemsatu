package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mofax/iso8583"
)

func printSortedDE(parsedMessage iso8583.IsoStruct) {
	dataElement := parsedMessage.Elements.GetElements()
	int64toSort := make([]int, 0, len(dataElement))
	for key := range dataElement {
		int64toSort = append(int64toSort, int(key))
	}
	sort.Ints(int64toSort)
	for _, key := range int64toSort {
		// log.Printf("[%v] : %v\n", int64(key), dataElement[int64(key)])
		fmt.Printf("[%v] : %v\n", int64(key), dataElement[int64(key)])
	}
}

func cetakDataDiTerminal(dataDiterima Pesan) {
	dataFormatJSON, err := json.MarshalIndent(dataDiterima, "", "  ")
	if err != nil {
		fmt.Println("Error mencetak JSON:", err)
		return
	}

	fmt.Println(string(dataFormatJSON))
}
