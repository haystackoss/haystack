![Test](https://github.com/nabaz-io/nabaz/actions/workflows/test-nabaz.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/nabaz-io/nabaz.svg)](https://pkg.go.dev/github.com/nabaz-io/nabaz) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)



# The change based test runner
Hate waiting for tests?    
Reduce up to **98%** of test wait time by skipping tests unaffected by code change.

**Under the hood:**
Gathers code coverage for individual tests, compares new code to code coverage.

![nabaz installation](https://raw.githubusercontent.com/nabaz-io/nabaz/main/docs/demo.gif)

Run it **locally** like so:

```bash
CMDLINE="pytest -v"
nabaz test --cmdline $CMDLINE ./
```
_**Note:** [Contact us](#contact-us) for early access to **CI/CD** version (remote storage, integration, system tests support)._

---
## Installation

### **From source** 🧙‍♂️
```bash
# Install nabaz binary.
export GOPATH=`go env GOPATH`
go install github.com/nabaz-io/nabaz/cmd/nabaz@latest
cp $GOPATH/src/github.com/nabaz-io/nabaz/bin/* /usr/local/bin
chmod +x /usr/local/bin/nabaz

# Required for go test support
mkdir -p $GOPATH/github.com/nabaz-io
cd $GOPATH/github.com/nabaz-io
git clone https://github.com/nabaz-io/go
cd go/src
./make.bash
mv $GOPATH/src/github.com/nabaz-io/go /usr/local/nabaz-go

# Required for pytest support
pip3 install pytest pytest-cov pytest-json pytest-json-report pytest-metadata pydantic

# Verify install
$ nabaz version
version 0.0
```

 ### **Linux** 🐧
```bash
# Ubuntu
wget -qO- https://nabaz.jfrog.io/artifactory/api/security/keypair/nabazgpg/public | apt-key add -
echo "deb [arch=amd64] https://nabaz.jfrog.io/artifactory/nabaz-debian-local stable main" >> /etc/apt/sources.list
sudo apt update
sudo apt install -y nabaz

# Debian
wget https://nabaz.jfrog.io/artifactory/nabaz-debian-local/pool/stable/nabaz-0.0-amd64.deb -O nabaz.deb
sudo dpkg -i ./nabaz.deb
```

---

## Running Tests
### pytest
```bash
export CMDLINE="pytest -v"
nabaz test --cmdline "$CMDLINE" .
```

### go test
```bash
export CMDLINE="go test"
export PKGS="./..." # IMPORTANT make sure packages are written SEPERATLY
nabaz test --cmdline $CMDLINE --pkgs $PKGS .
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
