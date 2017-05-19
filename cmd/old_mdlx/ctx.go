package main

import "fmt"
import "strconv"

//where all of the cases are stored
var AllCases = []Lcase{
	//the case to determine weather or not a *Letter given is the beginning of a header of any size
	{
		Test: func(ltr *Letter) bool {
			if ltr.Char == "#" {
				return true
			}
			return false
		},
		Truth: false,
		Stype: "bgin.header",
		Handler: func(ltr *Letter) string {
			//keeps track of how many #'s there are
			k := 0
			//for every #, k++, if you reach a character that is not a #, then break
			for i := ltr.Place; i < i+1; i++ {
				if f[i].Char != "#" {
					break
				}
				//devNull is the case for deleting a letter
				f[i].Set("devNull")
				k++
			}
			//<find next newline's place in f>
			//<and set it to end.header>
			var nl = Letter{
				Char:  "\n",
				Place: ltr.Place,
			}
			nk := strconv.Itoa(k)

			ks := "<h" + nk + ">"
			newlineid := findNextInBytes(nl)
			if newlineid != -1 {
				var nnl = &f[newlineid]
				nnl.Set("end.header")
			}
			appendString(ks)
			//</find next newline's place in f and set it to end.header>
			return "bgin.header"
		},
	},
	{
		Test: func(ltr *Letter) bool {
			//if there are four ---- in a row, delete the last three, make this one a sld.line
			var is int
			if ltr.Char == "-" {
				for i := ltr.Place; i < ltr.Place+4; i++ {
					if f[i].Char == "-" {
						is++
					}
				}
				if is == 4 {

					return true
				}
			}
			return false
		},
		Truth: false,
		Stype: "sld.line",
		Handler: func(ltr *Letter) string {
			for i := ltr.Place; i < ltr.Place+4; i++ {
				f[i].Set("devNull")
			}
			appendString("<br>")
			return "-sld.line-"
		},
	},
	{
		Test: func(ltr *Letter) bool {
			//if this character is set to not exist
			for i := 0; i < len(ltr.Blocked); i++ {
				if ltr.Blocked[i] == "plain text" {
					return true
				}
			}
			return false
		},
		Truth: false,
		Stype: "devNull",
		Handler: func(ltr *Letter) string {
			return "devNull"
		},
	},

	{
		Test: func(ltr *Letter) bool {
			//if this character is set to plain text
			for i := 0; i < len(ltr.Blocked); i++ {
				if ltr.Blocked[i] == "plain text" {
					return true
				}
			}
			return false
		},
		Truth: false,
		Stype: "end.header",
		Handler: func(ltr *Letter) string {
			//find last #, count all of the hashtags before it, then set this one to the appropriate header value
			var h1 = Letter{
				Place: ltr.Place,
				Char:  "#",
			}
			id := findLastInBytes(h1)
			//id = 2
			dll := 0
			fmt.Println(id)
			for i := id; i > i-1; i-- {
				if f[i].Char == "#" {
					for q := i; q < q+1; q-- {
						asdf := fcheckargs(f, q)
						if asdf == true {
							if f[q].Char != "#" {
								break
							}
						} else {
							break
						}
						dll++

					}
					break
				}
			}
			dl := strconv.Itoa(dll)

			appendString("</h" + dl + ">\n")
			return "end.header"
		},
	},

	//determine weather or not the letter given is a backtick
	{
		Test: func(ltr *Letter) bool {
			if ltr.Char == "`" {
				return true
			}
			return false
		},
		Truth: false,
		Stype: "backtick",
		Handler: func(ltr *Letter) string {
			//get the id of the last backtick
			fmt.Println(ltr.Place)
			last := findLastInBytes(f[ltr.Place])
			var lb *Letter
			//-1 is the signal for doesen't exist, there is no -1 byte in an array, there for it's impossible for this one to exist
			//if the last one exists
			if last != -1 {
				//create a pointer to the the last byte.
				lb = &f[last]
				//if the last byte was a beginning backtick
				if lb.Was == "-bgin-backtick-" {
					//make this one an nding backtick
					appendString("</code>")
					return "-nd-backtick-"
				} else {
					//if the last one was an ending backtick, make this one a beginning backtick
					appendString("<code>")
					return "-bgin-backtick-"
				}
			} else {
				//if the last one didn't exist, make this one a beginning backtick
				appendString("<code>")
				return "-bgin-backtick-"
			}
		},
	},
	{
		Test: func(ltr *Letter) bool {
			var isbeginning, isending bool
			if ltr.Char == "*" {
				//if this is a beginning astrick, then set is as a bold character
				//if there is a newline between this and the next *, than this one cannot be a beginning backtick
				for i := ltr.Place + 1; i < len(f); i++ {
					if string(f[i].Char) == "\n" {
						isbeginning = false
						break
					}
					if string(f[i].Char) == "*" {
						isbeginning = true
						break
					}
				}
				//if there is a newline inbetween the last * and this *, than this one is not an ending *
				last := findLastInBytes(f[ltr.Place])
				if last != -1 {
					lba := &f[last]
					var unns = true
					for i := lba.Place + 1; i < ltr.Place; i++ {
						if f[i].Char == "\n" {
							unns = false
							isending = false
							break
						}
					}
					if unns == true {
						isending = true
					}
				}
			}
			if isbeginning == true {
				return true
			}
			if isending == true {
				return true
			}
			return false
		},
		Truth: false,
		Stype: "bgin-bold",
		Handler: func(ltr *Letter) string {
			nid := findNextInBytes(f[ltr.Place])
			if nid != -1 {
				next := &f[nid]
				next.Block("bgin-bold")
			}
			appendString("<b>")
			return "-bgin.bold-"
		},
	},
	{
		Test: func(ltr *Letter) bool {
			//see if this is a ending bold
			if ltr.Char == "*" {
				var isending bool
				last := findLastInBytes(f[ltr.Place])
				if last != -1 {
					lba := &f[last]
					var unns = true
					for i := lba.Place + 1; i < ltr.Place; i++ {
						if f[i].Char == "\n" {
							unns = false
							isending = false
							break
						}
					}
					if unns == true {
						isending = true
					}
				}
				if isending == true {
					return true
				}
			}
			return false
		},
		Truth: false,
		Stype: "end-bold",
		Handler: func(ltr *Letter) string {
			nid := findNextInBytes(f[ltr.Place])
			if nid != -1 {
				next := &f[nid]
				next.Block("end-bold")
			}
			appendString("</b>")
			return "-end.bold-"
		},
	},
	{
		Test: func(ltr *Letter) bool {
			//see if this is a list item, or if this is an indented list item
			if _, err := strconv.Atoi(ltr.Char); err == nil && f[ltr.Place+1].Char == "." {
				if f[ltr.Place-1].Char == "\n" {
					return true
				}
				var brkn = false
				if spaces := findAllBeforeCharInThisLine(ltr, " "); len(spaces) != 0 {
					for i := 1; i < len(spaces); i++ {
						if f[spaces[i]].Char != " " {
							fmt.Println(f[spaces[i]].Char)
							brkn = true
							break
						}
					}
				}
				fmt.Println(brkn)
				if brkn != true {
					return true
				}
			}
			return false
		},
		Truth: false,
		Stype: "list item",
		Handler: func(ltr *Letter) string {
			var msg string
			//I like descriptive function names
			spaces := findAllBeforeCharInThisLine(ltr, " ")
			if len(spaces)-1 != -1 {
				for i := 0; i < len(spaces)-1; i = i + 4 {
					msg = msg + "-tab-"
				}
			}
			appendString("-list.item[" + msg + "]-")
			return "-list.item-"
		},
	},
	{
		Test: func(ltr *Letter) bool {
			//unorderd list item
			if ltr.Char == "*" || ltr.Char == "-" || ltr.Char == "+" {
				if f[ltr.Place-1].Char == "\n" {
					for i := ltr.Place + 1; i < len(f); i++ {
						if f[i].Char != " " && f[i].Char != "\n" {
							return true
						}
						if f[i].Char == "\n" {
							break
						}
					}
				}
			}
			return false
		},
		Truth: false,
		Stype: "unorderd item",
		Handler: func(ltr *Letter) string {
			appendString("-list.unorderedItem-")
			return "-list.unorderedItem-"
		},
	},
	//my first case, harkens back to github.com/17liamnaddell/lxr
	{
		Test: func(ltr *Letter) bool {
			if ltr.Char == ltr.Char {
				return true
			}
			return false
		},
		Truth: false,
		Stype: "plain text",
		Handler: func(ltr *Letter) string {
			log = append(log, ltr.Byte)
			return "ltr"
		},
	},
}
