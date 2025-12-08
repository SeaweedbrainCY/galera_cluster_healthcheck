package healthcheck

import (
	"database/sql"
	"strconv"

	"github.com/SeaweedbrainCY/galera_cluster_healthcheck/config"
	"go.uber.org/zap"
)

type HealthCheck struct {
	IsHealthy            bool
	ClusterSizeMsg       string
	ClusterStatusMsg     string
	NodeStatusMsg        string
	NodeConnectivityMsg  string
	IncomingAddressesMsg string
}

func (healthCheck *HealthCheck) CheckClusterSize(db *sql.DB, config *config.Config, clusterSize int, logger *zap.Logger) {
	var name string
	var size int
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'wsrep_cluster_size';").Scan(&name, &size)
	if err != nil {
		logger.Error("Failed to query wsrep_cluster_size", zap.Error(err))
		healthCheck.ClusterSizeMsg = "Error querying wsrep_cluster_size"
		healthCheck.IsHealthy = false
		return
	}

	if size < clusterSize {
		healthCheck.ClusterSizeMsg = "ERROR. Cluster size is less than expected. Was expecting " + strconv.Itoa(clusterSize) + " nodes, but only " + strconv.Itoa(size) + " are available"
		healthCheck.IsHealthy = false
	} else {
		healthCheck.ClusterSizeMsg = "OK. Cluster has " + strconv.Itoa(size) + " nodes. This is the expected size."
	}

}

func (healthCheck *HealthCheck) CheckClusterStatus(db *sql.DB, config *config.Config, clusterSize int, logger *zap.Logger) {
	var name string
	var status string
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'wsrep_cluster_status';").Scan(&name, &status)
	if err != nil {
		logger.Error("Failed to query wsrep_cluster_status", zap.Error(err))
		healthCheck.ClusterSizeMsg = "Error querying wsrep_cluster_status"
		healthCheck.IsHealthy = false
		return
	}

	if status != "Primary" {
		healthCheck.ClusterStatusMsg = "ERROR. Cluster is not in primary state. Got " + status + " state. This could mean that the cluster is partitioned and " + config.MariaDB_Node_Name + " is separated from the others."
		healthCheck.IsHealthy = false
	} else {
		healthCheck.ClusterStatusMsg = "OK. Cluster is in primary state. " + config.MariaDB_Node_Name + " is connected to the cluster."
	}
}

func (healthCheck *HealthCheck) CheckNodeStatus(db *sql.DB, config *config.Config, clusterSize int, logger *zap.Logger) {
	var name string
	var status string
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'wsrep_local_state_comment';").Scan(&name, &status)
	if err != nil {
		logger.Error("Failed to query wsrep_local_state_comment", zap.Error(err))
		healthCheck.ClusterSizeMsg = "Error querying wsrep_local_state_comment"
		healthCheck.IsHealthy = false
		return
	}

	switch status {
	case "Synced":
		healthCheck.NodeStatusMsg = "OK. " + config.MariaDB_Node_Name + " is in synced state."
	case "Initialized":
		healthCheck.NodeStatusMsg = "WARNING. " + config.MariaDB_Node_Name + " is in initialized state. This means that the node is not fully synced with the cluster yet."
		healthCheck.IsHealthy = false
	case "Joining":
		healthCheck.NodeStatusMsg = "WARNING. " + config.MariaDB_Node_Name + " is in joining state. This means that the node is in the process of joining the cluster."
		healthCheck.IsHealthy = false
	case "Donor/Desynced":
		healthCheck.NodeStatusMsg = "WARNING. " + config.MariaDB_Node_Name + " is in donor/desynced state. This means that the node is donating data to other nodes and is not fully synced."
		healthCheck.IsHealthy = false
	case "Desynced":
		healthCheck.NodeStatusMsg = "ERROR. " + config.MariaDB_Node_Name + " is in desynced state. This means that the node is out of sync with the cluster."
		healthCheck.IsHealthy = false
	default:
		healthCheck.NodeStatusMsg = "ERROR. " + config.MariaDB_Node_Name + " is in an unknown state: " + status + "."
		healthCheck.IsHealthy = false
	}
}

func (healthCheck *HealthCheck) CheckNodeConnectivity(db *sql.DB, config *config.Config, clusterSize int, logger *zap.Logger) {
	var name string
	var status string
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'wsrep_connected';").Scan(&name, &status)
	if err != nil {
		logger.Error("Failed to query wsrep_connected", zap.Error(err))
		healthCheck.ClusterSizeMsg = "Error querying wsrep_connected"
		healthCheck.IsHealthy = false
		return
	}
	if status != "ON" {
		healthCheck.NodeConnectivityMsg = "ERROR. " + config.MariaDB_Node_Name + " is not connected to the cluster. " + config.MariaDB_Node_Name + " is " + status
	} else {
		healthCheck.NodeConnectivityMsg = "OK. " + config.MariaDB_Node_Name + " is connected to the cluster."
	}
}

func (healthCheck *HealthCheck) CheckIncomingAddresses(db *sql.DB, config *config.Config, clusterSize int, logger *zap.Logger) {
	var name string
	var addresses string
	err := db.QueryRow("SHOW GLOBAL STATUS LIKE 'wsrep_incoming_addresses';").Scan(&name, &addresses)
	if err != nil {
		logger.Error("Failed to query wsrep_incoming_addresses", zap.Error(err))
		healthCheck.ClusterSizeMsg = "Error querying wsrep_incoming_addresses"
		healthCheck.IsHealthy = false
		return
	}

	healthCheck.IncomingAddressesMsg = "INFO. Connected addresses are: " + addresses
}
