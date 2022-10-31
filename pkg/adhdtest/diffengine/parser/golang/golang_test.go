package golang_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/adhdtest/diffengine/parser/golang"
)

func TestParse(t *testing.T) {
	// TestParse tests the parser
	parser, err := golang.NewGolangParser()
	if err != nil {
		t.Error(err)
	}

	goFile := []byte(`   
    func loginHandler(w http.ResponseWriter, r *http.Request) {
        var details LoginDetails
    }
    
    type MyInteger int
    func (a MyInteger) MyMethod(b int) int {
      return a + b
    }
    
    func main() {
        http.HandleFunc("/login", loginHandler)
        fmt.Printf("Starting server at port 8080\n")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatal(err)
        }
    }
    `)
	funcs := parser.GetFunctions(goFile)
	if len(funcs) != 3 {
		t.Errorf("Expected 3 functions, got %d", len(funcs))
	}

	// validate func names
	if _, ok := funcs["loginHandler"]; !ok {
		t.Errorf("Expected function \"loginHandler\" to be present")
	}
	if _, ok := funcs["MyMethod"]; !ok {
		t.Errorf("Expected function \"MyMethod\" to be present")
	}
	if _, ok := funcs["main"]; !ok {
		t.Errorf("Expected function \"main\" to be present")
	}
}
