package mdlx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type parserStateFn func(p *mParser) parserStateFn

type mParser struct {
	tree        *MdTree
	lex         *lexer
	last        *item
	oldparent   []*item
	current     *item
	tokenBuffer []item
}

type MdTree struct {
	Val    *item           `json:"Val"`
	Tree   map[int]*MdTree `json:"Tree"`
	parent []*item
}

func (t *MdTree) String() string {
	s, _ := json.Marshal(t)
	var buf bytes.Buffer
	json.Indent(&buf, s, "", "\t")
	return buf.String()
}

func (t *MdTree) Json() []byte {
	s, _ := json.Marshal(t)
	return s
}

func getPrefix(pos int) string {
	//tab := "   "
	dash := "---"
	//tabs := ""
	dashes := ""
	for i := 0; i < pos; i++ {
		//\n---\n\ \ \
		//tabs = tabs + tab
		dashes = dashes + dash
	}
	return "\n" + dashes
}

func PrintStruct(t *MdTree) {
	printLessStruct(t, -1)
	fmt.Println("")
}

func printLessMap(m map[int]*MdTree, pos int) {
	pos++
	for key, val := range m {
		vt := fmt.Sprintf("%T", val)
		if vt == "*mdlx.MdTree" {
			printLessStruct(val, pos)
		} else {
			fmt.Print(getPrefix(pos), key, ":MAPVAL:", val)
		}
	}
}

func printLessStruct(t *MdTree, pos int) {
	pos++
	v := reflect.Indirect(reflect.ValueOf(t))

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	for key, val := range values {
		vt := fmt.Sprintf("%T", val)
		if vt == "map[int]*mdlx.MdTree" {
			printLessMap(val.(map[int]*MdTree), pos)
		} else {
			if vt == "int" {
				fmt.Printf("%strueval:%d", getPrefix(pos), val)
			} else {
				fmt.Printf("%s%d  :  %#v", getPrefix(pos), key, val)
			}
		}
	}
}
func (t *MdTree) createSubTree(keys []int, val *item) {
	subTree := t
	for _, intKey := range keys {
		nextTree, exists := subTree.Tree[intKey]
		if !exists {
			tree := newMdTree()
			subTree.Tree[intKey] = tree
			nextTree = tree
		}
		_ = nextTree
	}
}

func (p *mParser) next() *item {
	if len(p.tokenBuffer) != 0 {
		token := p.tokenBuffer[0]
		p.tokenBuffer = p.tokenBuffer[1:]
		return &token
	}
	token := p.lex.nextItem()
	return &token
}

func (p *mParser) peek() *item {
	if len(p.tokenBuffer) != 0 {
		return &(p.tokenBuffer[0])
	}
	token := p.lex.nextItem()
	p.tokenBuffer = append(p.tokenBuffer, token)
	return &token
}

func (p *mParser) checkBack() *item {
	return p.last
}

func newMdTree() *MdTree {
	return &MdTree{
		Tree: make(map[int]*MdTree),
	}
}

func LoadString(input string) *MdTree {
	return parse(lex(input))
}

func LoadBytes(input []byte) *MdTree {
	return parse(lex(string(input)))
}

func (p *mParser) run() {
	for mstate := startParse(p); mstate != nil; {
		mstate = mstate(p)
	}
}

func parse(lexer *lexer) *MdTree {
	p := mParser{
		tree:        newMdTree(),
		lex:         lexer,
		tokenBuffer: []item{},
	}
	p.run()
	return p.tree
}

func startParse(p *mParser) parserStateFn {
	tkn := p.peek()
	if tkn.Typ != itemNewline {
		p.last = p.current
		p.current = tkn
	}
	switch {
	case tkn.Typ == itemHeader:
		return parseHeader
	case tkn.Typ == itemText:
		return parseText
	case tkn.Typ == itemNewline:
		return parseNewline
	case tkn.Typ == itemHeader2:
		return parseHeader2
	case tkn.Typ == itemHeader3:
		return parseHeader3
	case tkn.Typ == itemHeader4:
		return parseHeader4
	case tkn.Typ == itemAsterisk:
		return parseAsterisk
	case tkn.Typ == itemBacktick:
		return parseBacktick
	case tkn.Typ == itemUnderscore:
		return parseUnderscore
	case tkn.Typ == itemUl:
		return parseUl
	case tkn.Typ == itemTabUl:
		return parseTabUl
	}
	p.last = tkn
	return nil
}

