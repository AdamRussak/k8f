// https: //github.com/DataDog/datadog-api-client-go/blob/master/examples/v2/metrics/SubmitMetrics.go
// https://docs.datadoghq.com/metrics/custom_metrics/
// Submit metrics returns "Payload accepted" response

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"k8f/core"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/common"
	datadog "github.com/DataDog/datadog-api-client-go/v2/api/v2/datadog"
)

// FIXME: i need to get access to DD to test this feature
// DD_SITE="datadoghq.com" DD_API_KEY="<DD_API_KEY>"
func DdMain() {
	core.CheckEnvVarOrSitIt("DD_API_KEY", "")
	body := datadog.MetricPayload{
		Series: []datadog.MetricSeries{
			{
				Metric: "system.load.1",
				Type:   datadog.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadog.MetricPoint{
					{
						Timestamp: common.PtrInt64(time.Now().Unix()),
						Value:     common.PtrFloat64(0.7),
					},
				},
				Resources: []datadog.MetricResource{
					{
						Name: common.PtrString("dummyhost"),
						Type: common.PtrString("host"),
					},
				},
			},
		},
	}
	ctx := common.NewDefaultContext(context.Background())
	configuration := common.NewConfiguration()
	apiClient := common.NewAPIClient(configuration)
	api := datadog.NewMetricsApi(apiClient)
	resp, r, err := api.SubmitMetrics(ctx, body, *datadog.NewSubmitMetricsOptionalParameters())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.SubmitMetrics`: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.SubmitMetrics`:\n%s\n", responseContent)
}
