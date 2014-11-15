package grobot

import "net/http"

func init() {
	HttpClientProvider = &defaultHttpClient{http.DefaultClient}
}

type HttpClient interface {
	Send(*http.Request) (*http.Response, error)
}

var HttpClientProvider HttpClient

type defaultHttpClient struct {
	client *http.Client
}

func (c *defaultHttpClient) Send(request *http.Request) (*http.Response, error) {
	return c.client.Do(request)
}

func GetHTTP(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return HttpClientProvider.Send(request)
}
