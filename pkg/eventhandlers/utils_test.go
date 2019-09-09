package eventhandlers

import "testing"

func TestHasRole(t *testing.T) {
	tt := map[string]struct {
		roles []string
		role  string
		exp   bool
	}{
		"empty list should return false": {
			[]string{}, "W", false,
		},
		"list has values but not the one we want": {
			[]string{"W", "R", "M"}, "O", false,
		},
		"ok": {
			[]string{"W", "R"}, "W", true,
		},
	}
	for sc, tc := range tt {
		res := hasRole(tc.roles, tc.role)
		if res != tc.exp {
			t.Errorf("expected %v got %v in `%s`", tc.exp, res, sc)
		}
	}
}
