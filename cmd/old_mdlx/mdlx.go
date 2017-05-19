package main

import (
	//	"fmt"
	//for text normilization issues
	//	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"os"
)

//log is what the program prints out at the end of the running Go to line 25 for more info about letters
var log []byte

//f is all of the "letters" that the file md.md contains
var f []Letter

//byte archive or byte array is an array of all bytes in md.md
var ba []byte

//Liam case is the syntax for a case, it has a test, a handler(if the test is true) And a stype, a short description of the type. Truth is if the case is true about the Letter in question(letter refers to my letter, not the english letter). Capitilizd letter will mean the stuff declared on line 25.
type Lcase struct {
	Test    func(ltr *Letter) bool
	Truth   bool
	Stype   string
	Handler func(ltr *Letter) string
}

type Letter struct {
	//character that the byte stands for
	Char string
	//Place in f, used as an iding system.
	Place int
	//Which handler did the letter execute
	Was string
	//the byte that the letter represents
	Byte byte
	//blocked Cases that this letter cannot be
	Blocked []string
	//all of the cases this Letter could be
	Cases []Lcase
}

func main() {
	appendString("<!Doctype html>\n")
	bytearray, err := ioutil.ReadFile("mdexamples/md.md")
	ba = bytearray
	//checkerr is in msc.go
	checkerr(err)
	//create a copy of AllCases AllCases in ctx.go or context.go
	pcases := AllCases
	//populate f with all of the letters in ba. Create copies of all letters in ba, and make those copies Letters
	for i := 0; i < len(ba); i++ {
		CurrentByte := Letter{
			Char:  string(ba[i]),
			Place: i,
			Byte:  ba[i],
			Cases: pcases,
		}
		f = append(f, CurrentByte)
	}
	//Testing phase.
	for i := 0; i < len(f); i++ {
		//Creates CurrentByte, a pointer to the location of f[i]
		CurrentByte := &f[i]
		//Tests every case. Set the truth of each case to the result of the Test function
		for q := 0; q < len(CurrentByte.Cases); q++ {
			truth := CurrentByte.Cases[q].Test(CurrentByte)
			CurrentByte.Cases[q].Truth = truth
		}
		//Iterate over all cases, and run tests on each to determine weather or not the Handler can be used, or is blocked
		for w := 0; w < len(AllCases); w++ {
			//if the case is true
			if CurrentByte.Cases[w].Truth == true {
				//dummy variable
				listed := false
				//if the case is on the blocked list, than set listed to true
				for qw := 0; qw < len(CurrentByte.Blocked); qw++ {
					CurrentByte.Cases[w].Truth = false
					if CurrentByte.Cases[w].Stype == CurrentByte.Blocked[qw] {
						listed = true
					}
				}
				//if the case was not listed
				//do the handler
				if listed == false {
					ret := CurrentByte.Cases[w].Handler(&f[i])
					CurrentByte.Was = ret
					break
				}
			}
		}
	}
	//for everything in log
	//for q := 0; q < len(log); q++ {
	//norm.NFC.Bytes is a unnorming program, It allows emoji's and other unicode characters to be printed properly
	//fmt.Print(string(norm.NFC.Bytes(log[q : q+1])))
	//}
	if _, err := ioutil.ReadFile("log.log"); err != nil {
		os.Remove("log.log")
	}
	ioutil.WriteFile("log.log", log, 0666)

}
