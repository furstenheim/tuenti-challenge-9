package main

import (
	"log"
	"bufio"
	"strings"
	"strconv"
	"os"
	"sort"
	"fmt"
	"math/bits"
)


func main () {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	firstLineFields := strings.Fields(line)
	a := parseAlmanac()
	numberOfCases, err := strconv.Atoi(firstLineFields[0])
	if (err != nil) {
		log.Fatal(err)
	}
	for i := 0; i < numberOfCases; i ++ {
		log.Println(i)
		c := parseCase(reader, a)
		s := c.solve(&a)
		s.printResult(i)
	}
}

func parseCase (reader *bufio.Reader, a Almanac) Case {
	line, err := reader.ReadString('\n')
	handleError(err)
	fields := strings.Fields(line)
	gold, err := strconv.Atoi(fields[0])
	handleError(err)
	id := a.almanacCharactersMap[fields[1]]
	nSkills, err := strconv.Atoi(fields[2])
	handleError(err)
	skills := SkillMask{}
	for i := 0; i < nSkills; i++ {
		skills.addSkill(a.skillsMap[fields[3 + i]])
	}
	return Case {
		id: id,
		gold: gold,
		skills: skills,
	}
}

type Case struct {
	gold int
	id CharacterId
	skills SkillMask
}

type Solution struct {
	found bool
	gold int
}

func (s *Solution) printResult (i int) {
	var text string
	if !s.found {
		text = fmt.Sprintf("Case #%d: IMPOSSIBLE", i + 1)
	} else {
		text = fmt.Sprintf("Case #%d: %d", i + 1, s.gold)
	}
	fmt.Println(text)
}


func (c *Case) solve (a * Almanac) Solution {
	found, gold := a.branchSolve(c.gold, c.id, c.skills)
	return Solution{found: found, gold: gold}
}

func (a * Almanac) branchSolve (remainingGold int, id CharacterId, missingSkills SkillMask) (found bool, gold int) {
	almanacCharacter := a.almanacCharacters[id]
	remainingSkills := missingSkills.and(almanacCharacter.skills.neg())
	currentBestGold := 200000 // Max in problem is 100000
	// log.Println(remainingSkills, almanacCharacter, missingSkills, a.skills[68], almanacCharacter.expandedSkills)
	if remainingSkills.isZero() && almanacCharacter.gold <= remainingGold {
		found = true
		currentBestGold = almanacCharacter.gold
	}

	possibleCombinations := remainingSkills.split()
	for _, comb := range(almanacCharacter.combinations) {
		c1 := comb.char1
		c2 := comb.char2
		char1 := a.almanacCharacters[c1]
		char2 := a.almanacCharacters[c2]
		combinedMaxSkills := char1.expandedSkills.or(char2.expandedSkills)
		// Combination might hold all necessary skills
		if !combinedMaxSkills.neg().and(remainingSkills).isZero() {
			// log.Println("Could not reach combination", char1, char2, remainingSkills, combinedMaxSkills.neg().and(remainingSkills))
			continue // cannot possibly reach with this combination
		}
		for _, comb := range (possibleCombinations) {
			if !char1.expandedSkills.neg().and(comb[0]).isZero() ||
				!char2.expandedSkills.neg().and(comb[1]).isZero() {
				// log.Println("There was no sub combination")
				continue // This combination wouldn't work
			}
			branch1Found, branch1Gold := a.branchSolve(remainingGold, c1, comb[0])
			branch2Found, branch2Gold := a.branchSolve(remainingGold - branch1Gold, c2, comb[1])
			branchFound := branch1Found && branch2Found
			branchGold := branch1Gold + branch2Gold
			// log.Println("Branch", branch1Found, branch2Found, branch1Gold, branch2Gold)
			if branchFound && branchGold < currentBestGold {
				found = true
				currentBestGold = branchGold
			}
		}
	}
	return found, currentBestGold
}



func parseAlmanac () Almanac {
	f, err := os.Open("./almanac.data")
	handleError(err)
	reader := bufio.NewReader(f)
	line, err := reader.ReadString('\n')
	handleError(err)
	nCharacters, err := strconv.Atoi(strings.Fields(line)[0])
	handleError(err)
	nCombinations, err := strconv.Atoi(strings.Fields(line)[1])
	handleError(err)
	a := Almanac{
		currentSkillId: 1,
		currentCharacterId: 1,
		skillsMap: map[string]SkillId{},
		almanacCharactersMap: map[string]CharacterId{},
		rawCombinations: []RawCombination{},
	}

	for i := 0; i < nCharacters; i++ {
		a.parseCharacter(reader)
	}
	for i := 0; i< nCombinations; i++ {
		a.parseCombination(reader)
	}

	sort.Slice(a.rawCombinations, func (i, j int) bool {
		return a.almanacCharacters[a.rawCombinations[i].result].level < a.almanacCharacters[a.rawCombinations[j].result].level
	})
	for _, c := range(a.rawCombinations) {
		s1 := a.almanacCharacters[c.char1].expandedSkills
		s2 := a.almanacCharacters[c.char2].expandedSkills
		s3 := a.almanacCharacters[c.result].expandedSkills
		a.almanacCharacters[c.result].expandedSkills = s3.Combine(s1, s2)
		a.almanacCharacters[c.result].combinations = append(a.almanacCharacters[c.result].combinations, c)
	}

	return a
}


