# go-ilog10
Fast integer log10 in Go (number of decimal digits - 1)

Usage:

```go
import "github.com/knz/go-ilog10"

func main() {
   fmt.Println(ilog10.FastUint32Log10(1234)) // prints 3
   fmt.Println(ilog10.NumInt32DecimalDigits(1234)) // prints 4

   fmt.Println(ilog10.FastUint64Log10(1234)) // prints 3
   fmt.Println(ilog10.NumInt64DecimalDigits(1234)) // prints 4
}
```
