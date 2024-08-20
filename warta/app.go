package warta

import (
	"encoding/csv"
	"fmt"
	"gobot/pkg"
	"net/http"
	"strings"
	"time"
)

func Scouting() {
	pkg.EverySaturdaySundayThreeHour(func() {
	fmt.Println("Scouting warta")
	_, lastRow, err := getData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	filteredData, err := filterData(lastRow)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// fmt.Println("\nFiltered Data:")
	// for key, value := range filteredData {
	//     fmt.Printf("%s: %s\n", key, value)
	// }

	// Convert filteredData to SupabaseWarta format
	warta := SupabaseWarta{
		CreatedAt:    int(time.Now().Unix()),
		BulletinDate: filteredData["BulletinDate"],
		Preacher1:    filteredData["PreacherID1"],
		Preacher2:    filteredData["PreacherID2"],
	}

	// Print warta before insertion
	// fmt.Println("\nWarta Data to be Inserted:")
	// fmt.Printf("%+v\n", warta)

	// Insert data into Supabase and capture the result
	if err = UpdateOrInsert(warta); err != nil {
		fmt.Println("Error during insertion:", err)
		return
	}
	})
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
		return nil, nil, fmt.Errorf("bulletin Date column not found")
	}

	// Prepare the last row with non-empty Bulletin Date
	var lastRow map[string]string
	for i := len(records) - 1; i > 0; i-- {
		row := records[i]
		if bulletinDateIndex < len(row) && strings.TrimSpace(row[bulletinDateIndex]) != "" {
			lastRow = make(map[string]string)
			for j, value := range row {
				lastRow[header[j]] = strings.TrimSpace(value)
			}
			return header, lastRow, nil
		}
	}

	return header, nil, fmt.Errorf("no non-empty Bulletin Date found")
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
