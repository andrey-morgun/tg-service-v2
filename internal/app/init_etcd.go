package app

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func (a *App) initEtcd() {
	var err error
	a.etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{a.config.Etcd.ConnString},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		a.logger.Fatal(err.Error())
		return
	}
}
