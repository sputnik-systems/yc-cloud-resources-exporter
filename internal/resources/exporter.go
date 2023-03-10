package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	computev1 "github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	clickhousev1 "github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/clickhouse/v1"
	postgresqlv1 "github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/postgresql/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

var (
	instanceInfoMetric                                    *prometheus.GaugeVec
	coresMetric, memoryMetric                             *prometheus.GaugeVec
	diskSizeMetric                                        *prometheus.GaugeVec
	clickhouseClusterMetric, clickhouseHostDiskSizeMetric *prometheus.GaugeVec
	postgresClusterMetric, postgresHostDiskSizeMetric     *prometheus.GaugeVec

	wg sync.WaitGroup
)

func Run() error {
	var creds ycsdk.Credentials
	switch {
	case os.Getenv("YC_IAM_TOKEN") != "":
		creds = ycsdk.NewIAMTokenCredentials(os.Getenv("YC_IAM_TOKEN"))
	case os.Getenv("YC_OAUTH_TOKEN") != "":
		creds = ycsdk.OAuthToken(os.Getenv("YC_OAUTH_TOKEN"))
	case os.Getenv("YC_SA_KEY") != "":
		key, err := iamkey.ReadFromJSONBytes([]byte(os.Getenv("YC_SA_KEY")))
		if err != nil {
			return fmt.Errorf("failed to read sa key: %w", err)
		}
		creds, err = ycsdk.ServiceAccountKey(key)
		if err != nil {
			return fmt.Errorf("failed to get sa key based creds: %w", err)
		}
	case os.Getenv("YC_SA_KEY_FILE") != "":
		key, err := iamkey.ReadFromJSONFile(os.Getenv("YC_SA_KEY_FILE"))
		if err != nil {
			return fmt.Errorf("failed to read sa key file: %w", err)
		}
		creds, err = ycsdk.ServiceAccountKey(key)
		if err != nil {
			return fmt.Errorf("failed to get sa key based creds: %w", err)
		}
	default:
		creds = ycsdk.InstanceServiceAccount()
	}

	ctx := context.Background()
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{Credentials: creds})
	if err != nil {
		return fmt.Errorf("failed to init yc sdk: %w", err)
	}

	folderIds, err := getFolderIds()
	if err != nil {
		return err
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		wg.Add(4)

		updateComputeInstanceMetrics(ctx, sdk, folderIds)
		updateComputeDiskMetrics(ctx, sdk, folderIds)
		updateManagedClickhouseMetrics(ctx, sdk, folderIds)
		updateManagedPostgresMetrics(ctx, sdk, folderIds)

		wg.Wait()

		handler := promhttp.HandlerFor(
			prometheus.DefaultGatherer, promhttp.HandlerOpts{})
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {})
	return http.ListenAndServe(":8080", nil)
}

func getFolderIds() ([]string, error) {
	v := os.Getenv("YC_FOLDER_IDS")
	if v == "" {
		return nil, errors.New("you should specify YC_FOLDER_IDS env var")
	}

	in := strings.Split(v, ",")
	check := make(map[string]struct{})
	out := make([]string, 0)
	for _, id := range in {
		if _, ok := check[id]; !ok {
			out = append(out, id)
		}
	}

	return out, nil
}

func updateComputeInstanceMetrics(ctx context.Context, sdk *ycsdk.SDK, folderIds []string) {
	instances := make([]*computev1.Instance, 0)
	for _, folderId := range folderIds {
		result, err := listComputeInstances(ctx, sdk.Compute().Instance(), folderId)
		if err != nil {
			log.Printf("failed to list compute instances: %s", err)
		}

		instances = append(instances, result...)
	}

	if instanceInfoMetric != nil {
		prometheus.Unregister(instanceInfoMetric)
	}
	if coresMetric != nil {
		prometheus.Unregister(coresMetric)
	}
	if memoryMetric != nil {
		prometheus.Unregister(memoryMetric)
	}
	instanceInfoMetric = getComputeInstanceInfoGaugeVec()
	coresMetric = getComputeInstanceCoresGaugeVec()
	memoryMetric = getComputeInstanceMemoryGaugeVec()
	prometheus.MustRegister(instanceInfoMetric)
	prometheus.MustRegister(coresMetric)
	prometheus.MustRegister(memoryMetric)

	for _, instance := range instances {
		folderId := instance.GetFolderId()
		id := instance.GetId()
		name := instance.GetName()
		status := strings.ToLower(instance.GetStatus().String())
		platformId := instance.GetPlatformId()
		resources := instance.GetResources()
		coreFraction := strconv.FormatInt(int64(resources.GetCoreFraction()), 10)
		preemptible := strconv.FormatBool(instance.GetSchedulingPolicy().GetPreemptible())
		instanceInfoMetric.With(prometheus.Labels{
			"folder_id":     folderId,
			"id":            id,
			"name":          name,
			"status":        status,
			"platform_id":   platformId,
			"core_fraction": coreFraction,
			"preemptible":   preemptible,
		}).Set(1)
		coresMetric.With(prometheus.Labels{
			"folder_id":     folderId,
			"id":            id,
			"name":          name,
			"status":        status,
			"platform_id":   platformId,
			"core_fraction": coreFraction,
			"preemptible":   preemptible,
		}).Set(float64(resources.GetCores()))
		memoryMetric.With(prometheus.Labels{
			"folder_id":   folderId,
			"id":          id,
			"name":        name,
			"status":      status,
			"platform_id": platformId,
			"preemptible": preemptible,
		}).Set(float64(resources.GetMemory()))
	}

	wg.Done()
}

