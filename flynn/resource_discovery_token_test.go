package flynn

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "regexp"
    "testing"

    "github.com/hashicorp/terraform/helper/resource"
    "github.com/hashicorp/terraform/terraform"
)

type TestHTTPMock struct {
	server *httptest.Server
}

const testFlynnDiscoveryTokenConfigBasic = `
resource "flynn_discovery_token" "test" {
  server = "%s"
}
output "token" {
  value = flynn_discovery_token.test.token
}
`

func TestDiscoveryToken_http201(t *testing.T) {
	TestHTTPMock := setUpMockHTTPServer()

	defer TestHTTPMock.server.Close()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testFlynnDiscoveryTokenConfigBasic, TestHTTPMock.server.URL),
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["flynn_discovery_token.test"]
					if !ok {
						return fmt.Errorf("Missing Discovery Token resource")
					}

					outputs := s.RootModule().Outputs

					if outputs["token"].Value != fmt.Sprintf("%s/clusters/toke-toke-token", TestHTTPMock.server.URL) {
						return fmt.Errorf(
							`'token' output is %s; want %s`,
							outputs["token"].Value, fmt.Sprintf("%s/clusters/toke-toke-token", TestHTTPMock.server.URL),
						)
					}

					return nil
				},
			},
		},
	})
}

func TestDiscoveryToken_http404(t *testing.T) {
	TestHTTPMock := setUpMockHTTPServer()

	defer TestHTTPMock.server.Close()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testFlynnDiscoveryTokenConfigBasic, fmt.Sprintf("%s/404", TestHTTPMock.server.URL)),
				ExpectError: regexp.MustCompile("HTTP request error. Response code: 404"),
			},
		},
	})
}

func setUpMockHTTPServer() *TestHTTPMock {
	Server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path == "/clusters" {
				w.Header().Set("Location", "/clusters/toke-toke-token")
				w.WriteHeader(http.StatusCreated)
			} else if r.URL.Path != "/clusters" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)

	return &TestHTTPMock{
        server: Server,
	}
}