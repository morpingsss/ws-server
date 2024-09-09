package etcd

import (
	"context"
	"time"

	"github.com/morpingsss/ws-server/pkg/setting"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ClientDis struct {
	client *clientv3.Client
}

func NewClientDis(addr []string) (*ClientDis, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	return &ClientDis{
		client: client,
	}, nil
}

func (c *ClientDis) GetService(prefix string) ([]string, error) {
	resp, err := c.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	addrs := c.extractAddrs(resp)

	go c.watcher(prefix)
	return addrs, nil
}

func (c *ClientDis) watcher(prefix string) {
	rch := c.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				c.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				c.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (c *ClientDis) extractAddrs(resp *clientv3.GetResponse) []string {
	addrs := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return addrs
	}
	for _, kv := range resp.Kvs {
		if v := kv.Value; v != nil {
			c.SetServiceList(string(kv.Key), string(v))
			addrs = append(addrs, string(v))
		}
	}
	return addrs
}

func (c *ClientDis) SetServiceList(key, val string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	setting.GlobalSetting.ServerList[key] = val
	log.Info("发现服务：", key, " 地址:", val)
}

func (c *ClientDis) DelServiceList(key string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	delete(setting.GlobalSetting.ServerList, key)
	log.Println("服务下线:", key)
}
