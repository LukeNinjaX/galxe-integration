package biz

import (
	"fmt"
	"testing"

	"github.com/artela-network/galxe-integration/config"
)

func TestSend(t *testing.T) {

	body := &PostBody{
		ChannelCode:   "artela",
		ChannelTaskId: "891b8fbef81c43c7aec3e4bfeea2c752",
		CompleteTime:  "1697076853",
		UserAddress:   "0x1dcabfc8807beb9c2314508f561a9ef43c9a2b03",
	}
	config := &config.GoPlusConfig{
		ChannelCode: "artela",
		ManageId:    "100005",
		ManageKey:   "mqucjot7NBTBPSjEL95tCZ4HL3BtYllV",
		SecwarexUrl: "",
	}
	GoPlus_Config = config
	sign, s, err := createSign(body)
	fmt.Println(sign)
	fmt.Println(s)
	fmt.Println(err)
}
