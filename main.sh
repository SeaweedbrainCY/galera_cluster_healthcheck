#! /bin/bash 


size=$(mysql --user="$user" --password="$password" --execute='SHOW GLOBAL STATUS LIKE "wsrep_cluster_size";')
