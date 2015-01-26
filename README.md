# nut

Gophers love nuts.

![nut](https://raw.githubusercontent.com/gophergala/nut/master/nut.png)

`nut` is a tool that allows Go projects to declare dependencies, download dependencies, rewrite import paths and ensure that dependencies are properly vendored.

To accomplish this goal, `nut` does three things:

* Introduces a metadata file with dependency information.
* Fetches your project's dependencies and rewrite their import paths.
* Introduces conventions for managing vendored dependencies.

`nut` is specifically designed to manage dependencies for a binary program. If you are making a library, please follow [the standard Go way](https://golang.org/doc/code.html#Library).

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

github.com/octokit/go-octokit/octokit = ""
github.com/fhs/go-netrc/netrc = "4422b68c9c"
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

## Who is using `nut`?

* [nut](github.com/gophergala/nut)
* [hub](github.com/github/hub)(WIP)

## Other arts

* [godep](https://github.com/tools/godep)
* [party](https://github.com/mjibson/party)
