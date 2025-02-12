package metrics

import (
	"runtime"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/config"
)

var registerClusterManagerBaseMetrics sync.Once

// MetricMasterLeader identifies whether this instance of ovnkube-master is a leader or not
var MetricClusterManagerLeader = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: MetricOvnkubeNamespace,
	Subsystem: MetricOvnkubeSubsystemClusterManager,
	Name:      "leader",
	Help:      "Identifies whether the instance of ovnkube-cluster-manager is a leader(1) or not(0).",
})

var MetricClusterManagerReadyDuration = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: MetricOvnkubeNamespace,
	Subsystem: MetricOvnkubeSubsystemClusterManager,
	Name:      "ready_duration_seconds",
	Help:      "The duration for the cluster manager to get to ready state",
})

var metricV4HostSubnetCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: MetricOvnkubeNamespace,
	Subsystem: MetricOvnkubeSubsystemClusterManager,
	Name:      "num_v4_host_subnets",
	Help:      "The total number of v4 host subnets possible",
})

var metricV6HostSubnetCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: MetricOvnkubeNamespace,
	Subsystem: MetricOvnkubeSubsystemClusterManager,
	Name:      "num_v6_host_subnets",
	Help:      "The total number of v6 host subnets possible",
})

var metricV4AllocatedHostSubnetCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: MetricOvnkubeNamespace,
	Subsystem: MetricOvnkubeSubsystemClusterManager,
	Name:      "allocated_v4_host_subnets",
	Help:      "The total number of v4 host subnets currently allocated",
})

var metricV6AllocatedHostSubnetCount = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: MetricOvnkubeNamespace,
	Subsystem: MetricOvnkubeSubsystemClusterManager,
	Name:      "allocated_v6_host_subnets",
	Help:      "The total number of v6 host subnets currently allocated",
})

// RegisterClusterManagerBase registers ovnkube cluster manager base metrics with the Prometheus registry.
// This function should only be called once.
func RegisterClusterManagerBase() {
	registerClusterManagerBaseMetrics.Do(func() {
		prometheus.MustRegister(MetricClusterManagerLeader)
		prometheus.MustRegister(MetricClusterManagerReadyDuration)
		prometheus.MustRegister(prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Namespace: MetricOvnkubeNamespace,
				Subsystem: MetricOvnkubeSubsystemClusterManager,
				Name:      "build_info",
				Help: "A metric with a constant '1' value labeled by version, revision, branch, " +
					"and go version from which ovnkube was built and when and who built it",
				ConstLabels: prometheus.Labels{
					"version":    "0.0",
					"revision":   config.Commit,
					"branch":     config.Branch,
					"build_user": config.BuildUser,
					"build_date": config.BuildDate,
					"goversion":  runtime.Version(),
				},
			},
			func() float64 { return 1 },
		))
	})
}

// RegisterClusterManagerFunctional is a collection of metrics that help us understand ovnkube-cluster-manager functions. Call once after
// LE is won.
func RegisterClusterManagerFunctional() {
	prometheus.MustRegister(metricV4HostSubnetCount)
	prometheus.MustRegister(metricV6HostSubnetCount)
	prometheus.MustRegister(metricV4AllocatedHostSubnetCount)
	prometheus.MustRegister(metricV6AllocatedHostSubnetCount)
}

func UnregisterClusterManagerFunctional() {
	prometheus.Unregister(metricV4HostSubnetCount)
	prometheus.Unregister(metricV6HostSubnetCount)
	prometheus.Unregister(metricV4AllocatedHostSubnetCount)
	prometheus.Unregister(metricV6AllocatedHostSubnetCount)
}

// RecordSubnetUsage records the number of subnets allocated for nodes
func RecordSubnetUsage(v4SubnetsAllocated, v6SubnetsAllocated float64) {
	metricV4AllocatedHostSubnetCount.Set(v4SubnetsAllocated)
	metricV6AllocatedHostSubnetCount.Set(v6SubnetsAllocated)
}

// RecordSubnetCount records the number of available subnets per configuration
// for ovn-kubernetes
func RecordSubnetCount(v4SubnetCount, v6SubnetCount float64) {
	metricV4HostSubnetCount.Set(v4SubnetCount)
	metricV6HostSubnetCount.Set(v6SubnetCount)
}