func (a *Almanac) parseCombination (reader * bufio.Reader) {
	line, err := reader.ReadString('\n')
	handleError(err)
	fields := strings.Fields(line)
	c0, ok0 := a.almanacCharactersMap[fields[0]]
	c1, ok1 := a.almanacCharactersMap[fields[1]]
	c2, ok2 := a.almanacCharactersMap[fields[2]]
	if !(ok0 && ok1 && ok2) {
		log.Fatal("Unkonwn characters", c0, c1, c2)
	}

	a.rawCombinations = append(a.rawCombinations, RawCombination{
		result: c0,
		char1: c1,
		char2: c2,
	})
}


func (a *Almanac) parseCharacter (reader *bufio.Reader) {
	line, err := reader.ReadString('\n')
	handleError(err)
	fields := strings.Fields(line)
	name := fields[0]
	id := a.getNextCharacterId(name)
	level, err := strconv.Atoi(fields[1])
	handleError(err)
	gold, err := strconv.Atoi(fields[2])
	handleError(err)
	nSkills, err := strconv.Atoi(fields[3])
	handleError(err)
	skillsMask := SkillMask{}

	for i := 0; i< nSkills; i++ {
		skill := fields[i + 4]
		var id SkillId
		if skillId, ok := a.skillsMap[skill]; ok {
			id = skillId
		} else {
			id = a.registerSkill(skill)
		}
		skillsMask.addSkill(id)
	}
	c := AlmanacCharacter{
		id: id,
		skills: skillsMask,
		expandedSkills: skillsMask,
		name: name,
		gold: gold,
		combinations: []RawCombination{},
		level: level,
	}
	a.almanacCharacters[id] = c
}


func (a * Almanac) getNextCharacterId (name string) CharacterId {
	id := a.currentCharacterId
	a.currentCharacterId++
	a.almanacCharactersMap[name] = id
	return id
}
func (a * Almanac) registerSkill( skill string) SkillId {
	id := a.currentSkillId
	a.skillsMap[skill] = id
	a.skills[id] = Skill{
		id: id,
		name: skill,
	}
	a.currentSkillId++
	return id
}

type SkillId uint8
type CharacterId uint8

type Skill struct {
	id SkillId
	name string
}

type RawCombination struct {
	result, char1, char2 CharacterId
}

type CharacterType struct {
	gold int
	skills SkillMask
}
type ExpandedCharacter struct {
	id CharacterId
	characterTypes []CharacterType
}
type AlmanacCharacter struct {
	id CharacterId
	skills SkillMask
	name string
	expandedSkills SkillMask
	combinations []RawCombination
	gold int
	level int
}

type Almanac struct {
	currentSkillId SkillId
	currentCharacterId CharacterId
	skillsMap map[string]SkillId
	skills [256]Skill
	almanacCharactersMap map[string]CharacterId
	almanacCharacters [256]AlmanacCharacter
	rawCombinations []RawCombination
}

// There are 100 skills
type SkillMask [2]uint64

func (sm * SkillMask) addSkill (id SkillId) {
	if id < 64 {
		sm[0] = sm[0] | (1 << id)
	} else {
		sm[1] = sm[1] | (1 << (id - 64))
	}
}

func (sm SkillMask) Combine (sm1, sm2 SkillMask) (combination SkillMask) {
	return sm.or(sm1).or(sm2)
}

func (sm SkillMask) or (sm1 SkillMask) (sm2 SkillMask) {
	sm2[0] = sm1[0] | sm[0]
	sm2[1] = sm1[1] | sm[1]
	return
}
func (sm SkillMask) and (sm1 SkillMask) (sm2 SkillMask) {
	sm2[0] = sm1[0] & sm[0]
	sm2[1] = sm1[1] & sm[1]
	return
}
func (sm SkillMask) xor (sm1 SkillMask) (sm2 SkillMask) {
	sm2[0] = sm1[0] ^ sm[0]
	sm2[1] = sm1[1] ^ sm[1]
	return
}
func (sm SkillMask) neg () (sm2 SkillMask) {
	sm2[0] = ^ sm[0]
	sm2[1] = ^ sm[1]
	return
}

func (sm SkillMask) isZero() bool {
	return sm[0] == 0 && sm[1] == 0
}

func (sm SkillMask) popCount () int {
	return bits.OnesCount64(sm[0]) + bits.OnesCount64(sm[1])
}

func (sm SkillMask) split () [][2]SkillMask {
	p1 := partition64Mask(sm[0])
	p2 := partition64Mask(sm[1])
	partitions := [][2]SkillMask{}
	for _, subp1 := range (p1) {
		for _, subp2 := range(p2) {
			partitions = append(
				partitions, [2]SkillMask{{subp1[1], subp2[1]}, {subp1[0], subp2[0]}})
		}
	}
	return partitions
}

func partition64Mask (mask uint64) [][2]uint64 {
	combinations := [][2]uint64{{0, 0}}
	for true {
		v := mask & -mask;
		mask = mask & ^v
		if v == 0 {
			break
		}
		newCombinations := [][2]uint64{}
		for _, pair := range(combinations) {
			newPair1 := [2]uint64{pair[0] | v, pair[1]}
			newPair2 := [2]uint64{pair[0], pair[1] | v}
			newCombinations = append(newCombinations, newPair1, newPair2)
		}
		combinations = newCombinations
	}
	return combinations
}

func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}

