# Wastebasket

[![CI](https://github.com/Bios-Marcel/wastebasket/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/Bios-Marcel/wastebasket/actions/workflows/test.yml)

Wastebasket is a go library allowing you to move files into your trashbin.

## Dependencies

### Windows

The only dependency is `PowerShell`, which is installed by default. This will change at some point, as I am planning to reimplement the functionallity using the win32 api.

### Linux

Your either need to have `gio`, `gvfs-trash` or `trash-cli` installed.
At least one of these is usually installed by default.

### Mac OS

The only dependency is `Finder`, which is installed by default.

## How do i use it

Grab it via

```bash
go get github.com/Bios-Marcel/wastebasket
```

and you are ready to go.

Minimal Go example that creates a file, deletes it and empties the trashbin:

```GO
package main

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/Bios-Marcel/wastebasket"
)

func main() {
    ioutil.WriteFile("test.txt", []byte("Test"), os.ModePerm)
    fmt.Println(wastebasket.Trash("test.txt"))
    wastebasket.Empty()
}
```