# galera Cluster Healthcheck
This is a very basic docker container that helthcheck the nodes of a Mariadb Galera Cluster. 

## Features
### Monitoring a specific node
Use the healthcheck to monitor a specific and retrieve metrics from Galera. To determine the health of your cluster, the script will monitor : 
- The number of connected node in the cluster 
- The state of the cluster 
- The state of your node inside the cluster 
- The sync status of your cluster

**This metrics are fetch from a node, which mean that you can detect network partition, malfunction, lose of quorum or desync. Thus, to have a precise and global view of your cluster, it is recommanded to healthcheck every node of your cluster.**

### Discord notification
For now, the only mean of notification implemented is through discord webhook notification. 

Thus the script is made to be modulable with many other notification medium. 