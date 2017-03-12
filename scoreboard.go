package godge

import (
	"sort"
	"sync"
)

const (
	failedVerdict = "Failed"
	passedVerdict = "Passed"
)

type scoreboard struct {
	sync.RWMutex
	m map[string]map[string]string
}

func (r scoreboard) set(user, task string, verdict string) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.m[user]; !ok {
		r.m[user] = make(map[string]string)
	}
	r.m[user][task] = verdict
}

func (r scoreboard) get(user, task string) string {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.m[user]; !ok {
		return ""
	}
	return r.m[user][task]
}

// returns a 2D array of the results (including the tasks as the first row and
// the users as the first column). The results are rows are sorted by the score
// of each user.
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
