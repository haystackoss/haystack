# nabaz: a change based test runner
Hate waiting for test runs to finish? nabaz only runs tests impacted by changed code.
 

![An animation showcasing that nabaz transforms a text log into an interactive log with folding sections.](https://raw.githubusercontent.com/GoTestTools/.github/main/gotestfmt.svg)

Run it **localy** like so:

```bash
CMDLINE="pytest -v"
nabaz test --cmdline $CMDLINE ./
```

### Support
## Languages
- [x] Python
- [x] Go
- [ ] Java
- [ ] .NET/C#
- [ ] Javascript
## Frameworks
- [x] pytest
- [x] go test

**Note:** Contact us for early access to CI version.

---
## Installation

 ### **Linux** 🐧
```bash
sudo apt install nabaz
```

### **With `go install`**

You can install `nabaz` using the `go install` command:

```bash
go install github.com/nabaz-io/nabaz/cmd/nabaz@latest

# make sure PATH is set up
go env -w GOBIN=$GOPATH/bin
```

You can then use the `nabaz` command, provided that your Go `bin` directory is added to your system path.


## Building

```bash
 go build -o nabaz ./cmd/nabaz
 ```

## Contact us!
at hello@nabaz.io.
## License

Licensed under the [MIT license](LICENSE.md).
