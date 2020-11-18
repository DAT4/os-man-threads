package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Till struct {
	sync.Mutex
	execution chan Customer
	cashFlow  chan int
	cash      int
	line      chan Customer
}

type Shop struct {
	sync.Mutex
	wg                    sync.WaitGroup
	name                  string
	opening               int
	closing               int
	door                  chan bool
	customers             chan Customer
	customersDoneShopping chan Customer
	openingHours          time.Duration
	tills                 []*Till
	funds                 int
}

/**
 * This function will create a shop and all the necesary
 * channels and so on. it is FactoryStyle
 **/
func CreateAndOpenShop(
	name string,
	opening int,
	closing int,
	tills int,
	funds int,
	customerAmount int,
	lineSize int,
) (shop Shop) {
	shop = Shop{
		name:                  name,
		opening:               opening,
		closing:               closing,
		openingHours:          hours(closing - opening),
		door:                  make(chan bool),
		customers:             make(chan Customer, customerAmount),
		customersDoneShopping: make(chan Customer, customerAmount),
	}

	for i := 0; i < tills; i++ {
		till := Till{
			execution: make(chan Customer),
			cashFlow:  make(chan int),
			line:      make(chan Customer, lineSize),
		}
		fmt.Println("CREA Line named:", till.line)
		shop.tills = append(shop.tills, &till)
	}

	go func(door chan<- bool) {
		time.Sleep(shop.openingHours)
		close(door)
	}(shop.door)

	return shop
}

type Grocery struct {
	price int
}

/**
 * This function will create the groceries that a
 * customer will use.
 **/
func CreateGrocery() (grocery Grocery) {
	grocery.price = rand.Intn(100-20) + 20
	return grocery
}

type Customer struct {
	shoppingDuration time.Duration
	groceries        []Grocery
}

/**
 * This function is used to generate a unique customer
 * with some groceries and a shoppin duration.
 **/
func CreateCustomer() (customer Customer) {
	customer.shoppingDuration = minutes(rand.Intn(25-5) + 5)

	groceriesAmount := rand.Intn(30-5) + 5
	for i := 0; i < groceriesAmount; i++ {
		customer.groceries = append(customer.groceries, CreateGrocery())
	}
	return customer
}

/**
 * This method will get the total amount that a customer has
 * to pay for his groceries.
 **/
func (self *Customer) getTotal() int {
	var total int
	for _, grocery := range self.groceries {
		total += grocery.price
	}
	return total
}
