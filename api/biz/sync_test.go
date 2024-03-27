package biz

import (
	"fmt"
	"testing"

	"github.com/artela-network/galxe-integration/config"
)

func TestSend(t *testing.T) {

	body := &PostBody{
		ChannelCode:   "artela",
		ChannelTaskId: "83115",
		CompleteTime:  "1697076853",
		UserAddress:   "0x1dcabfc8807beb9c2314508f561a9ef43c9a2b03",
	}
	config := &config.GoPlusConfig{
		ChannelCode: "artela",
		ManageId:    "100005",
		ManageKey:   "mqucjot7NBTBPSjEL95tCZ4HL3BtYllV",
		SecwarexUrl: "",
	}

	sign, s, err := createSign(body, config)
	fmt.Print(sign, s, err)
}
