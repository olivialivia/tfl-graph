package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

type Link struct {
	Source string
	Target string
	Line   string
	Cost   int
}

type Node struct {
	ID string
}

type Graph struct {
	Nodes []Node
	Links []Link
}

func (g *Graph) findLinks(station string) (result []Link) {
	for i := 0; i < len(g.Links); i++ {
		l := g.Links[i]

		if station == l.Source || station == l.Target {
			result = append(result, l)
		}
	}
	return
}

func linkExists(links []Link, link Link) bool {
	for i := 0; i < len(links); i++ {
		if link.Source == links[i].Source && link.Target == links[i].Target && link.Line == links[i].Line {
			return true
		}
	}
	return false
}

func main() {
	var startStation = flag.String("start", "source", "a start station")
	var finalStation = flag.String("dest", "target", "an end station")
	flag.Parse()
	fmt.Println(*startStation, *finalStation)
	start := *startStation
	dest := *finalStation

	contents, err := ioutil.ReadFile("tfl-graph.json")
	if err != nil {
		log.Fatal(err)
	}
	graph := Graph{}
	err = json.Unmarshal(contents, &graph)
	if err != nil {
		log.Fatal(err)
	}

	var SSfound bool
	var FSfound bool

	for i := 0; i < len(graph.Nodes); i++ {
		l := graph.Nodes[i]

		if l.ID == start {
			SSfound = true
		}
		if l.ID == dest {
			FSfound = true
		}
	}

	if !SSfound || !FSfound {
		log.Fatal("station not found")
	}

	links := graph.findLinks(start)

	var visited []Link
	var unvisited []Link
	var route []Link

	unvisited = append(unvisited, links...)

Walk:
	for len(unvisited) > 0 {
		sort.Slice(unvisited, func(i, j int) bool {
			return unvisited[i].Cost < unvisited[j].Cost
		})
		currLink := unvisited[0]
		unvisited = unvisited[1:]
		visited = append(visited, currLink)

		nextLinks := append(graph.findLinks(currLink.Source), graph.findLinks(currLink.Target)...)

		for _, l := range nextLinks {
			if linkExists(visited, l) {
				continue
			}
			l.Cost = currLink.Cost + 1

			if currLink.Line != l.Line {
				l.Cost += 2
			}
			unvisited = append(unvisited, l)

			if dest == l.Source || dest == l.Target {
				route = append(route, l)
				break Walk
			}
		}
	}

	if len(route) == 0 {
		log.Fatal("no route exists")
	}

	station := dest

	for i := 0; i < len(route); i++ {
		l := route[i]

		if l.Source == station {
			station = l.Target
		} else {
			station = l.Source
		}

		for _, next := range visited {
			if next.Cost >= l.Cost {
				continue
			}
			if next.Source == station || next.Target == station {
				route = append(route, next)
				break
			}
		}
	}

	//fmt.Printf("%+v", visited)

	res, err := json.MarshalIndent(route, "", "    ")
	if err != nil {
		panic(err)
	}
	

	fmt.Printf("\n \n %v", string(res))
}

