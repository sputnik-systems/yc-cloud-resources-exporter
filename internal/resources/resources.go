package resources

import (
	"context"

	computev1 "github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	clickhousev1 "github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/clickhouse/v1"
	postgresqlv1 "github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/postgresql/v1"
	"github.com/yandex-cloud/go-sdk/gen/compute"
	"github.com/yandex-cloud/go-sdk/gen/mdb/clickhouse"
	"github.com/yandex-cloud/go-sdk/gen/mdb/postgresql"
)

func listComputeInstances(ctx context.Context, c *compute.InstanceServiceClient, folderId string) ([]*computev1.Instance, error) {
	var instances []*computev1.Instance
	req := &computev1.ListInstancesRequest{FolderId: folderId}
	resp, err := c.List(ctx, req)
	if err != nil {
		return nil, err
	}
	instances = append(instances, resp.GetInstances()...)

	for resp.GetNextPageToken() != "" {
		req.PageToken = resp.GetNextPageToken()
		resp, err = c.List(ctx, req)
		if err != nil {
			return nil, err
		}
		instances = append(instances, resp.GetInstances()...)
	}

	return instances, nil
}

func listComputeDisks(ctx context.Context, c *compute.DiskServiceClient, folderId string) ([]*computev1.Disk, error) {
	var disks []*computev1.Disk
	req := &computev1.ListDisksRequest{FolderId: folderId}
	resp, err := c.List(ctx, req)
	if err != nil {
		return nil, err
	}
	disks = append(disks, resp.GetDisks()...)

	for resp.GetNextPageToken() != "" {
		req.PageToken = resp.GetNextPageToken()
		resp, err = c.List(ctx, req)
		if err != nil {
			return nil, err
		}
		disks = append(disks, resp.GetDisks()...)
	}

	return disks, nil
}

func listManagedClickhouseClusters(ctx context.Context, c *clickhouse.ClusterServiceClient, folderId string) ([]*clickhousev1.Cluster, error) {
	var clusters []*clickhousev1.Cluster
	req := &clickhousev1.ListClustersRequest{FolderId: folderId}
	resp, err := c.List(ctx, req)
	if err != nil {
		return nil, err
	}
	clusters = append(clusters, resp.GetClusters()...)

	for resp.GetNextPageToken() != "" {
		req.PageToken = resp.GetNextPageToken()
		resp, err = c.List(ctx, req)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, resp.GetClusters()...)
	}

	return clusters, nil
}

func listManagedClickhouseHosts(ctx context.Context, c *clickhouse.ClusterServiceClient, cluster *clickhousev1.Cluster) ([]*clickhousev1.Host, error) {
	var hosts []*clickhousev1.Host
	req := &clickhousev1.ListClusterHostsRequest{ClusterId: cluster.GetId()}
	resp, err := c.ListHosts(ctx, req)
	if err != nil {
		return nil, err
	}
	hosts = append(hosts, resp.GetHosts()...)

	for resp.GetNextPageToken() != "" {
		req.PageToken = resp.GetNextPageToken()
		resp, err = c.ListHosts(ctx, req)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, resp.GetHosts()...)
	}

	return hosts, nil
}

func listManagedPostgresClusters(ctx context.Context, c *postgresql.ClusterServiceClient, folderId string) ([]*postgresqlv1.Cluster, error) {
	var clusters []*postgresqlv1.Cluster
	req := &postgresqlv1.ListClustersRequest{FolderId: folderId}
	resp, err := c.List(ctx, req)
	if err != nil {
		return nil, err
	}
	clusters = append(clusters, resp.GetClusters()...)

	for resp.GetNextPageToken() != "" {
		req.PageToken = resp.GetNextPageToken()
		resp, err = c.List(ctx, req)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, resp.GetClusters()...)
	}

	return clusters, nil
}

func listManagedPostgresHosts(ctx context.Context, c *postgresql.ClusterServiceClient, cluster *postgresqlv1.Cluster) ([]*postgresqlv1.Host, error) {
	var hosts []*postgresqlv1.Host
	req := &postgresqlv1.ListClusterHostsRequest{ClusterId: cluster.GetId()}
	resp, err := c.ListHosts(ctx, req)
	if err != nil {
		return nil, err
	}
	hosts = append(hosts, resp.GetHosts()...)

	for resp.GetNextPageToken() != "" {
		req.PageToken = resp.GetNextPageToken()
		resp, err = c.ListHosts(ctx, req)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, resp.GetHosts()...)
	}

	return hosts, nil
}