func updateComputeDiskMetrics(ctx context.Context, sdk *ycsdk.SDK, folderIds []string) {
	disks := make([]*computev1.Disk, 0)
	for _, folderId := range folderIds {
		result, err := listComputeDisks(ctx, sdk.Compute().Disk(), folderId)
		if err != nil {
			log.Printf("failed to list compute disks: %s", err)
		}

		disks = append(disks, result...)
	}

	if diskSizeMetric != nil {
		prometheus.Unregister(diskSizeMetric)
	}
	diskSizeMetric = getComputeDiskSizeGaugeVec()
	prometheus.MustRegister(diskSizeMetric)

	for _, disk := range disks {
		for _, instanceId := range disk.GetInstanceIds() {
			diskSizeMetric.With(prometheus.Labels{
				"folder_id":   disk.GetFolderId(),
				"id":          disk.GetId(),
				"instance_id": instanceId,
				"name":        disk.GetName(),
				"type_id":     disk.GetTypeId(),
			}).Set(float64(disk.GetSize()))
		}
	}

	wg.Done()
}

func updateManagedClickhouseMetrics(ctx context.Context, sdk *ycsdk.SDK, folderIds []string) {
	clusters := make([]*clickhousev1.Cluster, 0)
	for _, folderId := range folderIds {
		result, err := listManagedClickhouseClusters(ctx, sdk.MDB().Clickhouse().Cluster(), folderId)
		if err != nil {
			log.Printf("failed to list clickhouse clusters: %s", err)
		}

		clusters = append(clusters, result...)
	}

	if clickhouseClusterMetric != nil {
		prometheus.Unregister(clickhouseClusterMetric)
	}
	clickhouseClusterMetric = getClickhouseClusterInfoGaugeVec()
	prometheus.MustRegister(clickhouseClusterMetric)

	hosts := make([]*clickhousev1.Host, 0)
	for _, cluster := range clusters {
		clickhouseClusterMetric.With(prometheus.Labels{
			"folder_id": cluster.GetFolderId(),
			"id":        cluster.GetId(),
			"name":      cluster.GetName(),
			"status":    strings.ToLower(cluster.GetStatus().String()),
		}).Set(1)

		chosts, err := listManagedClickhouseHosts(ctx, sdk.MDB().Clickhouse().Cluster(), cluster)
		if err != nil {
			log.Printf("failed to list clickhouse cluster hosts: %s", err)
		}
		hosts = append(hosts, chosts...)
	}

	if clickhouseHostDiskSizeMetric != nil {
		prometheus.Unregister(clickhouseHostDiskSizeMetric)
	}
	clickhouseHostDiskSizeMetric = getClickhouseHostDiskSizeGaugeVec()
	prometheus.MustRegister(clickhouseHostDiskSizeMetric)
	for _, host := range hosts {
		resources := host.GetResources()
		clickhouseHostDiskSizeMetric.With(prometheus.Labels{
			"name":               host.GetName(),
			"type":               strings.ToLower(host.GetType().String()),
			"resource_preset_id": resources.GetResourcePresetId(),
			"disk_type_id":       resources.GetDiskTypeId(),
			"cluster_id":         strings.ToLower(host.GetClusterId()),
		}).Set(float64(resources.GetDiskSize()))
	}

	wg.Done()
}

func updateManagedPostgresMetrics(ctx context.Context, sdk *ycsdk.SDK, folderIds []string) {
	clusters := make([]*postgresqlv1.Cluster, 0)
	for _, folderId := range folderIds {
		result, err := listManagedPostgresClusters(ctx, sdk.MDB().PostgreSQL().Cluster(), folderId)
		if err != nil {
			log.Printf("failed to list postgres clusters: %s", err)
		}

		clusters = append(clusters, result...)
	}

	if postgresClusterMetric != nil {
		prometheus.Unregister(postgresClusterMetric)
	}
	postgresClusterMetric = getPostgresClusterInfoGaugeVec()
	prometheus.MustRegister(postgresClusterMetric)

	hosts := make([]*postgresqlv1.Host, 0)
	for _, cluster := range clusters {
		postgresClusterMetric.With(prometheus.Labels{
			"folder_id": cluster.GetFolderId(),
			"id":        cluster.GetId(),
			"name":      cluster.GetName(),
			"status":    strings.ToLower(cluster.GetStatus().String()),
		}).Set(1)

		chosts, err := listManagedPostgresHosts(ctx, sdk.MDB().PostgreSQL().Cluster(), cluster)
		if err != nil {
			log.Printf("failed to list postgres cluster hosts: %s", err)
		}
		hosts = append(hosts, chosts...)
	}

	if postgresHostDiskSizeMetric != nil {
		prometheus.Unregister(postgresHostDiskSizeMetric)
	}
	postgresHostDiskSizeMetric = getPostgresHostDiskSizeGaugeVec()
	prometheus.MustRegister(postgresHostDiskSizeMetric)
	for _, host := range hosts {
		resources := host.GetResources()
		postgresHostDiskSizeMetric.With(prometheus.Labels{
			"name":               host.GetName(),
			"resource_preset_id": resources.GetResourcePresetId(),
			"disk_type_id":       resources.GetDiskTypeId(),
			"cluster_id":         strings.ToLower(host.GetClusterId()),
		}).Set(float64(resources.GetDiskSize()))
	}

	wg.Done()
}
