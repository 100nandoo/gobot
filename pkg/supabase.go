package pkg

import (
	"github.com/nedpals/supabase-go"
	"gobot/config"
	"os"
)

var SupabaseClient = supabase.CreateClient(os.Getenv(config.SupabaseUrl), os.Getenv(config.SupabaseKey))
