package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
)

const Http = "http"
const Tcp = "tcp"
const Udp = "udp"

var waitGroup = sync.WaitGroup{}

type httpWay string
type tcpWay string
type udpWay string

type netAccess interface {
	access()
}

/*HTTP*/
func (httpDo httpWay) access() {
	resp, _ := doGet("http://127.0.0.1:8000")
	defer resp.Body.Close()
	defer waitGroup.Done()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			log.Panicln(err)
		}
		fmt.Println("http reply", string(body))
	}
}
func doGet(url string) (r *http.Response, e error) {
	resp, err := http.Get(url)
	if err != nil {
		//log.Panicln(resp.StatusCode)
		log.Panicln(err)
	}
	return resp, err
}

/*TCP*/
func (tcpDo tcpWay) access() {
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()
	defer waitGroup.Done()
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	fmt.Println("tcp reply:", string(buffer))
}

/*UDP*/
func (udpDo udpWay) access() {
	conn, err := net.Dial("udp", "127.0.0.1:8002")
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()
	defer waitGroup.Done()

	_, err = conn.Write([]byte(""))
	if err != nil {
		log.Panicln(err)
	}

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	fmt.Println("udp reply", string(buffer))
}

func main() {
	var httpDo httpWay = Http
	var tcpDo tcpWay = Tcp
	var udpDo udpWay = Udp
	accessWay := [...]netAccess{httpDo, tcpDo, udpDo}
	for _, Type := range accessWay {
		go Type.access()
	}
	waitGroup.Add(3)
	waitGroup.Wait()
}
