package main

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/config"
	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/discord"
	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/healthcheck"
	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/notification"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Version = "dev" // Will be set during build time

func DatabaseHealthCheck(db *sql.DB, config *config.Config, logger *zap.Logger) (*healthcheck.HealthCheck, error) {
	err := db.Ping()
	if err != nil {
		logger.Warn("Database ping failed. Attempting to reconnect", zap.Error(err))
		return nil, errors.New("DatabaseConnectionFailed")
	}
	var healthCheck healthcheck.HealthCheck

	healthCheck.IsHealthy = true

	healthCheck.CheckClusterSize(db, config, config.Galera_Cluster_Minimum_Size, logger)
	healthCheck.CheckClusterStatus(db, config, config.Galera_Cluster_Minimum_Size, logger)
	healthCheck.CheckNodeStatus(db, config, config.Galera_Cluster_Minimum_Size, logger)
	healthCheck.CheckNodeConnectivity(db, config, config.Galera_Cluster_Minimum_Size, logger)
	healthCheck.CheckIncomingAddresses(db, config, config.Galera_Cluster_Minimum_Size, logger)

	return &healthCheck, nil

}

func main() {
	loggerConfig := zap.NewProductionConfig()                          // Start with production defaults
	loggerConfig.Encoding = "console"                                  // Change the encoding to "console"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 time format
	logger, err := loggerConfig.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()
	logger.Info("Application started", zap.String("version", Version))

	config := config.LoadConfig(logger)

	logger.Info("Connecting to the database",
		zap.String("Db_Host", config.Db_Host),
		zap.Int("Db_Port", config.Db_Port),
		zap.String("Db_User", config.Db_User),
	)

	dsn := config.Db_User + ":" + config.Db_Password + "@tcp(" + config.Db_Host + ":" + strconv.Itoa(config.Db_Port) + ")/"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()
	if err != nil {
		logger.Fatal("Database ping failed", zap.Error(err))
	}

	logger.Info("Successfully connected to the database")

	ticker := time.NewTicker(time.Duration(config.Check_Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var healthCheck *healthcheck.HealthCheck
		healthCheck, err := DatabaseHealthCheck(db, config, logger)
		if err != nil {
			if err.Error() == "DatabaseConnectionFailed" {
				db.Close()
				db, err = sql.Open("mysql", dsn)
				if err != nil {
					logger.Error("Reconnection to database failed", zap.Error(err))
					continue
				}
				err = db.Ping()
				if err != nil {
					logger.Error("Database ping after reconnection failed", zap.Error(err))
					continue
				}
				logger.Info("Reconnected to the database successfully")
				healthCheck, err = DatabaseHealthCheck(db, config, logger)
				if err != nil {
					logger.Error("Health check failed after reconnection", zap.Error(err))
					continue
				}
			} else {
				logger.Error("Health check failed", zap.Error(err))
				continue
			}
		}
		logger.Info("Health performed",
			zap.Bool("IsHealthy", healthCheck.IsHealthy),
			zap.String("ClusterSizeMsg", healthCheck.ClusterSizeMsg),
			zap.String("ClusterStatusMsg", healthCheck.ClusterStatusMsg),
			zap.String("NodeStatusMsg", healthCheck.NodeStatusMsg),
			zap.String("NodeConnectivityMsg", healthCheck.NodeConnectivityMsg),
			zap.String("IncomingAddressesMsg", healthCheck.IncomingAddressesMsg),
		)
		should_trigger_new_notif, err := notification.ShouldSendNewNotification(healthCheck, config, logger)
		if err == nil && should_trigger_new_notif {
			err := discord.SendNotification(config, healthCheck, logger)
			if err != nil {
				notification.UpdateLastNotificationStatus(healthCheck, logger)
			}
		}
	}

}
