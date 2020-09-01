package healthcheck

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// These constants are the possible status for responses. Use them to prevent typos.
const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

const (
	defaultResolution = 3 * time.Second
	defaultTimeout    = 1 * time.Second
)

// New creates a new health check service that can be mounted as a HTTP handler
func New(service string) *HealthCheck {
	return &HealthCheck{
		service: service,
	}
}

// HealthCheck is a struct to hold health checks and creates a handler for this.
type HealthCheck struct {
	service      string
	dependencies []*Monitor
}

// Monitor will allow you to register a service to monitor. We will call
// the `MonitorFunc` every `Resolution`.
type Monitor struct {
	Name        string
	Category    string
	MonitorFunc func(context.Context) error
	Resolution  time.Duration
	Timeout     time.Duration
	result      *serviceState
	mu          sync.RWMutex
}

type serviceState struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Took    string `json:"took"`
	Details string `json:"details,omitempty"`
}

type response struct {
	CreatedAt time.Time                  `json:"created_at"`
	Service   string                     `json:"service"`
	Status    string                     `json:"status"`
	Took      string                     `json:"took"`
	Monitors  map[string][]*serviceState `json:"monitors"`
}

func (m *Monitor) status(ctx context.Context) *serviceState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.result == nil {
		m.mu.RUnlock()
		m.check(ctx)
		m.mu.RLock()
	}

	return m.result
}

func (m *Monitor) start(ctx context.Context) {
	ticker := time.NewTicker(m.Resolution)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf(`{"msg":"stopping monitoring","service":%q`, m.Name)
			return
		case <-ticker.C:
			m.check(ctx)
		}
	}
}

func (m *Monitor) check(ctx context.Context) error {
	var (
		start = time.Now()
		err   error
	)

	status := &serviceState{Name: m.Name, Status: StatusOK}
	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()

	if err = m.MonitorFunc(ctx); err != nil {
		status.Status = StatusError
		status.Details = err.Error()
	}

	status.Took = time.Now().Sub(start).String()

	m.mu.Lock()
	defer m.mu.Unlock()
	m.result = status
	return err
}

// AddMonitor adds a new check function.
func (hc *HealthCheck) AddMonitor(m *Monitor) {
	if m.Resolution == 0 {
		m.Resolution = defaultResolution
	}

	if m.Timeout == 0 {
		m.Timeout = defaultTimeout
	}

	hc.dependencies = append(hc.dependencies, m)
}

// check will check every dependency monitor and compile results. This will read
// the status for **every** dependency and create the response. But each dependency
// will have it's cached status.
func (hc *HealthCheck) check(ctx context.Context) *response {
	resp := &response{
		Service:   hc.service,
		Status:    StatusOK,
		Monitors:  make(map[string][]*serviceState),
		CreatedAt: time.Now(),
	}

	for _, service := range hc.dependencies {
		check := service.status(ctx)

		resp.Monitors[service.Category] = append(resp.Monitors[service.Category], check)
		if check.Status != StatusOK {
			resp.Status = StatusError
		}
	}

	resp.Took = time.Now().Sub(resp.CreatedAt).String()
	return resp
}

// Start will start monitoring all services and pass a channel that we
// can close letting them know it's time to exit.
func (hc *HealthCheck) Start(ctx context.Context) {
	for _, monitor := range hc.dependencies {
		go monitor.start(ctx)
	}
}

// ServeHTTP implements the http.Hanlder interface and allows us to mount the
// health check service in any endpoint.
func (hc *HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	resp := hc.check(r.Context())
	if resp.Status != StatusOK {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	encoder.Encode(resp)
}
