package user_maps

import (
	"context"
	"encoding/json"
	"github.com/andReyM228/lib/errs"
	"github.com/andReyM228/lib/log"
	"go.etcd.io/etcd/client/v3"
)

type Service struct {
	etcdClient *clientv3.Client
	log        log.Logger
}

func NewService(etcdClient *clientv3.Client, log log.Logger) Service {
	return Service{
		etcdClient: etcdClient,
		log:        log,
	}
}

func (s Service) Put(ctx context.Context, key string, values ...interface{}) error {
	value := key
	if len(values) != 0 {
		valueBytes, err := json.Marshal(values[0])
		if err != nil {
			return err
		}

		value = string(valueBytes)
	}

	_, err := s.etcdClient.Put(ctx, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Fetch(ctx context.Context, key string, receiver ...interface{}) bool {
	resp, err := s.etcdClient.Get(ctx, key)
	if err != nil {
		s.log.Error(err.Error())
		return false
	}

	if len(resp.Kvs) == 0 {
		s.log.Errorf("%s", errs.NotFoundError{What: key})
		return false
	}

	if len(receiver) == 0 {
		return true
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &receiver[0]); err != nil {
		s.log.Error(err.Error())
		return false
	}

	return true
}

func (s Service) Delete(ctx context.Context, key string) error {
	_, err := s.etcdClient.Delete(ctx, key)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}
