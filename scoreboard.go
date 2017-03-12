package godge

import "sort"

const (
	failedVerdict = "Failed"
	passedVerdict = "Passed"
)

type scoreboard map[string]map[string]string

func (r scoreboard) set(user, task string, verdict string) {
	if _, ok := r[user]; !ok {
		r[user] = make(map[string]string)
	}
	r[user][task] = verdict
}

func (r scoreboard) get(user, task string) string {
	if _, ok := r[user]; !ok {
		return ""
	}
	return r[user][task]
}

func (r scoreboard) toScoreboard(allUsers, allTasks []string) [][]string {

	var ret [][]string

	for _, u := range allUsers {
		row := []string{u}
		for _, t := range allTasks {
			row = append(row, r.get(u, t))
		}
		ret = append(ret, row)
	}

	score := func(row []string) int {
		var sc int
		for _, r := range row {
			if r == passedVerdict {
				sc++
			}
		}
		return sc
	}

	sort.Slice(ret, func(i, j int) bool {
		return score(ret[i][1:]) < score(ret[j][1:])
	})

	header := append([]string{}, append([]string{""}, allTasks...)...)
	ret = append([][]string{header}, ret...)

	return ret
}
