/*
Copyright © 2024 Joshua Benn sublimeibanez@protonmail.com
*/
package main

import (
	"fmt"

	"github.com/SublimeIbanez/todor/cmd"
	"github.com/SublimeIbanez/todor/common"
)

func main() {
	cmd.Execute()

	fmt.Println(common.Path(""))
}
