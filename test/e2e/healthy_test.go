package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bkielbasa/go-web-app/cmd/web/run"
)

func TestRunningApp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short tests")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	run, tearDown, err := run.App(ctx)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		tearCtx, cancelTear := context.WithTimeout(context.Background(), time.Second)
		defer cancelTear()

		_ = tearDown(tearCtx)
	}()
	go run()

	err = retry(checkReadyStatus, time.Second, 100*time.Millisecond)
	if err != nil {
		t.Errorf("cannot check the ready status: %s", err)
	}

	err = retry(checkHealthyStatus, time.Second, 100*time.Millisecond)
	if err != nil {
		t.Errorf("cannot check the healthy status: %s", err)
	}

	// the app is ready to go!
}

func checkHealthyStatus() bool {
	resp, err := http.Get("http://localhost:8080/healthyz")
	if err != nil {
		return false
	}

	if resp.StatusCode >= 500 {
		return false
	}

	return true
}

func checkReadyStatus() bool {
	resp, err := http.Get("http://localhost:8080/readyz")
	if err != nil {
		return false
	}

	if resp.StatusCode >= 500 {
		return false
	}

	return true
}

func retry(condition func() bool, waitFor time.Duration, tick time.Duration) error {
	ch := make(chan bool, 1)

	timer := time.NewTimer(waitFor)
	defer timer.Stop()

	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for tick := ticker.C; ; {
		select {
		case <-timer.C:
			return fmt.Errorf("condition never satisfied")
		case <-tick:
			tick = nil
			go func() { ch <- condition() }()
		case v := <-ch:
			if v {
				return nil
			}
			tick = ticker.C
		}
	}
}
