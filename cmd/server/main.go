package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var HEADERS = map[string][]string{
	"MATCH":             []string{"Date", "Time", "Season", "Competition", "Phase", "Type", "DayID", "MatchID", "Regulation", "?", "?", "?"},
	"TEAMS":             []string{"ID", "Name", "Result", "HeadCoach", "Assistant", "Color"},
	"MORE":              []string{"Referee", "Spectators", "Receipts", "City", "Hall", "Scout"},
	"SET":               []string{"IsTieBreak", "PartialScore1", "PartialScore2", "PartialScore3", "Score", "Duration"},
	"PLAYERS-H":         []string{"?", "Number", "Index", "RoleSet1", "RoleSet2", "RoleSet3", "RoleSet4", "RoleSet5", "ID", "LastName", "FirstName", "Nickname", "IsCaptain", "Role", "IsForeign"},
	"PLAYERS-V":         []string{"?", "Number", "Index", "RoleSet1", "RoleSet2", "RoleSet3", "RoleSet4", "RoleSet5", "ID", "LastName", "FirstName", "Nickname", "IsCaptain", "Role", "IsForeign"},
	"ATTACKCOMBINATION": []string{"Code", "Zone", "Ball", "Attack", "Description", "?", "?", "?", "?"},
	"SETTERCALL":        []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?"},
	"RESERVE":           []string{},
	"SCOUT":             []string{},
}

type Match struct {
	Date        string
	Time        string
	Season      string
	Competition string
	Phase       string
	Type        string
	DayID       string
	MatchID     string
}

type Team struct {
	ID        string
	Name      string
	Result    string
	HeadCoach string
	Assistant string
	Color     string
	DayID     string
	MatchID   string
}

type More struct {
	Referee    string
	Spectators string
	Receipts   string
	City       string
	Hall       string
	Scout      string
}

type Set struct {
	IsTieBreak    string
	PartialScore1 string
	PartialScore2 string
	PartialScore3 string
	Score         string
	Duration      string
}

type Player struct {
	Number    string
	Index     string
	RoleSet1  string
	RoleSet2  string
	RoleSet3  string
	RoleSet4  string
	RoleSet5  string
	ID        string
	LastName  string
	FirstName string
	Nickname  string
	IsCaptain string
	Role      string
	IsForeign string
}

type AttackCombination struct {
	Code        string
	Zone        string
	Ball        string
	Attack      string
	Description string
}

func main() {
	f, _ := os.Open("../../data/sample.dvw")
	defer f.Close()

	s := bufio.NewScanner(f)
	sectionRegex := regexp.MustCompile(`^\[3(.*)\]$`)
	section := ""
	rowIndex := 0

	for s.Scan() {
		line := s.Text()
		match := sectionRegex.FindStringSubmatch(line)

		if len(match) > 0 {
			section = match[1]
			rowIndex = 0
			fmt.Printf("--%s--\n", section)
			continue
		}

		// Non-CSV sections
		switch section {
		case "DATAVOLLEYSCOUT":
			continue
		case "WINNINGSYMBOLS":
			continue
		case "COMMENTS":
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

		header := HEADERS[section]
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
			match := Match{
				values["Date"],
				values["Time"],
				values["Season"],
				values["Competition"],
				values["Phase"],
				values["Type"],
				values["DayID"],
				values["MatchID"],
			}

			fmt.Println(match)
		case "TEAMS":
			team := Team{
				values["ID"],
				values["Name"],
				values["Result"],
				values["HeadCoach"],
				values["Assistant"],
				values["Color"],
				values["DayID"],
				values["MatchID"],
			}
			fmt.Println(team)
		case "MORE":
			more := More{
				values["Referee"],
				values["Spectators"],
				values["Receipts"],
				values["City"],
				values["Hall"],
				values["Scout"],
			}
			fmt.Println(more)
			// continue
		case "SET":
			set := Set{
				values["IsTieBreak"],
				values["PartialScore1"],
				values["PartialScore2"],
				values["PartialScore3"],
				values["Score"],
				values["Duration"],
			}

			fmt.Println(set)
			// continue
		case "PLAYERS-H", "PLAYERS-V":
			player := Player{
				values["Number"],
				values["Index"],
				values["RoleSet1"],
				values["RoleSet2"],
				values["RoleSet3"],
				values["RoleSet4"],
				values["RoleSet5"],
				values["ID"],
				values["LastName"],
				values["FirstName"],
				values["Nickname"],
				values["IsCaptain"],
				values["Role"],
				values["IsForeign"],
			}

			fmt.Println(player)
			// continue
		case "ATTACKCOMBINATION":
			combination := AttackCombination{
				values["Code"],
				values["Zone"],
				values["Ball"],
				values["Attack"],
				values["Description"],
			}

			fmt.Println(combination)
			// continue
		case "SETTERCALL":
			// continue
		case "RESERVE":
			// continue
		case "SCOUT":
			// continue
		}

		rowIndex++
	}
}
