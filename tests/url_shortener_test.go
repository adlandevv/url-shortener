package tests

import (
	"net/http"
	"net/url"
	"path"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	// Universal method for creating URLs
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	// Creating httpexpect client
	e := httpexpect.Default(t, u.String())

	e.POST("/url"). // Sending POST-request, path - '/url'
			WithJSON(save.Request{ // Constructing request body
			URL:   gofakeit.URL(),             // Generating a random URL
			Alias: random.NewRandomString(10), // Generating a random string
		}).
		WithBasicAuth("my_user", "my_pass"). // Add credentials to the request
		Expect().                            // Next, we list the expectations for the response
		Status(200).                         // Code need to be 200
		JSON().Object().                     // Get the JSON object from the response body
		ContainsKey("alias")                 // Check that it has the "alias" key
}

func TestURLShortener_SaveRedirectRemove(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("my_user", "my_pass").
				Expect().Status(200).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect
			testRedirect(t, alias, tc.url)

			// Remove
			reqDel := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("my_user", "my_pass").
				Expect().Status(http.StatusOK).JSON().Object()

			reqDel.Value("status").String().IsEqual("OK")

			// Redirect again
			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())
	require.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
