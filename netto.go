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

	shop.customers = make(chan Customer, 500)
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
				fmt.Println("Customer enter shop!")
			default:
				time.Sleep(time.Second)
				fmt.Println("Shop is full!")
			}
		}
	}(shop.customers, closed)

	doneShopping := make(chan Customer, 500)

	go shopping(shop.customers, doneShopping)

	go goingToLine(shop.tills, doneShopping)

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

func shopping(customersEnterShop <-chan Customer, lineArea chan Customer) {
	/**
	 * Customer in the shop will take his time shopping
	 * and then he will go to check for a line
	 */
	var customer Customer

	for {
		select {
		case customer = <-customersEnterShop:
			time.Sleep(time.Duration(customer.shoppingDuration) * time.Millisecond)
			go toLine(customer, lineArea)
		case <-time.After(time.Second):
			fmt.Println("No customers in shop")
			return
		}
	}
}

func toLine(customer Customer, lineArea chan<- Customer) {
	select {
	case lineArea <- customer:
	case <-time.After(time.Second):
		fmt.Println("Some error")
		return
	}
}

func goingToLine(tills []*Till, customersInShop <-chan Customer) {
	for {
		select {
		case tills[0].line <- <-customersInShop:
		case tills[1].line <- <-customersInShop:
		case tills[2].line <- <-customersInShop:
		case tills[3].line <- <-customersInShop:
		case tills[4].line <- <-customersInShop:
		case tills[5].line <- <-customersInShop:
		case <-time.After(time.Second):
			fmt.Println("All tills are closed!")
			return
		}
	}
}

func (self Customer) getTotal() int {
	var total int
	for _, grocery := range self.groceries {
		total += grocery.price
	}
	return total
}

func tillArea(wg *sync.WaitGroup, tills []*Till, closed chan bool) {
	var customer Customer
	defer wg.Done()
	for {
		select {
		case customer = <-tills[0].line:
			fmt.Println("Customer payed", customer.getTotal(), "DKK , at till 1")
			tills[0].cash += customer.getTotal()
		case customer = <-tills[1].line:
			fmt.Println("Customer payed", customer.getTotal(), "DKK , at till 2")
			tills[1].cash += customer.getTotal()
		case customer = <-tills[2].line:
			fmt.Println("Customer payed", customer.getTotal(), "DKK , at till 3")
			tills[2].cash += customer.getTotal()
		case customer = <-tills[3].line:
			fmt.Println("Customer payed", customer.getTotal(), "DKK , at till 4")
			tills[3].cash += customer.getTotal()
		case customer = <-tills[4].line:
			fmt.Println("Customer payed", customer.getTotal(), "DKK , at till 5")
			tills[4].cash += customer.getTotal()
		case customer = <-tills[5].line:
			fmt.Println("Customer payed", customer.getTotal(), "DKK , at till 6")
			tills[5].cash += customer.getTotal()
		case <-time.After(3 * time.Second):
			select {
			case <-closed:
				fmt.Println("Shop is closed 4 REALZ")
				fmt.Println("Till 1 made", tills[0].cash)
				fmt.Println("Till 2 made", tills[1].cash)
				fmt.Println("Till 3 made", tills[2].cash)
				fmt.Println("Till 4 made", tills[3].cash)
				fmt.Println("Till 5 made", tills[4].cash)
				fmt.Println("Till 6 made", tills[5].cash)
				return
			case <-time.After(time.Second):
				fmt.Println("Customers are slow... till is bored...")
			}
		}
	}
}
