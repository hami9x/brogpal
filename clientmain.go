package main

import (
	wd "github.com/phaikawl/wade"
)

type UserInfo struct {
	Name string
	Age  int
}

func main() {
	wade := wd.WadeUp("_wade")
	wade.RegisterElement("t-userinfo", &UserInfo{
		Name: "Hai Thanh Nguyen",
		Age:  18,
	})
}
