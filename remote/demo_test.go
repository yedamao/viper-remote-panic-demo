package remote

import (
	"fmt"
	"io/ioutil"
	"testing"
)

type defaultRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func (rp defaultRemoteProvider) Provider() string {
	return rp.provider
}

func (rp defaultRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp defaultRemoteProvider) Path() string {
	return rp.path
}

func (rp defaultRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}

func TestProvider(t *testing.T) {
	provider := New()
	reader, err := provider.Watch(&defaultRemoteProvider{})
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	respChan, _ := provider.WatchChannel(&defaultRemoteProvider{})
	for resp := range respChan {
		fmt.Println(string(resp.Value))
	}
}
