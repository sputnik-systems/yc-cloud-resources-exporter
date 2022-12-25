package resources

import (
	"github.com/prometheus/client_golang/prometheus"
)

func getComputeInstanceInfoGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_instance_info",
		Help: "The total size of requested disks",
	}, []string{"folder_id", "id", "name", "status", "platform_id", "core_fraction", "preemptible"})
}

func getComputeInstanceCoresGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_instance_cores",
		Help: "The total size of requested disks",
	}, []string{"folder_id", "id", "name", "status", "platform_id", "core_fraction", "preemptible"})
}

func getComputeInstanceMemoryGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_instance_memory",
		Help: "The total size of requested disks",
	}, []string{"folder_id", "id", "name", "status", "platform_id", "preemptible"})
}

func getComputeDiskSizeGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_disk_size",
		Help: "The total size of requested disks",
	}, []string{"folder_id", "id", "name", "type_id"})
}

func getClickhouseClusterInfoGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_clickhouse_clusterInfo",
		Help: "The used clickhouse clusters",
	}, []string{"folder_id", "id", "name", "status"})
}

func getClickhouseHostDiskSizeGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_clickhouse_host_disk_size",
		Help: "The used clickhouse cluster hosts",
	}, []string{"name", "type", "resource_preset_id", "disk_type_id", "cluster_id"})
}

func getPostgresClusterInfoGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_postgresql_cluster_info",
		Help: "The used postgresql clusters",
	}, []string{"folder_id", "id", "name", "status"})
}

func getPostgresHostDiskSizeGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_postgresql_host_disk_size",
		Help: "The used postgresql cluster hosts",
	}, []string{"name", "resource_preset_id", "disk_type_id", "cluster_id"})
}
