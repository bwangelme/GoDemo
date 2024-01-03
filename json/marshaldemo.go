package main

import (
	"bwdemo/wtypes"
	"encoding/json"
	"fmt"
	"log"
)

// 测试自定义类型的 Marshaler 接口和  UnMarshaler 接口
func main() {
	var args = struct {
		IsTv  wtypes.Bool `json:"is_tv,omitempty"`
		Empty wtypes.Bool `json:"empty,omitempty"`
		// Bool 类型应该设置上 ,omitempty ，否则在遇到空值时程序会报错
		//ErrEmpty wtypes.Bool `json:"err_empty"`
	}{
		IsTv:  wtypes.True,
		Empty: wtypes.Nil,
		//ErrEmpty: wtypes.Nil,
	}

	// 注意 Marshal 和 UnMarshal 函数只接受一个指针参数，不能传入值类型
	data, err := json.Marshal(&args)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))

	var newArgs struct {
		IsTv  wtypes.Bool `json:"is_tv"`
		Empty wtypes.Bool `json:"empty"`
	}
	// empty 值必须是 false 或 true，否则就会报错
	//var errData = []byte(`{"is_tv":true, "empty": 233}`)
	var newData = []byte(`{"is_tv":true}`)
	err = json.Unmarshal(newData, &newArgs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(newArgs)
}
