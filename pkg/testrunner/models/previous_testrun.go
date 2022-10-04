package models

// PreviousTestRun is the info representing a previously run test.
type PreviousTestRun struct {
	Name     string
	Success  bool
	TimeInMs float64
	Ran      bool
	RunID    string
}

// SetSkipped sets the test as skipped.
func (r *PreviousTestRun) SetSkipped() *PreviousTestRun {
	r.Ran = false
	return r
}