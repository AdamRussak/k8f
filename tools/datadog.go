// https: //github.com/DataDog/datadog-api-client-go/blob/master/examples/v2/metrics/SubmitMetrics.go
// https://docs.datadoghq.com/metrics/custom_metrics/
// Submit metrics returns "Payload accepted" response

package tools

type Test struct {
	Version string `json:"version,omitempty"`
}

// FIXME: i need to get access to DD to test this feature
// DD_SITE="datadoghq.com" DD_API_KEY="<DD_API_KEY>"
// func DdMain(dd_api string) {
// 	core.CheckEnvVarOrSitIt("DD_API_KEY", dd_api)
// 	body := datadogV2.MetricPayload{
// 		Series: []datadogV2.MetricSeries{
// 			{
// 				Metric: "adam.test.1",
// 				Type:   datadogV2.METRICINTAKETYPE_UNSPECIFIED.Ptr(),
// 				Points: []datadogV2.MetricPoint{
// 					{
// 						Timestamp: datadog.PtrInt64(time.Now().Unix()),
// 						Value:     datadog.PtrFloat64(1),
// 					},
// 				},
// 				Resources: []datadogV2.MetricResource{
// 					{
// 						Name: datadog.PtrString("1.2.3"),
// 						Type: datadog.PtrString("host"),
// 					},
// 				},
// 			},
// 		},
// 	}
// 	ctx := datadog.NewDefaultContext(context.Background())
// 	configuration := datadog.NewConfiguration()
// 	apiClient := datadog.NewAPIClient(configuration)
// 	api := datadogV2.NewMetricsApi(apiClient)
// 	resp, r, err := api.SubmitMetrics(ctx, body, *datadogV2.NewSubmitMetricsOptionalParameters())
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.SubmitMetrics`: %v\n", err)
// 		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
// 	}

// 	responseContent, _ := json.MarshalIndent(resp, "", "  ")
// 	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.SubmitMetrics`:\n%s\n", string(responseContent))
// }
