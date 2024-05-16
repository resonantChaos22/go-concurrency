package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const (
	NumberOfPizzas = 10
)

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNum int, r *rand.Rand) *PizzaOrder {
	pizzaNum++
	if pizzaNum > NumberOfPizzas {
		return &PizzaOrder{
			pizzaNumber: pizzaNum,
		}
	}

	delay := r.Intn(2) + 1
	color.Yellow("Received Order #%d!!", pizzaNum)
	rnd := rand.Intn(12) + 1
	msg := ""
	success := false

	switch {
	case rnd < 5:
		pizzasFailed++
	case rnd >= 5:
		pizzasMade++
	}
	total++

	color.Yellow("Making pizza #%d, it will take %d seconds.", pizzaNum, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	switch {
	case rnd <= 2:
		msg = fmt.Sprintf("***	We ran out of ingredients for pizza #%d! ***", pizzaNum)
	case rnd > 2 && rnd <= 4:
		msg = fmt.Sprintf("***	The cook quit while making pizza #%d! ***", pizzaNum)
	default:
		success = true
		msg = fmt.Sprintf("Pizza Order #%d is Ready.", pizzaNum)
	}

	p := PizzaOrder{
		pizzaNumber: pizzaNum,
		message:     msg,
		success:     success,
	}

	return &p
}

func pizzeria(pizzaMaker *Producer, r *rand.Rand) {
	i := 0

	for {
		currentPizza := makePizza(i, r)
		if currentPizza == nil {
			continue
		}
		i = currentPizza.pizzaNumber
		select {
		case pizzaMaker.data <- *currentPizza:

		case quitChan := <-pizzaMaker.quit:
			close(pizzaMaker.data)
			close(quitChan)
			return
		}
	}
}

func main() {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizzeria(pizzaJob, r)

	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!!")
			}
		} else {
			color.Cyan("Done making pizzas....")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("***	Error closing channels - ", err)
			}
		}
	}

	color.Cyan("----------------------------")
	color.Cyan("Done for the day")
	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed >= 6:
		color.Red("It was not a good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was an okay day...")
	case pizzasFailed >= 2:
		color.Yellow("It was a pretty good day...")
	default:
		color.Green("It was a great day...")
	}

}
