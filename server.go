package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net"
	"time"
)

type IpInfo struct {
	Ip				string		`json:"ip"`
	Title			string		`json:"info"`
}

type Item struct {
	ID     	 		uint   		`json:"id"`
	IpAdrClient   	string 		`json:"ipClient"`
	IpInfo			IpInfo		`json:"ipInfo"`
}

type udpSrv struct {
	ipPort			string
	udpAddr			*net.UDPAddr
	listener		*net.UDPConn
	err				error
}

func main() {

	//подключение к дб
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	collection := client.Database("test").Collection("test")

	//Создаем коннект
	var udpSrv udpSrv
	var item  Item

	udpSrv.ipPort = "127.0.0.1:6000"

	udpSrv.udpAddr, udpSrv.err = net.ResolveUDPAddr("udp4", udpSrv.ipPort)
	if udpSrv.err  != nil {
		log.Fatal(udpSrv.err )
	}
	//Слушаем UDP
	udpSrv.listener, udpSrv.err = net.ListenUDP("udp4", udpSrv.udpAddr)
	if udpSrv.err != nil {
		log.Fatal(udpSrv.err)
	}

	fmt.Println("UDP server up and listening on port 6000")

	//defer udpSrv.listener.Close()
	//Обработка коннекта
	for {
		// wait for UDP client to connect
		handleUDPConnection(udpSrv.listener, &item)

		// Добавление данных в бд
		_, err := collection.InsertOne(ctx,item)
		if err != nil {
			panic(err)
		}
	}
}
//====================================================//
//Decompress data
func zlUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = zlib.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

//HandlerUdpConnection
func handleUDPConnection(conn *net.UDPConn, item *Item) {

	buffer := make([]byte, 1024)
	//Запись данных в буфер
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}

	//	fmt.Println("UDP client : ", addr)
	//	fmt.Println("Received from UDP client :  ", buffer[:n])

	item.IpAdrClient = addr.String()


	//Decompress
	uncompressedData, uncompressedDataErr := zlUnzipData(buffer[:n])
	if uncompressedDataErr != nil {
		log.Fatal(uncompressedDataErr)
	}
	//	os.Stdout.Write(uncompressedData) //Вывести в консоль данные в формате JSON после Decompress

	//Unmarshal JSON и запись в структуру
	_ = json.Unmarshal(uncompressedData, item)
	log.Println(*item)
}


