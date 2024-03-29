package handle

import (

	"browser/utils"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"strings"
	"time"
)

func (f FabSdk)ListenBlock(){
	client, err := event.New(f.channeProvider,event.WithBlockEvents(),event.WithBlockNum(uint64(0)))
	if err != nil{
		fmt.Println(fmt.Errorf("listen block event err :%s",err.Error()))

	}
	register,notifier,err := client.RegisterBlockEvent()
	if err != nil{
		fmt.Println(fmt.Errorf("regist block event err :%s",err.Error()))
	}

	defer client.Unregister(register)

	for ;; {
		select {
		case ccEvent := <-notifier:
			fmt.Println("receive block event")
			fmt.Println(fmt.Sprintf("url:%s",ccEvent.SourceURL))

			tx,err := utils.UpdateBlockAndTx(*ccEvent.Block)
			if err != nil{
				fmt.Printf("received ledger event err :%s\n", err.Error())
			}

			/////// 判断是否发布token
			if tx.Args != nil {
				if  strings.ToLower(tx.Args[0]) == "issue" {
					fabsdk := InitSdk()
					defer fabsdk.Close()
					err = fabsdk.SyncToken()
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}


		case <-time.After(time.Second * 60):
			fmt.Println("timeout while waiting for block event")
		}
	}

}