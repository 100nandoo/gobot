package warta

import (
	"encoding/csv"
	"fmt"
	"gobot/pkg"
	"net/http"
	"strings"
	"time"
)
func Scouting(now bool) {
	scoutingLogic := func() {
		pkg.LogWithTimestamp("Scouting warta started")

		_, lastRow, err := getData()
		if err != nil {
			pkg.LogWithTimestamp("Error: %v", err)
			return
		}

		filteredData, err := filterData(lastRow)
		if err != nil {
			pkg.LogWithTimestamp("Error: %v", err)
			return
		}

		// Convert filteredData to SupabaseWarta format
		warta := SupabaseWarta{
			CreatedAt:    int(time.Now().Unix()),
			BulletinDate: filteredData["BulletinDate"],
			Preacher1:    filteredData["PreacherID1"],
			Preacher2:    filteredData["PreacherID2"],
		}

		// Insert data into Supabase and capture the result
		if err = UpdateOrInsert(warta); err != nil {
			pkg.LogWithTimestamp("Error during insertion: %v", err)
			return
		}
		pkg.LogWithTimestamp("Scouting warta finished successfully")
	}

	if now {
		pkg.LogWithTimestamp("Running scouting logic immediately")
		scoutingLogic()
	} else {
		pkg.LogWithTimestamp("Scheduling scouting logic for Saturday and Sunday")
		pkg.EverySaturdayDayAtThisHour(scoutingLogic, "12:12")
	}
}

var indonesianMonths = map[string]string{
	"Januari":   "01",
	"JANUARI":   "01",
	"Februari":  "02",
	"FEBRUARI":  "02",
	"Maret":     "03",
	"MARET":     "03",
	"April":     "04",
	"APRIL":     "04",
	"Mei":       "05",
	"MEI":       "05",
	"Juni":      "06",
	"JUNI":      "06",
	"Juli":      "07",
	"JULI":      "07",
	"Agustus":   "08",
	"AGUSTUS":   "08",
	"September": "09",
	"SEPTEMBER": "09",
	"Oktober":   "10",
	"OKTOBER":   "10",
	"November":  "11",
	"NOVEMBER":  "11",
	"Desember":  "12",
	"DESEMBER":  "12",
}

// Function to parse Indonesian date format (e.g., "1 Januari 2024")
func parseIndonesianDate(dateStr string) (time.Time, error) {
	for month, num := range indonesianMonths {
		dateStr = strings.Replace(dateStr, month, num, 1)
	}
	return time.Parse("2 01 2006", dateStr)
}

// Function to find the next Sunday from the current date
func nextSunday() time.Time {
	now := time.Now().Truncate(24 * time.Hour) // Reset the time part to midnight

	// Since the function runs on Saturday, we know the next Sunday is tomorrow
	nextSunday := now.AddDate(0, 0, 1)
	pkg.LogWithTimestamp("Next Sunday: %v", nextSunday.Format("2006-01-02"))
	
	return nextSunday
}

func getData() ([]string, map[string]string, error) {
	// The URL of the published Google Sheets as a CSV
	url := "https://docs.google.com/spreadsheets/d/e/2PACX-1vS9mVW8Ld_E_wkt7IF65StvyGiZ_LpSxL_Uaoop4qqIsatlIuIqj38V8wnZCy3k5Clo22lQPJhs4qA5/pub?output=csv"

	// Fetch the CSV data from the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Read the CSV data
	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	// Check if there are any records
	if len(records) == 0 {
		return nil, nil, fmt.Errorf("no records found")
	}

	// Find the index of the "Bulletin Date" column (assuming it's in the header)
	header := records[0]
	bulletinDateIndex := -1
	for i, colName := range header {
		if strings.TrimSpace(colName) == "BulletinDate" {
			bulletinDateIndex = i
			break
		}
	}

	if bulletinDateIndex == -1 {
		return nil, nil, fmt.Errorf("BulletinDate column not found")
	}

	// Find the row with the upcoming Sunday's date
	upcomingSunday := nextSunday()

	for _, row := range records[1:] {
		if bulletinDateIndex < len(row) && strings.TrimSpace(row[bulletinDateIndex]) != "" {
			dateStr := strings.TrimSpace(row[bulletinDateIndex])
			parsedDate, err := parseIndonesianDate(dateStr)
			// pkg.LogWithTimestamp("Parsed date: %v", parsedDate.Format("2006-01-02"))
			if err != nil {
				pkg.LogWithTimestamp("Error: %v", err)
				continue // Skip if parsing fails
			}

			// Check if the parsed date matches the upcoming Sunday
			if parsedDate.Year() == upcomingSunday.Year() &&
				parsedDate.YearDay() == upcomingSunday.YearDay() {
				// Return the row with the matching Bulletin Date
				lastRow := make(map[string]string)

				for j, value := range row {
					lastRow[header[j]] = strings.TrimSpace(value)
				}
				preacherID1 := lastRow["PreacherID1"]
				preacherID2 := lastRow["PreacherID2"]
				pkg.LogWithTimestamp("PreacherID1: %s, PreacherID2: %s", preacherID1, preacherID2)

				return header, lastRow, nil
			}
		}
	}

	return header, nil, fmt.Errorf("no matching Bulletin Date for upcoming Sunday found")
}

func formatMap(m map[string]string) string {
	var result string
	for key, value := range m {
		result += fmt.Sprintf("%s: %s\n", key, value)
	}
	return result
}

func filterData(row map[string]string) (map[string]string, error) {
	// Define the keys you're interested in
	desiredKeys := []string{"BulletinDate", "PreacherID1", "PreacherID2"}

	// Create a map to store the filtered data
	filteredData := make(map[string]string)

	// Extract the desired fields from the row
	for _, key := range desiredKeys {
		if value, exists := row[key]; exists {
			filteredData[key] = value
		} else {
			return nil, fmt.Errorf("field %s not found in the data", key)
		}
	}

	return filteredData, nil
}

func UpdateOrInsert(warta SupabaseWarta) error {
	result := GetAllSupabaseWarta()
	existingWarta := findWarta(result, warta)

	if existingWarta != nil {
		if !isIdentical(*existingWarta, warta) {
			// fmt.Println("Update")
			Update(warta)
			SendWarta(warta)
		}
	} else {
		// fmt.Println("Insert")
		Insert(warta)
		SendWarta(warta)
	}

	return nil
}

func findWarta(wartas []SupabaseWarta, target SupabaseWarta) *SupabaseWarta {
	for _, warta := range wartas {
		if warta.BulletinDate == target.BulletinDate {
			return &warta
		}
	}
	return nil
}

func isIdentical(warta SupabaseWarta, target SupabaseWarta) bool {
	return warta.Preacher1 == target.Preacher1 && warta.Preacher2 == target.Preacher2
}
