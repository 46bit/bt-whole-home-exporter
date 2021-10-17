package main

import (
	"io"
	"strconv"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type Metrics = map[string]*dto.MetricFamily

func stationMetrics(stations []Station, nodesByDeviceID map[string]Node) Metrics {
	metrics := Metrics{
		"wifi_device_signal_strength_db": &dto.MetricFamily{
			Name:   derefS("wifi_device_signal_strength_db"),
			Type:   derefT(dto.MetricType_GAUGE),
			Metric: []*dto.Metric{},
		},
	}

	for _, station := range stations {
		if !station.IsOnline() {
			continue
		}

		signalStrength, err := strconv.ParseFloat(station.SignalStrength, 64)
		if err != nil {
			continue
		}

		metric := &dto.Metric{
			Gauge: &dto.Gauge{Value: &signalStrength},
			Label: []*dto.LabelPair{
				{
					Name:  derefS("name"),
					Value: derefS(station.BestName()),
				},
				{
					Name:  derefS("mac"),
					Value: derefS(station.StationMac),
				},
				{
					Name:  derefS("ip"),
					Value: derefS(station.StationIP),
				},
				{
					Name:  derefS("freq"),
					Value: derefS(station.ConnectType),
				},
				{
					Name:  derefS("link_rate"),
					Value: derefS(station.LinkRate),
				},
				{
					Name:  derefS("ap_name"),
					Value: derefS(station.AccessPointName(nodesByDeviceID)),
				},
				{
					Name:  derefS("ap_id"),
					Value: derefS(station.ParentID),
				},
			},
		}

		metrics["wifi_device_signal_strength_db"].Metric = append(
			metrics["wifi_device_signal_strength_db"].Metric,
			metric,
		)
	}

	return metrics
}

func nodeMetrics(nodes []Node) Metrics {
	metrics := Metrics{
		"wifi_ap_cpu_ratio": &dto.MetricFamily{
			Name:   derefS("wifi_ap_cpu_ratio"),
			Type:   derefT(dto.MetricType_GAUGE),
			Metric: []*dto.Metric{},
		},
		"wifi_ap_mem_ratio": &dto.MetricFamily{
			Name:   derefS("wifi_ap_mem_ratio"),
			Type:   derefT(dto.MetricType_GAUGE),
			Metric: []*dto.Metric{},
		},
	}

	for _, node := range nodes {
		cpuUser, err := strconv.ParseFloat(node.CPUUser, 64)
		if err != nil {
			continue
		}
		cpuSystem, err := strconv.ParseFloat(node.CPUSystem, 64)
		if err != nil {
			continue
		}
		cpuRatio := (cpuUser + cpuSystem) / 100
		cpuRatioMetric := &dto.Metric{
			Gauge: &dto.Gauge{Value: &cpuRatio},
			Label: []*dto.LabelPair{
				{
					Name:  derefS("name"),
					Value: derefS(node.DeviceName),
				},
				{
					Name:  derefS("id"),
					Value: derefS(node.DeviceID),
				},
				{
					Name:  derefS("ip"),
					Value: derefS(node.DeviceIP),
				},
			},
		}

		memUsed, err := strconv.ParseFloat(node.MemUsed, 64)
		if err != nil {
			continue
		}
		memTotal, err := strconv.ParseFloat(node.MemTotal, 64)
		if err != nil {
			continue
		}
		memRatio := memUsed / memTotal
		memRatioMetric := &dto.Metric{
			Gauge: &dto.Gauge{Value: &memRatio},
			Label: []*dto.LabelPair{
				{
					Name:  derefS("name"),
					Value: derefS(node.DeviceName),
				},
				{
					Name:  derefS("id"),
					Value: derefS(node.DeviceID),
				},
				{
					Name:  derefS("ip"),
					Value: derefS(node.DeviceIP),
				},
			},
		}

		metrics["wifi_ap_cpu_ratio"].Metric = append(
			metrics["wifi_ap_cpu_ratio"].Metric,
			cpuRatioMetric,
		)
		metrics["wifi_ap_mem_ratio"].Metric = append(
			metrics["wifi_ap_mem_ratio"].Metric,
			memRatioMetric,
		)
	}

	return metrics
}

func renderMetricsInPromFormat(metrics Metrics, out io.Writer) int {
	totalBytesWritten := 0
	for _, metricFamily := range metrics {
		bytesWritten, err := expfmt.MetricFamilyToText(out, metricFamily)
		totalBytesWritten += bytesWritten
		if err != nil {
			continue
		}
	}
	return totalBytesWritten
}

func derefS(s string) *string {
	return &s
}

func derefT(i dto.MetricType) *dto.MetricType {
	return &i
}
