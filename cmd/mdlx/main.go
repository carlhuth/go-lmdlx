package main

import "github.com/essential-web/mdlx"
import "io/ioutil"
import "log"
import "fmt"

func main() {
	var ba, err = ioutil.ReadFile("mdexamples/mkdn.md")
	if err != nil {
		log.Fatal(err)
	}
	b := mdlx.LoadBytes(ba)
	_ = b.Json()
	fmt.Println(b)
}
