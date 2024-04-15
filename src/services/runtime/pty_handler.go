package runtime

import (
	"encoding/json"
	"fmt"
	"github.com/igm/sockjs-go/v3/sockjs"
	"github.com/pkg/errors"
	"io"
	"k8s.io/client-go/tools/remotecommand"
)

const EndOfTransmission = "\u0004"

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

type Session struct {
	Id            string
	SockJSSession sockjs.Session
	SizeChan      chan remotecommand.TerminalSize
	Bound         chan error
	DoneChan      chan struct{}
	ClusterName   string
	Namespace     string
	Service       string
	PodName       string
	Container     string
}

type Message struct {
	Op        string `json:"op"`
	SessionID string `json:"sessionId"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
	Data      string `json:"data"`
}

type Result struct {
	SessionId   string `json:"sessionId,omitempty"`
	Token       string `json:"token,omitempty"`
	Cluster     string `json:"cluster,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	ServiceName string `json:"serviceName,omitempty"`
	PodName     string `json:"podName,omitempty"`
	Container   string `json:"container,omitempty"`
	Cmd         string `json:"cmd,omitempty"`
}

func (t Session) Read(p []byte) (int, error) {
	m, err := t.SockJSSession.Recv()
	if err != nil {
		err = errors.Wrap(err, "sockJSSession.Recv")
		return 0, err
	}

	var msg Message
	if err := json.Unmarshal([]byte(m), &msg); err != nil {
		_ = t.Toast(fmt.Sprintf("read msg (%s) form client error.%v", string(p), err))
		return 0, nil
	}
	t.SizeChan <- remotecommand.TerminalSize{msg.Cols, msg.Rows}
	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.SizeChan <- remotecommand.TerminalSize{msg.Cols, msg.Rows}
		return 0, nil
	default:
		return copy(p, EndOfTransmission), fmt.Errorf("unknown message type '%s'", msg.Op)
	}
}

func (t Session) Write(p []byte) (int, error) {
	msg, err := json.Marshal(Message{
		Op:   "stdout",
		Data: string(p),
	})

	if err != nil {
		return 0, err
	}

	if err = t.SockJSSession.Send(string(msg)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t Session) Close(status uint32, reason string) {
	err := t.SockJSSession.Close(status, reason)
	fmt.Println(fmt.Sprintf("close socket (%s). %d, %s, %v", t.Id, status, reason, err))
}

// Next TerminalSize handles pty->process resize events
// Called in a loop from remotecommand as long as the process is running
func (t Session) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.SizeChan:
		return &size
	case <-t.DoneChan:
		return nil
	}
}

// Toast can be used to send the user any OOB messages
// hterm puts these in the center of the terminal
func (t Session) Toast(p string) error {
	msg, err := json.Marshal(Message{
		Op:   "toast",
		Data: p,
	})
	if err != nil {
		return err
	}

	if err = t.SockJSSession.Send(string(msg)); err != nil {
		return err
	}
	return nil
}
