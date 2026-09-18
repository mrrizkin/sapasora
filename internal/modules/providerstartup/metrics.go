package providerstartup

import (
	"sort"
	"sync"
	"time"
)

// LifecycleMetrics records provider connection lifecycle events without
// retaining device identifiers. Keeping the cardinality at provider level
// makes the counters safe to export to a metrics backend.
type LifecycleMetrics struct {
	mu        sync.RWMutex
	providers map[string]*providerMetrics
}

type providerMetrics struct {
	startupAttempts  uint64
	startupSuccesses uint64
	startupFailures  uint64
	disconnects      uint64
	active           uint64
	lastFailureAt    time.Time
}

// LifecycleSnapshot is a point-in-time view of one provider's lifecycle
// counters.
type LifecycleSnapshot struct {
	Provider          string
	StartupAttempts   uint64
	StartupSuccesses  uint64
	StartupFailures   uint64
	Disconnects       uint64
	ActiveConnections uint64
	LastFailureAt     time.Time
}

// NewLifecycleMetrics creates an isolated metrics collector. Each provider
// owns one collector so tests and embedders can inspect lifecycle state without
// relying on global process state.
func NewLifecycleMetrics() *LifecycleMetrics {
	return &LifecycleMetrics{providers: make(map[string]*providerMetrics)}
}

func (m *LifecycleMetrics) provider(name string) *providerMetrics {
	if m.providers == nil {
		m.providers = make(map[string]*providerMetrics)
	}
	metrics := m.providers[name]
	if metrics == nil {
		metrics = &providerMetrics{}
		m.providers[name] = metrics
	}
	return metrics
}

// RecordStartupAttempt records one actual provider connection attempt.
func (m *LifecycleMetrics) RecordStartupAttempt(provider string) {
	m.record(provider, func(metrics *providerMetrics) {
		metrics.startupAttempts++
	})
}

// RecordStartupSuccess records a connection that became usable.
func (m *LifecycleMetrics) RecordStartupSuccess(provider string) {
	m.record(provider, func(metrics *providerMetrics) {
		metrics.startupSuccesses++
		metrics.active++
	})
}

// RecordStartupFailure records a failed provider connection attempt.
func (m *LifecycleMetrics) RecordStartupFailure(provider string) {
	m.record(provider, func(metrics *providerMetrics) {
		metrics.startupFailures++
		metrics.lastFailureAt = time.Now().UTC()
	})
}

// RecordDisconnect records an active provider connection that stopped. The
// active gauge is saturating so cleanup paths can safely be retried without
// counting cleanup for a connection that never became usable.
func (m *LifecycleMetrics) RecordDisconnect(provider string) {
	m.record(provider, func(metrics *providerMetrics) {
		if metrics.active == 0 {
			return
		}
		metrics.disconnects++
		metrics.active--
	})
}

func (m *LifecycleMetrics) record(provider string, update func(*providerMetrics)) {
	if m == nil || provider == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	update(m.provider(provider))
}

// Snapshot returns a consistent view of one provider. Unknown providers have
// zero counters and retain the requested provider name.
func (m *LifecycleMetrics) Snapshot(provider string) LifecycleSnapshot {
	snapshot := LifecycleSnapshot{Provider: provider}
	if m == nil {
		return snapshot
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	metrics := m.providers[provider]
	if metrics == nil {
		return snapshot
	}
	return snapshotFrom(provider, metrics)
}

// Snapshots returns all provider snapshots in stable provider-name order.
func (m *LifecycleMetrics) Snapshots() []LifecycleSnapshot {
	if m == nil {
		return nil
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	providers := make([]string, 0, len(m.providers))
	for provider := range m.providers {
		providers = append(providers, provider)
	}
	sort.Strings(providers)

	snapshots := make([]LifecycleSnapshot, 0, len(providers))
	for _, provider := range providers {
		snapshots = append(snapshots, snapshotFrom(provider, m.providers[provider]))
	}
	return snapshots
}

func snapshotFrom(provider string, metrics *providerMetrics) LifecycleSnapshot {
	return LifecycleSnapshot{
		Provider:          provider,
		StartupAttempts:   metrics.startupAttempts,
		StartupSuccesses:  metrics.startupSuccesses,
		StartupFailures:   metrics.startupFailures,
		Disconnects:       metrics.disconnects,
		ActiveConnections: metrics.active,
		LastFailureAt:     metrics.lastFailureAt,
	}
}
