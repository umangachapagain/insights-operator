package clusterconfig

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"github.com/openshift/insights-operator/pkg/gatherers/common"
	"github.com/openshift/insights-operator/pkg/record"
)

// GatherOpenShiftAPIServerOperatorLogs Collects logs from `openshift-apiserver-operator` with following substrings:
// - "the server has received too many requests and has asked us"
// - "because serving request timed out and response had been started"
//
// ### API Reference
// - https://github.com/kubernetes/client-go/blob/master/kubernetes/typed/core/v1/pod_expansion.go#L48
// - https://docs.openshift.com/container-platform/4.6/rest_api/workloads_apis/pod-core-v1.html#apiv1namespacesnamespacepodsnamelog
//
// ### Sample data
// - docs/insights-archive-sample/config/pod/openshift-apiserver-operator/logs/openshift-apiserver-operator-6ddb679b87-4kn55/errors.log
//
// ### Location in archive
// - `config/pod/{namespace-name}/logs/{pod-name}/errors.log`
//
// ### Config ID
// `clusterconfig/openshift_apiserver_operator_logs`
//
// ### Released version
// - 4.7.0
//
// ### Backported versions
// None
//
// ### Changes
// None
func (g *Gatherer) GatherOpenShiftAPIServerOperatorLogs(ctx context.Context) ([]record.Record, []error) {
	containersFilter := common.LogContainersFilter{
		Namespace:     "openshift-apiserver-operator",
		LabelSelector: "app=openshift-apiserver-operator",
	}

	gatherKubeClient, err := kubernetes.NewForConfig(g.gatherProtoKubeConfig)
	if err != nil {
		return nil, []error{err}
	}

	coreClient := gatherKubeClient.CoreV1()

	records, err := common.CollectLogsFromContainers(
		ctx,
		coreClient,
		containersFilter,
		getAPIServerOperatorLogsMessagesFilter(),
		nil,
	)
	if err != nil {
		return nil, []error{err}
	}

	return records, nil
}

func getAPIServerOperatorLogsMessagesFilter() common.LogMessagesFilter {
	return common.LogMessagesFilter{
		MessagesToSearch: []string{
			"the server has received too many requests and has asked us",
			"because serving request timed out and response had been started",
		},
		IsRegexSearch: false,
		SinceSeconds:  logDefaultSinceSeconds,
		LimitBytes:    logDefaultLimitBytes,
	}
}
