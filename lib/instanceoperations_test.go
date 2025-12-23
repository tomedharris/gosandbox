package lib

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type DoFunc func(*http.Request) (*http.Response, error)

func (f DoFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

func loadFixture(t *testing.T, fname string) string {
	fn := "testdata/" + fname
	if _, err := os.Stat(fn); err != nil {
		t.Errorf("fixture %s not found", fn)
	}
	b, err := os.ReadFile(fn)
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	return string(b)
}

func Test_ListInstances(t *testing.T) {
	tests := []struct {
		title             string
		body              string
		numberOfInstances uint8
		instanceIds       []string
	}{
		{"Single Instance", loadFixture(t, "describe-instances-single.xml"), 1, []string{"i-1234567890abcdef0"}},
		{"Multiple Instances", loadFixture(t, "describe-instances-multiple.xml"), 2, []string{"i-aaa", "i-bbb"}},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			httpClient := DoFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(tt.body)),
					Header:     http.Header{"Content-Type": []string{"text/xml"}},
				}, nil
			})
			cfg, err := config.LoadDefaultConfig(t.Context(), config.WithHTTPClient(httpClient))
			if err != nil {
				t.Errorf("error setting up config: %s", err)
			}
			client := ec2.NewFromConfig(cfg)
			result, err := ListInstances(client)
			if err != nil {
				t.Fatalf("error listing instances: %s", err)
			}
			if len(result) != int(tt.numberOfInstances) {
				t.Fatalf("expected %d instance, got %d", tt.numberOfInstances, len(result))
			}
			for idx, expectedInstanceId := range tt.instanceIds {
				if result[idx].InstanceId != expectedInstanceId {
					t.Fatalf("expected instance id %s, got %s", expectedInstanceId, result[0].InstanceId)
				}
			}
		})
	}
}
