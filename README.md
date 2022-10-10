# nabaz: a change based test runner
Hate waiting for test runs to finish? nabaz only runs tests impacted by changed code.
 

![nabaz installation](https://raw.githubusercontent.com/nabaz-io/nabaz/main/docs/goinstall.gif)

Run it **localy** like so:

```bash
CMDLINE="pytest -v"
nabaz test --cmdline $CMDLINE ./
```
_**Note:** [Contact us](#contact-us) for early access to **CI/CD** version (remote storage, integration, system tests support)._

---
# Support
## Languages
- [x] Python 🐍
- [x] Go 
- [ ] Java (Coming soon)
- [ ] .NET/C# (Coming soon)
- [ ] Javascript (TBD)
- [ ] C/C++ (not planned currently)
- [ ] [Suggest here](#contact-us)
## Frameworks
- [x] pytest
- [x] go test
- [ ] JUnit
- [ ] XUnit
- [ ] Cypress
- [ ] [Suggest here](#contact-us)

---
## Installation

 ### **Linux** 🐧
```bash
sudo apt install nabaz
```

### **With `go install`**

You can install `nabaz` using the `go install` command:

```bash
# make sure PATH is set up
go env -w GOBIN=$GOPATH/bin
go install github.com/nabaz-io/nabaz/cmd/nabaz@latest
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
