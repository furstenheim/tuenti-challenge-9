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

type Case struct {
	tableSize int
	restrictions []Restriction
}
type Solution struct {
	positions [8][]string
	tableSize int
}

type Restriction struct {
	u1, u2 int
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
			u1: restriction1,
			u2: restriction2,
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
	log.Fatal("Unknown")
	return Solution{}

}


func handleError (err error){
	if err != nil {
		log.Fatal(err)
	}
}

