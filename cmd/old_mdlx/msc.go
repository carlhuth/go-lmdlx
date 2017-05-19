package main

import (
	"fmt"
	"strings"
)

//the name says it all
func findAllBeforeCharInThisLine(ltr *Letter, char string) (ids []int) {
	var nlid = -1
	for i := ltr.Place; i < i+1; i-- {
		if f[i].Char == "\n" {
			nlid = i
			break
		}
	}
	if nlid != -1 {
		fmt.Println(nlid, "2:", ltr.Place)
		for i := nlid; i < ltr.Place; i++ {
			//if f[i].Char == char {
			ids = append(ids, f[i].Place)
			//}
		}
	}
	return
}

//Blocks a string given. The string should represent a case's stype
func (ltr *Letter) Block(str string) {
	ltr.Blocked = append(ltr.Blocked, str)
	id := getCaseByStype(f[ltr.Place], str)
	ltr.Cases[id].Truth = false
}

func (ltr *Letter) Set(stype string) {
	for i := 0; i < len(ltr.Cases); i++ {
		dont := false
		if ltr.Cases[i].Stype == stype {
			dont = true
		}
		if ltr.Cases[i].Stype != stype {
			if dont != true {
				ltr.Block(ltr.Cases[i].Stype)
			}
		}
	}
}

//gets a case's id by it's stype
func getCaseByStype(ltr Letter, Stype string) (id int) {
	for i := 0; i < len(ltr.Cases); i++ {
		if ltr.Cases[i].Stype == Stype {
			id = i
			return
		}
	}
	id = -1
	return
}

//find all ltr.Char's in f
func findAllInBytes(ltr Letter) (last []int) {
	for i := 0; i < ltr.Place; i++ {
		if f[i].Char != ltr.Char {
			last = append(last, i)
		}
	}
	return
}

//find last ltr.Char in bytes
func findLastInBytes(ltr Letter) (lastid int) {
	lastid = -1
	for i := ltr.Place; i < i+1; i-- {
		if asf := fcheckargs(f, i); asf == true {
			if f[i].Char == ltr.Char && f[i].Place != ltr.Place {
				lastid = i
				break
			}
		} else {
			return
		}
	}
	return
}

//find next ltr.Char in bytes
func findNextInBytes(ltr Letter) (lastid int) {
	lastid = -1
	for i := ltr.Place; i < len(f); i++ {
		if f[i].Char == ltr.Char && f[i].Place != ltr.Place {
			lastid = i
			return
		}
	}
	return
}

//if args[i] exists, than return true, else return false
func checkargs(args []string, i int) (truth bool) {
	defer func() {
		if r := recover(); r != nil {
			truth = false
		} else {
			truth = true
		}
	}()
	_ = args[i]
	return
}

//check if arg in f exists
func fcheckargs(args []Letter, i int) (truth bool) {
	defer func() {
		if r := recover(); r != nil {
			truth = false
		} else {
			truth = true
		}
	}()
	_ = args[i]
	return
}

//I think you can guess what is going on here.
//this is the most complicated part of my code, 😂
func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

//i didn't write inline docs for this, I will never.
func appendString(str string) {
	split := strings.Split(str, "")
	for i := 0; i < len(split); i++ {
		bt := []byte(split[i])[0]
		log = append(log, bt)
	}
}
