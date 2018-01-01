
The Tabular Data (TDAT) Data Interchange Format
=============================================================================

TDAT is a data interchange format for tabular data. It is derived from CSV
(Comma Separated Values). Its main characteristics are:

* well-specified
* type-safe
* machine readable (and writable)
* human readable (and writable)

A sample looks like this:

	persons
	|id:i   |name:s      |male:b     |birth:t
	|1      |"John Doe"  |true       |1972-05-03T10:11:12.193
	|2      |"Jane Doe"  |false      |1973-04-21T04:12:54.677
	|3      |"Paul Doe"  |true       |2004-01-02T08:04:12.677
	|4      |"Baby Doe"  |           |

For a full specification of TDAT, see rfc.txt included in this repository.

What is included
-----------------------------------------------------------------------------

This repository is the refernce implementation for TDAT. It provides a parser
and a renderer (generator) for TDAT models. Additionally, it provides the tdat
tool. The tdat tool can be used to validate TDAT models and to convert TDAT
models into JSON and CSV.


Performance
-----------------------------------------------------------------------------

TDAT compares well with JSON. This reference implementation includes benchmarks
that compare TDAT performance against JSON performance for sample data. The
benchmarks can be executed with the go tool:

    go test -bench github.com/cvilsmeier/tdat/...


Author
-----------------------------------------------------------------------------
C.Vilsmeier


