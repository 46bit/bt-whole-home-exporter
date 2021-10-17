package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	password := os.Args[1]

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		topologyInfo, err := fetchTopologyInfoFromBtWholeHomeWifi(password)
		if err != nil {
			panic(err)
		}
		nodesByDeviceID := mapByNodeDeviceID(topologyInfo.Nodes)

		output := &bytes.Buffer{}
		renderMetricsInPromFormat(nodeMetrics(topologyInfo.Nodes), output)
		renderMetricsInPromFormat(stationMetrics(topologyInfo.Stations, nodesByDeviceID), output)
		outputBytes := output.Bytes()
		fmt.Printf("%s", string(outputBytes))
		w.Write(outputBytes)
	})
	log.Fatal(http.ListenAndServe(":8994", mux))
}
