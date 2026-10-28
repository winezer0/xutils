package logging_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/winezer0/xutils/cmdutils"
	"github.com/winezer0/xutils/logging"
	"github.com/winezer0/xutils/poolwriter"
)

func TestSilentLogger(t *testing.T) {
	// The process starts with a silent logger.
	logging.Debugf("default %d", 1)
	logging.Warnf("default %d", 2)
	logging.Errorf("default %d", 3)
	logging.SetLogger(logging.FuncLogger{})
	logging.Debugf("empty")
	logging.Warnf("empty")
	logging.Errorf("empty")
	logging.SetLogger(nil)
	logging.Errorf("reset")
}

func TestLevelAndArgumentForwarding(t *testing.T) {
	t.Cleanup(func() { logging.SetLogger(nil) })
	var got []string
	record := func(level string) func(string, ...any) {
		return func(format string, args ...any) {
			if format != "value=%d %s" {
				t.Errorf("format changed: %q", format)
			}
			got = append(got, level+":"+fmt.Sprintf(format, args...))
		}
	}
	logging.SetLogger(logging.FuncLogger{
		DebugFunc: record("debug"),
		WarnFunc:  record("warn"),
		ErrorFunc: record("error"),
	})
	logging.Debugf("value=%d %s", 1, "a")
	logging.Warnf("value=%d %s", 2, "b")
	logging.Errorf("value=%d %s", 3, "c")
	want := []string{"debug:value=1 a", "warn:value=2 b", "error:value=3 c"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	logging.SetLogger(nil)
	logging.Errorf("discarded")
	if len(got) != 3 {
		t.Fatal("nil did not restore silent logging")
	}
}

func TestDiagnosticsReachConfiguredLogger(t *testing.T) {
	t.Cleanup(func() { logging.SetLogger(nil) })
	var mu sync.Mutex
	warnings, errors := 0, 0
	logging.SetLogger(logging.FuncLogger{
		WarnFunc: func(string, ...any) {
			mu.Lock()
			warnings++
			mu.Unlock()
		},
		ErrorFunc: func(string, ...any) {
			mu.Lock()
			errors++
			mu.Unlock()
		},
	})
	headers := cmdutils.ParseHeaders([]string{"invalid", "Accept=text/plain"})
	if headers["Accept"] != "text/plain" {
		t.Fatal("valid header was lost")
	}
	pool := poolwriter.NewPool(4, 16)
	for i := 0; i < 16; i++ {
		pool.Submit(poolwriter.StoreTask{}) // Empty path deterministically fails.
	}
	pool.StopAndWait()
	mu.Lock()
	defer mu.Unlock()
	if warnings != 1 || errors != 16 || pool.GetFailCount() != 16 {
		t.Fatalf("warnings=%d errors=%d failures=%d", warnings, errors, pool.GetFailCount())
	}
}
