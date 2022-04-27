package main

import (
	"container/list"
)

const (
	MaxDim =	80
	MinTime =	800
	VarTime =	1500
)

type (
	pending struct {
		todo	*list.List
		done	*list.List
		roots	int
	}
	link struct {
		edge	*node
		count	int
	}
	node struct {
		infodim	int
		ccount	int
		hooks	[]link
	}
)

var	(
	treelist	[MaxDim]pending
	plCount		int
)

func newPlayer() bool {
	if rnd.Float64() < float64(rate.x) / float64(rate.y) {
		return true
	}
	return false
}

func initTreelist() {
	for i := 0; i < MaxDim; i++ {
		treelist[i] = pending{list.New(), list.New(), 0}
	}
}

func newNode(d int) *node {
	var edges int
	var res node
	res.infodim = d
	if d == 0 {
		edges = 1 + 0
		res.hooks = make([]link, 1)
		res.hooks[0] = link{nil, MinTime + rnd.Intn(VarTime)}
	} else {
		switch m := d % 5; {
			case m == 0:
				edges = 1 + 3
			case m == 1:
				edges = 1 + 1
			default:
				edges = 1 + 2
		}
		for i := 0; i < edges; i++ {
			res.hooks = append(res.hooks, link{nil, MinTime + rnd.Intn(VarTime)})
		}
	}
	res.ccount = len(res.hooks) - 1
	return &res
}

func (n *node) doCycle() int {
	if n.isDone() {
		return Ignore
	}
	if n.ccount < 0 {
		return Ignore
	}
	if n.hooks[n.ccount].count > 0 {
		n.hooks[n.ccount].count--
		if n.hooks[n.ccount].count == 0 {
			if n.ccount == 0 {
				return OpenUp
			} else {
				return OpenDown
			}
		}
	}
	return Ignore
}

func (n *node) isDone() bool {
	done := true
	for _, l := range n.hooks {
		done = done && (l.edge != nil) 
	}
	return done
}

func (n *node) putMom() {
	var stop bool
	var nmom *node
	var found *list.Element
	for mom := treelist[n.infodim + 1].todo.Front(); (mom != nil) && !stop; mom = mom.Next() {
		if mom.Value.(*node).ccount > 0 {
			stop = true
			found = mom
			nmom = found.Value.(*node)
		}
	}
	if !stop {
		nmom = newNode(n.infodim + 1)
		treelist[n.infodim + 1].todo.PushBack(nmom)
		treelist[nmom.infodim].roots++
	}
	n.hooks[0].edge = nmom
	nmom.hooks[nmom.ccount].edge = n
	nmom.ccount--
	treelist[n.infodim].roots--
	if nmom.isDone() {
		treelist[nmom.infodim].done.PushBack(treelist[nmom.infodim].todo.Remove(found))
	}
}

func (n *node) putChild() {
	var stop bool
	var nchild *node
	var found *list.Element
	for child := treelist[n.infodim - 1].todo.Front(); (child != nil) && !stop; child = child.Next() {
		if child.Value.(*node).hooks[0].edge == nil {
			stop = true
			found = child
			nchild = found.Value.(*node)
			treelist[nchild.infodim].roots--
		}
	}
	if !stop {
		nchild = newNode(n.infodim - 1)
		treelist[n.infodim - 1].todo.PushBack(nchild)
	}
	n.hooks[n.ccount].edge = nchild
	nchild.hooks[0].edge = n
	n.ccount--
	if nchild.isDone() {
		treelist[nchild.infodim].done.PushBack(treelist[nchild.infodim].todo.Remove(found))
	}
}

const (
	Ignore = iota
	OpenUp
	OpenDown
)

func doCycle() {
	if newPlayer() {
		n := newNode(0)
		treelist[0].todo.PushBack(n)
		treelist[0].roots++
		plCount++
		showPlayers()
	}
	for i := 0; i < MaxDim - 5; i++ {
		var next *list.Element
		for e := treelist[i].todo.Front(); e != nil; e = next {
			next = e.Next()
			n := e.Value.(*node)
			res := n.doCycle()
			switch res {
			case OpenUp:
				n.putMom()
			case OpenDown:
				if n.infodim > 2 {
					n.putChild()
				}
			}
			if n.isDone() {
				treelist[i].todo.Remove(e)
				treelist[i].done.PushBack(n)
			}
		}
	}
}
//	-- Leer canal
//	-- Si es 1,
//		- crear un nuevo espacio personal
//			* (Nodo D0 + Nodo D1)
//			* Poner contador en D1
//		- añadirlo a treelist
//	-- Ejecutar ciclo
//		Para cada dimensión
//			Para cada elemento e de la lista todo
//			Invocar e.doCycle (Decrementar el contador adecuado y reportar status)
//			Si el resultado es,
//			OpenUp:
//				Buscarle una madre
//				Pasarlo a la lista done
//			OpenDown:
//				Buscarle un hijo
