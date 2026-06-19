package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL          string  `mapstructure:"DATABASE_URL"`
	JWTSecret            string  `mapstructure:"JWT_SECRET"`
	JWTExpiryHours       int     `mapstructure:"JWT_EXPIRY_HOURS"`
	SupabaseURL          string  `mapstructure:"SUPABASE_URL"`
	SupabaseServiceKey   string  `mapstructure:"SUPABASE_SERVICE_KEY"`
	SupabaseBucketProofs string  `mapstructure:"SUPABASE_BUCKET_PROOFS"`
	SupabaseBucketSlips  string  `mapstructure:"SUPABASE_BUCKET_SLIPS"`
	FCMCredentialsPath   string  `mapstructure:"FCM_CREDENTIALS_PATH"`
	RadarRadiusMeters    float64 `mapstructure:"RADAR_RADIUS_METERS"`
	RadarStaleMinutes    int     `mapstructure:"RADAR_STALE_MINUTES"`
	Reward               RewardConfig
	CronOverdueExpense   string `mapstructure:"CRON_OVERDUE_EXPENSE"`
	CronOverdueMission   string `mapstructure:"CRON_OVERDUE_MISSION"`
	CronScraper          string `mapstructure:"CRON_SCRAPER"`
	Port                 string `mapstructure:"PORT"`
}

type RewardConfig struct {
	SpeedBonus7dPct    int `mapstructure:"REWARD_SPEED_BONUS_7D_PCT"`
	SpeedBonus1dPct    int `mapstructure:"REWARD_SPEED_BONUS_1D_PCT"`
	Streak3Pct         int `mapstructure:"REWARD_STREAK_3_PCT"`
	Streak7Pct         int `mapstructure:"REWARD_STREAK_7_PCT"`
	FirstCompleterFlat int `mapstructure:"REWARD_FIRST_COMPLETER_FLAT"`
}

func MustLoad() *Config {
	log.Println("debugprint: entering MustLoad")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, using environment variables: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Unmarshal reward config separately if needed or ensure mapstructure tags match
	if err := viper.Unmarshal(&cfg.Reward); err != nil {
		log.Fatalf("Unable to decode reward config, %v", err)
	}

	// Validation
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return &cfg
}

func (c *Config) JWTExpiry() time.Duration {
	log.Println("debugprint: entering (*Config).JWTExpiry")
	return time.Duration(c.JWTExpiryHours) * time.Hour
}
