package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const supabaseUrlKey, supabaseKeyKey string = "Supabase.url", "Supabase.key"
const telegramAaronKey string = "Telegram.bot.aaron"
const telegramFreeGamesDebugKey string = "Telegram.channel.free-games-debug"
const telegramFreeGamesKey string = "Telegram.channel.free-games"

var (
	SupabaseKey            string
	SupabaseUrl            string
	TelegramAaron          string
	TelegramFreeGamesDebug int64
	TelegramFreeGames      int64
)

// Read config.yml content and put the content into respective global variables
func init() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	SupabaseUrl = viper.GetString(supabaseUrlKey)
	SupabaseKey = viper.GetString(supabaseKeyKey)
	TelegramAaron = viper.GetString(telegramAaronKey)
	TelegramFreeGamesDebug = viper.GetInt64(telegramFreeGamesDebugKey)
	TelegramFreeGames = viper.GetInt64(telegramFreeGamesKey)
	//fmt.Println(TelegramAaron)
	if err != nil {
		fmt.Println("Error calling init in config.go", err)
		return
	}
}
