package main

import (
	"fmt"
	"github.com/spf13/viper"
)

const supabaseUrlKey, supabaseKeyKey string = "Supabase.url", "Supabase.key"

var (
	SupabaseKey string
	SupabaseUrl string
)

func Init() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	SupabaseUrl = viper.GetString(supabaseUrlKey)
	SupabaseKey = viper.GetString(supabaseKeyKey)
	if err != nil {
		fmt.Println("Error calling init in config.go", err)
		return
	}
}
