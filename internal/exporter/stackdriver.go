package exporter

import (
	"contrib.go.opencensus.io/exporter/stackdriver"

	"go.opencensus.io/trace"
)

func InitStackdriver(projectID string) {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
    panic(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(), 
	})
}
