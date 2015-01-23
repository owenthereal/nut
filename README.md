# nut

Gophers love nuts.

`nut` is a tool that allows Go projects to declare dependencies, download dependencies, rewrite import paths and ensure that dependencies are properly vendored.

To accomplish this goal, `nut` does two things:

* Introduces a metadata file with dependency information.
* Fetches your project's dependencies and rewrite their import paths.
* Introduces conventions for managing vendored dependencies.

`nut` is specifically designed to manage dependencies for a binary program. If you are making a library, please follow [the standard Go way](https://golang.org/doc/code.html#Library).

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
├── LICENSE
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

github.com/octokit/go-octokit =
github.com/fhs/go-netrc/netrc = "4422b68c9c"
```

The format of declaring a dependency is `PACKAGE = COMMIT-ISH`.
`PACKAGE` is valid package path that passes into `go get`.
`COMMIT-ISH` can be any tag, sha or branch.
A blank `COMMIT-ISH` means the latest version of the dependency.

## Downloading dependencies

Run `nut install` to download dependencies. `nut` puts dependencies to `vendor/_nuts/`.
Let's take a look this folder:

```
$ tree vendor
vendor
└── _nuts
    ├── go-netrc
    |   ...
    └── go-octokit
    |   ...

8 directories, 73 files
```

The package name of a dependency is rewritten to the `PACKAGE_NAME` definite in `Nut.toml`.
The import paths of a dependency is rewritten to refer directly to the project folder.

## Import a dependency

To import to a dependency, refer to it by its full import path:

```go
package main

import (
  github.com/jingweno/hello_world/vendor/_nuts/go-octokit
  github.com/jingweno/hello_world/vendor/_nuts/go-netrc
)

func main() {
    c := octokit.NewClient()
    ...
}
```

## Build a project

All dependencies are properly vendored to `vendor/_nuts` and your program is referring to import paths relative to this folder.
`go build` and `go test` should just work.

## Prior art

* [godep](https://github.com/tools/godep)
* [party](https://github.com/mjibson/party)
