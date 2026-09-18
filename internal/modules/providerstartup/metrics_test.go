package providerstartup

import (
	"sync"
	"testing"
)

func TestLifecycleMetricsTracksProviderState(t *testing.T) {
	metrics := NewLifecycleMetrics()

	metrics.RecordStartupAttempt("telegram")
	metrics.RecordStartupAttempt("telegram")
	metrics.RecordStartupFailure("telegram")
	metrics.RecordStartupSuccess("telegram")
	metrics.RecordDisconnect("telegram")
	metrics.RecordDisconnect("telegram")

	snapshot := metrics.Snapshot("telegram")
	if snapshot.StartupAttempts != 2 {
		t.Errorf("startup attempts = %d, want 2", snapshot.StartupAttempts)
	}
	if snapshot.StartupFailures != 1 {
		t.Errorf("startup failures = %d, want 1", snapshot.StartupFailures)
	}
	if snapshot.StartupSuccesses != 1 {
		t.Errorf("startup successes = %d, want 1", snapshot.StartupSuccesses)
	}
	if snapshot.Disconnects != 1 {
		t.Errorf("disconnects = %d, want 1", snapshot.Disconnects)
	}
	if snapshot.ActiveConnections != 0 {
		t.Errorf("active connections = %d, want 0", snapshot.ActiveConnections)
	}
	if snapshot.LastFailureAt.IsZero() {
		t.Fatal("last failure timestamp is zero")
	}
}

func TestLifecycleMetricsIsSafeForConcurrentUpdates(t *testing.T) {
	metrics := NewLifecycleMetrics()
	const workers = 8
	const eventsPerWorker = 100

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for range eventsPerWorker {
				metrics.RecordStartupAttempt("whatsapp")
				metrics.RecordStartupFailure("whatsapp")
			}
		}()
	}
	wg.Wait()

	snapshot := metrics.Snapshot("whatsapp")
	want := uint64(workers * eventsPerWorker)
	if snapshot.StartupAttempts != want || snapshot.StartupFailures != want {
		t.Fatalf("snapshot = %+v, want %d attempts and failures", snapshot, want)
	}
}

func TestLifecycleMetricsSnapshotsAreSortedAndUnknownIsEmpty(t *testing.T) {
	metrics := NewLifecycleMetrics()
	metrics.RecordStartupAttempt("whatsapp")
	metrics.RecordStartupAttempt("telegram")

	snapshots := metrics.Snapshots()
	if len(snapshots) != 2 || snapshots[0].Provider != "telegram" || snapshots[1].Provider != "whatsapp" {
		t.Fatalf("snapshots = %+v, want telegram then whatsapp", snapshots)
	}
	unknown := metrics.Snapshot("sms")
	if unknown.Provider != "sms" || unknown.StartupAttempts != 0 || unknown.ActiveConnections != 0 {
		t.Fatalf("unknown snapshot = %+v, want zero counters", unknown)
	}
}