func setPath(p *mParser, ntkn *item) {
	subtree := p.tree
	for _, i := range p.tree.parent {
		nextTree := subtree.Tree[i.Id]
		subtree = nextTree
	}
	subtree.Tree[ntkn.Id] = newMdTree()
	subtree.Tree[ntkn.Id].Val = ntkn
}

func setHeaderPath(p *mParser, ntkn *item, wantedpos int) {
	var basearray []*item
	if wantedpos-1 > len(p.tree.parent) {
		basearray = p.tree.parent[0 : wantedpos-2]
	} else {
		basearray = p.tree.parent[0 : wantedpos-1]

	}
	p.tree.parent = basearray
	setPath(p, ntkn)
	basearray = append(basearray, ntkn)
	p.tree.parent = basearray
}

func setHeader(p *mParser, ntkn *item, wantedpos int) {
	var basearray []*item
	if wantedpos-1 > len(p.tree.parent) {
		basearray = p.tree.parent[0 : wantedpos-2]
	} else {
		basearray = p.tree.parent[0 : wantedpos-1]

	}
	p.tree.parent = basearray
	basearray = append(basearray, ntkn)
	p.tree.parent = basearray
}

func parseText(p *mParser) parserStateFn {
	ntkn := p.next()
	setPath(p, ntkn)
	return startParse
}

func parseUl(p *mParser) parserStateFn {
	ntkn := p.next()
	fmt.Println("xd, ", p.tree.parent[len(p.tree.parent)-1].Id)
	if p.tree.parent[len(p.tree.parent)-1].Typ == itemUl {
		fmt.Println("GOING")
		p.tree.parent = p.tree.parent[:len(p.tree.parent)-1]
	}
	setPath(p, ntkn)
	return startParse

}

func parseTabUl(p *mParser) parserStateFn {
	ntkn := p.next()
	oldntkn := p.last
	nextntkn := p.peek()
	if oldntkn.Typ == itemUl {
		p.oldparent = p.tree.parent
		p.tree.parent = append(p.tree.parent, oldntkn)
	}
	setPath(p, ntkn)
	if nextntkn.Typ != itemTabUl {
		if nextntkn.Typ != itemNewline {
			p.tree.parent = p.oldparent
		}
	}
	return startParse
}

//parse header needs to set the parse tree one layer deeper, then instruct all other funcs to use that parse tree location
func parseHeader(p *mParser) parserStateFn {
	//get token
	ntkn := p.next()
	setHeaderPath(p, ntkn, 1)
	return startParse
}

func parseUnderscore(p *mParser) parserStateFn {
	ntkn := p.next()
	setPath(p, ntkn)
	return startParse
}

func parseHeader2(p *mParser) parserStateFn {
	ntkn := p.next()
	setHeaderPath(p, ntkn, 2)
	return startParse
}

func parseHeader3(p *mParser) parserStateFn {
	ntkn := p.next()
	setHeaderPath(p, ntkn, 3)
	return startParse
}

func parseHeader4(p *mParser) parserStateFn {
	ntkn := p.next()
	setHeaderPath(p, ntkn, 4)
	return startParse
}

func parseAsterisk(p *mParser) parserStateFn {
	ntkn := p.next()
	setPath(p, ntkn)
	return startParse
}

func parseBacktick(p *mParser) parserStateFn {
	ntkn := p.next()
	setPath(p, ntkn)
	return startParse
}

func parseNewline(p *mParser) parserStateFn {
	ntkn := p.next()
	setPath(p, ntkn)
	return startParse
}
