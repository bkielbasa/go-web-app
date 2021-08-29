package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bkielbasa/go-web-app/links/infra"
	"github.com/matryer/is"
)

func TestAddNewLinks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short tests")
	}

	is := is.New(t)

	runApp(func() {
		shortLink, err := addNewLink("https://google.com", []string{})
		is.NoErr(err) // creating a new link shouldn't return any error
		is.True(strings.HasPrefix(shortLink, shortBaseURL))
	})
}

func TestLinksShouldBeUnique(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short tests")
	}

	is := is.New(t)

	runApp(func() {
		shortLink, err := addNewLink("https://google.com", []string{})
		is.NoErr(err)
		is.True(strings.HasPrefix(shortLink, shortBaseURL))

		shortLink2, err := addNewLink("https://google.com", []string{})
		is.NoErr(err)
		is.True(strings.HasPrefix(shortLink, shortBaseURL))
		is.True(shortLink != shortLink2)
	})
}

func TestClickingLinkShouldRedirectToTargetPage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short tests")
	}

	is := is.New(t)

	runApp(func() {
		shortLink, err := addNewLink("https://google.com", []string{})
		is.NoErr(err)
		is.True(strings.HasPrefix(shortLink, shortBaseURL))

		t.Log(shortLink)
		resp, err := clickLink(shortLink)
		is.Equal(http.StatusMovedPermanently, resp.StatusCode)
		is.NoErr(err)
	})
}

const baseUrl = "http://localhost:8080"

func clickLink(target string) (*http.Response, error) {
	c := http.Client{
		// if CheckRedirect returns ErrUseLastResponse,
		// then the most recent respons∑ is returned with its body
		// unclosed, along with a nil error.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest(http.MethodGet, target, nil)
	return c.Do(req)
}

func addNewLink(target string, tags []string) (string, error) {
	c := http.Client{}

	addReq := infra.AddLinkRequest{
		Target: target,
		Tags:   tags,
	}
	reqBody, _ := json.Marshal(addReq)
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/link", bytes.NewReader(reqBody))
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot send the request: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("problem while reading the body: %w", err)
	}

	r := infra.AddLinkResponse{}
	if err = json.Unmarshal(respBody, &r); err != nil {
		return "", fmt.Errorf("problem while unmarshaling the body %s: %w", respBody, err)
	}

	return r.ShortURL, nil
}
