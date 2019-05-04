package main

import (
	"net"
	"log"
	"fmt"
	"bufio"
	"strconv"
	"strings"
	"sort"
)


func main () {
	conn, err := net.Dial("tcp", "52.49.91.111:1888")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()


	query := "SUBMIT\n"
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
	for true {
		line, err := reader.ReadString('\n')
		log.Println(line)
		handleError(err)
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
	restrictedPeople map[PersonId]bool
	partitionCache map[PartitionKey]PartitionEl
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
		visitedRestrictions: map[Restriction]bool{},
		restrictedPeople: map[PersonId]bool{},
		person2Chain: map[PersonId]ChainId{},
		partitionCache: map[PartitionKey]PartitionEl{},
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
	c.joinRestrictions()
	freePeople := c.getFreePeople()
	chains := c.getChains()
	for _ , p := range(freePeople) {
		chains = append(chains, Chain{p})
	}
	lengths := make([]int, len(chains))
	for i, c := range(chains) {
		lengths[i] = len(c)
	}
	distribution := c.partitionSum(lengths)


	log.Println("Start of loop", freePeople, chains)
	sittings := [8][]string{}
	for i, _ := range(sittings) {
		currentSittings := []string{}
		for _, index := range (distribution[i]) {
			for _, person := range(chains[index]) {
				currentSittings = append(currentSittings, strconv.Itoa(int(person)))
			}
		}
		sittings[i] = currentSittings
		/*

		remainingSits := c.tableSize
		for remainingSits > 0 {
			found, nextChain, newChains := c.getBiggestChain(chains, remainingSits)
			log.Println("Chain size", remainingSits, found, len(nextChain), len(chains), len(newChains))
			if !found {
				break
			}
			chains = newChains
			remainingSits -= len(nextChain)
			for _, p := range (nextChain) {
				currentSittings = append(currentSittings, strconv.Itoa(int(p)))
			}
		}
		for remainingSits > 0 {
			var nextPerson PersonId
			log.Println("free people", freePeople, remainingSits, sittings, chains)
			nextPerson, freePeople = freePeople[len(freePeople) - 1], freePeople[:len(freePeople) - 1]
			currentSittings = append(currentSittings, strconv.Itoa(int(nextPerson)))
			remainingSits--
		}
		*/
	}

	return Solution{
		positions: sittings,
		tableSize: c.tableSize,
	}

}

type Visits [8 * 24]int
type PartitionKey struct {
	remainingSpace int
	filling [8]int
	visits Visits
}

type PartitionEl struct {
	found bool
	distribution Visits
}

func (c * Case) partitionSum (lengths []int ) [8][]int {
	visits := Visits{}
	for i, _ :=range(visits) {
		visits[i] = - 1
	}
	el := c.partitionBranch(PartitionKey{
		remainingSpace: c.tableSize * 8,
		filling: [8]int{},
		visits: visits,
	}, lengths)
	if !el.found {
		log.Fatal("Could not distribute")
	}
	result := [8][]int{}
	for i, _ := range(result) {
		result[i] = []int{}
	}
	log.Println(lengths, el.distribution)
	for i, _ := range(lengths) {
		table := el.distribution[i]
		result[table] = append(result[table], i)
	}
	return result
}

func (c * Case) partitionBranch (partitionStatus PartitionKey, lengths []int) PartitionEl {
	if cache, ok := c.partitionCache[partitionStatus]; ok {
		return cache
	}

	// log.Println(partitionStatus, lengths)
	if partitionStatus.remainingSpace == 0 {
		return PartitionEl{
			found: true,
			distribution: partitionStatus.visits,
		}
	}
	for i, v := range(lengths) {
		if partitionStatus.visits[i] > -1 {
			continue // already assigned
		}
		for table := 0; table < 8; table ++ {
			if partitionStatus.filling[table] + v <= c.tableSize {
				newVisits := cloneVisits(partitionStatus.visits)
				newVisits[i] = table
				newFillings := cloneFillings(partitionStatus.filling)
				newFillings[table] = newFillings[table] + v
				branchEl := c.partitionBranch(PartitionKey{
					remainingSpace: partitionStatus.remainingSpace - v,
					visits: newVisits,
					filling: newFillings,
				}, lengths)
				if branchEl.found {
					return branchEl
				}

			}
		}

	}
	return PartitionEl{}

}
func cloneLengths (lengths []int) []int {
	clone := make([]int, len(lengths))
	for i, v := range(lengths) {
		clone[i] = v
	}
	return clone
}
func cloneVisits(distribution Visits) Visits {
	distributionCopy := distribution
	return distributionCopy
}
func cloneFillings (fillings [8]int) [8]int {
	clone := fillings
	return clone
}


func (c * Case) getBiggestChain (chains []Chain, maxPossibleSize int) (found bool, ch Chain, newChains []Chain) {
	for i, ch := range(chains) {
		log.Println("Chain size", len(ch), maxPossibleSize)
		if len(ch) <= maxPossibleSize {
			return true, ch, append(chains[:i], chains[i + 1: ]...)
		}
	}
	return false, Chain{}, chains
}
func (c * Case) getFreePeople () []PersonId {
	freePeople := []PersonId{}
	for i := 1; i <= c.tableSize * 8; i++ {
		if _, ok := c.restrictedPeople[PersonId(i)]; ! ok {
			freePeople = append(freePeople, PersonId(i))
		}
	}
	return freePeople
}

func (c * Case) getChains () []Chain {
	chains := []Chain{}
	for _, c := range (c.chains) {
		if c != nil {
			chains = append(chains, c)
		}
	}
	// bigger slices at the beginning
	sort.Slice(chains, func (i, j int) bool {
		return len(chains[i]) > len(chains[j])
	})
	return chains
}

func (c * Case) joinRestrictions () {
	for _, r := range (c.restrictions) {
		if _, handledRestriction := c.visitedRestrictions[r]; handledRestriction {
			continue
		}
		p1, p2 := r.u1, r.u2
		c.restrictedPeople[p1] = true
		c.restrictedPeople[p2] = true
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
			if c1 == c2 {
				continue
			}
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

