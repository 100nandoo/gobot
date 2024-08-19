package warta

import (
    "encoding/csv"
    "fmt"
    "net/http"
    "strings"
)

func Scouting(){
	// pkg.EveryHourOnSaturday(func() {
	// 	fmt.Println("Scouting warta")
		getData()
	// })
}

func getData() {
    // The URL of the published Google Sheets as a CSV
    url := "https://docs.google.com/spreadsheets/d/e/2PACX-1vS9mVW8Ld_E_wkt7IF65StvyGiZ_LpSxL_Uaoop4qqIsatlIuIqj38V8wnZCy3k5Clo22lQPJhs4qA5/pub?output=csv"

    // Fetch the CSV data from the URL
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Read the CSV data
    reader := csv.NewReader(resp.Body)
    records, err := reader.ReadAll()
    if err != nil {
        panic(err)
    }

    // Check if there are any records
    if len(records) == 0 {
        fmt.Println("No records found")
        return
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
        fmt.Println("Bulletin Date column not found")
        return
    }

    // Print the column names with their indices
    fmt.Println("Column Names and Indices:")
    for i, colName := range header {
        fmt.Printf("Index %d: %s\n", i, strings.TrimSpace(colName))
    }

    // Iterate over the rows from the end to find the last row where "Bulletin Date" is not empty
    for i := len(records) - 1; i > 0; i-- {
        row := records[i]
        if bulletinDateIndex < len(row) && strings.TrimSpace(row[bulletinDateIndex]) != "" {
            fmt.Println("\nLast row with non-empty Bulletin Date:")
            for j, value := range row {
                fmt.Printf("%s: %s\n", header[j], strings.TrimSpace(value))
            }
            return
        }
    }

    fmt.Println("No non-empty Bulletin Date found")
}
