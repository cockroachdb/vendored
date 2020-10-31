# gcassert

gcassert is a program for making assertions about compiler decisions in
Golang programs, via inline comment directives like `//gcassert:inline`.

## Example

Given a file `foo.go`:

```
package foo

func addOne(i int) int {
    return i+1
}

//gcassert:inline
func addTwo(i int) int {
    return i+1
}

func a(ints []int) int {
    var sum int
    for i := range ints {
        //gcassert:bce,inline
        sum += addOne(ints[i])

        sum += addTwo(ints[i]) //gcassert:bce

        sum += ints[i] //gcassert:bce
    }
    return sum
}
```

The inline `//gcassert` directive will cause `gcassert` to fail if the line
`sum += addOne(ints[i])` is either not inlined or contains bounds checks.

A `//gcassert:inline` directive on a function will cause `gcassert` to fail
if any of the callers of that function do not get inlined.

`//gcassert` comments expect a comma-separated list of directives after
`//gcassert:`. They can be included above the line in question or after, as an
inline comment.

## Installation

```
go get github.com/jordanlewis/gcassert/cmd/gcassert
```

## Usage

Run gcassert on packages containing gcassert directives, like this:

```
gcassert ./package/path
```

The program will output all lines that had a gcassert directive that wasn't
respected by the compiler.

For example, running on the testdata directory in this library will produce the
following output:

```
$ gcassert ./testdata
testdata/bce.go:8:	fmt.Println(ints[5]): Found IsInBounds
testdata/bce.go:16:	sum += notInlinable(ints[i]): call was not inlined
testdata/inline.go:22:	sum += notInlinable(i): call was not inlined
```

Inspecting each of the listed lines will show a `//gcassert` directive
that wasn't upheld when running the compiler on the package.

## Directives


```
//gcassert:inline
```

The inline directive on a CallExpr asserts that the following statement
contains a function that is inlined by the compiler. If the function does not
get inlined, gcassert will fail.

The inline directive on a FuncDecl asserts that every caller of that function
is actually inlined by the compiler

```
//gcassert:bce
```

The bce directive asserts that the following statement contains a slice index
that has no necessary bounds checks. If the compiler adds bounds checks,
gcassert will fail.
