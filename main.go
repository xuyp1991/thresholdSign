package main

import (
	"fmt"
	"github.com/xuyp1991/thresholdSign/ethutil"
)

const (
	strPwd = "password"
	testKeyFile = "/home/xuyapeng/go_workspace/src/github.com/xuyp1991/thresholdSign/testdata/keystore/UTC--2020-08-07T06-40-39.809100731Z--8f59d26233d58c03069f3d4169c4b8ade7edce7d"
)

func main() {
	key, err := ethutil.DecryptKeyFile(testKeyFile,strPwd)
	if  err != nil {
		fmt.Println("DecryptKeyFile error ",err)
	}
	fmt.Println("DecryptKeyFile key succ ",key)
}