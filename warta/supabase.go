package warta

import (
	"fmt"
	"gobot/pkg"
)

type SupabaseWarta struct {
	CreatedAt    int    `json:"created_at"`
	BulletinDate string `json:"bulletin_date"`
	Preacher1    string `json:"preacher1"`
	Preacher2    string `json:"preacher2"`
}

const dbName = "Warta"

/*
GetAllSupabaseWarta

Get All rows from Warta Database, return arrays of SupabaseWarta objects
*/
func GetAllSupabaseWarta() []SupabaseWarta {
	var result []SupabaseWarta
	err := pkg.SupabaseClient.DB.From(dbName).Select("*").Execute(&result)
	if err != nil {
		fmt.Println("Error calling GetAllSupabaseWarta", err)
		return nil
	}
	return result
}

// Insert inserts a row into the Warta Database and returns the inserted record with ID
func Insert(warta SupabaseWarta) (*SupabaseWarta, error) {
	var result SupabaseWarta
	err := pkg.SupabaseClient.DB.From(dbName).Insert(warta).Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("Error calling Insert: %w", err)
	}
	return &result, nil
}

// Update a row in the Warta Database
func Update(updatedWarta SupabaseWarta) (*SupabaseWarta, error) {
	var result SupabaseWarta
	err := pkg.SupabaseClient.DB.From(dbName).Update(updatedWarta).Eq("bulletin_date", updatedWarta.BulletinDate).Execute(&result)
	if err != nil {
		return nil, fmt.Errorf("Error calling Update: %w", err)
	}
	return &result, nil
}

// Delete deletes a row in the Warta Database by ID
func Delete(id string) error {
	var result []SupabaseWarta
	err := pkg.SupabaseClient.DB.From(dbName).Delete().Eq("id", id).Execute(&result)
	if err != nil {
		return fmt.Errorf("Error calling Delete: %w", err)
	}
	return nil
}