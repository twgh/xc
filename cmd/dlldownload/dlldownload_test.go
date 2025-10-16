package dlldownload

import (
	"fmt"
	"testing"
)

func Test_getDll(t *testing.T) {
	// 获取不存在的文件, 看输出是什么
	bs, err := getDll("https://cnb.cool/twgh521/xcguidll/-/releases/download/3.3.3.1/xcgui.dll")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bs))
}

func Test_getLatestVersion(t *testing.T) {
	fmt.Println(getLatestVersion())
}
