package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/senyabanana/go-mcaa-service/internal/models"
	"github.com/senyabanana/go-mcaa-service/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewAgent(t *testing.T) {
	type want struct {
		expectedURL    string
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
				expectedURL:    "http://example.com",
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
			assert.Equal(t, agent.ServerURL, tt.want.expectedURL)
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
		metric        models.Metrics
		expectedError bool
	}{
		{
			name: "Send Gauge Metric",
			metric: models.Metrics{
				ID:    "TestGauge",
				MType: "gauge",
				Value: func() *float64 { v := 123.45; return &v }(),
			},
			expectedError: false,
		},
		{
			name: "Send Counter Metric",
			metric: models.Metrics{
				ID:    "TestCounter",
				MType: "counter",
				Delta: func() *int64 { v := int64(123); return &v }(),
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)

				var receivedMetric models.Metrics
				err := json.NewDecoder(r.Body).Decode(&receivedMetric)
				assert.NoError(t, err)
				assert.Equal(t, tt.metric.ID, receivedMetric.ID)
				assert.Equal(t, tt.metric.MType, receivedMetric.MType)
				if tt.metric.MType == storage.Gauge {
					assert.Equal(t, *tt.metric.Value, *receivedMetric.Value)
				} else if tt.metric.MType == storage.Counter {
					assert.Equal(t, *tt.metric.Delta, *receivedMetric.Delta)
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			agent := NewAgent(server.URL, time.Second, time.Second)
			err := agent.sendMetric(tt.metric)

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
