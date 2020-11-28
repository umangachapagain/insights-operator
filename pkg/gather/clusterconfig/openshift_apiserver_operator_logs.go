package clusterconfig

import (
	"bufio"
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog"
	"strings"

	"github.com/openshift/insights-operator/pkg/record"
)

// GatherOpenShiftAPIServerOperatorLogs collects logs from openshift-apiserver-operator with following substrings:
//   - "the server has received too many requests and has asked us"
//   - "because serving request timed out and response had been started"
//
// The Kubernetes API https://github.com/kubernetes/client-go/blob/master/kubernetes/typed/core/v1/pod_expansion.go#L48
// Response see https://docs.openshift.com/container-platform/4.6/rest_api/workloads_apis/pod-core-v1.html#apiv1namespacesnamespacepodsnamelog
//
// Location in archive: logs/openshift-api-server-operator
func GatherOpenShiftAPIServerOperatorLogs(g *Gatherer) func() ([]record.Record, []error) {
	return func() ([]record.Record, []error) {
		gatherKubeClient, err := kubernetes.NewForConfig(g.gatherProtoKubeConfig)
		if err != nil {
			return nil, []error{err}
		}

		client := gatherKubeClient.CoreV1()

		records, err := gatherOpenShiftAPIServerOperatorLogs(g.ctx, client)
		if err != nil {
			return nil, []error{err}
		}

		return records, nil
	}
}

func gatherOpenShiftAPIServerOperatorLogs(ctx context.Context, coreClient corev1client.CoreV1Interface) ([]record.Record, error) {
	const namespace = "openshift-apiserver-operator"
	var (
		messagesToSearch = []string{
			"the server has received too many requests and has asked us",
			"because serving request timed out and response had been started",
		}
		sinceSeconds int64 = 86400     // last day
		limitBytes   int64 = 1024 * 64 // maximum 64 kb of logs
	)

	pods, err := coreClient.Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var records []record.Record

	for _, pod := range pods.Items {
		request := coreClient.Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
			SinceSeconds: &sinceSeconds,
			LimitBytes:   &limitBytes,
		})

		logs, err := filterLogs(ctx, request, messagesToSearch)
		if err != nil {
			return nil, err
		}

		records = append(records, record.Record{
			Name:        "logs/openshift-api-server-operator",
			Item:        Raw{logs},
		})
	}

	return records, nil
}

func filterLogs(ctx context.Context, request *restclient.Request, messagesToSearch []string) (string, error) {
	stream, err := request.Stream(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		err := stream.Close()
		if err != nil {
			klog.Errorf("error during closing a stream: %v", err)
		}
	}()

	scanner := bufio.NewScanner(stream)

	var result string

	for scanner.Scan() {
		line := scanner.Text()
		for _, messageToSearch := range messagesToSearch {
			if strings.Contains(strings.ToLower(line), strings.ToLower(messageToSearch)) {
				result += line + "\n"
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return result, nil
}
