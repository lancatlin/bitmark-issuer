package main

import (
	"errors"
	"time"

	"github.com/bitmark-inc/bitmark-sdk-go/asset"
	"github.com/bitmark-inc/bitmark-sdk-go/bitmark"
)

var (
	// ErrAlreadyRegister mean a file was register in the blockchain
	ErrAlreadyRegister = errors.New("this file already has been register")
	// ErrTimeout mean the register action was timeout
	ErrTimeout = errors.New("Register timeout")
)

func (a *Asset) register(data []byte) (err error) {
	// register an asset to the network but not issue yet
	// need a.issue after register
	acct := a.User.Account()
	params, err := asset.NewRegistrationParams(
		a.Name,
		map[string]string{"Creater": a.User.Name},
	)
	params.SetFingerprint(data)
	params.Sign(acct)
	chErr := make(chan error, 1)
	chID := make(chan string, 1)
	go func(a *Asset) {
		id, err := asset.Register(params)
		chErr <- err
		chID <- id
	}(a)
	select {
	case err = <-chErr:
		a.AssetID = <-chID
		if err != nil {
			if err.Error() == "[1000] message: invalid parameters reason: asset already registered" {
				return ErrAlreadyRegister
			}
			return err
		}
	case <-time.After(10 * time.Second):
		// Timeout
		// The asset already has been register
		return ErrTimeout
	}
	return nil
}

func (a *Asset) issue() (err error) {
	// issue an asset to the network
	// use quantity to issue multiple asset
	params := bitmark.NewIssuanceParams(a.AssetID, a.Amount)
	params.Sign(a.User.Account())
	bitmarkIDs, err := bitmark.Issue(params)
	if err != nil {
		return errors.New("bitmark issue fatal: " + err.Error())
	}
	a.Issues = make([]Issue, a.Amount)
	for i, id := range bitmarkIDs {
		issue := Issue{
			ID:      id,
			AssetID: a.ID,
		}
		a.Issues[i] = issue
	}
	return nil
}
