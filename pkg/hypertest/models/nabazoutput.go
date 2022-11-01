package models

type FailedTest struct {
	Name string
	FileLink string
	Err string
	IsSubTest bool
}

type NabazOutput struct {
	IsThinking bool
	IsRunningTests bool
	IsFirstRun bool
	Err string
	FailedTests []FailedTest
}

type OutputState struct {
	PreviousTestsFailedOutput string
	FailedTests []FailedTest
}

func (o *OutputState) FailedTestIndex(failedTest string) (test *FailedTest, index int) {
	for index, test := range o.FailedTests {
		if test.Name == failedTest {
			return &test, index
		}
	}
	return nil, -1
}

func (o *OutputState) RemoveRottonTest(index int) {
	if index >= len(o.FailedTests) -1 {
		o.FailedTests = o.FailedTests[:len(o.FailedTests) -1]
	} else {
		o.FailedTests = append(o.FailedTests[:index], o.FailedTests[index+1:]...)
	}
}

func (o *OutputState) AddFailedTest(failedTest FailedTest) {
	o.FailedTests = append(o.FailedTests, failedTest)
}

func (o *OutputState) UpdateFailedTestError(index int, newError string) {
	o.FailedTests[index].Err = newError
}


