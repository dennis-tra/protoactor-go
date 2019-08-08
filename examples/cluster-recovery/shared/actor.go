package shared

import (
	fmt "fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/cluster"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Actor struct {
	Partner string
	seqNum  int32
}

func (state *Actor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		ctx.SetReceiveTimeout(5 * time.Second)

	case *actor.ReceiveTimeout:
		defer ctx.SetReceiveTimeout(5 * time.Second)

		req := &Request{
			SequenceNumber: state.seqNum,
		}

		partners := cluster.GetMemberPIDs(state.Partner)
		if partners.Len() == 0 {
			return
		}

		partners.ForEach(func(i int, pid actor.PID) {
			fmt.Printf("--> Sending REQUEST with Seq Num %d\n", state.seqNum)
			ctx.Request(&pid, req)
		})
		state.seqNum++

	case *Request:
		fmt.Printf("<-- Received REQUEST with Seq Num %d\n", msg.SequenceNumber)
		resp := &Response{
			SequenceNumber: msg.SequenceNumber,
		}
		ctx.Respond(resp)

	case *Response:
		fmt.Printf("<-- Received RESPONSE with Seq Num %d\n", msg.SequenceNumber)

	case *Panic:
		panic(ctx.Self().String())
	}
}
