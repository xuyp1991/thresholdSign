package tss

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/xuyp1991/thresholdSign/operator"
	tssKeygen "github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/binance-chain/tss-lib/tss"
	"time"
	"math/big"
	"context"
	"github.com/keep-network/keep-ecdsa/pkg/ecdsa"

)

const (
		defaultPreParamsGenerationTimeout = 1 * time.Minute
)
//两个chan 一个获取party,一个获取message
var (
	Partychan = make(chan tss.Party)
	tssGlobalMessageChan = make(chan tss.Message)
)

func generatePartiesIDs(
	thisMemberID MemberID,
	groupMemberIDs []MemberID,
) (
	*tss.PartyID,
	[]*tss.PartyID,
	error,
) {
	var thisPartyID *tss.PartyID
	groupPartiesIDs := []*tss.PartyID{}

	for _, memberID := range groupMemberIDs {
		if memberID.bigInt().Cmp(big.NewInt(0)) <= 0 {
			return nil, nil, fmt.Errorf("member ID must be greater than 0, but found [%v]", memberID.bigInt())
		}

		newPartyID := tss.NewPartyID(
			memberID.String(), // id - unique string representing this party in the network
			"",                // moniker - can be anything (even left blank)
			memberID.bigInt(), // key - unique identifying key
		)

		if thisMemberID.Equal(memberID) {
			thisPartyID = newPartyID
		}

		groupPartiesIDs = append(groupPartiesIDs, newPartyID)
	}

	return thisPartyID, groupPartiesIDs, nil
}

func GenerateGroup(key1,key2,key3 *keystore.Key) {
	_, operatorPublicKey1 := operator.EthereumKeyToOperatorKey(key1)
	_, operatorPublicKey2 := operator.EthereumKeyToOperatorKey(key2)
	_, operatorPublicKey3 := operator.EthereumKeyToOperatorKey(key3)

	memberID1 := MemberIDFromPublicKey(operatorPublicKey1)
	memberID2 := MemberIDFromPublicKey(operatorPublicKey2)
	memberID3 := MemberIDFromPublicKey(operatorPublicKey3)

	memberIDs := make([]MemberID, 0)
	memberIDs = append(memberIDs, memberID1)
	memberIDs = append(memberIDs, memberID2)
	memberIDs = append(memberIDs, memberID3)

	groupInfo := GroupInfo{
		GroupID:"justfortest",
		MemberID:memberID1,
		GroupMemberIDs:memberIDs,
		DishonestThreshold:3,
	}

//	GenerateParty(groupInfo)

	//----------------------test   需要使用同一个message chan--------------------------------------
		groupInfo1 := GroupInfo{
		GroupID:"justfortest",
		MemberID:memberID2,
		GroupMemberIDs:memberIDs,
		DishonestThreshold:3,
	}
	groupInfo2 := GroupInfo{
		GroupID:"justfortest",
		MemberID:memberID3,
		GroupMemberIDs:memberIDs,
		DishonestThreshold:3,
	}
	GenerateParty(groupInfo,groupInfo1,groupInfo2)
}

func GenerateParty(groupInfo1,groupInfo2,groupInfo3 GroupInfo) {


	_, groupPartiesIDs, err := generatePartiesIDs(
		groupInfo1.MemberID,
		groupInfo1.GroupMemberIDs,
	)
	if err != nil {
		fmt.Println("generatePartiesIDs error")
		return 
	}
	sortPartIDS := tss.SortPartyIDs(groupPartiesIDs)
	go manageParty(sortPartIDS)
	go generateforparty(groupInfo1)
	go generateforparty(groupInfo2)
	go generateforparty(groupInfo3)
	time.Sleep(time.Duration(2100)*time.Second)
	// tssMessageChan := make(chan tss.Message, len(groupInfo1.GroupMemberIDs))
	// endChan := make(chan tssKeygen.LocalPartySaveData)
	
	// currentPartyID, groupPartiesIDs, err := generatePartiesIDs(
	// 	groupInfo1.MemberID,
	// 	groupInfo1.GroupMemberIDs,
	// )
	// if err != nil {
	// 	fmt.Println("generatePartiesIDs error")
	// 	return 
	// }

	// params := tss.NewParameters(
	// 	tss.NewPeerContext(tss.SortPartyIDs(groupPartiesIDs)),
	// 	currentPartyID,
	// 	len(groupPartiesIDs),
	// 	groupInfo1.DishonestThreshold,
	// )

	// party := tssKeygen.NewLocalParty(params, tssMessageChan, endChan, *tssPreParams)
	// fmt.Println("test  GenerateParty  ",party)
	// generateKey(party,endChan,tssMessageChan)
}

