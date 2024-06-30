package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemStorage(t *testing.T) {
	ms := NewMemStorage()
	assert.NotNil(t, ms, "expected MemStorage instance, got nil")
	assert.NotNil(t, ms.gauges, "expected gauges map to be initialized, got nil")
	assert.NotNil(t, ms.counters, "expected counters map to be initialized, got nil")
}

func TestUpdateGauge(t *testing.T) {
	tests := []struct {
		name      string
		gaugeName string
		value     float64
		want      float64
	}{
		{
			name:      "UpdateGauge with positive value",
			gaugeName: "testGauge1",
			value:     123.45,
			want:      123.45,
		},
		{
			name:      "UpdateGauge with negative value",
			gaugeName: "testGauge2",
			value:     -2.4,
			want:      -2.4,
		},
		{
			name:      "UpdateGauge with zero value",
			gaugeName: "testGauge3",
			value:     0,
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			ms.UpdateGauge(tt.gaugeName, tt.value)
			val, ok := ms.gauges[tt.gaugeName]
			assert.True(t, ok, "expected gauge '%s' to be present", tt.gaugeName)
			assert.Equal(t, tt.want, val, "expected gauge '%s' to be %f, got %f", tt.gaugeName, tt.want, val)
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	tests := []struct {
		name        string
		counterName string
		value       int64
		want        int64
	}{
		{
			name:        "UpdateCounter with positive value",
			counterName: "testCounter1",
			value:       33,
			want:        33,
		},
		{
			name:        "UpdateCounter with negative value",
			counterName: "testCounter2",
			value:       -22,
			want:        -22,
		},
		{
			name:        "UpdateCounter with zero value",
			counterName: "testCounter3",
			value:       0,
			want:        0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			ms.UpdateCounter(tt.counterName, tt.value)
			val, ok := ms.counters[tt.counterName]
			assert.True(t, ok, "expected counter '%s' to be present", tt.counterName)
			assert.Equal(t, tt.want, val, "expected counter '%s' to be %d, got %d", tt.counterName, tt.want, val)
		})
	}
}

func TestGetGauge(t *testing.T) {
	tests := []struct {
		name       string
		gaugeName  string
		initial    map[string]float64
		wantValue  float64
		wantExists bool
	}{
		{
			name:       "Existing gauge",
			gaugeName:  "testGauge1",
			initial:    map[string]float64{"testGauge1": 42.0},
			wantValue:  42.0,
			wantExists: true,
		},
		{
			name:       "Non-existing gauge",
			gaugeName:  "nonExistingGauge",
			initial:    map[string]float64{"testGauge2": 33.3},
			wantValue:  0,
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{gauges: tt.initial}
			val, ok := ms.GetGauge(tt.gaugeName)
			assert.Equal(t, tt.wantExists, ok, "expected gauge '%s' existence to be %v", tt.gaugeName, tt.wantExists)
			assert.Equal(t, tt.wantValue, val, "expected gauge '%s' value to be %f, got %f", tt.gaugeName, tt.wantValue, val)
		})
	}
}

func TestGetCounter(t *testing.T) {
	tests := []struct {
		name        string
		counterName string
		initial     map[string]int64
		wantValue   int64
		wantExists  bool
	}{
		{
			name:        "Existing counter",
			counterName: "testCounter1",
			initial:     map[string]int64{"testCounter1": 100},
			wantValue:   100,
			wantExists:  true,
		},
		{
			name:        "Non-existing counter",
			counterName: "nonExistingCounter",
			initial:     map[string]int64{"testCounter2": 200},
			wantValue:   0,
			wantExists:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{counters: tt.initial}
			val, ok := ms.GetCounter(tt.counterName)
			assert.Equal(t, tt.wantExists, ok, "expected counter '%s' existence to be %v", tt.counterName, tt.wantExists)
			assert.Equal(t, tt.wantValue, val, "expected counter '%s' value to be %d, got %d", tt.counterName, tt.wantValue, val)
		})
	}
}

func TestGetAllMetrics(t *testing.T) {
	tests := []struct {
		name     string
		initial  *MemStorage
		expected string
	}{
		{
			name:     "Empty storage",
			initial:  NewMemStorage(),
			expected: "<html><body></body></html>",
		},
		{
			name: "Storage with metrics",
			initial: &MemStorage{
				gauges:   map[string]float64{"gauge1": 42.0},
				counters: map[string]int64{"counter1": 100},
			},
			expected: "<html><body><br>gauge/gauge1: 42</br><br>counter/counter1: 100</br></body></html>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.initial.GetAllMetrics(), "expected GetAllMetrics() output to match")
		})
	}
}

func TestSetGauge(t *testing.T) {
	tests := []struct {
		name       string
		gaugeName  string
		value      float64
		wantValue  float64
		wantExists bool
	}{
		{
			name:       "Set gauge",
			gaugeName:  "gauge1",
			value:      42.0,
			wantValue:  42.0,
			wantExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			ms.SetGauge(tt.gaugeName, tt.value)
			val, ok := ms.gauges[tt.gaugeName]
			assert.True(t, ok, "expected gauge '%s' to be present", tt.gaugeName)
			assert.Equal(t, tt.wantValue, val, "expected gauge '%s' to be %f, got %f", tt.gaugeName, tt.wantValue, val)
		})
	}
}

func TestSetCounter(t *testing.T) {
	tests := []struct {
		name        string
		counterName string
		value       int64
		wantValue   int64
		wantExists  bool
	}{
		{
			name:        "Set counter",
			counterName: "counter1",
			value:       100,
			wantValue:   100,
			wantExists:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			ms.SetCounter(tt.counterName, tt.value)
			val, ok := ms.counters[tt.counterName]
			assert.True(t, ok, "expected counter '%s' to be present", tt.counterName)
			assert.Equal(t, tt.wantValue, val, "expected counter '%s' to be %d, got %d", tt.counterName, tt.wantValue, val)
		})
	}
}
