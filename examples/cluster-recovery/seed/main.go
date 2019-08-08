package main

import (
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/examples/cluster-recovery/shared"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func main() {

	props := actor.PropsFromProducer(func() actor.Actor {
		return &shared.Actor{
			Partner: "member",
		}
	})
	remote.Register("seed", props)

	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("mycluster", "127.0.0.1:9090", cp)

	rootCtx := actor.NewRootContext(nil)
	rootCtx.SpawnNamed(props, "seed")

	console.ReadLine()

	cluster.Shutdown(true)
}
