package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const supabaseUrlKey, supabaseKeyKey string = "Supabase.url", "Supabase.key"

var (
	SupabaseKey string
	SupabaseUrl string
)

// Read config.yml content and put the content into respective global variables
func init() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	SupabaseUrl = viper.GetString(supabaseUrlKey)
	SupabaseKey = viper.GetString(supabaseKeyKey)
	//fmt.Println(SupabaseUrl, SupabaseKey)
	if err != nil {
		fmt.Println("Error calling init in config.go", err)
		return
	}
}
