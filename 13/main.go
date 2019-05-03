package main

import (
	"log"
)

/*
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
		skillsMap: map[string]SkillId{},
		almanacCharactersMap: map[string]CharacterId{},
	}

	for i := 0; i < nCharacters; i++ {
		a.parseCharacter(reader)
	}


}
*/

/*
func (a *Almanac) parseCharacter (reader *bufio.Reader) {
	line, err := reader.ReadString('\n')
	handleError(err)
	id := a.getNextCharacterId()
	fields := strings.Fields(line)
	name := fields[0]
	a.almanacCharactersMap[name] = id
	level, err := strconv.Atoi(fields[1])
	handleError(err)
	gold, err := strconv.Atoi(fields[2])
	handleError(err)
	nSkills, err := strconv.Atoi(fields[3])
	handleError(err)
	c := AlmanacCharacter{
		id: id,
	}

}
*/

func (a * Almanac) getNextCharacterId () CharacterId {
	id := a.currentCharacterId
	a.currentCharacterId++
	return id
}

type SkillId uint8
type CharacterId uint8

type Skill struct {
	id SkillId
	name string
}

type AlmanacCharacter struct {
	id CharacterId
	skills SkillMask
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


func handleError (e error) {
	if e != nil {
		log.Fatal(e)
	}
}






