package ao3

import "testing"

type hostTest struct {
	parts string
	host  *Host
}

var hostTests = []hostTest{
	hostTest{
		parts: "www.audible.com",
		host: &Host{
			sub:    "www",
			tld:    "com",
			domain: "audible",
		},
	},
	hostTest{
		parts: "audible.com",
		host: &Host{
			tld:    "com",
			domain: "audible",
		},
	},
	hostTest{
		parts: "audible",
		host:  &Host{},
	},
}

func TestParseHost(t *testing.T) {
	for _, test := range hostTests {
		h, err := ParseHost(test.parts)
		if err != nil {
			t.Log(err)
		}
		if h.String() != test.host.String() {
			t.Errorf("input %v: output %#v != expected %v", test.parts, h, test.host)
		}
	}
}

type newTest struct {
	input *Query
	want  string
}

var newTests = []newTest{
	newTest{
		input: New().PrependPath("poot").SetHost("www.audible.com").SetParam("query", "toot").AppendPath("root"),
		want:  "https://www.audible.com/poot/root?query=toot",
	},
	newTest{
		input: New("https://www.audible.com/poot?query=toot"),
		want:  "https://www.audible.com/poot?query=toot",
	},
}

func TestNew(t *testing.T) {
	for _, test := range newTests {
		if output := test.input.String(); output != test.want {
			t.Errorf("input %v: output %#v != expected %v", test.input, output, test.want)
		}
	}
}

func TestPoot(t *testing.T) {
	println("poot")
}
