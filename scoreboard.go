package godge

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	failedVerdict = "Failed"
	passedVerdict = "Passed"
)

func saveToScoreboard(db *sqlx.DB, user, task string, verdict string) error {
	_, err := db.Exec("INSERT INTO scoreboard (username, task_name, verdict, submitted_at) VALUES (?,?,?,?)", user, task, verdict, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save scoreboard record: %v", err)
	}
	return nil
}

func getFromScoreboard(db *sqlx.DB, user, task string) (*string, error) {
	var res string
	err := db.Get(&res, "SELECT verdict FROM scoreboard WHERE username=? AND task_name=? ORDER BY ID DESC LIMIT 1", user, task)
	if err == sql.ErrNoRows {
		e := ""
		return &e, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from scoreboard: %v", err)
	}
	return &res, nil
}

// returns a 2D array of the results (including the tasks as the first row and
// the users as the first column). The results are rows are sorted by the score
// of each user.
func buildScoreboard(db *sqlx.DB, allUsers, allTasks []string) ([][]string, error) {

	var ret [][]string

	for _, u := range allUsers {
		row := []string{u}
		for _, t := range allTasks {
			v, err := getFromScoreboard(db, u, t)
			if err != nil {
				return nil, fmt.Errorf("failed to build scoreboard: %v", err)
			}
			row = append(row, *v)
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

	return ret, nil
}
