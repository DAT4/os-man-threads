package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/**
 * The main flow of the program is defined here.
 * The flow is to create and open the shop, and then
 * start runnig the concurent tasks, and waiting for
 * everything to finish until it will terminate.
 */
func main() {
	// The shop is created with a factory method
	shop := CreateAndOpenShop(
		"Netto", // Name is Netto
		8,       // Opening at 8
		22,      // Closing at 22
		6,       // The shop has 6 tills
		0,       // The shop has no funds
		50,      // Maximum customer amount
		10,      // Maximum line size pr till
	)

	// The program will concurrently spawn customers to the channel in the shop.
	go spawnCustomers(shop.customers, shop.door)

	// The tills will open sequentially.
	openTills(&shop.wg, shop.tills, shop.customersDoneShopping, shop.door)

	// Each customer in the shop will go to do shopping in his own time.
	go shopping(shop.customers, shop.customersDoneShopping)

	// The shop will wait for all tills to close before terminating.
	shop.wg.Wait()
}

/**
 * This method will pass generated customers into the
 * channel of customers in the shop. if the door of the
 * shop is not closed. Else it will print and terminate
 **/
func spawnCustomers(customers chan<- Customer, door chan bool) {
	for {
		/**
		 * Wait for a random interval between 10 and 2 minutes
		 * Time is scaled down 1:3600 (so one hour is one second)
		 **/
		time.Sleep(minutes(rand.Intn(10-2) + 2))
		select {
		case <-door:
			fmt.Println("Door is closed")
			return
		case customers <- CreateCustomer():
		case <-time.After(minutes(15)):
			/**
			 * The shop has a max capacity and customers channel
			 * will block if it is reached. Here it will be
			 * announced if the shop has been full for 15 minutes
			 **/
			fmt.Println("Shop is full!")
		}
	}
}

/**
 * This function will pass customers into their own individual flow
 **/
func shopping(customersEnterShop <-chan Customer, lineArea chan<- Customer) {
	for {
		select {
		case customer := <-customersEnterShop:
			// customer is thrown into his own routine
			go toLine(customer, lineArea)
		case <-time.After(minutes(15)):
			fmt.Println("No customers in shop")
			return
		}
	}
}

/**
 * The customer will take his time shopping and then he will go to
 * stand in line to checkout his groceries.
 **/
func toLine(customer Customer, lineArea chan<- Customer) {
	time.Sleep(customer.shoppingDuration)
	select {
	case lineArea <- customer:
	case <-time.After(minutes(15)):
		fmt.Println("Some error")
		return
	}
}

/**
 * This function will create a worker for dispaching customers
 * from the shop channel to each of the channels in the lines
 * also it will create another worker for executing the customers
 * in the lines.
 **/
func openTills(
	wg *sync.WaitGroup,
	tills []*Till,
	customers <-chan Customer,
	door <-chan bool,
) {
	for _, till := range tills {
		// The till adds a worker to the waitgroup of the shop
		wg.Add(1)
		/**
		 * This is the dispacher function, and it will collect
		 * the customers from the customer stream and put them
		 * into the lines.
		 *
		 * NOTE:	Something is wrong where I think and the
		 *			concurency is not running well.
		 **/
		go func(
			wait *sync.WaitGroup,
			line chan<- Customer,
			stream <-chan Customer,
			closed <-chan bool,
		) {
			/**
			 * When the function returns, then the worker will
			 * notify the shop that he is done.
			 **/
			defer wait.Done()
			for {
				select {
				case customer := <-customers:
					line <- customer
				case _, ok := <-closed:
					/**
					 * If the shop door is closed and there is no
					 * more customers, then the line will close.
					 **/
					if ok {
						fmt.Println("Someone is kidding you.. you still have to work!")
					} else {
						close(line)
						return
					}
				}
			}
		}(wg, till.line, customers, door)

		/**
		 * This is a worker function which will take one customer each
		 * from the line of a till and execute them each in one minute.
		 **/
		go func(t *Till) {
			for {
				select {
				case customer, ok := <-t.line:
					/**
					 * Mutual exclution is used to update the till cash register
					 * in a threadsafe way.
					 **/
					t.Lock()
					if ok {
						// If open then take cash
						t.cash += customer.getTotal()
					} else {
						// If closed, then count the total cash.
						fmt.Println("This till has earned", till.cash, "DKK")
						return
					}
					t.Unlock()
				}
			}
		}(till)
	}
}
