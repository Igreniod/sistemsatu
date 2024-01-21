package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mofax/iso8583"
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

func IsoHandler(dataDariTCP string) string {

	//-----Menyimpan pesan ISO dalam variabel-----//
	pesanDalamISO, err := bacaHeader(dataDariTCP)
	if err != nil {
		fmt.Println("Error Baca : ", err)
	}
	//-----Memparsing pesan ISO ke struktur yang telah ditentukan-----//
	strukturISO := iso8583.NewISOStruct("spec1987.yml", true)
	parsed, err := strukturISO.Parse(pesanDalamISO)
	if err != nil {
		fmt.Println(err)
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to unpack valid isomsg")
	}
	if isomsgUnpacked != pesanDalamISO {
		fmt.Printf("%s should be %s \n", isomsgUnpacked, pesanDalamISO)
	}

	printSortedDE(parsed) //-----mencetak data elements
	responsKeTCP := response(parsed)
	packed, _ := responsKeTCP.ToString()
	fmt.Println("Hasil response : ")
	printSortedDE(responsKeTCP)
	strNum1 := fmt.Sprintf("%04d", len(packed))
	rspTCP := strNum1 + packed

	// fmt.Println("response ke TCP : ", rspTCP)
	// bitmapHex, err := iso8583.BitMapArrayToHex(responsKeTCP.Bitmap)
	// fmt.Println("Bitmap Hex: ", bitmapHex)

	return rspTCP
}

func bacaHeader(pesanDalamISO string) (string, error) {
	header := pesanDalamISO[:4]

	intHeader, err := strconv.Atoi(header)
	if err != nil {
		fmt.Println("Error : ", err)
		return "", err
	}

	isoMSG := pesanDalamISO[4 : 4+intHeader]

	// fmt.Println("nilai integer : ", intHeader)
	// fmt.Println(isoMSG)

	return isoMSG, nil
}

func response(dataISO iso8583.IsoStruct) (responseISO iso8583.IsoStruct) {

	// bit := dataISO.Elements.GetElements()
	mtiReq := dataISO.Mti.String()
	// responseCode := ""
	// switch bit[3] {
	// case "310101":
	// 	responseCode = "00"
	// default:
	// 	responseCode = "13"
	// }

	//-----Memasukan pesan ke dalam format JSON dan mendapatkan respon dari app responder-----//
	responseFromResponder := kirimPesanResponder(dataISO)

	responseISO = iso8583.NewISOStruct("spec1987.yml", true)
	dataElement := map[int64]string{
		2:  responseFromResponder.Dua,
		3:  responseFromResponder.Tiga,
		7:  responseFromResponder.Tujuh,
		11: responseFromResponder.Sebelas,
		12: responseFromResponder.Duabelas,
		13: responseFromResponder.Tigabelas,
		18: responseFromResponder.Delapanbelas,
		32: responseFromResponder.Tigapuluhdua,
		33: responseFromResponder.Tigapuluhtiga,
		37: responseFromResponder.Tigapuluhtujuh,
		39: responseFromResponder.Tigapuluhsembilan,
		41: responseFromResponder.Empatpuluhsatu,
		// 42:  bit[42],
		// 43:  bit[43],
		47:  responseFromResponder.Empatpuluhtujuh,
		49:  responseFromResponder.Empatpuluhsembilan,
		62:  responseFromResponder.Enampuluhdua,
		120: responseFromResponder.Seratusduapuluh,
	}
	//-----memasukan data element ke responseISO-----//
	switch mtiReq {
	case "0200":
		responseISO.AddMTI("0210")
	case "0400":
		responseISO.AddMTI("0410")
	}
	fmt.Println("bit : ", mtiReq)

	for field, value := range dataElement {
		err := responseISO.AddField(field, value)
		if err != nil {
			log.Println("Error adding value to ISO:", err)
			return iso8583.IsoStruct{}
		}
	}
	// fmt.Println(responseISO.)

	return responseISO
}

func kirimPesanResponder(pesanIso iso8583.IsoStruct) (responseFromResponder Pesan) {
	var newPesan Pesan
	bat := pesanIso.Elements.GetElements()

	newPesan.Dua = bat[2]
	newPesan.Tiga = bat[3]
	newPesan.Tujuh = bat[7]
	newPesan.Sebelas = bat[11]
	newPesan.Duabelas = bat[12]
	newPesan.Tigabelas = bat[13]
	newPesan.Delapanbelas = bat[18]
	newPesan.Tigapuluhdua = bat[32]
	newPesan.Tigapuluhtiga = bat[33]
	newPesan.Tigapuluhtujuh = bat[37]
	// newPesan.Tigapuluhsembilan = bat[39]
	newPesan.Empatpuluhsatu = bat[41]
	newPesan.Empatpuluhtujuh = bat[47]
	newPesan.Empatpuluhsembilan = bat[49]
	newPesan.Enampuluhdua = bat[62]
	newPesan.Seratusduapuluh = bat[120]

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
	// fmt.Println("Response dari app responder :")
	// fmt.Println(newPesan.Tiga)
	// cetakDataDiTerminal(newPesan)
	defer respon.Body.Close()

	return newPesan
}
