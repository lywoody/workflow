package workflow

import "fmt"

// Work is the definition of one Node
type Work struct {
	ID    uint
	Name  string
	Start *WorkNode // Workflow started here
}

// WorkNode is one node of a Work job
type WorkNode struct {
	ID   uint
	Name string
	// Start End
	// NPass means how many nodes should finished before this one
	// "" means normal
	Type string
	// Execute status
	// timeout - executed but timeout
	// pass - not executed but not nessasery, continue anyway
	// hangup - waiting by some reason, manual click needed
	status string // waiting executing done donewitherror timeout pass hangup
	// Execute result
	result      string
	resultState bool

	NextNodes []*WorkNode
}

// IsEnd is a func return
func (wn WorkNode) IsEnd() bool {
	// return len(wn.NextNodes) == 0
	return wn.NextNodes == nil
}

// Exec is a func return
func (wn *WorkNode) Exec() (string, error) {
	wn.status = WorkNodeStateDone
	return fmt.Sprintf("%s Executed OK", wn), nil
}

// ExecutedAndPass is
func (wn WorkNode) ExecutedAndPass() bool {
	if wn.status == WorkNodeStateDone || wn.status == WorkNodeStatePass {
		return true
	}
	return false
}

// CanBeScheduled is
func (wn WorkNode) CanBeScheduled() bool {
	if wn.status == WorkNodeStateWaiting {
		return true
	}
	return false
}

// String is
func (wn WorkNode) String() string {
	return fmt.Sprintf("[%d]%s(%s)/%s", wn.ID, wn.Name, wn.Type, wn.status)
}
