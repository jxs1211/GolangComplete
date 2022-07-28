
<b>[Go 101](https://go101.org)</b> is a series of books on Go programming.
Currently, the following books are avaliable:

* [Go (Fundamentals) 101](https://go101.org/article/101.html), which focuses on Go syntax/semantics (except custom generics related) and all kinds of runtime related things.
* [Go Generics 101](https://go101.org/generics/101.html), which explains Go custom generics in detail.
* [Go Optimizations 101](https://go101.org/optimizations/101.html), which provides some code performance optimization tricks, tips, and suggestions.
* [Go Details & Tips 101](https://go101.org/details-and-tips/101.html), which collects many details and provides several tips in Go programming.

These books are expected to help gophers gain a deep and thorough understanding of Go
and be helpful for both beginner and experienced Go programmers.

To get latest news of Go 101 books, please follow the official twitter account [@go100and1](https://twitter.com/go100and1)
and join [the Go 101 Slack space](https://go-101.slack.com).

### Install, Update, and Read Locally

If you use Go toolchain v1.16+, then you don't need to clone the project respository:

```shell
### Install or update.

$ go install go101.org/go101@latest

### Read. (GOBIN path, defaulted as GOPATH/bin, should be set in PATH)

$ go101
Server started:
   http://localhost:55555 (non-cached version)
   http://127.0.0.1:55555 (cached version)
```

If you use Go toolchain v1.15-, or you would make some modifications (for contribution, etc.):

```shell
### Install.

$ git clone https://github.com/go101/go101.git

### Update. Enter the Go 101 project directory (which
# contains the current `README.md` file), then run

$ git pull

### Read. Enter the Go 101 project directory, then run

$ go run .
Server started:
   http://localhost:55555 (non-cached version)
   http://127.0.0.1:55555 (cached version)
```

The start page should be opened in a browser automatically.
If it is not opened, please visit http://localhost:55555.

Options:
```
-port=1234
-theme=light # or dark (default is auto)
```

Some HTML files are generated from their corresponding markdown files.
If a markdown file is modified, we can run `go run . -gen`
to synchronize its corresponding HTML file.

### Contributing

Welcome to improve Go 101 by:
* Submitting corrections for all kinds of mistakes, such as typos, grammar errors, wording inaccuracies, description flaws, code bugs and broken links.
* Suggesting interesting Go related contents.

Current contributors are listed on [this page](https://go101.org/article/acknowledgements.html).

Translations are also welcome. Here is a list of the ongoing translation projects:
* [中文版](https://github.com/golang101/golang101)

### License

Please read the [LICENSE](LICENSE) for more details.
