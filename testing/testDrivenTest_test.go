package testing

// Source: https://speakerdeck.com/mitchellh/advanced-testing-with-go

import "testing"

func TestAddNumber(t *testing.T) {
	cases := map[string]struct{ A, B, Expected int }{
		"test case 01": {1, 1, 2},
		"test case 02": {1, -1, 0},
		"test case 03": {1, 0, 1},
	}

	for k, tc := range cases {
		actual := tc.A + tc.B
		if actual != tc.Expected {
			t.Errorf("%s: %d + %d = %d, expected %d",
				k, tc.A, tc.B, actual, tc.Expected)
		}
	}
}
