package main

import (
	"errors"
	"sync"
)

type pool struct {
	idle     []iPoolObject
	active   []iPoolObject
	capacity int
	mux      *sync.Mutex
}

type connection struct {
	id string
}

type iPoolObject interface {
	getID() string
}

func initPool(connections []iPoolObject) (*pool, error) {
	if len(connections) == 0 {
		return nil, errors.New("Cannot create pool of length 0")
	}

	return &pool{
		idle:     connections,
		active:   make([]iPoolObject, 0),
		capacity: len(connections),
		mux:      &sync.Mutex{},
	}, nil
}

func (p *pool) loan() (iPoolObject, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if len(p.idle) == 0 {
		return nil, errors.New("no have idle")
	}

	obj := p.idle[0]
	p.idle = p.idle[1:]
	p.active = append(p.active, obj)

	return obj, nil
}

func (p *pool) receive(obj iPoolObject) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	err := p.remove(obj)
	if err != nil {
		return err
	}

	p.idle = append(p.idle, obj)

	return nil
}

func (p *pool) remove(target iPoolObject) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	for i, obj := range p.active {
		if obj.getID() == target.getID() {
			p.active[0], p.active[i] = p.active[i], p.active[0]
			p.active = p.active[1:]
			return nil
		}
	}

	return errors.New("target is not exist in pool")
}
