package domain_test

import (
	"strconv"
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/domain"
	"github.com/alcortesm/demo-mongodb-transactions/testhelp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillLevel_Ctor(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		// WHEN we create a valid skill level
		_, err := domain.NewSkillLevel(domain.MinSkillLevel)

		// THEN we get no error
		require.NoError(t, err)
	})

	t.Run("invalid value", func(t *testing.T) {
		t.Parallel()

		invalidValues := []int{-42, -1, 11, 42}

		for _, value := range invalidValues {
			t.Run(strconv.Itoa(value), func(t *testing.T) {
				t.Parallel()

				// WHEN we try to create a skill level with an invalid value
				_, err := domain.NewSkillLevel(value)

				// THEN we get an error
				require.Error(t, err)
			})
		}
	})
}

func TestSkillLevel_String(t *testing.T) {
	t.Parallel()

	subtests := []struct {
		skillLevel domain.SkillLevel
		want       string
	}{
		{
			skillLevel: domain.SkillLevel{},
			want:       "unknown_skill_level",
		},
		{
			skillLevel: testhelp.SkillLevel(t, domain.MinSkillLevel),
			want:       strconv.Itoa(domain.MinSkillLevel), // "0"
		},
		{
			skillLevel: testhelp.SkillLevel(t, domain.MaxSkillLevel),
			want:       strconv.Itoa(domain.MaxSkillLevel), // "10"
		},
	}

	for _, test := range subtests {
		t.Run(test.skillLevel.String(), func(t *testing.T) {
			t.Parallel()

			// WHEN we ask for the value
			got := test.skillLevel.String()

			// THEN we get the expected string
			require.Equal(t, test.want, got)
		})
	}
}

func TestSkillLevel_Value(t *testing.T) {
	t.Parallel()

	subtests := []struct {
		skillLevel domain.SkillLevel
		wantValue  int
		wantOK     bool
	}{
		{
			skillLevel: domain.SkillLevel{},
			wantValue:  0,
			wantOK:     false,
		},
		{
			skillLevel: testhelp.SkillLevel(t, domain.MinSkillLevel),
			wantValue:  0,
			wantOK:     true,
		},
		{
			skillLevel: testhelp.SkillLevel(t, domain.MaxSkillLevel),
			wantValue:  10,
			wantOK:     true,
		},
	}

	for _, test := range subtests {
		t.Run(test.skillLevel.String(), func(t *testing.T) {
			t.Parallel()

			// WHEN we ask for the value
			value, ok := test.skillLevel.Value()

			// THEN we get the expected ok
			assert.Equal(t, test.wantOK, ok)

			// THEN we get the expected value
			assert.Equal(t, test.wantValue, value)
		})
	}
}

func TestSkillLevel_IsKnown(t *testing.T) {
	t.Parallel()

	subtests := []struct {
		skillLevel domain.SkillLevel
		want       bool
	}{
		{
			skillLevel: domain.SkillLevel{},
			want:       false,
		},
		{
			skillLevel: testhelp.SkillLevel(t, domain.MinSkillLevel),
			want:       true,
		},
		{
			skillLevel: testhelp.SkillLevel(t, domain.MaxSkillLevel),
			want:       true,
		},
	}

	for _, test := range subtests {
		t.Run(test.skillLevel.String(), func(t *testing.T) {
			t.Parallel()

			// WHEN we ask if the value is known
			isKnown := test.skillLevel.IsKnown()

			// THEN we get the expected response
			require.Equal(t, test.want, isKnown)
		})
	}
}

// Tests that identical skill levels are equal to each other.
func TestSkillLevel_Comparable(t *testing.T) {
	// two skill levels representing the minimum skill level
	// are equal.
	a, err := domain.NewSkillLevel(domain.MinSkillLevel)
	require.NoError(t, err)
	b, err := domain.NewSkillLevel(domain.MinSkillLevel)
	require.NoError(t, err)
	require.True(t, a == b)

	// two skill levels representing the maximum skill level
	// are equal.
	a, err = domain.NewSkillLevel(domain.MaxSkillLevel)
	require.NoError(t, err)
	b, err = domain.NewSkillLevel(domain.MaxSkillLevel)
	require.NoError(t, err)
	require.True(t, a == b)

	// two skill levels representing an unknown skill level
	// are equal.
	a = domain.SkillLevel{}
	b = domain.SkillLevel{}
	require.True(t, a == b)
}