func generateforparty(groupInfo GroupInfo) {
	tssPreParams, err := tssKeygen.GeneratePreParams(defaultPreParamsGenerationTimeout)
	if err != nil {
		fmt.Println("GeneratePreParams error")
		return 
	}

	tssMessageChan := make(chan tss.Message, len(groupInfo.GroupMemberIDs))
	endChan := make(chan tssKeygen.LocalPartySaveData)

	currentPartyID, groupPartiesIDs, err := generatePartiesIDs(
		groupInfo.MemberID,
		groupInfo.GroupMemberIDs,
	)
	if err != nil {
		fmt.Println("generatePartiesIDs error")
		return 
	}
	sortPartIDS := tss.SortPartyIDs(groupPartiesIDs)
	params := tss.NewParameters(
		tss.NewPeerContext(sortPartIDS),
		currentPartyID,
		len(groupPartiesIDs),
		groupInfo.DishonestThreshold,
	)

	party := tssKeygen.NewLocalParty(params, tssMessageChan, endChan, *tssPreParams)
	Partychan<- party
	time.Sleep(time.Duration(15)*time.Second)
	generateKey(party,endChan,tssMessageChan,sortPartIDS)
}

// generateKey executes the protocol to generate a signing key. This function
// needs to be executed only after all members finished the initialization stage.
// As a result it will return a Signer who has completed key generation, or error
// if the key generation failed.
func  generateKey(party tss.Party ,endChan <-chan tssKeygen.LocalPartySaveData,tssMessageChan <-chan tss.Message,sortedPartyIDs tss.SortedPartyIDs)  {//(*ThresholdSigner, error)
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)
	defer cancel()
	if err := party.Start(); err != nil {
		// return nil, fmt.Errorf(
		// 	"failed to start key generation: [%v]",
		// 	party.WrapError(err),
		// )
		fmt.Println("start error ",err)
		return 
	}
	for {
		select {
		case keygenData := <-endChan:
			// signer := &ThresholdSigner{
			// 	groupInfo:    s.groupInfo,
			// 	thresholdKey: ThresholdKey(keygenData),
			// }
			fmt.Println("get keygenData ",keygenData)
			pkX, pkY := keygenData.ECDSAPub.X(), keygenData.ECDSAPub.Y()

			curve := tss.EC()
			publicKey := ecdsa.PublicKey{
				Curve: curve,
				X:     pkX,
				Y:     pkY,
			}
			fmt.Println("get publicKey ",publicKey)
		//	return (*ecdsa.PublicKey)(&publicKey)
			return //signer, nil
		case tmp := <-tssMessageChan://如何把message路由给其他的chan?
		//	fmt.Println("get message ",tmp)
			 _, routing, _ := tmp.WireBytes()
			 senderPartyID := sortedPartyIDs.FindByKey(MemberID(routing.From.GetKey()).bigInt())

			// _, err := party.UpdateFromBytes(
			// 	bytes,
			// 	senderPartyID,
			// 	true,
			// )
			// if err != nil {
			// 	fmt.Println("UpdateFromBytes error ",err)
			// 	continue
			// }
			if senderPartyID == party.PartyID() {
				tssGlobalMessageChan<-tmp
//				fmt.Println("senderPartyID equal partyID ",senderPartyID)
			}
			// fmt.Println("UpdateFromBytes succ ")
			continue
		case <-ctx.Done():
			memberIDs := []MemberID{}

			if party.WaitingFor() != nil {
				for _, partyID := range party.WaitingFor() {
					memberID, err := MemberIDFromString(partyID.GetId())
					if err != nil {
						fmt.Println("get error ",err)
						continue
					}
					fmt.Println("waiting for ",partyID,"---",memberID,"---",party.String())
					memberIDs = append(memberIDs, memberID)
				}
			}
			fmt.Println("done return ")
			return //nil, timeoutError{KeyGenerationProtocolTimeout, "key generation", memberIDs}
		}
	}
}

