package main

func wrapLine(s string, l int) []string {
	rs := []rune(s)
	if len(rs) < l {
		return []string{s}
	}

	var lines []string
	i := 0
	for ; i+l < len(rs); i += l {
		ss := string(rs[i:(i + l)])
		lines = append(lines, ss)
	}
	lines = append(lines, string(rs[i:]))

	return lines
}
