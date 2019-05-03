package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
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