func manageParty(sortedPartyIDs tss.SortedPartyIDs) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Second)
	defer cancel()
	AllPartys := []tss.Party{}
	for {
		select {
			case newParty := <-Partychan:
				AllPartys = append(AllPartys,newParty)
				fmt.Println("test new newParty ",newParty)
				continue
			case newMessage := <-tssGlobalMessageChan:
				fmt.Println("test new message ",newMessage)
				bytes, routing, _ := newMessage.WireBytes()
				senderPartyID := sortedPartyIDs.FindByKey(MemberID(routing.From.GetKey()).bigInt())
				if routing.To != nil {
					for _, destination := range routing.To {
						destPartyID := sortedPartyIDs.FindByKey(MemberID(destination.GetKey()).bigInt())
						for _, party := range AllPartys {
							if destPartyID.Index != party.PartyID().Index {
								//fmt.Println("Id equals ",senderPartyID.Index,"---",party.PartyID().Index)
								continue
							} else {
								fmt.Println("Id not equals ",senderPartyID.Index,"---",party.PartyID().Index)
							}
							go partyUpdateFromBytes(party,bytes,senderPartyID,routing.IsBroadcast)
							// _, err := party.UpdateFromBytes(
							// 	bytes,
							// 	senderPartyID,
							// 	routing.IsBroadcast,
							// )
							// if err != nil {
							// 	fmt.Println("UpdateFromBytes error ",err,"---",senderPartyID,"---",party.PartyID())
							// } else {
							// 	time.Sleep(time.Duration(1)*time.Second)
							// 	fmt.Println("updateFromByte succ",senderPartyID,"---",party.PartyID())
							// }
						}
					}
				} else {
					for _, party := range AllPartys {
						if senderPartyID.Index == party.PartyID().Index {
							//fmt.Println("Id equals ",senderPartyID.Index,"---",party.PartyID().Index)
							continue
						} else {
							fmt.Println("Id not equals ",senderPartyID.Index,"---",party.PartyID().Index)
						}

						go partyUpdateFromBytes(party,bytes,senderPartyID,routing.IsBroadcast)
						// _, err := party.UpdateFromBytes(
						// 	bytes,
						// 	senderPartyID,
						// 	routing.IsBroadcast,
						// )
						// if err != nil {
						// 	fmt.Println("UpdateFromBytes error ",err,"---",senderPartyID,"---",party.PartyID())
						// } else {
						// 	time.Sleep(time.Duration(1)*time.Second)
						// 	fmt.Println("updateFromByte succ",senderPartyID,"---",party.PartyID())
						// }
					}
				}
				continue
			case <-ctx.Done():
				fmt.Println("manageParty Done")
				return 
		}
	}
}

func partyUpdateFromBytes(party tss.Party ,wireBytes []byte, from *tss.PartyID, isBroadcast bool) {
	_, err := party.UpdateFromBytes(
		wireBytes,
		from,
		isBroadcast,
	)
	if err != nil {
		fmt.Println("UpdateFromBytes error ",err,"---",from,"---",party.PartyID())
	} else {
	//	time.Sleep(time.Duration(1)*time.Second)
		fmt.Println("updateFromByte succ",from,"---",party.PartyID())
	}
}