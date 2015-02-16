# nut

Gophers love nuts.

![nut](https://dl.dropboxusercontent.com/u/1079131/nut.png)

`nut` is a tool that allows Go projects to declare dependencies, download dependencies, rewrite import paths and ensure that dependencies are properly vendored.

To accomplish this goal, `nut` does three things:

* Introduces a metadata file with dependency information.
* Fetches your project's dependencies and rewrite their import paths.
* Introduces conventions for managing vendored dependencies.

`nut` is specifically designed to manage dependencies for a binary program. If you are making a library, please follow [the standard Go way](https://golang.org/doc/code.html#Library).

## Voting for `nut`

I made `nut` as part of the [Gopher Gala](http://www.gophergala.com/) contest.
If you like what you see, please [vote](http://gophergala.com/voting/) for it.

## Installation

Make sure you have a working Go environment. [See the install instructions](http://golang.org/doc/install.html).

To install `nut`, simply run:
```
$ go get github.com/gophergala/nut
```

Make sure your `PATH` includes to the `$GOPATH/bin` directory so your
commands can be easily used:

```
export PATH=$PATH:$GOPATH/bin
```

## Creating A New Project

Inside your `GOPATH`, run:

```
$ nut new hello_world
```

Let's checkout what `nut` has generated for us:

```
$ cd hello_world
$ tree .
.
├── .gitignore
├── Nut.toml
├── README.md
└── main.go

0 directories, 4 files
```

First let's check out `Nut.toml`:

```toml
[application]

name = "hello_world"
version = "0.0.1"
authors = ["Your Name <you@example.com>"]
```

This piece of application information is reserved for future usage.

## Adding a dependency

It's quite simple to add a dependency. Simply add it to your `Nut.toml` file:

```toml
[dependencies]

"github.com/octokit/go-octokit/octokit" = ""
"github.com/fhs/go-netrc/netrc" = "4422b68c9c"
```

The format of declaring a dependency is `PACKAGE = COMMIT-ISH`.
`PACKAGE` is a valid package path that passes into `go get`.
`COMMIT-ISH` can be any tag, sha or branch.
An empty `COMMIT-ISH` means the latest version of the dependency.

## Downloading dependencies

Run `nut install` to download dependencies. `nut` puts dependencies to `vendor/_nuts/`.
Let's take a look this folder:

```
$ tree vendor
vendor
└── _nuts
    └── github.com
        ├── fhs
        │   └── go-netrc
        │       ...
        └── octokit
            └── go-octokit
                ...

20 directories, 98 files
```

The import paths of all dependencies are rewritten to be relative to `vendor/_nuts`.
For example, assuming `github.com/octokit/go-octokit` depends on `github.com/sawyer/go-sawyer`,
all the import paths of `go-octokit` to `go-sawyer` will be relative to `vendor/_nuts` since they're vendored.

## Importing a dependency

To import to a dependency, refer to it by its full import path:

```go
package main

import (
  github.com/jingweno/hello_world/vendor/_nuts/github.com/octokit/go-octokit/octokit
)

func main() {
    c := octokit.NewClient()
    ...
}
```

## Building

All dependencies are properly vendored to `vendor/_nuts` and your program is referring to import paths relative to this folder.
`go build` and `go test` should just work.

## Demo

![demo](https://dl.dropboxusercontent.com/u/1079131/nut_demo.gif)

## FAQ

### What makes `nut` different than other dependency management tools?

`nut` allows you to lock dependencies, vendor them and rewrite their import paths.
The dependencies vendored by `nut` are "self contained" and are ready for use without overriding `GOPATH`.
Most existing dependency management tools in Go override `GOPATH` and you need an extra tool to build your project.
With `nut`, you can build your project with just `go build`.
`nut` properly vendors dependencies so that existing `go` commands work as a standard Go project.

### Is `nut` the same as `godep save -r`?

`godep` allow rewriting the import paths of dependencies with [`godep save -r`](https://github.com/tools/godep/blob/master/save.go#L46-L47).
`nut` does exactly the same thing in this regard.
However, `godep` doesn't allow updating any dependency with rewritten import paths.
`nut` doesn't currently support update of dependencies but it's a high priority item that will be implemented next.
The workflow will be as straightforward as `nut update FOO` to update an individual dependency specified in `Nut.toml`.

Besides, `nut`'s design philosophy is very different from `godep`:
`nut` is explicit about dependency management with a manifest file (`Nut.toml`) and allows locking of dependencies, whereas `godep` isn't.

## Who is using `nut`?

* [nut](https://github.com/gophergala/nut)
* [hub](https://github.com/github/hub) (WIP)

## Other arts

* [godep](https://github.com/tools/godep)
* [party](https://github.com/mjibson/party)
