# Galera Cluster Healthcheck
<div align="center">
 <img src="https://github.com/SeaweedbrainCY/galera_cluster_healthcheck/actions/workflows/deploy.yml/badge.svg/"> <img src="https://github.com/SeaweedbrainCY/galera_cluster_healthcheck/actions/workflows/security_scan.yml/badge.svg/"> <img src="https://img.shields.io/github/v/tag/SeaweedbrainCY/galera_cluster_healthcheck"/> <img src="https://img.shields.io/github/license/seaweedbraincy/galera_cluster_healthcheck"/>
</div>
<br>
This is a very basic docker container that helthcheck the nodes of a Mariadb Galera Cluster. 


**Contributors :**

![GitHub Contributors Image](https://contrib.rocks/image?repo=seaweedbraincy/galera_cluster_healthcheck)
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

**When the cluster is healthy :**

<img width="40%" src=https://github.com/SeaweedbrainCY/galera_cluster_healthcheck/assets/42048771/c66ee07b-427d-4a6e-ab2e-3dd0741eceee)/>

**When an error is detected :**

<img width="40%" src=https://github.com/SeaweedbrainCY/galera_cluster_healthcheck/assets/42048771/f292048f-a1bb-4c66-8fd5-6e3f806460fb)/>


Thus the script is made to be modulable with many other notification medium. 

## Installation
### Docker compose
You can use the ![docker-compose.yml](https://raw.githubusercontent.com/SeaweedbrainCY/galera_cluster_healthcheck/main/docker-compose.yml) file from the repo, or directly this config : 
```yml
version: '3.7'

services:
  galera_health_check:
    container_name: galera_health_check
    image: ghcr.io/seaweedbraincy/galera_cluster_healthcheck:b0.3
    user: 1000:1000
    security_opt:
      - no-new-privileges:true
    read_only: true
    environment:
       DB_USER: <REPLACE_ME> # required
       DB_PASSWORD: <REPLACE_ME> # required
       DB_HOST: <REPLACE_ME> # required
       DB_PORT: 3306 # optional. Default is 3306
       CHECK_INTERVAL: 60 # optional. Default is 60. Healthcheck interval in seconds
       ALERT_THROTTLE: 10800 # optional. Default is 10800s (3h). Amount of time in senconds between 2 consecutive alerts
       NODE_NAME: <REPLACE_ME> # required. Used to identify the node in the alert message
       CLUSTER_MIN_SIZE: <REPLACE_ME> # required. Minimum number of nodes in the cluster, used to determine if the cluster is in a degraded state
       WEBHOOK_URL: <REPLACE_ME> # required. Discord webhook URL for notifications
    volumes:
      - ./docker_logs:/app/logs 
    restart: always
```
### Configuration
You can use several env variable to configure the script : 

| Variable    | Is required | Defintion| 
| -------- | ------- | ------- |
| DB_USER  |   **REQUIRED**  | Mariadb user used to connect to the database. The user should only be able to read mariadb GLOBAL STATUS variables |
| DB_PASSWORD |  **REQUIRED**     |Password of the user used to connect to the database.|
| DB_HOST    |  **REQUIRED**    | Mariadb host|
| DB_PORT    | *optional*    |Port used by mariadb. Default value : 3306|
| ALERT_THROTTLE    | *optional* | Amount of time in senconds between 2 consecutive alerts. Default value : 10800s (3h).|
| CHECK_INTERVAL  | *optional* | Healthcheck interval in seconds. Default value: 60|
| NODE_NAME    |  **REQUIRED**    |Used to identify the node in the alert message. |
| CLUSTER_MIN_SIZE    |  **REQUIRED**    |Minimum number of nodes in the cluster. This is used to determine if the cluster is in a degraded state|
| WEBHOOK_URL    |  **REQUIRED**    |Discord webhook URL used to send notifications|


## Contribution

All contributions are welcome ! Feel free to open a merge request 

## Licence 

The code is under MIT Licence
