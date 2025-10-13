package xc

import (
	"fmt"
	"testing"
)

func Test_GetDll(t *testing.T) {
	// 获取不存在的文件, 看输出是什么
	bs, err := GetDll("https://cnb.cool/twgh521/xcguidll/-/releases/download/3.3.3.1/xcgui.dll")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bs))
}

func Test_GetLatestVersion(t *testing.T) {
	fmt.Println(GetLatestVersion())
}
