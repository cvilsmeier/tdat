
The Tabular Data (TDAT) Data Interchange Format
=============================================================================

[![GoDoc](https://godoc.org/github.com/cvilsmeier/tdat?status.svg)](https://godoc.org/github.com/cvilsmeier/tdat)
[![Build Status](https://travis-ci.org/cvilsmeier/tdat.svg?branch=master)](https://travis-ci.org/cvilsmeier/tdat)
[![Go Report Card](https://goreportcard.com/badge/github.com/cvilsmeier/tdat)](https://goreportcard.com/report/github.com/cvilsmeier/tdat)


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

This repository is the reference implementation for TDAT. It provides

* A specification (see rcf.txt)
* A parser and a renderer (generator) for TDAT models
* The tdat tool for validating TDAT files and converting them into various other formats (JSON, CSV)


Use it
-----------------------------------------------------------------------------

Get it with

    go get github.com/cvilsmeier/tdat



Performance
-----------------------------------------------------------------------------

TDAT compares well with JSON. This reference implementation includes benchmarks
that compare TDAT performance against JSON performance for sample data. The
benchmarks can be executed with the go tool:

    go test -bench github.com/cvilsmeier/tdat/...


Author
-----------------------------------------------------------------------------
C.Vilsmeier


