/*
 这个文件实现了 ServiceDiscover 服务

 它从 zk 读取 app 和 service 信息，将变化的数据通过 channel 返回给调用方

 zk key 的格式: `/xiaohei/services/rpc_client/{appname}/{service_list}`
*/
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
	xerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ServiceDiscover struct {
	basicPath string
	servers   []string

	conn     *zk.Conn
	zkEvents <-chan zk.Event

	appMap     map[string]uint8
	serviceMap map[string]uint8

	stop      chan struct{}
	cleanUpWG sync.WaitGroup

	sync.Mutex
}

func (sd *ServiceDiscover) watchService(appName, serviceName string) {
	servicePath := fmt.Sprintf("%v/%v", appName, serviceName)

	sd.Lock()
	if _, ok := sd.serviceMap[servicePath]; ok {
		sd.Unlock()
		return
	}
	sd.serviceMap[servicePath] = 1
	sd.Unlock()

	zkPath := fmt.Sprintf("%v/%v", sd.basicPath, servicePath)
	defer func() {
		sd.Lock()
		delete(sd.serviceMap, servicePath)
		sd.Unlock()

		logrus.Infof("[service] Cancel watch service path %s", zkPath)
		sd.cleanUpWG.Done()
	}()

	sd.cleanUpWG.Add(1)
	logrus.Infof("[service] begin to watch service %s", servicePath)

	for {
		var ch <-chan zk.Event
		var err error
		var entrypoints []string

		for retry := 0; retry < 3; retry++ {
			entrypoints, _, ch, err = sd.conn.ChildrenW(zkPath)
			if err == nil {
				break
			}

			if err == zk.ErrNoNode {
				logrus.Infof("[service] %s not exist, cancel watch", zkPath)
				return
			}

			logrus.Errorf("watch app failed: %s, retry after 3 seconds", err)
			time.Sleep(3 * time.Second)
		}

		if err != nil {
			logrus.Errorf("watch service failed %s, stop watching", zkPath)
			return
		}

		for _, entrypoint := range entrypoints {
			zkPath := fmt.Sprintf("%v/%v/%v/%v", sd.basicPath, appName, serviceName, entrypoint)
			data, _, err := sd.conn.Get(zkPath)
			if err != nil {
				logrus.Errorf("Get %v data failed: %v", zkPath, err)
			}
			fmt.Println(appName, serviceName, entrypoint)
			fmt.Println(string(data))
		}

		select {
		case <-ch:
			logrus.Debugf("[service] %s changed, refresh watch info", zkPath)
			continue
		case <-sd.stop:
			logrus.Infof("Stop goroutine watch service %v", servicePath)
			return
		}
	}
}

func (sd *ServiceDiscover) watchApp(appname string) {
	sd.Lock()
	if _, ok := sd.appMap[appname]; ok {
		sd.Unlock()
		return
	}
	sd.appMap[appname] = 1
	sd.Unlock()

	zkPath := fmt.Sprintf("%v/%v", sd.basicPath, appname)
	defer func() {
		sd.Lock()
		delete(sd.appMap, appname)
		sd.Unlock()

		logrus.Infof("Cancel watch app %s", appname)
		sd.cleanUpWG.Done()
	}()

	sd.cleanUpWG.Add(1)
	logrus.Infof("[app] begin to watch %s", appname)

	for {
		var ch <-chan zk.Event
		var err error
		var serviceNames []string

		for retry := 0; retry < 3; retry++ {
			serviceNames, _, ch, err = sd.conn.ChildrenW(zkPath)
			if err == nil {
				break
			}

			if err == zk.ErrNoNode {
				logrus.Infof("[app] %s is not exist, cancel watch", zkPath)
				return
			}
			logrus.Errorf("watch app failed: %s", err)
			time.Sleep(3 * time.Second)
		}

		if err != nil {
			logrus.Errorf("watch app failed: %s, stop watching", zkPath)
			return
		}

		for _, serviceName := range serviceNames {
			go sd.watchService(appname, serviceName)
		}

		select {
		case <-ch:
			logrus.Debugf("[app] %s changed, refresh watch info", zkPath)
			continue
		case <-sd.stop:
			logrus.Infof("Stop goroutine watchapp %v", appname)
			return
		}

	}
}

func (sd *ServiceDiscover) watchAll() {
	var err error
	var ch <-chan zk.Event
	var appNames []string
	sd.cleanUpWG.Add(1)

	defer func() {
		sd.cleanUpWG.Done()
	}()

	for {
		for retry := 0; retry < 3; retry++ {
			appNames, _, ch, err = sd.conn.ChildrenW(sd.basicPath)
			if err == nil {
				// watch path 成功，跳出 retry loop
				break
			}

			if err == zk.ErrNoNode {
				sd.stop <- struct{}{}
				logrus.Fatalf("watch path %s not exist: %s", sd.basicPath, err)
				return
			}

			logrus.Errorf("[all] watch path failed: %v, retry after 3 seconds", err)
			time.Sleep(3 * time.Second)
		}

		if err != nil {
			sd.stop <- struct{}{}
			logrus.Fatalf("[all] watch path %s failed: %s", sd.basicPath, err)
			return
		}

		for _, appName := range appNames {
			go sd.watchApp(appName)
		}

		select {
		case <-ch:
			logrus.Debugf("[all] %s changed, refresh watch info", sd.basicPath)
			continue
		case <-sd.stop:
			logrus.Info("[all] stop watchall goroutine")
			return
		}

	}
}

//startZKSentinel
// 启动一个刷新 zk 连接的哨兵
func (sd *ServiceDiscover) startZKSentinel() {
	// TODO
	sd.cleanUpWG.Add(1)
	logrus.Info("Start ZK Sentinel")

	for {
		select {
		case <-sd.stop:
			logrus.Info("Stop zk sentinel")
			return
		}
	}
	sd.cleanUpWG.Done()
}

//StartServer
// 启动服务发现 Server
func (sd *ServiceDiscover) StartServer() {
	go sd.startZKSentinel()
	go sd.watchAll()
}

func newServiceDiscover(zkServers []string, pathPrefix string) (*ServiceDiscover, error) {
	var conn *zk.Conn
	var ch <-chan zk.Event
	var err error
	for retry := 0; retry < 3; retry++ {
		conn, ch, err = zk.Connect(zkServers, time.Second)
		if err == nil {
			break
		}

		if err != nil {
			logrus.Infof("Connect to zk %v failed, retry after 1 seconds", zkServers)
			time.Sleep(1 * time.Second)
		}
	}
	if err != nil {
		return nil, xerrors.WithMessagef(err, "connect to %s failed", zkServers)
	}

	return &ServiceDiscover{
		basicPath:  pathPrefix,
		servers:    zkServers,
		zkEvents:   ch,
		appMap:     make(map[string]uint8, 0),
		serviceMap: make(map[string]uint8, 0),
		conn:       conn,
		stop:       make(chan struct{}),
	}, nil
}

func NewServiceDiscover(zkServers []string, pathPrefix string) (*ServiceDiscover, error) {
	sd, err := newServiceDiscover(zkServers, pathPrefix)
	if err != nil {
		return nil, err
	}
	return sd, nil
}

func main() {
	var ch = make(chan struct{}, 0)

	sd, err := NewServiceDiscover([]string{"127.0.0.1:2181"}, "/xiaohei/services/rpc_client")
	if err != nil {
		logrus.Fatalf("init service discover failed: %v", err)
	}
	sd.StartServer()
	logrus.Info("Start Service Discover")

	// TODO, TERM信号停掉服务
	<-ch
}
