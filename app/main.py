import mariadb
import sys
import os
from time import sleep
from datetime import datetime
import send_discord as send_discord

try:
    DB_USER = os.environ['DB_USER']
    DB_PASSWORD = os.environ['DB_PASSWORD']
    DB_HOST = os.environ['DB_HOST'] 
    DB_PORT = os.environ['DB_PORT'] if 'DB_PORT' in os.environ else 3306
    CHECK_INTERVAL = os.environ['CHECK_INTERVAL'] if 'CHECK_INTERVAL' in os.environ else 60
    ALERT_THROTTLE = int(os.environ['ALERT_THROTTLE']) if 'ALERT_THROTTLE' in os.environ else 10800 #3h in seconds
    NODE_NAME = os.environ['NODE_NAME'] 
    CLUSTER_MIN_SIZE = os.environ['CLUSTER_MIN_SIZE'] 
except Exception as e :
    print("Error while fetching environment variables. Exiting.")
    print("Required variables are : DB_USER, DB_PASSWORD, DB_HOST, NODE_NAME, CLUSTER_MIN_SIZE")
    print(f"Error : {e}")
    sys.exit(1)


STATUS = {
        "cluster_size":{
        "msg": "",
    },
        "cluster_status": {
        "msg": "",
    },
        "node_status":{
        "msg": "",
    },
        "node_connectivity":{
        "msg": "",
    },
        "incoming_addresses":{
        "msg": "",
    }
}

IS_ON_ERROR = False

def log_message(error_msg, is_error=True,fatal=False, include_status=False):
    fire_available = can_fire_alert(is_error=is_error)
    if fire_available:
        log_last_notif_date(is_error=is_error)
        send_discord.send(error_msg, include_status, STATUS, is_error)
        print(f"{datetime.now()} : Notification sent")
        print(f"{datetime.now()} : New status : {STATUS}")
    else:
        if  is_error:
            print(f"{datetime.now()} : Cluster still in error state. No notification sent because of treshold.")
        else :
            print(f"{datetime.now()} : Cluster health is normal. All checks passed. Next check in {CHECK_INTERVAL} seconds.")
    if fatal:
        sys.exit(1)



def log_last_notif_date(is_error):
        status = "error" if is_error else "normal"
        with open('logs/last_error_date.log', 'w') as f:
            f.write(f"{datetime.now()},{status}")

def get_last_notif_date():
    if not os.path.isfile('logs/last_error_date.log'):
        return "2001-01-01 01:01:01.0000" , "normal" # default value

    with open('logs/last_error_date.log', 'r') as f:
        log = f.read()
    return log.split(',')[0], log.split(',')[1]

def can_fire_alert(is_error):
    last_notif_date,status = get_last_notif_date()
    was_last_notif_error = status == "error"
    last_notif_date = datetime.strptime(last_notif_date, '%Y-%m-%d %H:%M:%S.%f')
    current_date = datetime.now()
    return ((current_date - last_notif_date).seconds > ALERT_THROTTLE and is_error) or was_last_notif_error != is_error 


