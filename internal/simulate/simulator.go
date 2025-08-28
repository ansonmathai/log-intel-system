package simulate

import (
	"fmt"
	"math/rand"
	"time"
)

type LogEntry struct {
	Timestamp time.Timestamp
	Category string
	Amount int
	Note string
	Location string
}

var categories = []string{"Food", "Travel", "Groceries", "Entertainment", "Utilities"}
var notes = []string{"Lunch", "Bus fare", "Movie", "Electricity bill", "Snacks"}
var locations = []string{"Ernakulam", "Kochi", "Thrissur", "Aluva", "Perumbavoor"}

func GenerateLog() LogEntry {
	return LogEntry{
		Timestamp: time.Now(),
		Category: categories[rand.Intn(len(categories))],
		Amount: rand.Intn(500)+50,
		Note: notes[rand.Intn(len(notes))],
		Location: locations[rand.Intn(len(locations))], 
	}
}

func FormatLog(entry LogEntry) string {
	return fmt.sprintf("%s\t%s\t₹%d\t%s\t%s",
			entry.Timestamp.Format(time.RFC3339),
			entry.Category,
			entry.Amount,
			entry.Note,
			entry.Location,
		)
}