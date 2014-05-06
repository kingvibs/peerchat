package dht

import "time"
import "crypto/sha1"
import "net/rpc"
import "fmt"

const (
	Online = "Online"
	Offline = "Offline"
)
type Status string

const (
	IDLen = 64
	K = 8
)

type RoutingEntry struct {
	ipAddr string
	nodeId ID
}

type ID uint64

type SendMessageArgs struct {
	Content string
	Timestamp time.Time
	ToUsername string
	FromUsername string
}

type SendMessageReply struct {
	
}

type AnnouceUserArgs struct {
	QueryingNodeId ID
	QueryingIpAddr string
	AnnoucedUsername string
}

type AnnouceUserReply struct {
	QueriedNodeId ID
}

type FindNodeArgs struct {
	QueryingNodeId ID
	TargetNodeId ID
}

type FindNodeReply struct {
	QueriedNodeId ID
	TryNodes []string // if list is of length 1, then we found it
}

type GetUserArgs struct {
	QueryingNodeId ID
	TargetUsername ID
}

type GetUserReply struct {
	QueriedNodeId ID
	TryNodes []string // if list is of length 1, then we found it
}

type PingArgs struct {
	PingingNodeId ID
}

type PingReply struct {
	PingedNodeId ID
}

func Sha1(s string) ID {
	/*
		Returns a 160 bit integer based on a
		string input. 
	*/
    h := sha1.New()
    h.Write([]byte("hi"))
    bs := h.Sum(nil)
    l := len(bs)
    var a ID
	for i, b := range bs {
	    shift := ID((l-i-1) * 8)	
	    a |= ID(b) << shift
   	}
   	return a
}

func Xor(a, b ID) ID {
	/*
		Zors together two big.Ints and
		returns the result.
	*/
	return a ^ b
}

func find_n(a, b uint64) uint{
	var IDLen uint
	IDLen = 64
	var d, diff uint64
	diff = a ^ b
	var i uint
	for i = 0; i < IDLen; i++{
		d = 1<<(IDLen - 1 - i)
		if d & diff != 0 { // if true, return i
			break
		}
	}
	return i
}

// call() sends an RPC to the rpcname handler on server srv
// with arguments args, waits for the reply, and leaves the
// reply in reply. the reply argument should be a pointer
// to a reply structure.
//
// the return value is true if the server responded, and false
// if call() was not able to contact the server. in particular,
// the reply's contents are only valid if call() returned true.
//
// you should assume that call() will time out and return an
// error after a while if it doesn't get a reply from the server.
//
// please use call() to send all RPCs, in client.go and server.go.
// please don't change this function.
//
func call(srv string, rpcname string, args interface{}, reply interface{}) bool {
	c, errx := rpc.Dial("unix", srv)
	if errx != nil {
		return false
	}
	defer c.Close()
		
	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	} 

	fmt.Println(err)
	return false
}