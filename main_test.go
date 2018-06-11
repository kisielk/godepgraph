package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	args := os.Args
	stdout := output
	testCases := []struct {
		args []string
		want string
	}{
		{
			[]string{"godepgraph", "github.com/kisielk/godepgraph"},
			`digraph godep {
_0 [label="flag" style="filled" color="palegreen"];
_1 [label="fmt" style="filled" color="palegreen"];
_2 [label="github.com/kisielk/godepgraph" style="filled" color="paleturquoise" URL="https://github.com/kisielk/godepgraph" target="_top"];
_2 -> _0;
_2 -> _1;
_2 -> _3;
_2 -> _4;
_2 -> _5;
_2 -> _6;
_2 -> _7;
_2 -> _8;
_3 [label="go/build" style="filled" color="palegreen"];
_4 [label="io" style="filled" color="palegreen"];
_5 [label="log" style="filled" color="palegreen"];
_6 [label="os" style="filled" color="palegreen"];
_7 [label="sort" style="filled" color="palegreen"];
_8 [label="strings" style="filled" color="palegreen"];
}
`},
		{
			[]string{"godepgraph", "-l", "1", "github.com/kisielk/godepgraph"},
			`digraph godep {
_0 [label="github.com/kisielk/godepgraph" style="filled" color="paleturquoise" URL="https://github.com/kisielk/godepgraph" target="_top"];
}
`},
		{
			[]string{"godepgraph", "-i", "io,fmt,os", "github.com/kisielk/godepgraph"},
			`digraph godep {
_0 [label="flag" style="filled" color="palegreen"];
_1 [label="github.com/kisielk/godepgraph" style="filled" color="paleturquoise" URL="https://github.com/kisielk/godepgraph" target="_top"];
_1 -> _0;
_1 -> _2;
_1 -> _3;
_1 -> _4;
_1 -> _5;
_2 [label="go/build" style="filled" color="palegreen"];
_3 [label="log" style="filled" color="palegreen"];
_4 [label="sort" style="filled" color="palegreen"];
_5 [label="strings" style="filled" color="palegreen"];
}
`},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			defer func() {
				os.Args = args
				output = stdout
			}()
			buf := bytes.NewBuffer([]byte{})
			output = buf
			os.Args = tc.args
			main()
			actual := string(buf.Bytes())
			if tc.want != actual {
				t.Log(actual)
				t.Log(tc.want)
				t.Fatal()
			}
		})
	}
}

func TestGetUrl(t *testing.T) {
	testCases := []struct {
		name string
		want string
	}{
		{"github.com/kisielk/godepgraph", "https://github.com/kisielk/godepgraph"},
		{"github.ibm.com/logging/beates", "https://github.ibm.com/logging/beates"},
		{"github.com/collectd", ""},
		{"collectd.org/api", "https://godoc.org/collectd.org/api"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			actual := getUrl(tc.name)
			if actual != tc.want {
				t.Log(actual)
				t.Log(tc.want)
				t.Fatal()
			}
		})
	}
}
