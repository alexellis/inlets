package client

import "testing"

func Test_AllowsAll_Allows(t *testing.T) {
	f := makeAllowsAllFilter()

	got := f("tcp", "example.com")
	want := true

	if got != want {
		t.Fatalf("connection allowed: want %v, got %v", want, got)
	}
}

func Test_StrictFilteringPassesOnMatch(t *testing.T) {
	u := map[string]string{
		"tcp": "http://example.com",
	}

	f := makeFilter(u)

	got := f("tcp", "example.com")
	want := true

	if got != want {
		t.Fatalf("connection allowed: want %v, got %v", want, got)
	}
}

func Test_StrictFilteringFails(t *testing.T) {
	u := map[string]string{
		"tcp": "example.com",
	}

	f := makeFilter(u)

	got := f("tcp", "test.com")
	want := false

	if got != want {
		t.Fatalf("connection allowed: want %v, got %v", want, got)
	}
}

func Test_StrictFilteringFailsUDP(t *testing.T) {
	u := map[string]string{}

	f := makeFilter(u)

	got := f("udp", "example.com")
	want := false

	if got != want {
		t.Fatalf("connection allowed: want %v, got %v", want, got)
	}
}
