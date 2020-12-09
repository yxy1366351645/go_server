package main

import (
	"cluster-broadcast/shared"
	"fmt"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
)

func main() {
	cluster := startNode(8081)

	fmt.Print("\nBoot other nodes and press Enter\n")
	console.ReadLine()

	fmt.Print("\nAdding 1 Egg - Enter\n")
	console.ReadLine()
	calcAdd("Eggs", 1)

	fmt.Print("\nAdding 10 Egg - Enter\n")
	console.ReadLine()
	calcAdd("Eggs", 10)

	fmt.Print("\nAdding 100 Bananas - Enter\n")
	console.ReadLine()
	calcAdd("Bananas", 100)

	fmt.Print("\nAdding 2 Meat - Enter\n")
	console.ReadLine()
	calcAdd("Meat", 3)
	calcAdd("Meat", 9000)

	getAll()

	console.ReadLine()

	cluster.Shutdown(true)
}

func startNode(port int64) *cluster.Cluster {
	// how long before the grain poisons itself
	timeout := 10 * time.Minute

	system := actor.NewActorSystem()
	config := remote.Configure("localhost", 0)
	remote := remote.NewRemote(system, config)

	provider, _ := consul.New()
	clusterConfig := cluster.Configure("my-cluster", provider, config)
	cluster := cluster.New(system, clusterConfig)

	// this node knows about Hello kind
	remote.Register("Calculator", actor.PropsFromProducer(func() actor.Actor {
		return &shared.CalculatorActor{
			Timeout: &timeout,
		}
	}))

	// this node knows about Hello kind
	remote.Register("Tracker", actor.PropsFromProducer(func() actor.Actor {
		return &shared.TrackerActor{
			Timeout: &timeout,
		}
	}))

	shared.CalculatorFactory(func() shared.Calculator {
		return &shared.CalcGrain{}
	})

	shared.TrackerFactory(func() shared.Tracker {
		return &shared.TrackGrain{}
	})

	cluster.Start()
	return cluster
}

func calcAdd(grainId string, addNumber int64) {
	calcGrain := shared.GetCalculatorGrain(grainId)
	total1, err := calcGrain.Add(&shared.NumberRequest{Number: addNumber})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Grain: %v - Total: %v \n", calcGrain.ID, total1.Number)
}

func getAll() {
	trackerGrain := shared.GetTrackerGrain("singleTrackerGrain")
	totals, err := trackerGrain.BroadcastGetCounts(&shared.Noop{})
	if err != nil {
		panic(err)
	}

	fmt.Println("--- Totals ---")
	for grainId, count := range totals.Totals {
		fmt.Printf("%v : %v\n", grainId, count)
	}
}
