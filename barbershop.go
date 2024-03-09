package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type BarberShop struct {
	clients   chan struct{}
	barbers   chan struct{}
	waiting   chan struct{}
	closeShop chan struct{}
	wg        sync.WaitGroup
}

const (
	noOfBarbers       = 2
	noOfWaitingChairs = 5
	openingHours      = 5 * time.Second
)

func (bs *BarberShop) closeBarberShop(closingTime time.Duration) {
	time.Sleep(closingTime)
	close(bs.closeShop)
}

func (bs *BarberShop) OpenShop() {
	fmt.Println("Welcome to Harley Saloon!")
	go bs.closeBarberShop(openingHours)

	for i := 0; i < noOfBarbers; i++ {
		bs.wg.Add(1)
		go bs.StartJob()
	}

	go bs.clientVisits()
	bs.wg.Wait()
	fmt.Println("We are closed! Please visit tomorrow.")
}

func (bs *BarberShop) clientVisits() {
	for {
		select {
		case <-bs.closeShop:
			return
		default:
			visitInterval := time.Duration(rand.Intn(2) * int(time.Second))
			time.Sleep(visitInterval)
			bs.clients <- struct{}{}
			select {
			case bs.waiting <- struct{}{}:
				fmt.Println("Client arrived and is waiting.")
			default:
				fmt.Println("Client left, no space available.")
			}

		}
	}
}

func (b *BarberShop) StartJob() {
	defer b.wg.Done()
	for {
		select {
		case <-b.closeShop:
			return
		case <-b.clients:
			select {
			case b.barbers <- struct{}{}:
				select {
				case <-b.waiting:
					fmt.Println("Haircut started.")
					time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
					fmt.Println("Haircut done.")
				default:
					fmt.Println("Saloon is Empty.")
				}
				<-b.barbers
			default:
				fmt.Println("Barber is sleeping.")
			}
		}
	}
}

func main() {
	bs := &BarberShop{
		clients:   make(chan struct{}, noOfWaitingChairs),
		barbers:   make(chan struct{}, noOfBarbers),
		waiting:   make(chan struct{}, noOfWaitingChairs),
		closeShop: make(chan struct{}),
	}
	bs.OpenShop()

}
