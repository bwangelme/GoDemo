package main

import (
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"time"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

type Controller struct {
	// 本地缓存对象，从k8s集群通过来的数据缓存在这里
	indexer cache.Indexer
	// 工作队列，用于解耦informer与controller之间的强关联
	// 为什么需要一个工作队列呢？我猜测注册的函数会以同步的方式运行
	// 如果监听函数执行时间过长会对informer的主业务逻辑造成影响
	// 但是我还没看informer的代码，所以只是猜测
	queue workqueue.RateLimitingInterface
	// informer对象，用于k8s集群状态的数据到本地缓存对象
	informer cache.Controller
}

// New*属于go代码风格，因为没有构造函数
func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

// 注册到队列的构造函数
// 每次有新的数据进入队列，都会执行
func (c *Controller) processNextItem() bool {
	// 从队列中拿出一个key, 这个key是字符串，通常是"{namespace}-{name}"的形式
	// 工作队列在shutdown之后，quit的值会是true
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// 告诉工作队列我们已经处理过这个key了，这也是一种解锁的方法，当一个key从队列中取出来之后就会加上锁，那么其他监听函数就会同时处理这个key，这保证了队列的并行安全
	defer c.queue.Done(key)

	// 这里放controller的主要业务逻辑
	// 将key交由controller业务代码处理
	err := c.syncToStdout(key.(string))
	// 处理错误的主要入口
	c.handleErr(err, key)
	return true
}

// 因为是演示作用，所以这个函数的业务逻辑就是简单的将pod的相关信息打印到标准输出
// 如果出现了错误也不处理，原样抛到上一层即可
func (c *Controller) syncToStdout(key string) error {
	// 通过key获取最新的obj,这里是pod对象
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		fmt.Printf("Pod %s does not exist anymore\n", key)
	} else {
		fmt.Printf("Sync/Add/Update for Pod %s\n", obj.(*v1.Pod).GetName())
	}
	return nil
}

// 错误处理函数
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// 为什么既要Done又要Forget呢？
		// Done方法告诉队列，这个数据项已经处理过了，但是是否应该移除，它不管，并且每个key在被调用前，这个数据项不能同时被其他进程处理
		// Forget方法告诉队列，这个数据项忘了吧，无论是已经用完了还是没办法处理，队列会把这个数据项从队列中移除
		c.queue.Forget(key)
		return
	}

	// 看看这个key重新入队了多少次，超过五次就不在放进队列了
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)

		// 再次放入队列中，等待下次调用
		// AddRateLimited会控制队列的流速，控制并发数量
		c.queue.AddRateLimited(key)
		return
	}

	// 丢弃这个数据项
	c.queue.Forget(key)
	// 实在处理不了，把错误往上抛
	runtime.HandleError(err)
	klog.Infof("Dropping pod %q out of the queue: %v", key, err)
}

func (c *Controller) Run(workers int, stopCh chan struct{}) {
	// 捕获未知错误
	defer runtime.HandleCrash()

	// 关闭工作队列，工作进程在get的时候会判断工作队列关闭没有
	defer c.queue.ShutDown()

	// 启动informer
	// informer通常作为一个独立的gorouting运行
	go c.informer.Run(stopCh)

	// 等待第一此数据同步完成，全量数据同步
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	// 启动n个工作进程
	for i := 0; i < workers; i++ {
		// 直到stopCh收到退出信号才退出，每秒运行runWorker
		// 因为runWorker其实是一个for的死循环，其实会一直等待runWorker运行完成
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	// 等待关闭信号
	<-stopCh
	klog.Info("Stopping Pod controller")
}

// 工作进程
func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func main() {
	var kubeconfig string
	var master string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.Parse()
	config, _ := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	clientset, _ := kubernetes.NewForConfig(config)
	// 上面就是为了创建可以与k8s交互的客户端

	//配置监听器监听的对象及参数
	// 这监听默认namespace的pod资源
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())

	// 创建一个带流速控制的工作队列
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// 正常情况下，一般创建共享的informer，即SharedInformer，共享的informer可以保证同一主进程下的子进程复用底层连接及资源，比如k8s的controller-manager里面有多个controller，比如deploymentController，daemonsetController,他们都需要监听pod资源，如果每个controller都使用独立的informer就太浪费资源了，这里的代码是为了简单
	// 一个informer可以注册AddFunc，UpdateFunc，DeleteFunc等回调函数，当对应的监听事件触发之后就会调用。
	// 这里的代码很简单，当事件了之后将放入controller的工作队列
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// 基于pod对象创建一个key，格式是{namespace}-{name}
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	controller := NewController(queue, indexer, informer)

	// 手动创建一个pod
	indexer.Add(&v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "mypod",
			Namespace: v1.NamespaceDefault,
		},
	})

	// 启动
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// 一直等待
	select {}
}
