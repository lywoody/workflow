package workflow

import (
	"context"
	"fmt"
	"time"
)

// Flow is a control of a work flow series
type Flow struct {
}

// Start is
func (f *Flow) Start(ctx context.Context, wn *WorkNode) (<-chan string, error) {
	if wn == nil {
		return nil, fmt.Errorf("Nil Work")
	}
	c := make(chan string, 1)
	go func(ctx context.Context) {
		tk := time.NewTicker(time.Second) // Driven flow every sec
		for {
			select {
			case <-ctx.Done():
				c <- "terminated"
				return
			case <-tk.C:
				err := checkFlow(c, wn)
				if err != nil {
					c <- fmt.Sprintf("error:%s", err.Error())
					return
				}
			}
		}
	}(ctx)
	return c, nil
}

// check and run one step every time
func checkFlow(rc chan<- string, wn *WorkNode) error {
	// deal with the first node
	if wn == nil {
		return fmt.Errorf("started with nil node")
	}
	// if n.IsEnd() {
	// 	return "[node]finished", nil
	// }
	// find the node need to be executed
	wns, err := findAvailableNodes(wn)
	if err != nil {
		return err
	}
	for _, n := range wns {
		if n == nil {
			return fmt.Errorf("started with nil node")
		}
		if n.IsEnd() {
			rc <- fmt.Sprintf("[node]%s done", n)
			close(rc)
			return nil
		}
		// state check
		// if n.Status == WorkNodeStateDenied {
		// 	return fmt.Sprintf("[node]%s %s", n, n.Status), "", nil
		// }
		// execute the node
		if res, err := n.Exec(); err != nil {
			rc <- fmt.Sprintf("[node]%s exe err:%s", n, err)
		} else {
			rc <- fmt.Sprintf("[node]%s executed, res:%s", n, res)
		}
	}
	return nil
}

func findAvailableNodes(wn *WorkNode) (ws []*WorkNode, err error) {
	if wn == nil {
		err = fmt.Errorf("nil node")
		return
	}
	if wn.ExecutedAndPass() {
		for _, i := range wn.NextNodes {
			var found []*WorkNode
			found, err = findAvailableNodes(i)
			if err != nil {
				return
			}
			ws = append(ws, found...)
		}
	} else if wn.CanBeScheduled() {
		ws = append(ws, wn)
		return
	} else {
		return
	}
	return
}
