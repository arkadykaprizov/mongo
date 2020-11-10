package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"log"
	"net"
	"os"
)

type Item struct {
	ID     	 		uint   		`json:"id"`
	IpInfo			IpInfo		`json:"ipInfo"`
}
type IpInfo struct {
	Ip				string		`json:"ip"`
	Title			string		`json:"info"`
}

func main() {
	//Проверка входхого аргумента
	if len(os.Args) != 2 {
		log.Println("Usage: ", os.Args[0], "host:port")
		os.Exit(1)
	}

	//Создание структуры для передачи
	item := Item{
		ID: 1,
		IpInfo: IpInfo{
			Ip: "192.168.3.12/24",
			Title: "use",
		},

	}

	//Преобразование структуры в JSON
	b, err := json.Marshal(item)
	checkError(err)

	//Compress JSON
	compressedData, compressedDataErr := zlibData(b)
	if compressedDataErr != nil {
		log.Fatal(compressedDataErr)
	}

	//	log.Println("compressed data:", compressedData)
	//	log.Println("compressed data len:", len(compressedData))

	//Создние Коннекта
	service := os.Args[1]

	conn, err := net.Dial("udp", service)
	checkError(err)

	//Отправка по UDP
	/*for i :=0; i<2000; i++{
		_, err = conn.Write(compressedData)
	}*/
	_, err = conn.Write(compressedData)
	checkError(err)


	os.Exit(0)
}

//========================================================//
//HandlerErr
func checkError(err error) {
	if err != nil {
		log.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

//Compress Data
func zlibData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	zl := zlib.NewWriter(&b)

	_, err = zl.Write(data)
	if err != nil {
		return
	}

	if err = zl.Flush(); err != nil {
		return
	}

	if err = zl.Close(); err != nil {
		return
	}

	compressedData = b.Bytes()

	return
}


