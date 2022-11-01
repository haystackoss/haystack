package python_test

import (
	"testing"

	"github.com/nabaz-io/nabaz/pkg/hypertest/diffengine/parser/python"
)

func TestParse(t *testing.T) {
	// TestParse tests the parser
	parser, err := python.NewPythonParser()
	if err != nil {
		t.Error(err)
	}

	pyFile := []byte(`   
	def hello():
		print("Hello World!")
		x = 5
		def p():
			print("Hello Worldp")
	
	def hello2():
		print("HEY World!")
    `)
	funcs := parser.GetFunctions(pyFile)
	if len(funcs) != 3 {
		t.Errorf("Expected 3 functions, got %d", len(funcs))
	}

	// validate func names
	if _, ok := funcs["hello"]; !ok {
		t.Errorf("Expected function \"hello\" to be present")
	}
	if _, ok := funcs["p"]; !ok {
		t.Errorf("Expected function \"p\" to be present")
	}
	if _, ok := funcs["hello2"]; !ok {
		t.Errorf("Expected function \"hello2\" to be present")
	}
}
