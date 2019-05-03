package main

import (
	"log"
	"bufio"
	"strings"
	"strconv"
	"os"
	"sort"
)



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
		s1 := a.almanacCharacters[c.char1].skills
		s2 := a.almanacCharacters[c.char2].skills
		s3 := a.almanacCharacters[c.result].skills
		a.almanacCharacters[c.result].expandedSkills = s3.Combine(s1, s2)
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
		gold: gold,
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
	expandedSkills SkillMask
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
	combination[0] = sm[0] | sm1[0] | sm2[0]
	combination[1] = sm[1] | sm1[1] | sm2[1]
	return
}


func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}






