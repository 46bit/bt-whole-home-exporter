package main

import (
	"encoding/json"
	"os/exec"
)

type Station struct {
	StationMac     string `json:"station_mac"`
	StationName    string `json:"station_name"`
	AliasName      string `json:"alias_name"`
	StationIP      string `json:"station_ip"`
	ParentID       string `json:"parent_id"`
	ConnectType    string `json:"connect_type"`
	LinkRate       string `json:"link_rate"`
	SignalStrength string `json:"signal_strength"`
	Online         string `json:"online"`
}

func (s *Station) IsOnline() bool {
	return s.Online == "1"
}

func (s *Station) BestName() string {
	if s.AliasName == "NULL" || len(s.AliasName) == 0 {
		return s.StationName
	}
	return s.AliasName
}

func (s *Station) AccessPointName(nodesByDeviceID map[string]Node) string {
	node, ok := nodesByDeviceID[s.ParentID]
	if ok {
		return node.DeviceName
	}
	return ""
}

type Node struct {
	DeviceName string `json:"device_name"`
	DeviceID   string `json:"device_id"`
	DeviceIP   string `json:"device_ip"`
	CPUUser    string `json:"cpuU"`
	CPUSystem  string `json:"cpuS"`
	CPUIdle    string `json:"cpuI"`
	MemTotal   string `json:"memT"`
	MemFree    string `json:"memF"`
	MemUsed    string `json:"memU"`
}

type TopologyInfo struct {
	Nodes    []Node    `json:"nodes"`
	Stations []Station `json:"stations"`
}

func fetchTopologyInfoFromBtWholeHomeWifi(password string) (*TopologyInfo, error) {
	cmd := exec.Command("bash", "./fetch-toplogy-info-from-bt-whole-home-wifi.bash", password)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var topologyInfo TopologyInfo
	if err := json.Unmarshal([]byte(stdout), &topologyInfo); err != nil {
		return nil, err
	}

	return &topologyInfo, nil
}

func mapByNodeDeviceID(nodes []Node) map[string]Node {
	nodesByDeviceID := map[string]Node{}
	for _, node := range nodes {
		nodesByDeviceID[node.DeviceID] = node
	}
	return nodesByDeviceID
}
