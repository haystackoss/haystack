![Test](https://github.com/nabaz-io/nabaz/actions/workflows/test-nabaz.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/nabaz-io/nabaz.svg)](https://pkg.go.dev/github.com/nabaz-io/nabaz) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)



# The focused test runner
`hypertest` is an open-source CLI tool designed to help you maintain focus on the current coding subtask, built specifically for the easily distracted.  

**A better way to run tests**  
By analyzing code coverage, hypertest runs only the tests affected by a code change. Run tests in the background while coding, and get immediate **live** feedback (<250ms).

![hypertest](https://raw.githubusercontent.com/nabaz-io/nabaz/main/docs/hyper.gif)

### Try it for yourself:

```bash
nabaz hypertest --cmdline "pytest -v"
```
---
## Get Started

### **From source (Cross-platform)** 🧙‍♂️
```bash
# Install nabaz binary.
export GOPATH=`go env GOPATH`
go install github.com/nabaz-io/nabaz/cmd/nabaz@latest

# go test support
mkdir -p $GOPATH/github.com/nabaz-io
cd $GOPATH/github.com/nabaz-io
git clone https://github.com/nabaz-io/go
cd go/src/github.com/nabaz-io
./make.bash
mv $GOPATH/src/github.com/nabaz-io/go /usr/local/nabaz-go

# pytest support
pip3 install pytest pytest-cov pytest-json pytest-json-report pytest-metadata pydantic

# Verify install
$ nabaz version
version 0.0
```

 ### **Linux (Ubuntu)** 🐧
```bash
# Ubuntu
wget -qO- https://nabaz.jfrog.io/artifactory/api/security/keypair/nabazgpg/public | apt-key add -
echo "deb [arch=amd64] https://nabaz.jfrog.io/artifactory/nabaz-debian-local stable main" >> /etc/apt/sources.list
sudo apt update
sudo apt install -y nabaz
```

---

## Running Tests
### pytest
```bash
nabaz hypertest --cmdline "pytest -v" .
```

### go test
```bash
nabaz hypertest --cmdline "go test ./..."
```

---
# Language Support
## Languages
- [x] Python 🐍
- [x] Go 
- [ ] Java (Coming soon)
- [ ] .NET/C# (Coming soon)
- [ ] Javascript (TBD)
- [ ] C/C++ (not planned currently)
- [ ] [Request here](https://github.com/nabaz-io/nabaz/issues/new?assignees=&labels=&template=feature_request.md&title=)
## Frameworks
- [x] pytest
- [x] go test
- [ ] JUnit
- [ ] XUnit
- [ ] Cypress
- [ ] [Request here](https://github.com/nabaz-io/nabaz/issues/new?assignees=&labels=&template=feature_request.md&title=)

---
  
## Building
```bash
 go build -o nabaz ./cmd/nabaz
 ```

## Contact us!
at hello@nabaz.io.
## License

Licensed under the [MIT license](LICENSE.md).
