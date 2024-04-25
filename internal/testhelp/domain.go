package testhelp

import (
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
)

func SkillLevel(t *testing.T, n int) domain.SkillLevel {
	t.Helper()

	sl, err := domain.NewSkillLevel(n)
	if err != nil {
		t.Fatal(err)
	}

	return sl
}
