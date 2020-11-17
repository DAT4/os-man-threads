package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Shop struct {
	done      sync.WaitGroup
	customers chan Customer
	tills     []*Till
}

type Till struct {
	sync.Mutex
	cash int
	line chan Customer
}

type Grocery struct {
	price int
}

type Customer struct {
	shoppingDuration int
	groceries        []Grocery
}

func main() {
	closed := make(chan bool)
	var shop Shop

	shop.tills = []*Till{
		&Till{
			cash: 0,
			line: make(chan Customer, 10),
		},
		&Till{
			cash: 0,
			line: make(chan Customer, 10),
		},
		&Till{
			cash: 0,
			line: make(chan Customer, 10),
		},
		&Till{
			cash: 0,
			line: make(chan Customer, 10),
		},
		&Till{
			cash: 0,
			line: make(chan Customer, 10),
		},
		&Till{
			cash: 0,
			line: make(chan Customer, 10),
		},
	}

	shop.customers = make(chan Customer, 50)
	shop.done.Add(1)

	go func(closed chan<- bool) {
		time.Sleep(5 * time.Second)
		close(closed)
	}(closed)

	go func(customers chan<- Customer, closed chan bool) {
		for {
			select {
			case <-closed:
				fmt.Println("Door is closed")
				return
			case customers <- customerSpawner():
			}
		}
	}(shop.customers, closed)

	go shoppingScenario(shop.tills, shop.customers)

	go tillArea(&shop.done, shop.tills, closed)
	shop.done.Wait()
}

func grocerySpawner() (grocery Grocery) {
	grocery.price = rand.Intn(100-20) + 20
	return grocery
}

func customerSpawner() (customer Customer) {
	customer.shoppingDuration = rand.Intn(20-5) + 5

	groceriesAmount := rand.Intn(30-5) + 5
	for i := 0; i < groceriesAmount; i++ {
		customer.groceries = append(customer.groceries, grocerySpawner())
	}
	return customer
}

func shoppingScenario(tills []*Till, customersInShop chan Customer) {
	for {
		select {
		case tills[0].line <- <-customersInShop:
		case tills[1].line <- <-customersInShop:
		case tills[2].line <- <-customersInShop:
		case tills[3].line <- <-customersInShop:
		case tills[4].line <- <-customersInShop:
		case tills[5].line <- <-customersInShop:
		case <-time.After(3 * time.Second):
			return
		}
	}
}

func tillArea(wg *sync.WaitGroup, tills []*Till, closed chan bool) {
	defer wg.Done()
	for {
		select {
		case customer := <-tills[0].line:
			fmt.Println("Customer payed", customer)
		case customer := <-tills[1].line:
			fmt.Println("Customer payed", customer)
		case customer := <-tills[2].line:
			fmt.Println("Customer payed", customer)
		case customer := <-tills[3].line:
			fmt.Println("Customer payed", customer)
		case customer := <-tills[4].line:
			fmt.Println("Customer payed", customer)
		case customer := <-tills[5].line:
			fmt.Println("Customer payed", customer)
		case <-time.After(3 * time.Second):
			select {
			case <-closed:
				fmt.Println("Shop is closed 4 REALZ")
				return
			case <-time.After(time.Second):
				fmt.Println("Customers are slow... till is bored...")
			}
		}
	}
}
