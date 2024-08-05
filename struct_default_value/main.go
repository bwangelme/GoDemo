package main

import (
	"fmt"
	"reflect"
)

type NoExchangeListReq struct {
	Uid              int64  `json:"uid" form:"uid" bindings:"required"`
	SystemId         string `json:"systemId" form:"systemId" bindings:"required"`
	ActivityTypeList []int  `json:"activityTypeList" form:"activityTypeList"`
	ExchangeTypeList []int  `json:"exchangeTypeList" form:"exchangeType"`
}

func main() {
	// ActivityTypeList 在创建结构体时没有传入，默认初始化成一个 int 数组
	req := NoExchangeListReq{
		Uid:              23,
		SystemId:         "homework",
		ExchangeTypeList: []int{1, 2},
	}
	fmt.Print(req.ActivityTypeList, reflect.TypeOf(req.ActivityTypeList))
}
