#MY UNDERSTANDING OF A-EBNF


##MY TAKE
So far, I would like to go with ABNF.
It seems more strict, and people modify it less.
EBNF can become full of differences from one spec to another.
ABNF also has nice syntax.
It is easier for people/programs to interpret.
----
>The EBNF repetiton VS ABNF repetition
>&lt;args&gt; ::= &lt;args&gt; { "," &lt;args&gt; } vs. &lt;n&gt;element = &lt;n&gt;&lt;n&gt;element 

Grouping and Concatenation are different
##EBNF
##NOT IN COMMON
##ABNF
###Rules
####Rule naming
Rule names are case insensitive 
HTML style <>'s are not required around rules, can be used tho
####Rule form
name = elements crlf

name=name of rule
elements=rule names or terminal specifications
crlf=end of line indicator
left alligned
indented soft returns

left alligned to start of ABNF, not ABNF document
####Terminal Values
rules = terminal values(at heart)
or called characters 
character = non negative intager
base interpretation is how a terminal interprets characters 
bases can be
b = binary
d = decimal
x = hexadecimal 
         CR          =  %d13

         CR          =  %x0D
values are represented as such with "."'s to separate strings of text
literal text strings are escaped with quotation marks as such
command = "command string"
abnf strings are case insensitive written in US-ASCII

Use space to indicate a case sensitive list of characters

####External encodings
External encodings are separate for syntax

###Operators
####Concatenation
A rule can be a list of other rules
for example
mundu = foo bat foo wow
####Linear White Space
You can put Linear White Space in your program, but it has to be
explicitly stated
####Alternatives
elements separated by "/" are alternatives 
foo / bar will accept foo or bar
####Incremental alternatives 
You can specify a list of alternatives in fragments by doing :
ex.1
         ruleset     =  alt1 / alt2

         ruleset     =/ alt3

         ruleset     =/ alt4 / alt5
ex.2
         ruleset     =  alt1 / alt2 / alt3 / alt4 / alt5
ex one and two are the same
####Value range Alternatives

				 DIGIT       =  %x30-39

   is equivalent to:

         DIGIT       =  "0" / "1" / "2" / "3" / "4" / "5" / "6" /

                        "7" / "8" / "9"

####Concatenated Numeric Values
A Concatenated Numeric Value is a list of values that cannot be
specified in the same string.
For example
char-line = %x0D.0A %x20-7E %x0D.0A

####Sequence Group
Elements enclosed in () are like single characters With ordered contents
Grouping is heavily advised here
It is easier to read
####Variable Repetition 
asterick = repetition
form=&lt;a&gt;*&lt;b&gt;element

Default values are 0 and infinity so that *&lt;element&gt; allows any
   number, including zero; 1*&lt;element&gt; requires at least one;
   3*3&lt;element&gt; allows exactly 3; and 1*2&lt;element&gt; allows one or two.
####Specific Repetition
  &lt;n&gt;element = &lt;n&gt;&lt;n&gt;element
####Optional Sequence 
[] enclose an optional escape sequence
*1(foo bar).
####Comment
; = comment
####Operator Precedence
Rule name, prose-val, Terminal value

      Comment

      Value range

      Repetition

      Grouping, Optional

      Concatenation

      Alternative
####ABNF OF ABNF
			   rulelist       =  1*( rule / (*c-wsp c-nl) )

         rule           =  rulename defined-as elements c-nl
                                ; continues if next line starts
                                ;  with white space

         rulename       =  ALPHA *(ALPHA / DIGIT / "-")

				 defined-as     =  *c-wsp ("=" / "=/") *c-wsp
                                ; basic rules definition and
                                ;  incremental alternatives

         elements       =  alternation *c-wsp

         c-wsp          =  WSP / (c-nl WSP)

         c-nl           =  comment / CRLF
                                ; comment or newline

         comment        =  ";" *(WSP / VCHAR) CRLF

         alternation    =  concatenation
                           *(*c-wsp "/" *c-wsp concatenation)

         concatenation  =  repetition *(1*c-wsp repetition)

         repetition     =  [repeat] element

         repeat         =  1*DIGIT / (*DIGIT "*" *DIGIT)

         element        =  rulename / group / option /
                           char-val / num-val / prose-val

         group          =  "(" *c-wsp alternation *c-wsp ")"

         option         =  "[" *c-wsp alternation *c-wsp "]"

         char-val       =  DQUOTE *(%x20-21 / %x23-7E) DQUOTE
                                ; quoted string of SP and VCHAR
                                ;  without DQUOTE

         num-val        =  "%" (bin-val / dec-val / hex-val)

         bin-val        =  "b" 1*BIT
                           [ 1*("." 1*BIT) / ("-" 1*BIT) ]
                                ; series of concatenated bit values
                                ;  or single ONEOF range

         dec-val        =  "d" 1*DIGIT
                           [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]

         hex-val        =  "x" 1*HEXDIG
                           [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]
				 prose-val      =  "<" *(%x20-3D / %x3F-7E) ">"
                                ; bracketed string of SP and VCHAR
                                ;  without angles
                                ; prose description, to be used as
                                ;  last resort

##OTHER
Golang was written in ebnf form

##OR IN ABNF
  /
