## 目的

本程序的目的是，检测 zk 中子节点添加/删除时，会输出什么事件

## ChildrenW 的特殊写法

`ChildrenW` 返回的第三个返回值是一个 channel, 当子节点的数据发生变化时，会向 channel 中写入 event

__注意:__ 此 channel 只能够用一次，在触发了一次后，必须重新调用 `ChildrenW`, 如果重复在此 channel 上接收时间，则会重复接收一个事件

```shell
{Unknown StateDisconnected  <nil> }
```

+ ChildrenW 示例代码

```go
for {
    sub, _, ch, err := conn.ChildrenW(path)
    if err != nil {
        log.Fatalf("watch %s failed %v\n", path, err)
    }
    
    select {
    case e := <-ch:
        //TODO
    }
}
```

## 子节点的操作以及对应的 Event

| 操作     | 事件                       |
|--------|--------------------------|
| 添加子节点  | EventNodeChildrenChanged |
| 删除子节点  | EventNodeChildrenChanged |
| 修改子节点  | 无                        |
| 删除节点自身 | EventNodeDeleted         |