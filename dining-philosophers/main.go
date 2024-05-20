package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

// some variables
var (
	eatTime   = 0 * time.Second
	thinkTime = 0 * time.Second
	sleepTime = 1 * time.Second
	hunger    = 3 //	how many times a philosopher eats
	eatChan   = make(chan string)
)

// var hunger = 3 //	how many times a philosopher eats
// eatTime := 1 * time.Second
// thinkTime := 3 * time.Second
// sleepTime := 1 * time.Second

func eat(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	color.Blue("%s is seated at the table.", philosopher.name)
	seated.Done()

	seated.Wait()

	for i := hunger; i > 0; i-- {
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			color.Yellow("\t%s takes the right fork.", philosopher.name)
			forks[philosopher.leftFork].Lock()
			color.Yellow("\t%s takes the left fork.", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			color.Yellow("\t%s takes the left fork.", philosopher.name)
			forks[philosopher.rightFork].Lock()
			color.Yellow("\t%s takes the right fork.", philosopher.name)
		}

		color.Yellow("\t%s has both forks and is eating.", philosopher.name)
		time.Sleep(eatTime)

		color.Yellow("\t%s has both forks and is thinking.", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		color.Blue("\t%s has put down the forks", philosopher.name)
	}

	color.Magenta("%s has left the table", philosopher.name)
	eatChan <- philosopher.name
}

func dine() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	forks := make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go eat(philosophers[i], wg, forks, seated)
	}

	wg.Wait()
	close(eatChan)
}

func main() {
	color.Cyan("==============================")
	color.Cyan(" Dining Philosophers Problem ")
	color.Cyan("==============================")
	color.HiMagenta("The table is empty")

	go func() {

		dine()
	}()

	eatingOrder := []string{}
	var mu sync.RWMutex

	for phil := range eatChan {
		mu.Lock()
		eatingOrder = append(eatingOrder, phil)
		mu.Unlock()
	}

	for _, phil := range eatingOrder {
		fmt.Printf("%s", color.HiGreenString("%s ", phil))
	}
	fmt.Println()

	color.HiMagenta("The table is empty")

}