if __name__ == "__main__":
    while True:
        # Connect to MariaDB Platform
        try:
            conn = mariadb.connect(
                user=DB_USER,
                password=DB_PASSWORD,
                host=DB_HOST,
                port=DB_PORT
            )
        except mariadb.Error as e:
            log_message(f"Error connecting to MariaDB : {e}. Docker is going to restart now.", True)


        # Get Cursor
        cur = conn.cursor()
        try:
            cur.execute("SHOW GLOBAL STATUS LIKE 'wsrep_cluster_size';")
            size = cur.fetchone()
            cur.execute("SHOW GLOBAL STATUS LIKE 'wsrep_cluster_status';")
            quorum_status = cur.fetchone()
            cur.execute("SHOW GLOBAL STATUS LIKE 'wsrep_local_state_comment';")
            node_status = cur.fetchone()
            cur.execute("SHOW GLOBAL STATUS LIKE 'wsrep_connected';")
            node_connectivity = cur.fetchone()
            cur.execute("SHOW GLOBAL STATUS LIKE 'wsrep_incoming_addresses';")
            incoming_addresses = cur.fetchone()
            conn.close()
        except Exception as e:
            log_message(f"Error while fetching data: {e}", fatal=True)


        if size is None :
            STATUS["cluster_size"]["msg"] = "Error while fetching cluster size. Got None value."
            IS_ON_ERROR = True
        else:
            try:
                size = size[1]
                if int(size) < int(CLUSTER_MIN_SIZE):
                    STATUS["cluster_size"]["msg"] = f"ERROR. Cluster size is less than expected. Was expecting {CLUSTER_MIN_SIZE} nodes, but only {size} are available"
                    IS_ON_ERROR = True
                else:
                    STATUS["cluster_size"]["msg"] = f"Cluster has {size} nodes. This is the expected size."
            except Exception as e:
                STATUS["cluster_size"]["msg"] = f"Error while parsing cluster size: {e}"
                IS_ON_ERROR = True
        
        if quorum_status is None:
            STATUS["cluster_status"]["msg"] = "Error while fetching cluster status. Got None value"
            IS_ON_ERROR = True
        else:
            try:
                quorum_status = quorum_status[1]
            except Exception as e:
                STATUS["cluster_status"]["msg"] = f"Error while parsing cluster status: {e}"
                IS_ON_ERROR = True
            if quorum_status != "Primary":
                STATUS["cluster_status"]["msg"] = f"ERROR. Cluster is not in primary state. Got {quorum_status} state. This means that the cluster is partitioned and {NODE_NAME} separated of others."
                IS_ON_ERROR = True
            else:
               STATUS["cluster_status"]["msg"] = f"Cluster is in primary state. {NODE_NAME} is connected to the others."
        
        if node_status is None:
            STATUS["node_status"]["msg"] = "Error while fetching node status. Got None value."
            IS_ON_ERROR = True
        else:
            try:
                node_status = node_status[1]
            except Exception as e:
                STATUS["node_status"]["msg"] = f"Error while parsing node status: {e}"
                IS_ON_ERROR = True
            if node_status != "Synced":
                STATUS["node_status"]["msg"] = f"ERROR. {NODE_NAME} is not synced, and in a transitionary state. {NODE_NAME} is {node_status}"
                IS_ON_ERROR = True
            elif node_status == "Initialized":
                STATUS["node_status"]["msg"] = f"ERROR. {NODE_NAME} is not synced. {NODE_NAME} is {node_status}. {NODE_NAME} is NOT operational."
                IS_ON_ERROR = True
            else:
                STATUS["node_status"]["msg"] = f"{NODE_NAME} is synced."
        
        if node_connectivity is None:
            STATUS["node_connectivity"]["msg"] = "Error while fetching node connectivity. Got None value."
            IS_ON_ERROR = True
        else:
            try:
                node_connectivity = node_connectivity[1]
            except Exception as e:
                STATUS["node_connectivity"]["msg"] = f"Error while parsing node connectivity: {e}"
                IS_ON_ERROR = True
            if node_connectivity != "ON":
                STATUS["node_connectivity"]["msg"] = f"ERROR. {NODE_NAME} is not connected to the cluster. {NODE_NAME} is {node_connectivity}"
                IS_ON_ERROR = True
            else:
                STATUS["node_connectivity"]["msg"] = f"{NODE_NAME} is connected to the cluster."
        
        if incoming_addresses is None:
            STATUS["incoming_addresses"]["msg"] = "Error while fetching incoming addresses. Got None value."
            IS_ON_ERROR = True
        else:
            try:
                incoming_addresses = incoming_addresses[1]
            except Exception as e:
                STATUS["incoming_addresses"]["msg"] = f"Error while parsing incoming addresses: {e}"
                IS_ON_ERROR = True
            STATUS["incoming_addresses"]["msg"] = f"Incoming addresses are {incoming_addresses}"

        
        log_message("Checking cluster status", is_error=IS_ON_ERROR, include_status=True, fatal=False)

            

        sleep(int(CHECK_INTERVAL))
    
    
    



