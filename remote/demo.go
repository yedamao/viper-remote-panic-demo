package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/spf13/viper"
)

type demoRemoteConfigProvider struct {
	sem    chan struct{} // 数据更新信号
	config map[string]interface{}
}

func (rc *demoRemoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	data, err := json.Marshal(rc.config)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

func (rc *demoRemoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	<-rc.sem
	return rc.Get(rp)
}

func (rc *demoRemoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	quit := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)

	go func(vr chan<- *viper.RemoteResponse) {
		for {
			select {
			case <-rc.sem:
				reader, err := rc.Get(rp)
				val, err := ioutil.ReadAll(reader)

				vr <- &viper.RemoteResponse{
					Error: err,
					Value: val,
				}
			}
		}
	}(viperResponsCh)
	return viperResponsCh, quit
}

// updateWorker 定期更新配置
func (rc *demoRemoteConfigProvider) updateWorker() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		ts := time.Now().Unix()
		rc.config[fmt.Sprintf("%d", ts)] = ts
		fmt.Println("config updated ...")
		rc.sem <- struct{}{}
	}
}

func New() *demoRemoteConfigProvider {
	provider := &demoRemoteConfigProvider{
		sem:    make(chan struct{}, 3),
		config: make(map[string]interface{}),
	}
	go provider.updateWorker()
	return provider
}

func init() {
	viper.RemoteConfig = New()
}
