package main

import (
	"errors"
	"time"

	"github.com/bitmark-inc/bitmark-sdk-go/bitmark"
)

var (
	// ErrAlreadyTransfer mean the account has got the asset already
	ErrAlreadyTransfer = errors.New("This account has got this asset already")
)

func (issue *Issue) transfer(receiver string) (err error) {
	// transfer this issue to an account number
	if err = db.Where("id = ?", issue.ID).Preload("Asset").Preload("Asset.User").First(&issue).Error; err != nil {
		panic(err)
	}
	if issue.Receiver != "" {
		return ErrAlreadyTransfer
	}
	/*
		if !account.IsValidAccountNumber(receiver) {
			return err
		}
	*/
	var count int
	err = db.Where("receiver = ?", receiver).Where("asset_id = ?", issue.AssetID).Find(&[]Issue{}).Count(&count).Error
	if err != nil {
		panic(err)
	}
	if count != 0 {
		return ErrAlreadyTransfer
	}

	params := bitmark.NewTransferParams(receiver)
	if err = params.FromBitmark(issue.ID); err != nil {
		panic(err)
	}
	if err = params.Sign(issue.Asset.User.Account()); err != nil {
		panic(err)
	}
	issue.TxID, err = bitmark.Transfer(params)
	if err != nil {
		return err
	}
	issue.Receiver = receiver
	issue.TransferredAt = time.Now()
	if err = db.Save(&issue).Error; err != nil {
		panic(err)
	}
	return nil
}
