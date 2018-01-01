/*
Package tdat implements TDAT for Go.

Sample File:

	$ cat sample.tdat
	customers
	|id:i    |name:s                         |customer_number:s  |registered:t
	|1       |"John \"Smiley \u263A\" Smith" |"88234"            |2010-03-04T09:34:55.882
	|2       |"Marco Polo"                   |"12399"            |2004-11-01T12:16:49.013
	|3       |"ႫႬႭႮႯ ႰႱႲႳ"               |"00233"            |2011-01-04T23:08:17.130

Parsing sample.tdat:

	import "github.com/cvilsmeier/tdat"

	func main() {
		tables, err := tdat.ParseFile("sample.tdat")
	}

*/
package tdat
