package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"

	_ "github.com/yedamao/viper-remote-panic-demo/remote"
)

// newViper create a new viper and copy all settings of old viper
func newViper(old *viper.Viper) *viper.Viper {
	v := viper.New()
	v.MergeConfigMap(old.AllSettings())
	return v
}

func main() {
	// 用来读取配置的viper
	var readOnlyViper atomic.Value

	// 用来读取/接收配置中心数据的viper
	remoteViper := viper.New()
	remoteViper.SetConfigType("json")
	if err := remoteViper.AddRemoteProvider("etcd", "xxxxx", ""); err != nil {
		panic(err)
	}
	remoteViper.ReadRemoteConfig()

	// set read only viper
	readOnlyViper.Store(newViper(remoteViper))

	// open a goroutine to watch remote changes forever
	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := remoteViper.WatchRemoteConfig()
			if err != nil {
				fmt.Printf("unable to read remote config: %v", err)
				continue
			}

			// unmarshal new config into our runtime config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			// runtime_viper.Unmarshal(&runtime_conf)
			readOnlyViper.Store(newViper(remoteViper))
		}
	}()

	for i := 0; i < 3; i++ {
		go func() {
			for {
				v := readOnlyViper.Load().(*viper.Viper)
				v.Get("dummy")
			}
		}()
	}

	for {
		time.Sleep(time.Second)
		v := readOnlyViper.Load().(*viper.Viper)
		fmt.Println(v.AllKeys())
	}
}
