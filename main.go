package main

import (
	"bufio"
	"fmt"
	"os"
)

var wordSize int

type node struct {
	value  string
	next   []*node
	prev   []*node
	used   bool
	merged bool
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	wordSize = len(lines[0])
	return lines, nil
}

func nodes(lines []string) []*node {
	prevMap := make(map[string][]*node)
	nextMap := make(map[string][]*node)
	nodes := []*node{}
	for _, line := range lines {
		new_node := &node{value: line, used: false, merged: false}
		nodes = append(nodes, new_node)
		if _, ok := prevMap[line[1:]]; !ok {
			prevMap[line[1:]] = []*node{}
		}
		if _, ok := nextMap[line[:len(line)-1]]; !ok {
			nextMap[line[:len(line)-1]] = []*node{}
		}
		prevMap[line[1:]] = append(prevMap[line[1:]], new_node)
		nextMap[line[:len(line)-1]] = append(nextMap[line[:len(line)-1]], new_node)
	}
	for _, n := range nodes {
		n.next = nextMap[n.value[1:]]
		n.prev = prevMap[n.value[:len(n.value)-1]]
	}
	return nodes
}

func merge(a *node) {
	if len(a.prev) == 1 {
		if len(a.prev[0].next) == 1 {
			merge(a.prev[0])
			return
		}
	}
	if len(a.next) == 1 {
		if len(a.next[0].prev) == 1 {
			next := a.next[0]
			a.value = a.value + next.value[wordSize-1:]
			a.next = next.next
			next.used = true
			merge(a)
			return
		}
	}
	if len(a.next) > 0 {
		for _, n := range a.next {
			for i, p := range n.prev {
				if p.value[:wordSize] == a.value[len(a.value)-wordSize:] {
					n.prev[i] = a
				}
			}
		}
	}
	a.merged = true
}

func mergeAll(nodes []*node, new_nodes []*node) []*node {
	if len(nodes) == 0 {
		return new_nodes
	}
	if nodes[0].merged {
		new_nodes = append(new_nodes, nodes[0])
		nodes = nodes[1:]
		return mergeAll(nodes, new_nodes)
	}
	if nodes[0].used {
		nodes = nodes[1:]
		return mergeAll(nodes, new_nodes)
	}
	merge(nodes[0])
	return mergeAll(nodes, new_nodes)
}

func resetUsed(nodes []*node) {
	for _, n := range nodes {
		n.used = false
	}
}

func getStarted(nodes []*node) []*node {
	started := []*node{}
	for _, n := range nodes {
		if len(n.prev) == 0 {
			started = append(started, n)
		}
	}
	return started
}

func getSolutions(n node) []node {
	nodes := []node{}
	if len(n.next) == 0 {
		return append(nodes, n)
	}
	for _, next := range n.next {
		next.value = n.value + next.value[len(next.value)-1:]
		nodes = append(nodes, getSolutions(*next)...)
	}
	return nodes
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>", os.Args[0])
		return
	}
	filename := os.Args[1]
	lines, err := readLines(filename)
	if err != nil {
		fmt.Println("Error reading file")
		return
	}
	nodes := nodes(lines)
	for _, n := range nodes {
		fmt.Printf("%s -> ", n.value)
		for _, next := range n.next {
			fmt.Printf("%s ", next.value)
		}
		fmt.Print(" <- ")
		for _, prev := range n.prev {
			fmt.Printf("%s ", prev.value)
		}
		fmt.Println()
	}
	fmt.Println("Merging...")
	merged_nodes := mergeAll(nodes, []*node{})
	for i, n := range merged_nodes {
		fmt.Printf("%d: %s -> ", i, n.value)
		for _, next := range n.next {
			fmt.Printf("%s ", next.value)
		}
		fmt.Print(" <- ")
		for _, prev := range n.prev {
			fmt.Printf("%s ", prev.value)
		}
		fmt.Println()
	}
	started := getStarted(merged_nodes)
	for _, n := range started {
		fmt.Println(n.value)
	}
	fmt.Println(len(merged_nodes))
	fmt.Println(len(started))
	for _, n := range started {
		solutions := getSolutions(*n)
		for _, s := range solutions {
			fmt.Printf("%s\n", s.value)
		}
	}
}
