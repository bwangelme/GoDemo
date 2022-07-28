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

var (
	SERVICE_DISCOVER_START_TIMEOUT = 10 * time.Second
)

type ServiceDiscover struct {
	zkPath  string
	servers []string

	conn     *zk.Conn
	zkEvents <-chan zk.Event

	appMap map[string]int

	stop      chan struct{}
	cleanUpWG sync.WaitGroup
	err       chan error

	sync.Mutex
}

func (sd *ServiceDiscover) watchService(servicePath string) {
	// TODO
}

func (sd *ServiceDiscover) watchApp(appname string) {
	sd.Lock()
	if _, ok := sd.appMap[appname]; ok {
		sd.Unlock()
		return
	}
	sd.appMap[appname] = 1
	sd.Unlock()

	zkPath := fmt.Sprintf("%v/%v", sd.zkPath, appname)
	defer func() {
		sd.Lock()
		delete(sd.appMap, appname)
		sd.Unlock()

		logrus.Infof("Cancel watch app path %s", zkPath)
		sd.cleanUpWG.Done()
	}()

	sd.cleanUpWG.Add(1)
	logrus.Infof("begin to watch %s", appname)

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
				logrus.Infof("app %s is not exist, cancel watch", zkPath)
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
			servicePath := fmt.Sprintf("%v/%v/%v", sd.zkPath, appname, serviceName)
			go sd.watchService(servicePath)
		}

		select {
		case <-ch:
			logrus.Infof("%s changed, refresh watch info", zkPath)
			continue
		case <-sd.stop:
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
			appNames, _, ch, err = sd.conn.ChildrenW(sd.zkPath)
			if err == nil {
				// watch path 成功，跳出 retry loop
				break
			}

			if err == zk.ErrNoNode {
				sd.stop <- struct{}{}
				sd.err <- xerrors.Errorf("watch path %s not exist: %s", sd.zkPath, err)
				return
			}

			logrus.Errorf("watch path failed: %v, retry after 3 seconds", err)
			time.Sleep(3 * time.Second)
		}

		if err != nil {
			sd.stop <- struct{}{}
			sd.err <- xerrors.Errorf("watch path %s failed: %s", sd.zkPath, err)
			return
		}

		for _, appName := range appNames {
			go sd.watchApp(appName)
		}
		sd.err <- nil

		select {
		case <-ch:
			continue
		case <-sd.stop:
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
	sd.err <- nil

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
func (sd *ServiceDiscover) StartServer() error {
	workerCnt := 2

	go sd.startZKSentinel()
	go sd.watchAll()

	startTimeout := SERVICE_DISCOVER_START_TIMEOUT / time.Duration(workerCnt)
	for cnt := 0; cnt < workerCnt; cnt++ {
		select {
		case <-time.After(startTimeout):
			return xerrors.New("Start discover server timeout")
		case err := <-sd.err:
			if err != nil {
				return err
			}
		}
	}

	return nil
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
		zkPath:   pathPrefix,
		servers:  zkServers,
		zkEvents: ch,
		appMap:   make(map[string]int, 0),
		conn:     conn,
		err:      make(chan error),
		stop:     make(chan struct{}),
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
	err = sd.StartServer()
	if err != nil {
		logrus.Fatalf("start service discover server failed: %v", err)
	}
	logrus.Info("Start Service Discover")

	fmt.Println(sd)
	<-ch
}
