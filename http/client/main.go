package main

import (
	"log"
	"net/http"
	"time"

	_ "net/http/pprof"
)

func main() {

	go func() {
		// serve pprof
		if err := http.ListenAndServe(":8001", nil); err != nil {
			panic(err)
		}
	}()

	go func() {
		client := &http.Client{}
		for {
			doWork(client)
		}
	}()

	//time.Sleep(100 * time.Hour)

	client := &http.Client{}
	for {
		resp, err := client.Get("http://localhost:8000")
		if err != nil {
			log.Println("Error:", err)
		} else {
			log.Println("Response status:", resp.Status)
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
}

//go:noinline
func doWork(client *http.Client) {
	actuallyDoWork(client)
}

//go:noinline
func actuallyDoWork(client *http.Client) {
	//resp, err := client.Get("http://localhost:8000")
	//if err != nil {
	//	fmt.Println("Error:", err)
	//} else {
	//	fmt.Println("Response status:", resp.Status)
	//	resp.Body.Close()
	//}
	time.Sleep(1 * time.Second)
}
