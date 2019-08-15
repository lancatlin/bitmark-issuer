package main

import (
	"testing"

	"github.com/bitmark-inc/bitmark-sdk-go/account"
)

func TestTransfer(t *testing.T) {
	var asset Asset
	for asset.Free() < 2 {
		if err := db.Take(&asset).Error; err != nil {
			t.Error(err)
		}
	}
	var issue Issue
	if err := db.Where("receiver = ''").Take(&issue).Error; err != nil {
		t.Error(err)
		return
	}
	receiver, err := account.New()
	if err != nil {
		t.Error(err)
		return
	}
	if err := issue.transfer(receiver.AccountNumber()); err != nil {
		t.Error("transfer fatal: ", err)
	}
	if err := issue.transfer(receiver.AccountNumber()); err != ErrAlreadyTransfer {
		t.Error("transfer fatal: ", err)
	}
}
