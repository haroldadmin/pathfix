# pathfix

Package pathfix provides the ability to retrieve the PATH environment variable from the user's login shell and append its value to the PATH of the current process.

This is helpful when your Go program's binary is bundled in an application which is started from an OS's GUI. The GUI shell does not have access to the custom PATHs a user may have set in their terminal shell, which leads to problems when trying to find executables from the Go program. `pathfix` package helps solve this problem.

## Usage

The public API contains only one method.

```go
pathfix.Fix()
```

## Installation

Add this module to your project using `go get`:

```shell
go get github.com/haroldadmin/pathfix
```

## Contribution

Please report bugs and problems by opening issues. External contributions through PRs are also welcome.

## License

See [License](./LICENSE)
