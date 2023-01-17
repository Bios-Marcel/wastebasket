# Wastebasket

[![CI](https://github.com/Bios-Marcel/wastebasket/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/Bios-Marcel/wastebasket/actions/workflows/test.yml)

Wastebasket is a go library allowing you to move files into your trashbin.

## Dependencies

## Golang

The library supports at least the 3 latest major Golang versions. Depending on
your OS it might still work on an older version, but there are no guarantees.

### Windows

There are no dependencies, it depends on the Shell32 API built into Windows.

**No CGO required**

### Linux (Unix)

There are two (well, four) options you've got here. Wastebasket offers a native
golang implementation of the [FreeDesktop Trash specification](https://specifications.freedesktop.org/trash-spec/trashspec-latest.html).
This implementation is used by default. Alternatively, you can fall back to
using wrapper code for binaries on your path. The supported binaries are
the CLI interfaces for `gio`, `gvfs-trash` and `trash-cli`. At least one of
these is usually installed by default on desktop systems.

Additionally, the custom implementation should also work for systems such
as BSD and its derivatives. However, this has not been tested and I do not
plan on doing so, simply because GitHub does not currently support running
tests on any BSD derivatives.

If anyone is willing to host a custom runner (which I think is possible), then
I'd be open to this though.

**No CGO required**

### Mac OS

The only dependency is `Finder`, which is installed by default.

There are plans for a better implementation, that uses the Objective-C API
provided by Mac OS, resulting in most likely much better performance.

**No CGO required (Might change in the future)**

## How do i use it

Grab it via

```bash
go get github.com/Bios-Marcel/wastebasket
```

and you are ready to go.

Minimal Go example that creates a file, deletes it and empties the trashbin:

```go
package main

import (
    "fmt"
    "os"

    "github.com/Bios-Marcel/wastebasket"
)

func main() {
    os.WriteFile("test.txt", []byte("Test"), os.ModePerm)
    fmt.Println(wastebasket.Trash("test.txt"))
    wastebasket.Empty()
}
```

## Benchmarks

Run benchmarks using:

```go
go test -tags=nix_wrapper_impl -bench=.
```

//FIXME Maybe supply bench tag for clarity?
