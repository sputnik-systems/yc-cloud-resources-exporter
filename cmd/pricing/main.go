package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cpuCost := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_instance_cpu_cost_per_hour",
		Help: "The cpu cost of compute instance per hour",
	}, []string{"platform_id", "core_fraction", "preemptible"})
	memoryCost := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_instance_memory_cost_per_hour",
		Help: "The memory cost of compute instance per hour",
	}, []string{"platform_id", "preemptible"})
	diskCost := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_disk_cost_per_month",
		Help: "The disk cost of compute instance per month",
	}, []string{"type_id"})
	// standard-v1
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v1", "core_fraction": "5", "preemptible": "false"}).Set(0.31)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v1", "core_fraction": "5", "preemptible": "true"}).Set(0.19)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v1", "core_fraction": "20", "preemptible": "false"}).Set(0.88)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v1", "core_fraction": "20", "preemptible": "true"}).Set(0.27)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v1", "core_fraction": "100", "preemptible": "false"}).Set(1.12)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v1", "core_fraction": "100", "preemptible": "true"}).Set(0.34)
	memoryCost.With(prometheus.Labels{"platform_id": "standard-v1", "preemptible": "false"}).Set(0.39)
	memoryCost.With(prometheus.Labels{"platform_id": "standard-v1", "preemptible": "true"}).Set(0.12)
	// standard-v2
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "5", "preemptible": "false"}).Set(0.16)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "5", "preemptible": "true"}).Set(0.01)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "20", "preemptible": "false"}).Set(0.49)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "20", "preemptible": "true"}).Set(0.16)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "50", "preemptible": "false"}).Set(0.72)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "50", "preemptible": "true"}).Set(0.22)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "100", "preemptible": "false"}).Set(1.19)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v2", "core_fraction": "100", "preemptible": "true"}).Set(0.32)
	memoryCost.With(prometheus.Labels{"platform_id": "standard-v2", "preemptible": "false"}).Set(0.31)
	memoryCost.With(prometheus.Labels{"platform_id": "standard-v2", "preemptible": "true"}).Set(0.07)
	// standard-v3
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v3", "core_fraction": "20", "preemptible": "false"}).Set(0.44)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v3", "core_fraction": "20", "preemptible": "true"}).Set(0.14)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v3", "core_fraction": "50", "preemptible": "false"}).Set(0.64)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v3", "core_fraction": "50", "preemptible": "true"}).Set(0.2)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v3", "core_fraction": "100", "preemptible": "false"}).Set(1.05)
	cpuCost.With(prometheus.Labels{"platform_id": "standard-v3", "core_fraction": "100", "preemptible": "true"}).Set(0.29)
	memoryCost.With(prometheus.Labels{"platform_id": "standard-v3", "preemptible": "false"}).Set(0.28)
	memoryCost.With(prometheus.Labels{"platform_id": "standard-v3", "preemptible": "true"}).Set(0.07)
	// disks
	diskCost.With(prometheus.Labels{"type_id": "network-ssd"}).Set(11.91)
	diskCost.With(prometheus.Labels{"type_id": "network-hdd"}).Set(2.92)
	diskCost.With(prometheus.Labels{"type_id": "non-replicated-ssd"}).Set(8.8)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
