package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
)

type pool struct {
	mux      *sync.Mutex
	idle     []iPoolObject
	active   []iPoolObject
	capacity int
}

type iPoolObject interface {
	getID() string
}

type connection struct {
	id string
}

func (c *connection) getID() string {
	return c.id
}

func initPool(poolObject []iPoolObject) (*pool, error) {
	if len(poolObject) == 0 {
		return nil, errors.New("Cannot create a pool of 0 length")
	}

	active := make([]iPoolObject, 0)
	pool := &pool{
		idle:     poolObject,
		active:   active,
		capacity: len(poolObject),
		mux:      new(sync.Mutex),
	}

	return pool, nil
}

func (p *pool) loan() (iPoolObject, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.idle) == 0 {
		return nil, errors.New("No pool object free. Please request after sometime")
	}

	obj := p.idle[0]
	p.idle = p.idle[1:]
	p.active = append(p.active, obj)
	fmt.Printf("Loan Pool object with ID: %s\n", obj.getID())

	return obj, nil
}

func (p *pool) receive(target iPoolObject) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	err := p.remove(target)
	if err != nil {
		return nil
	}

	p.idle = append(p.idle, target)
	fmt.Printf("Return Pool Object with ID: %s\n", target.getID())

	return nil
}

func (p *pool) remove(target iPoolObject) error {
	currentActiveLength := len(p.active)
	for i, obj := range p.active {
		if obj.getID() == target.getID() {
			p.active[currentActiveLength-1], p.active[i] = p.active[i], p.active[currentActiveLength-1]
			p.active = p.active[:currentActiveLength]
			return nil
		}
	}

	return errors.New("Target pool object doesn't belong to the pool")
}

func main() {
	connections := make([]iPoolObject, 0)
	for i := 0; i < 3; i++ {
		c := &connection{id: strconv.Itoa(i)}
		connections = append(connections, c)
	}

	pool, err := initPool(connections)
	if err != nil {
		log.Fatalf("Init Pool Error: %s", err)
	}

	conn1, err := pool.loan()
	if err != nil {
		log.Fatalf("Pool Loan Error: %s", err)
	}

	conn2, err := pool.loan()
	if err != nil {
		log.Fatalf("Pool loan Error: %s", err)
	}

	pool.receive(conn1)
	pool.receive(conn2)
}
