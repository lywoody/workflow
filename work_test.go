package workflow

import (
	"context"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	t.Logf("work flow started @ %v", time.Now())
	f := Flow{}
	// ctx, cancel := context.WithCancel(context.Background())
	ctx := context.Background()
	resc, err := f.Start(ctx, genFlow())
	if err != nil {
		t.Fatalf("%v", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msg, done := <-resc:
			if !done {
				return
			}
			t.Logf(msg)
		}
	}
}

func genFlow() *WorkNode {
	w := &WorkNode{ID: 1, Name: "start", status: WorkNodeStateWaiting}
	w.NextNodes = []*WorkNode{
		&WorkNode{ID: 2, Name: "approve1", status: WorkNodeStateTimeout},
		&WorkNode{ID: 3, Name: "approve2", status: WorkNodeStateTimeout},
	}
	e := &WorkNode{ID: 10, Name: "end", Type: WorkStateFinished}
	w.NextNodes[0].NextNodes = []*WorkNode{e}
	w.NextNodes[1].NextNodes = []*WorkNode{e}
	return w
}

func genFlow2() *WorkNode {
	start := &WorkNode{ID: 1, Name: "start"}
	fill := &WorkNode{ID: 2, Name: "fill"}
	check := &WorkNode{ID: 3, Name: "check"}
	approval1 := &WorkNode{ID: 11, Name: "app1"}
	approval21 := &WorkNode{ID: 121, Name: "app2-1"}
	approval22 := &WorkNode{ID: 122, Name: "app2-2"}
	process := &WorkNode{ID: 4, Name: "process"}
	validate := &WorkNode{ID: 5, Name: "validate"}
	end := &WorkNode{ID: 10, Name: "end", Type: WorkStateFinished}

	start.NextNodes = []*WorkNode{fill}
	fill.NextNodes = []*WorkNode{check}
	check.NextNodes = []*WorkNode{approval1, approval21}
	approval21.NextNodes = []*WorkNode{approval22}
	approval1.NextNodes = []*WorkNode{process}
	approval22.NextNodes = []*WorkNode{process}
	process.NextNodes = []*WorkNode{validate}
	validate.NextNodes = []*WorkNode{end}

	return start
}
