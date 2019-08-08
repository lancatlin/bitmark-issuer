package main

func newAsset(c *gin.Context) {
	user := getUser()
	if !user.IsLogin {
		c.Redirect(303, "/login")
		return
	}
	var a Asset 
	if err := c.ShouldBind(&a); err != nil {
		panic(err)
	}
	header, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}
	file, err := header.Open()
	if err := nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	a.User = user
}

func (a *Asset) register(data []byte) (err error) {
	acct, err := account.FromSeed(a.User.Wallet)
	if err != nil {
		panic(err)
	}
	params, err := asset.NewRegistrationParams(
		a.Name,
		map[string]string{"Creater": user.Name}
	)
	params.SetFingerprint(data)
	params.Sign(acct)
	a.ID, err = asset.Register(params)
	if err != nil {
		panic(err)
	}
}

func (a *Asset) issue() (err error) {
	params := bitmark.NewIssuanceParams(
		a.ID,
		bitmark.QuantityOptions{
			Quantity: 3,
		},
	)
	params.Sign()
}