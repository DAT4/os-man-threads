package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	shop := CreateAndOpenShop(
		"Netto", // Name
		8,       // Opening
		10,      // Closing
		6,       // Tills
		0,       // Funds
		50,      // Customer amount
		10,      // Line size pr till
	)

	go spawnCustomers(shop.customers, shop.door)

	openTills(&shop.wg, shop.tills, shop.customersDoneShopping, shop.door)

	go shopping(shop.customers, shop.customersDoneShopping)

	shop.wg.Wait() // Waiting for the guy to close the shop
}

func spawnCustomers(customers chan<- Customer, door chan bool) {
	for {
		time.Sleep(minutes(rand.Intn(10-2) + 2))
		select {
		case <-door:
			fmt.Println("Door is closed")
			return
		case customers <- CreateCustomer():
		case <-time.After(minutes(15)):
			fmt.Println("Shop is full!")
		}
	}
}

func shopping(customersEnterShop <-chan Customer, lineArea chan<- Customer) {
	var customer Customer
	for {
		select {
		case customer = <-customersEnterShop:
			time.Sleep(customer.shoppingDuration)
			go toLine(customer, lineArea)
		case <-time.After(minutes(15)):
			fmt.Println("No customers in shop")
			return
		}
	}
}

func toLine(customer Customer, lineArea chan<- Customer) {
	select {
	case lineArea <- customer:
	case <-time.After(minutes(15)):
		fmt.Println("Some error")
		return
	}
}

func openTills(wg *sync.WaitGroup, tills []*Till, customers <-chan Customer, door <-chan bool) {
	for _, till := range tills {
		wg.Add(1)
		go func(wait *sync.WaitGroup, line chan<- Customer, stream <-chan Customer, closed <-chan bool) {
			//THIS IS THE DISPACHER
			fmt.Println("RECV Line named:", line)
			defer wait.Done()
			for {
				select {
				case cus := <-customers:
					line <- cus
				case _, ok := <-closed:
					if ok {
						fmt.Println("some wird stff")
					} else {
						close(line)
						return
					}
				}
			}
		}(wg, till.line, customers, door)

		go func(t *Till) {
			//THIS IS THE WORKER
			fmt.Println("SEND Line named:", t.line)
			for {
				select {
				case customer, ok := <-t.line:
					t.Lock()
					if ok {
						t.cash += customer.getTotal()
					} else {
						fmt.Println("This till has earned", till.cash, "DKK")
						return
					}
					t.Unlock()
				}
			}
		}(till)
	}
}
