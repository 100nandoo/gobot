package reddit

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

const browserUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:139.0) Gecko/20100101 Firefox/139.0"

var redditHTTPClient = &http.Client{
	Transport: &http2.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			tcpConn, err := (&net.Dialer{}).DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}

			uconn := utls.UClient(tcpConn, &utls.Config{ServerName: host}, utls.HelloFirefox_Auto)
			if err := uconn.HandshakeContext(ctx); err != nil {
				tcpConn.Close()
				return nil, err
			}

			return uconn, nil
		},
	},
}

var (
	redditChallengePattern = regexp.MustCompile(`await\(async \w+\s*=>\s*\w+\s*\+\s*\w+\)\("([^"]+)"\)`)
	redditTokenPattern     = regexp.MustCompile(`name="token"\s+value="([^"]+)"`)
)

var getRedditLoidCookie = func() func() (string, error) {
	var lastUpdate time.Time
	var cachedLoid string

	return newSingleflight(func() (string, error) {
		if time.Since(lastUpdate) < 6*time.Hour && cachedLoid != "" {
			return cachedLoid, nil
		}

		loid, err := fetchRedditLoidCookie()
		if err != nil {
			if cachedLoid != "" {
				return cachedLoid, nil
			}
			return "", err
		}

		lastUpdate = time.Now()
		cachedLoid = loid
		return loid, nil
	})
}()

func setBrowserHeaders(req *http.Request) {
	req.Header.Set("User-Agent", browserUserAgent)
	req.Header.Set("Accept", "application/json,text/plain,*/*")
}

func fetchRedditLoidCookie() (string, error) {
	req, err := http.NewRequest("GET", "https://www.reddit.com/", nil)
	if err != nil {
		return "", err
	}

	setBrowserHeaders(req)

	resp, err := redditHTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	challengeBody, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if readErr != nil {
		return "", readErr
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d when requesting challenge page", resp.StatusCode)
	}

	challengeMatches := redditChallengePattern.FindSubmatch(challengeBody)
	if challengeMatches == nil {
		return "", fmt.Errorf("no JS challenge found")
	}

	tokenMatches := redditTokenPattern.FindSubmatch(challengeBody)
	if tokenMatches == nil {
		return "", fmt.Errorf("no token found in challenge page")
	}

	challengeStr := string(challengeMatches[1])
	params := url.Values{
		"solution":     {challengeStr + challengeStr},
		"js_challenge": {"1"},
		"token":        {string(tokenMatches[1])},
	}

	req, err = http.NewRequest("GET", "https://www.reddit.com/?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	setBrowserHeaders(req)

	resp, err = redditHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d when submitting challenge solution", resp.StatusCode)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "loid" {
			return cookie.Value, nil
		}
	}

	return "", fmt.Errorf("no loid cookie found")
}
