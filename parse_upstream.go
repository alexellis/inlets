package main

import "strings"

type UpstreamParser interface {
	Parse(input string) map[string]string
}

type ArgsUpstreamParser struct {
}

func (a *ArgsUpstreamParser) Parse(input string) map[string]string {
	upstreamMap := buildUpstreamMap(input)

	return upstreamMap
}

func buildUpstreamMap(args string) map[string]string {
	items := make(map[string]string)

	entries := strings.Split(args, ",")
	for _, entry := range entries {
		kvp := strings.Split(entry, "=")
		if len(kvp) == 1 {
			items[""] = strings.TrimSpace(kvp[0])
		} else {
			items[strings.TrimSpace(kvp[0])] = strings.TrimSpace(kvp[1])
		}
	}
	return items
}
