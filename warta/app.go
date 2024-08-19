package warta

import (
    "encoding/csv"
    "fmt"
    "net/http"
    "strings"
	"time"

)

func Scouting() {
    // pkg.EveryHourOnSaturday(func() {
    // 	fmt.Println("Scouting warta")
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

    fmt.Println("\nFiltered Data:")
    for key, value := range filteredData {
        fmt.Printf("%s: %s\n", key, value)
    }

    // Convert filteredData to SupabaseWarta format
    warta := SupabaseWarta{
        CreatedAt:    int(time.Now().Unix()),
        BulletinDate: filteredData["BulletinDate"],
        Preacher1:    filteredData["PreacherID1"],
        Preacher2:    filteredData["PreacherID2"],
    }

    // Print warta before insertion
    fmt.Println("\nWarta Data to be Inserted:")
    fmt.Printf("%+v\n", warta)

    // Insert data into Supabase and capture the result
    if err = UpdateOrInsert(warta); err != nil {
        fmt.Println("Error during insertion:", err)
        return
    }
    // })
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
        return nil, nil, fmt.Errorf("No records found")
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
        return nil, nil, fmt.Errorf("Bulletin Date column not found")
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

    return header, nil, fmt.Errorf("No non-empty Bulletin Date found")
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
            return nil, fmt.Errorf("Field %s not found in the data", key)
        }
    }

    return filteredData, nil
}

func UpdateOrInsert(warta SupabaseWarta) error {
	result := GetAllSupabaseWarta()

	// Check if any records were returned
	if len(result) > 0 {
		// Record exists, update logic here
		fmt.Println("Record exists, updating...")
		_, err := Update(warta)
		if err != nil {
			return fmt.Errorf("Error updating record: %w", err)
		}
	} else {
		// Record does not exist, insert new record
		_, err := Insert(warta)
		if err != nil {
			return fmt.Errorf("Error inserting new record: %w", err)
		}
	}

	return nil
}
