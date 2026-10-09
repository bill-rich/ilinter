# ilinter

A `go/analysis` analyzer that reports variable and constant names shorter
than a minimum length (3 characters by default), with an allow list for the
handful of short names that are idiomatic Go (`ok`, `i`, `tx`, `wg`).

It checks every variable and constant a package declares: package level and
local `var` and `const` declarations, short variable declarations, `range`
and `select` clause variables, and type switch guards. Function parameters,
results, receivers and struct fields are not checked, and the blank
identifier is never reported.

## Usage

Run it on packages like any analyzer:

    go run github.com/bill-rich/ilinter/cmd/ilinter@latest ./...

Or pin it as a tool dependency of your module and run it through `go vet`:

    go get -tool github.com/bill-rich/ilinter/cmd/ilinter@latest
    go vet -vettool="$(go tool -n ilinter)" ./...

Flags:

    -min N          minimum number of characters in a name (default 3)
    -allow a,b,c    names that are always allowed (default ok,i,tx,wg)

For example, `ilinter -min 2 -allow ok,i,id ./...`.

## Development

    go test ./...
