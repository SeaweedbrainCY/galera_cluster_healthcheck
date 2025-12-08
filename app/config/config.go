package config

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

type Config struct {
	Db_User                     string
	Db_Password                 string
	Db_Host                     string
	Db_Port                     int
	Check_Interval              int
	Alert_Throttle              int
	Discord_Webhook_Url         string
	Discord_Role_To_Mention     string
	MariaDB_Node_Name           string
	Galera_Cluster_Minimum_Size int
}

func LoadConfig(logger *zap.Logger) *Config {
	var config Config
	if os.Getenv("DB_USER") != "" {
		config.Db_User = os.Getenv("DB_USER")
	} else {
		panic("DB_USER environment variable is required")
	}

	if os.Getenv("DB_PASSWORD") != "" {
		config.Db_Password = os.Getenv("DB_PASSWORD")
	} else {
		panic("DB_PASSWORD environment variable is required")
	}

	if os.Getenv("DB_HOST") != "" {
		config.Db_Host = os.Getenv("DB_HOST")
	} else {
		panic("DB_HOST environment variable is required")
	}

	if os.Getenv("DB_PORT") != "" {
		db_port_int, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			panic("DB_PORT must be a valid integer")
		}
		config.Db_Port = db_port_int
	} else {
		config.Db_Port = 3306 // default port
		logger.Info("DB_PORT not set, using default 3306")
	}

	if os.Getenv("CHECK_INTERVAL") != "" {
		check_interval_int, err := strconv.Atoi(os.Getenv("CHECK_INTERVAL"))
		if err != nil {
			panic("CHECK_INTERVAL must be a valid integer")
		}
		config.Check_Interval = check_interval_int
	} else {
		config.Check_Interval = 60
		logger.Info("CHECK_INTERVAL not set, using default 60 seconds")
	}

	if os.Getenv("ALERT_THROTTLE") != "" {
		alert_throttle_int, err := strconv.Atoi(os.Getenv("ALERT_THROTTLE"))
		if err != nil {
			panic("ALERT_THROTTLE must be a valid integer")
		}
		config.Alert_Throttle = alert_throttle_int
	} else {
		config.Alert_Throttle = 10800
		logger.Info("ALERT_THROTTLE not set, using default 10800 seconds (3 hours)")
	}

	if os.Getenv("DISCORD_WEBHOOK_URL") != "" {
		config.Discord_Webhook_Url = os.Getenv("DISCORD_WEBHOOK_URL")
	} else {
		panic("DISCORD_WEBHOOK_URL environment variable is required")
	}

	if os.Getenv("DISCORD_ROLE_TO_MENTION") != "" {
		config.Discord_Role_To_Mention = os.Getenv("DISCORD_ROLE_TO_MENTION")
	} else {
		config.Discord_Role_To_Mention = ""
		logger.Info("DISCORD_ROLE_TO_MENTION not set, no role will be mentioned in alerts")
	}

	if os.Getenv("NODE_NAME") != "" {
		config.MariaDB_Node_Name = os.Getenv("NODE_NAME")
	} else {
		panic("NODE_NAME environment variable is required")
	}

	if os.Getenv("GALERA_CLUSTER_MIN_SIZE") != "" {
		min_size_int, err := strconv.Atoi(os.Getenv("GALERA_CLUSTER_MIN_SIZE"))
		if err != nil {
			logger.Fatal("Invalid GALERA_CLUSTER_MIN_SIZE", zap.Error(err))
		}
		config.Galera_Cluster_Minimum_Size = min_size_int
	} else {
		config.Galera_Cluster_Minimum_Size = 3
		logger.Info("GALERA_CLUSTER_MIN_SIZE not set, using default 3")
	}

	logger.Info("Configuration loaded successfully")

	return &config
}
