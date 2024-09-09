package etcd

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 注册租约服务
type ServiceReg struct {
	client        *clientv3.Client                        // etcd客户端实例
	lease         clientv3.Lease                          // 租约实例
	leaseResp     *clientv3.LeaseGrantResponse            // 租约响应，包含租约ID等信息
	cancelFunc    func()                                  // 取消续租的函数
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 续租响应通道，用于接收续租结果
}

func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
	var (
		err    error
		client *clientv3.Client
	)

	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}); err != nil {
		return nil, err
	}

	ser := &ServiceReg{
		client: client,
	}

	if err := ser.setLease(timeNum); err != nil {
		return nil, err
	}
	go ser.ListenLeaseRespChan()
	return ser, nil
}

// 设置租约
func (s *ServiceReg) setLease(timeNum int64) error {
	lease := clientv3.NewLease(s.client)

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	leaseResp, err := lease.Grant(ctx, timeNum)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)

	if err != nil {
		return err
	}

	s.lease = lease
	s.leaseResp = leaseResp
	s.cancelFunc = cancelFunc
	s.keepAliveChan = leaseRespChan
	return nil
}

// 监听续租情况
func (s *ServiceReg) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-s.keepAliveChan:
			if leaseKeepResp == nil {
				log.Error("已经关闭续租功能")
				return
			} else {
				// log.Info("续租成功")
			}
		}
	}
}

// 注册租约
func (s *ServiceReg) PutService(key, val string) error {
	kv := clientv3.NewKV(s.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(s.leaseResp.ID))
	return err
}

// 撤销租约
func (s *ServiceReg) RevokeLease() error {
	s.cancelFunc()
	time.Sleep(2 * time.Second)
	_, err := s.lease.Revoke(context.TODO(), s.leaseResp.ID)
	return err
}
