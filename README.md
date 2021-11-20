# go-satel

go-satel is a Go library to integrate
with [Satel ETHM-1/ETHM-1 Plus](https://www.satel.pl/produkty/sswin/komunikacja-i-powiadamianie/komunikacja-tcp-ip/ethm-1-plus/)
module

[![Build](https://github.com/probakowski/go-satel/actions/workflows/build.yml/badge.svg)](https://github.com/probakowski/go-satel/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/probakowski/go-satel)](https://goreportcard.com/report/github.com/probakowski/go-satel)

## Installation
go-satel is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/probakowski/go-satel
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/probakowski/go-satel"
```

and run `go get` without parameters.

Finally, to use the top-of-trunk version of this repo, use the following command:

```bash
go get github.com/probakowski/go-satel@master
```

## Usage
```go
s := satel.NewConfig("<ip:port>", satel.Config{EventsQueueSize: 1000})
go func() {
    value := true
    for {
        s.SetOutput("<user code>", <port number> value)
        time.Sleep(5 * time.Second)
        value = !value
    }
}()
for e, ok := <-s.Events; ok; e, ok = <-s.Events {
    logger.Print("change from satel", "type", e.Type, "index", e.Index, "value", e.Value)
}
```