package resources

import (
	"github.com/prometheus/client_golang/prometheus"
)

func getComputeInstanceCoresGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_instance_cores",
		Help: "The total size of requested disks",
	}, []string{"folder_id", "id", "name", "status", "platform_id", "core_fraction"})
}

func getComputeInstanceMemoryGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_isntance_memory",
		Help: "The total size of requested disks",
	}, []string{"folder_id", "id", "name", "status", "platform_id"})
}

func getComputeDiskSizeGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_compute_disk_size",
		Help: "The total size of requested disks",
	}, []string{"folder_id", "id", "name", "type_id"})
}

func getClickhouseClusterGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_clickhouse_cluster",
		Help: "The used clickhouse clusters",
	}, []string{"folder_id", "id", "name", "status"})
}

func getClickhouseHostDiskSizeGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_clickhouse_host_disk_size",
		Help: "The used clickhouse cluster hosts",
	}, []string{"folder_id", "name", "type", "resource_preset_id", "disk_type_id", "cluster_id"})
}

func getPostgresClusterGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_postgresql_cluster",
		Help: "The used postgresql clusters",
	}, []string{"folder_id", "id", "name", "status"})
}

func getPostgresHostDiskSizeGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "yc_mdb_postgresql_host_disk_size",
		Help: "The used postgresql cluster hosts",
	}, []string{"folder_id", "name", "resource_preset_id", "disk_type_id", "cluster_id"})
}
