package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	_ "github.com/yedamao/viper-remote-panic-demo/remote"
)

func main() {
	v := viper.New()
	v.SetConfigType("json")

	if err := v.AddRemoteProvider("etcd", "xxxxx", ""); err != nil {
		panic(err)
	}

	if err := v.WatchRemoteConfigOnChannel(); err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		go func() {
			for {
				v.Get("dummy")
			}
		}()
	}

	for {
		time.Sleep(time.Second)
		fmt.Println(v.AllKeys())
	}
}
