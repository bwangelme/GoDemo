/*
 这个文件实现了 ServiceDiscover 服务

 它从 zk 读取 app 和 service 信息，将变化的数据通过 channel 返回给调用方

 zk key 的格式: `/xiaohei/services/rpc_client/{appname}/{service_list}`
*/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-zookeeper/zk"
	xerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ServiceDiscover struct {
	basicPath   string
	servers     []string
	connTimeout time.Duration

	conn     *zk.Conn
	zkEvents <-chan zk.Event

	appMap     map[string]uint8
	serviceMap map[string]uint8
	cancel     context.CancelFunc

	cleanUpWG *sync.WaitGroup

	sync.Mutex
}

func (sd *ServiceDiscover) watchService(ctx context.Context, appName, serviceName string) {
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
		case <-ctx.Done():
			logrus.Infof("Stop goroutine watch service %v", servicePath)
			return
		}
	}
}

func (sd *ServiceDiscover) Close() {
	sd.cancel()
	sd.cleanUpWG.Wait()
}

func (sd *ServiceDiscover) watchApp(ctx context.Context, appname string) {
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
			go sd.watchService(ctx, appname, serviceName)
		}

		select {
		case <-ch:
			logrus.Debugf("[app] %s changed, refresh watch info", zkPath)
			continue
		case <-ctx.Done():
			logrus.Infof("Stop goroutine watchapp %v", appname)
			return
		}

	}
}

func (sd *ServiceDiscover) watchAll(ctx context.Context) {
	var err error
	var ch <-chan zk.Event
	var appNames []string

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
				logrus.Errorf("watch path %s not exist: %s", sd.basicPath, err)
				sd.cancel()
				return
			}

			logrus.Errorf("[all] watch path failed: %v, retry after 3 seconds", err)
			time.Sleep(3 * time.Second)
		}

		if err != nil {
			logrus.Errorf("[all] watch path %s failed: %s", sd.basicPath, err)
			sd.cancel()
			return
		}

		for _, appName := range appNames {
			go sd.watchApp(ctx, appName)
		}

		select {
		case <-ch:
			logrus.Debugf("[all] %s changed, refresh watch info", sd.basicPath)
			continue
		case <-ctx.Done():
			logrus.Info("[all] stop watch all goroutine")
			return
		}
	}
}

//startZKSentinel
// 启动一个刷新 zk 连接的哨兵
func (sd *ServiceDiscover) startZKSentinel(ctx context.Context) {
	logrus.Info("Start ZK Sentinel")

	defer func() {
		sd.cleanUpWG.Done()
	}()

	for {
		select {
		case event := <-sd.zkEvents:
			if event.Type == zk.EventSession && event.State == zk.StateExpired {
				logrus.Infof("zk expire, start refresh connection")
				newConn, newEvents, err := connectToZK(sd.servers, sd.connTimeout, 60)
				if err != nil {
					logrus.Errorf("Connect to zk failed: %v", err)
					sd.cancel()
				}

				oldConn := sd.conn
				sd.zkEvents, sd.conn = newEvents, newConn
				oldConn.Close()
				logrus.Infof("success create new zk connection")
			}
		case <-ctx.Done():
			logrus.Info("Stop zk sentinel")
			return
		}
	}
}

//StartServer
// 启动服务发现 Server
func (sd *ServiceDiscover) StartServer(ctx context.Context) {
	sd.cleanUpWG.Add(2)

	go sd.startZKSentinel(ctx)
	go sd.watchAll(ctx)
}

func connectToZK(servers []string, connTimeout time.Duration, retry int) (*zk.Conn, <-chan zk.Event, error) {
	var conn *zk.Conn
	var ch <-chan zk.Event
	var err error

	for i := 0; i < retry; i++ {
		conn, ch, err = zk.Connect(servers, connTimeout)
		if err == nil {
			return conn, ch, err
		}

		if err != nil {
			logrus.Infof("Connect to zk %v failed, i after 1 seconds", servers)
			time.Sleep(1 * time.Second)
		}
	}

	return nil, nil, xerrors.WithMessagef(err, "connect to %s failed", servers)
}

func newServiceDiscover(zkServers []string, pathPrefix string, cancel context.CancelFunc) (*ServiceDiscover, error) {
	sd := &ServiceDiscover{
		basicPath:   pathPrefix,
		servers:     zkServers,
		connTimeout: 10 * time.Second,
		cancel:      cancel,
		appMap:      make(map[string]uint8, 0),
		serviceMap:  make(map[string]uint8, 0),
		cleanUpWG:   &sync.WaitGroup{},
	}

	conn, ch, err := connectToZK(sd.servers, sd.connTimeout, 3)
	if err != nil {
		return nil, err
	}

	sd.conn = conn
	sd.zkEvents = ch

	return sd, nil
}

func NewServiceDiscover(ctx context.Context, zkServers []string, pathPrefix string, cancel context.CancelFunc) (*ServiceDiscover, error) {
	sd, err := newServiceDiscover(zkServers, pathPrefix, cancel)
	if err != nil {
		return nil, err
	}
	return sd, nil
}

func main() {
	var sigC = make(chan os.Signal, 0)
	signal.Notify(sigC, syscall.SIGTERM)
	var ctx, cancel = context.WithCancel(context.Background())

	sd, err := NewServiceDiscover(ctx, []string{"127.0.0.1:2181"}, "/xiaohei/services/rpc_client", cancel)
	if err != nil {
		logrus.Fatalf("init service discover failed: %v", err)
	}
	sd.StartServer(ctx)

	pid := os.Getpid()
	logrus.Infof("Start Service Discover on %v", pid)

	for {
		select {
		case sig := <-sigC:
			if sig == syscall.SIGTERM {
				logrus.Info("Receive term signal, stop discover")
				sd.Close()
				return
			}
		}
	}
}
