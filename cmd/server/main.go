package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Match struct {
	MatchMeta
	MatchAdditional
	Teams  [2]Team
	Sets   [5]Set
	Events []Event
}

type Team struct {
	TeamMeta
	Players map[int]Player
}

type MatchMeta struct {
	Date        time.Time
	Season      string
	Competition string
	Phase       string
	Type        string
	ID          int
}

type TeamMeta struct {
	ID        int
	Name      string
	Result    int
	HeadCoach string
	Assistant string
	Color     string
}

type MatchAdditional struct {
	Referee    string
	Spectators string
	City       string
	Hall       string
}

type Set struct {
	IsTieBreak    bool
	PartialScores [3]string
	Score         string
	Duration      int // Minutes
}

type Player struct {
	Number     int
	RoleInSets [5]string
	ID         int
	LastName   string
	FirstName  string
	Nickname   string
	IsCaptain  bool
	Role       string
	IsForeign  bool
}

type Event struct{}

func main() {
	f, _ := os.Open("../../data/sample.dvw")
	defer f.Close()

	s := bufio.NewScanner(f)
	sectionRegex := regexp.MustCompile(`^\[3(.*)\]$`)
	section := ""
	rowIndex := 0

	match := Match{}

	csvHeaders := map[string][]string{
		"MATCH":             []string{"Date", "Time", "Season", "Competition", "Phase", "Type", "DayID", "ID", "Regulation", "?", "?", "?"},
		"TEAMS":             []string{"ID", "Name", "Result", "HeadCoach", "Assistant", "Color"},
		"MORE":              []string{"Referee", "Spectators", "Receipts", "City", "Hall", "Scout"},
		"SET":               []string{"IsTieBreak", "PartialScore1", "PartialScore2", "PartialScore3", "Score", "Duration"},
		"PLAYERS-H":         []string{"TeamIndex", "Number", "Index", "RoleSet1", "RoleSet2", "RoleSet3", "RoleSet4", "RoleSet5", "ID", "LastName", "FirstName", "Nickname", "IsCaptain", "Role", "IsForeign"},
		"PLAYERS-V":         []string{"TeamIndex", "Number", "Index", "RoleSet1", "RoleSet2", "RoleSet3", "RoleSet4", "RoleSet5", "ID", "LastName", "FirstName", "Nickname", "IsCaptain", "Role", "IsForeign"},
		"ATTACKCOMBINATION": []string{"Code", "Zone", "Ball", "Attack", "Description", "?", "?", "?", "?"},
		"SETTERCALL":        []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
		"RESERVE":           []string{},
		"SCOUT":             []string{},
	}

	for s.Scan() {
		line := s.Text()
		regexMatches := sectionRegex.FindStringSubmatch(line)

		if len(regexMatches) > 0 {
			section = regexMatches[1]
			rowIndex = 0
			continue
		}

		// Skip these sections
		switch section {
		case
			"DATAVOLLEYSCOUT",
			"COMMENTS",
			"ATTACKCOMBINATION",
			"SETTERCALL",
			"WINNINGSYMBOLS",
			"RESERVE":
			continue
		}

		// Checks for irregular sections
		switch section {
		case "MATCH":
			if rowIndex != 0 {
				continue
			}
		case "MORE":
			if rowIndex != 0 {
				continue
			}
		}

		csvReader := csv.NewReader(strings.NewReader(line))
		csvReader.Comma = ';'

		header := csvHeaders[section]
		row, _ := csvReader.Read()

		values := make(map[string]string)
		for i, name := range header {
			if name == "?" {
				continue
			}

			values[name] = row[i]
		}

		switch section {
		case "MATCH":
			date, _ := time.Parse(
				"01/02/2006T15.04.05",
				fmt.Sprintf("%sT%s", values["Date"], values["Time"]),
			)
			id, _ := strconv.Atoi(values["ID"])
			meta := MatchMeta{
				date,
				values["Season"],
				values["Competition"],
				values["Phase"],
				values["Type"],
				id,
			}
			match.MatchMeta = meta

		case "TEAMS":
			id, _ := strconv.Atoi(values["ID"])
			result, _ := strconv.Atoi(values["Result"])
			meta := TeamMeta{
				id,
				values["Name"],
				result,
				values["HeadCoach"],
				values["Assistant"],
				values["Color"],
			}
			match.Teams[rowIndex].TeamMeta = meta

		case "MORE":
			additional := MatchAdditional{
				values["Referee"],
				values["Spectators"],
				values["City"],
				values["Hall"],
			}
			match.MatchAdditional = additional

		case "SET":
			// Ignore non-played sets
			if values["Score"] == "" {
				continue
			}

			duration, _ := strconv.Atoi(values["Duration"])
			set := Set{
				values["IsTieBreak"] == "True",
				[3]string{
					values["PartialScore1"],
					values["PartialScore2"],
					values["PartialScore3"],
				},
				values["Score"],
				duration,
			}

			match.Sets[rowIndex] = set

		case "PLAYERS-H", "PLAYERS-V":
			index, _ := strconv.Atoi(values["TeamIndex"])
			number, _ := strconv.Atoi(values["Number"])
			id, _ := strconv.Atoi(values["ID"])

			player := Player{
				number,
				[5]string{
					values["RoleSet1"],
					values["RoleSet2"],
					values["RoleSet3"],
					values["RoleSet4"],
					values["RoleSet5"],
				},
				id,
				values["LastName"],
				values["FirstName"],
				values["Nickname"],
				values["IsCaptain"] == "C",
				values["Role"],
				values["IsForeign"] == "True",
			}
			if match.Teams[index].Players == nil {
				match.Teams[index].Players = make(map[int]Player)
			}
			match.Teams[index].Players[player.ID] = player

		case "SCOUT":
			// TODO
		}

		rowIndex++
	}

	fmt.Printf("%+v\n", match)
}
