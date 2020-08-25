package main

import (
	"fmt"
	"github.com/xuyp1991/thresholdSign/ethutil"
	"github.com/xuyp1991/thresholdSign/tss"
)

const (
	strPwd = "password"
	testKeyFile1 = "/home/xuyapeng/go_workspace/src/github.com/xuyp1991/thresholdSign/testdata/keystore/UTC--2020-08-07T06-40-39.809100731Z--8f59d26233d58c03069f3d4169c4b8ade7edce7d"
	testKeyFile2 = "/home/xuyapeng/go_workspace/src/github.com/xuyp1991/thresholdSign/testdata/keystore/UTC--2020-08-10T07-56-58.217384982Z--72bfd5482329c670ef3465f1a1bde2265d3d375f"
	testKeyFile3 = "/home/xuyapeng/go_workspace/src/github.com/xuyp1991/thresholdSign/testdata/keystore/UTC--2020-08-10T07-57-17.817005235Z--bb242d23b282301a84b3597278fa49f6061a1bee"
)

func main() {
	key1, err := ethutil.DecryptKeyFile(testKeyFile1,strPwd)
	if  err != nil {
		fmt.Println("DecryptKeyFile error ",err)
	}
	fmt.Println("DecryptKeyFile key succ ",key1)
	key2, err := ethutil.DecryptKeyFile(testKeyFile2,strPwd)
	if  err != nil {
		fmt.Println("DecryptKeyFile error ",err)
	}
	fmt.Println("DecryptKeyFile key succ ",key2)
	key3, err := ethutil.DecryptKeyFile(testKeyFile3,strPwd)
	if  err != nil {
		fmt.Println("DecryptKeyFile error ",err)
	}
	fmt.Println("DecryptKeyFile key succ ",key3)

	tss.GenerateGroup(key1,key2,key3)
	// _, operatorPublicKey1 := operator.EthereumKeyToOperatorKey(key1)
	// _, operatorPublicKey2 := operator.EthereumKeyToOperatorKey(key1)
	// _, operatorPublicKey3 := operator.EthereumKeyToOperatorKey(key1)

	// memberID1 := tss.MemberIDFromPublicKey(operatorPublicKey1)
	// memberID2 := tss.MemberIDFromPublicKey(operatorPublicKey2)
	// memberID3 := tss.MemberIDFromPublicKey(operatorPublicKey3)

	// memberIDs := make([]tss.MemberID, 0)
	// memberIDs = append(memberIDs, memberID1)
	// memberIDs = append(memberIDs, memberID2)
	// memberIDs = append(memberIDs, memberID3)

	// groupInfo := tss.GroupInfo{
	// 	GroupID:"justfortest",
	// 	MemberID:memberID1,
	// 	GroupMemberIDs:memberIDs,
	// 	DishonestThreshold:3,
	// }
	// fmt.Println("test groupinfo ",groupInfo)
}