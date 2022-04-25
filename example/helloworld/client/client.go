package main

import (
	"log"
	"net"

	"github.com/muque7/my-rpc/client"
)

func main() {

	addr := "0.0.0.0:2333"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("dial error: %v\n", err)
	}
	cli := client.NewClient(conn)

	var correctQuery func(int) (map[string]interface{}, error)
	var wrongQuery func(int) (map[string]interface{}, error)

	cli.Call("myService", "QueryUser", &correctQuery)
	u, err := correctQuery(1)
	if err != nil {
		log.Printf("query1 error: %v\n", err)
	} else {
		log.Printf("query1 result: %v \n", u)
	}
	u, err = correctQuery(2)
	if err != nil {
		log.Printf("query2 error: %v\n", err)
	} else {
		log.Printf("query2 result: %v \n", u["Name"])
	}

	cli.Call("myService", "QueryUser", &wrongQuery)
	u, err = wrongQuery(1)
	if err != nil {
		log.Printf("query3 error: %v\n", err)
	} else {
		log.Println(u)
	}

	conn.Close()
}
