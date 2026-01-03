package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/config"
	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/healthcheck"
	"go.uber.org/zap"
)

func SendNotification(config *config.Config, healthCheck *healthcheck.HealthCheck, logger *zap.Logger) error {
	embed := Embed{}
	embed.Title = config.MariaDB_Node_Name + ": The galera cluster is in error state"
	embed.Color = 16528693
	embed.Timestamp = time.Now().Format(time.RFC3339)

	embed.Footer = &EmbedFooter{
		Text: "Galera Cluster Healthcheck - " + config.Version,
	}

	if healthCheck.IsHealthy {
		embed.Title = config.MariaDB_Node_Name + ": The galera cluster is healthy again"
		embed.Color = 1220903
	}

	embed.Description = "Current galera cluster health at " + time.Now().Format(time.RFC1123)
	embed.Fields = []EmbedField{
		{
			Name:   "Cluster Size Check",
			Value:  healthCheck.ClusterSizeMsg,
			Inline: false,
		},
		{
			Name:   "Cluster Status Check",
			Value:  healthCheck.ClusterStatusMsg,
			Inline: false,
		},
		{
			Name:   "Node Status Check",
			Value:  healthCheck.NodeStatusMsg,
			Inline: false,
		},
		{
			Name:   "Node Connectivity Check",
			Value:  healthCheck.NodeConnectivityMsg,
			Inline: false,
		},
		{
			Name:   "Incoming Addresses",
			Value:  healthCheck.IncomingAddressesMsg,
			Inline: false,
		},
	}
	content := ""
	if config.Discord_Role_To_Mention != "" {
		content = "<@&" + config.Discord_Role_To_Mention + ">"
	}
	webhook := Webhook{
		Username: "Galera Cluster Healthcheck",
		Content:  content,
		Embeds:   []Embed{embed},
	}
	payload, err := json.Marshal(webhook)
	if err != nil {
		logger.Error("Failed to marshal Discord webhook payload", zap.Error(err))
		return err
	}
	resp, err := http.Post(config.Discord_Webhook_Url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Failed to send Discord webhook", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body := new(bytes.Buffer)
		_, _ = body.ReadFrom(resp.Body)
		logger.Error("Error while sending Discord webhook: received non-204/200", zap.Int("StatusCode", resp.StatusCode), zap.String("ResponseBody", body.String()))
	}

	logger.Info("Discord notification sent successfully")
	return nil
}
