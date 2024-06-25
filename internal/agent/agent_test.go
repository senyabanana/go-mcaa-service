package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAgent(t *testing.T) {
	type want struct {
		expectedUrl    string
		expectedPoll   time.Duration
		expectedReport time.Duration
	}
	tests := []struct {
		name           string
		pollInterval   time.Duration
		reportInterval time.Duration
		url            string
		want           want
	}{
		{
			name:           "Basic Agent Creation",
			pollInterval:   2 * time.Second,
			reportInterval: 10 * time.Second,
			url:            "http://example.com",
			want: want{
				expectedUrl:    "http://example.com",
				expectedPoll:   2 * time.Second,
				expectedReport: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewAgent(tt.url, tt.pollInterval, tt.reportInterval)

			assert.Equal(t, agent.PollInterval, tt.want.expectedPoll)
			assert.Equal(t, agent.ReportInterval, tt.want.expectedReport)
			assert.Equal(t, agent.ServerUrl, tt.want.expectedUrl)
			assert.NotNil(t, agent.gauges)
			assert.NotNil(t, agent.counters)
		})
	}
}

func TestCollectRuntimeMetrics(t *testing.T) {
	tests := []struct {
		name      string
		pollCount int64
	}{
		{
			name:      "Basic Metric Collection",
			pollCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewAgent("", time.Second, time.Second)
			agent.collectRuntimeMetrics()

			agent.mu.Lock()
			defer agent.mu.Unlock()

			assert.Greater(t, len(agent.gauges), 0)
			assert.Equal(t, tt.pollCount, agent.counters["PollCount"])
		})
	}
}

func TestSendMetric(t *testing.T) {
	tests := []struct {
		name          string
		metricType    string
		metricName    string
		metricValue   interface{}
		expectedError bool
	}{
		{
			name:          "Send Gauge Metric",
			metricType:    "gauge",
			metricName:    "TestGauge",
			metricValue:   123.45,
			expectedError: false,
		},
		{
			name:          "Send Counter Metric",
			metricType:    "counter",
			metricName:    "TestCounter",
			metricValue:   123,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			agent := NewAgent(server.URL, time.Second, time.Second)
			err := agent.sendMetric(tt.metricType, tt.metricName, tt.metricValue)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSendAllMetrics(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Send All Metrics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			agent := NewAgent(server.URL, time.Second, time.Second)
			agent.collectRuntimeMetrics()
			agent.sendAllMetrics()
		})
	}
}
