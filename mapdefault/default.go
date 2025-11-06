package main

import "fmt"

type AnchorLevelBriefInfo struct {
	LevelName string `protobuf:"bytes,2,opt,name=level_name,json=levelName,proto3" json:"level_name,omitempty"` // 等级名字
}

func (a *AnchorLevelBriefInfo) GetLevelName() string {
	if a == nil {
		return ""
	}
	return a.LevelName
}

func main() {
	var m1 map[uint]*AnchorLevelBriefInfo = nil
	var m2 = make(map[uint64]*AnchorLevelBriefInfo)

	fmt.Println("m1 get", m1[1], m1[1].GetLevelName())
	fmt.Println("m2 get", m2[1], m2[1].GetLevelName())
}
