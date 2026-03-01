package reconstruct_otel_http

import (
	"sort"
	"testing"
)

func TestBuildResourceAttributes(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]string
		want     []attribute
	}{
		{
			name:     "empty metadata",
			metadata: map[string]string{},
			want:     nil,
		},
		{
			name: "maps service_name to OTEL convention",
			metadata: map[string]string{
				"service_name": "my-api",
			},
			want: []attribute{
				stringAttribute("service.name", "my-api"),
			},
		},
		{
			name: "maps environment_name to OTEL convention",
			metadata: map[string]string{
				"environment_name": "production",
			},
			want: []attribute{
				stringAttribute("deployment.environment.name", "production"),
			},
		},
		{
			name: "maps project_name to service.namespace",
			metadata: map[string]string{
				"project_name": "my-project",
			},
			want: []attribute{
				stringAttribute("service.namespace", "my-project"),
			},
		},
		{
			name: "passes through unknown keys unchanged",
			metadata: map[string]string{
				"project_id": "abc-123",
			},
			want: []attribute{
				stringAttribute("project_id", "abc-123"),
			},
		},
		{
			name: "all mapped keys plus passthrough",
			metadata: map[string]string{
				"service_name":     "my-api",
				"environment_name": "production",
				"project_name":     "my-project",
				"deployment_id":    "deploy-456",
				"log_type":         "environment",
			},
			want: []attribute{
				stringAttribute("service.name", "my-api"),
				stringAttribute("deployment.environment.name", "production"),
				stringAttribute("service.namespace", "my-project"),
				stringAttribute("deployment_id", "deploy-456"),
				stringAttribute("log_type", "environment"),
			},
		},
		{
			name: "mapped keys do not appear twice",
			metadata: map[string]string{
				"service_name": "my-api",
				"project_name": "my-project",
			},
			want: []attribute{
				stringAttribute("service.name", "my-api"),
				stringAttribute("service.namespace", "my-project"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildResourceAttributes(tt.metadata)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d attributes, want %d\ngot:  %v\nwant: %v", len(got), len(tt.want), got, tt.want)
			}

			// Sort both slices by key for stable comparison since map iteration is unordered
			sort.Slice(got, func(i, j int) bool { return got[i].Key < got[j].Key })
			sort.Slice(tt.want, func(i, j int) bool { return tt.want[i].Key < tt.want[j].Key })

			for i := range got {
				if got[i].Key != tt.want[i].Key {
					t.Errorf("attribute[%d] key = %q, want %q", i, got[i].Key, tt.want[i].Key)
				}
				if *got[i].Value.StringValue != *tt.want[i].Value.StringValue {
					t.Errorf("attribute[%d] value = %q, want %q", i, *got[i].Value.StringValue, *tt.want[i].Value.StringValue)
				}
			}
		})
	}
}
