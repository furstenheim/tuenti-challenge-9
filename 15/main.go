package main

import (
	"net"
	"log"
	"fmt"
	"bufio"
	"strconv"
	"strings"
)


func main () {
	conn, err := net.Dial("tcp", "52.49.91.111:1888")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()


	query := "TEST\n"
	fmt.Fprint(conn, query)
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	handleError(err)
	nCases, err := strconv.Atoi(strings.Fields(line)[0])
	handleError(err)
	for i := 0; i < nCases; i++ {
		c := parseCase(reader)
		log.Println(c)
		s := c.solve()
		s.printSolution(conn)

	}
}

type ChainId int
type PersonId int
type Chain []PersonId
type Case struct {
	tableSize int
	restrictions []Restriction
	chains []Chain
	person2Chain map[PersonId]ChainId
	visitedRestrictions map[Restriction]bool
}
type Solution struct {
	positions [8][]string
	tableSize int
}

type Restriction struct {
	u1, u2 PersonId
}

func parseCase (reader *bufio.Reader) Case {
	sizeLine, err := reader.ReadString('\n')
	log.Println("Received", sizeLine)
	handleError(err)
	size, err := strconv.Atoi(strings.Fields(sizeLine)[0])
	if err != nil {
		for true {
			validatorLine, err := reader.ReadString('\n')
			handleError(err)
			log.Println(validatorLine)
		}
	}
	handleError(err)
	restrictionsLine, err := reader.ReadString('\n')
	log.Println("Received", restrictionsLine)
	handleError(err)
	nRestrictions, err := strconv.Atoi(strings.Fields(restrictionsLine)[0])
	restrictions := make([]Restriction, nRestrictions)
	for i := 0; i < nRestrictions; i ++ {
		restrictionLine, err := reader.ReadString('\n')
		log.Println("Received", restrictionLine)
		handleError(err)
		restriction1, err := strconv.Atoi(strings.Fields(restrictionLine)[0])
		handleError(err)
		restriction2, err := strconv.Atoi(strings.Fields(restrictionLine)[1])
		restrictions[i] = Restriction{
			u1: PersonId(restriction1),
			u2: PersonId(restriction2),
		}
	}
	return Case{
		tableSize: size,
		restrictions: restrictions,
	}
}

func (s *Solution) printSolution (conn net.Conn) {
	for _, t := range (s.positions) {
		text := fmt.Sprintf("%s", strings.Join(t, ","))
		log.Println("Returning", text)
		fmt.Fprintln(conn, text)
	}
}

func (c * Case) solve () Solution {
	if len(c.restrictions) == 0 {
		sittings := [8][]string{}
		for i, _ := range(sittings) {
			for j := 0; j < c.tableSize; j++ {
				sittings[i] = append(sittings[i], strconv.Itoa(8 * i + j + 1))
			}
		}
		return Solution{
			positions: sittings,
			tableSize: c.tableSize,
		}
	}
	c.joinRestrictions()
	log.Fatal("Unknown")
	return Solution{}

}

func (c * Case) joinRestrictions () {
	for _, r := range (c.restrictions) {
		if _, handledRestriction := c.visitedRestrictions[r]; handledRestriction {
			continue
		}
		p1, p2 := r.u1, r.u2
		c1, ok1 := c.person2Chain[p1]
		c2, ok2 := c.person2Chain[p2]
		if !ok1 && !ok2 {
			nextId := ChainId(len(c.chains))
			c.chains = append(c.chains, Chain{p1, p2})
			c.person2Chain[p1] = nextId
			c.person2Chain[p2] = nextId
		} else if ok1 && !ok2 {
			c.appendToChain(c1, p1, p2)
		} else if !ok1 && ok2 {
			c.appendToChain(c2, p2, p1)
		} else {
			if c.chains[c1][0] == p1 {
				c.moveEndChain(c2, c1, p2, p1)
			} else if c.chains[c2][0] == p2 {
				c.moveEndChain(c1, c2, p1, p2)
			} else {
				log.Fatal("Could not move chains", c1, c2, p1, p2)
			}
		}
		c.visitedRestrictions[r] = true
		c.visitedRestrictions[Restriction{u1: r.u2, u2: r.u1}] = true
	}
}

func (c * Case) appendToChain (ch ChainId, attachTo PersonId, attach PersonId) {
	c.person2Chain[attach] = ch
	if c.chains[ch][0] == attachTo {
		c.chains[ch] = append([]PersonId{attach}, c.chains[ch]...)
	} else if c.chains[ch][len(c.chains[ch]) - 1] == attachTo {
		c.chains[ch] = append(c.chains[ch], attach)
	} else {
		log.Fatal("Id was not correct", c.chains[ch], attachTo)
	}

}

func (c * Case) moveEndChain (chRemains ChainId, chMoves ChainId, pRemains PersonId, pMoves PersonId) {
	c.chains[chRemains] = append(c.chains[chRemains], c.chains[chMoves]...)
	delete(c.person2Chain, pMoves)
	delete(c.person2Chain, pRemains)
	// reassign it could be pmove or premains
	c.person2Chain[c.chains[chRemains][0]] = chRemains
	c.person2Chain[c.chains[chRemains][len(c.chains[chRemains]) - 1]] = chRemains
	c.chains[chMoves] = nil
}


func handleError (err error){
	if err != nil {
		log.Fatal(err)
	}
}

