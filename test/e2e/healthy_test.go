package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bkielbasa/go-web-app/internal/application"
	"github.com/bkielbasa/go-web-app/links"
	"github.com/bkielbasa/go-web-app/links/app"
	"github.com/bkielbasa/go-web-app/links/infra"
	_ "github.com/lib/pq"
	"github.com/matryer/is"
)

func TestRunningApp(t *testing.T) {

	is := is.New(t)

	runApp(func() {
		err := retry(checkReadyStatus, time.Second, 100*time.Millisecond)
		is.NoErr(err)

		err = retry(checkHealthyStatus, time.Second, 100*time.Millisecond)
		is.NoErr(err)
	})
}

const shortBaseURL = "http://localhost:8080"

func runApp(f func()) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var storage app.Storage

	if testing.Short() {
		storage = infra.NewInMemoryStorage()
	} else {
		db, err := sql.Open("postgres", os.Getenv("POSTGRESCONNSTRING"))
		if err != nil {
			panic(err)
		}
		storage = infra.NewPostgresStorage(db)
	}

	app := application.New(ctx, 8080)
	app.AddModule(links.New(shortBaseURL, storage))

	defer func() {
		tearCtx, cancelTear := context.WithTimeout(context.Background(), time.Second)
		defer cancelTear()

		_ = app.Close(tearCtx)
	}()

	go func() {
		_ = app.Run()
	}()

	f()
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
