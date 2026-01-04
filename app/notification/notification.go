package notification

import (
	"os"
	"strings"
	"time"

	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/config"
	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/healthcheck"
	"go.uber.org/zap"
)

func getLastNotificationStatus(logger *zap.Logger) (string, time.Time, error) {
	last_notification_file_path := "./logs/last_notification_date"
	if _, err := os.Stat(last_notification_file_path); err != nil {
		logger.Info("Last notification file does not exist. Assuming no previous notifications sent.", zap.Error(err))
		return "", time.Now(), ErrorWhileOpeningFile
	}

	last_notification_content, err := os.ReadFile(last_notification_file_path)
	if err != nil {
		logger.Warn("Last notification file can't be open. Assuming no previous notifications sent.", zap.Error(err))
		return "", time.Now(), ErrorWhileReadingFile
	}

	file_parts := strings.Split(string(last_notification_content), "|")
	if len(file_parts) < 2 {
		logger.Warn("Last notification file is malformed. Assuming no previous notifications sent.")
		return "", time.Now(), NotificationFileMalformed
	}

	last_notification_status := file_parts[0]
	last_notification_date_srt := file_parts[1]

	last_notification_date, err := time.Parse(time.RFC3339, last_notification_date_srt)

	if err != nil {
		logger.Warn("Last notification date is malformed. Assuming no previous notifications sent.", zap.Error(err))
		return "", time.Now(), NotificationFileMalformed
	}

	return last_notification_status, last_notification_date, nil
}

func ShouldSendNewNotification(healthCheck *healthcheck.HealthCheck, config *config.Config, logger *zap.Logger) (bool, error) {
	last_notification_status, last_notification_date, err := getLastNotificationStatus(logger)
	if err != nil {
		return healthCheck.IsHealthy == false, nil
	}

	time_since_last_notification := time.Since(last_notification_date).Seconds()
	if (last_notification_status == "OK" && !healthCheck.IsHealthy) || (last_notification_status == "KO" && healthCheck.IsHealthy) {
		// Status changed, send notification
		return true, nil
	} else {
		if !healthCheck.IsHealthy && time_since_last_notification >= float64(config.Alert_Throttle) {
			// Still unhealthy, but throttle time passed, send notification
			return true, nil
		} else if !healthCheck.IsHealthy {
			// Still unhealthy, but throttle time not passed, do not send notification
			logger.Info("Not sending notification due to alert throttle",
				zap.Float64("time_since_last_notification_seconds", time_since_last_notification),
				zap.Int("alert_throttle_seconds", config.Alert_Throttle),
			)
			return false, nil
		}
		// Healthy and no status change, do not send notification
		return false, nil
	}
}

func UpdateLastNotificationStatus(healthCheck *healthcheck.HealthCheck, logger *zap.Logger) error {
	last_notification_file_path := "./logs/last_notification_date"
	status := "OK"
	if !healthCheck.IsHealthy {
		status = "KO"
	}
	current_time := time.Now().Format(time.RFC3339)
	content := status + "|" + current_time
	err := os.WriteFile(last_notification_file_path, []byte(content), 0644)
	if err != nil {
		logger.Error("Failed to update last notification file", zap.Error(err))
		return err
	}
	return nil
}

func InitLastNotificationFile(logger *zap.Logger) {
	last_notification_file_path := "./logs/last_notification_date"
	file, err := os.Open(last_notification_file_path)
	if os.IsNotExist(err) {
		epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		err := os.WriteFile(last_notification_file_path, []byte("OK|"+epoch.Format(time.RFC3339)), 0644)
		if err != nil {
			logger.Fatal("Failed to initialize last notification file", zap.Error(err))
		} else {
			logger.Info("Initialized last notification file with default values")
		}
	} else if err != nil {
		logger.Fatal("Failed to open last notification file", zap.Error(err))
	}
	defer file.Close()
}
