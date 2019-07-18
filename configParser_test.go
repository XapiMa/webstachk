package webstachk

import "testing"

func TestParseStatuses(t *testing.T) {
	type data struct {
		input  string
		expect []int
	}
	testDatas := []data{
		{
			"",
			[]int{},
		},
		{
			"|",
			[]int{},
		},
		{
			"||",
			[]int{},
		},
		{
			"200",
			[]int{200},
		},
		{
			"200|",
			[]int{200},
		},
		{
			"200|300",
			[]int{200, 300},
		},
		{
			"200|300|",
			[]int{200, 300},
		},
		{
			"200|300|301",
			[]int{200, 300, 301},
		},
	}

	for i, test := range testDatas {
		out, err := parseStatuses(test.input)
		if err != nil {
			t.Errorf("[ParseStatuses]in %dth test: %s", i, err)
		}
		if len(out) != len(test.expect) {
			t.Errorf("[ParseStatuses]in %dth test: %v found but expect %v", i, out, test.expect)
		}
		for j := range out {
			if out[j] != test.expect[j] {
				t.Errorf("[ParseStatuses]in %dth test %dth item : %d found but expect %d", i, j, out[j], test.expect[j])
			}
		}
	}

}

func TestParseTime(t *testing.T) {
	type data struct {
		input  string
		expect int64
	}
	testDatas := []data{
		{"0s", 0},
		{"10000s", 10000},
		{"10m", 10 * 60},
		{"10h", 10 * 60 * 60},
		{"10d", 10 * 60 * 60 * 24},
	}

	for i, test := range testDatas {
		out, err := parseTime(test.input)
		if err != nil {
			t.Errorf("[ParseStatuses]in %dth test: %s", i, err)
		}
		if out != test.expect {
			t.Errorf("[ParseStatuses]in %dth test : %d found but expect %d", i, out, test.expect)
		}
	}
}
