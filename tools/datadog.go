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

type Test struct {
	Version string `json:"version,omitempty"`
}

// FIXME: i need to get access to DD to test this feature
// DD_SITE="datadoghq.com" DD_API_KEY="<DD_API_KEY>"
func DdMain(dd_api string) {
	t := Test{Version: "1.2.0"}
	jsonStr, _ := json.Marshal(t)
	var mapData map[string]interface{}
	if err := json.Unmarshal(jsonStr, &mapData); err != nil {
		fmt.Println(err)
	}
	core.CheckEnvVarOrSitIt("DD_API_KEY", dd_api)
	body := datadog.MetricPayload{
		Series: []datadog.MetricSeries{
			{
				Metric: "adam.test.1",
				Type:   datadog.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
				Points: []datadog.MetricPoint{
					{
						Timestamp:      common.PtrInt64(time.Now().Unix()),
						UnparsedObject: mapData,
						Value:          common.PtrFloat64(1),
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
	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.SubmitMetrics`:\n%s\n", string(responseContent))
}
