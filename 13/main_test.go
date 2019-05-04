package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestSkillMask (t * testing.T) {
	testCases := []struct{
		sm SkillMask
		si SkillId
		expected SkillMask
	}{
		{
			sm: SkillMask{0, 0},
			si: 0,
			expected: SkillMask{1, 0},

		},
		{
			sm: SkillMask{1, 0},
			si: 1,
			expected: SkillMask{3, 0},

		},
		{
			sm: SkillMask{3, 0},
			si: 64,
			expected: SkillMask{3, 1},

		},
		{
			sm: SkillMask{3, 0},
			si: 65,
			expected: SkillMask{3, 2},

		},
	}
	for _, tc := range(testCases) {
		tc.sm.addSkill(tc.si)
		assert.Equal(t, tc.expected, tc.sm)
	}
}

func TestMissingSkills (t * testing.T) {
	testCases := []struct{
		sm, sm2 SkillMask
		expected SkillMask
	}{
		{
			sm: SkillMask{5, 0},
			sm2: SkillMask{2, 0},
			expected: SkillMask{5, 0},

		},
		{
			sm: SkillMask{7, 0},
			sm2: SkillMask{2, 0},
			expected: SkillMask{5, 0},

		},
		{
			sm: SkillMask{3, 0},
			sm2: SkillMask{2, 0},
			expected: SkillMask{1, 0},

		},
	}
	for _, tc := range(testCases) {
		expected := tc.sm.and(tc.sm2.neg())
		assert.Equal(t, tc.expected, expected)
	}
}

func TestPartitions (t * testing.T) {
	testCases := []struct{
		sm SkillMask
		expected [][2]SkillMask
	}{
		{
			sm: SkillMask{0, 0},
			expected: [][2]SkillMask{{{0, 0}, {0, 0}}},

		},
		{
			sm: SkillMask{1, 0},
			expected: [][2]SkillMask{{{0, 0}, {1, 0}}, {{1, 0}, {0, 0}}},

		},
		{
			sm: SkillMask{3, 0},
			expected: [][2]SkillMask{
				{{0, 0}, {3, 0}},
				{{2, 0}, {1, 0}},
				{{1, 0}, {2, 0}},
				{{3, 0}, {0, 0}},
			},
		},
		{
			sm: SkillMask{1, 1},
			expected: [][2]SkillMask{
				{{0, 0}, {1, 1}},
				{{0, 1}, {1, 0}},
				{{1, 0}, {0, 1}},
				{{1, 1}, {0, 0}},
			},
		},
	}
	for _, tc := range(testCases) {
		expected := tc.sm.split()
		assert.Equal(t, tc.expected, expected)
	}
}

func TestParseAlmanac (t *testing.T) {
	a := parseAlmanac()
	assert.Equal(t, 1, a.almanacCharacters[a.almanacCharactersMap["shiva"]].level)
	assert.Equal(t,
		a.almanacCharacters[a.almanacCharactersMap["character2"]].skills,
		SkillMask{(1 << a.skills[a.skillsMap["skill75"]].id) | ( 1 << a.skills[a.skillsMap["wind"]].id), 0})
	log.Println(a.almanacCharacters[a.almanacCharactersMap["ifrit"]])
}
