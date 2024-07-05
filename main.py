import mariadb
import sys
import os
from time import sleep
from datetime import datetime
import send_discord

DB_USER = os.environ['DB_USER']
DB_PASSWORD = os.environ['DB_PASSWORD']
DB_HOST = os.environ['DB_HOST'] 
DB_PORT = os.environ['DB_PORT'] if 'DB_PORT' in os.environ else 3306
CHECK_INTERVAL = os.environ['CHECK_INTERVAL'] if 'CHECK_INTERVAL' in os.environ else 60
ERROR_THRESHOLD = os.environ['ERROR_THRESHOLD'] if 'ERROR_THRESHOLD' in os.environ else 10800 #3h in seconds
NODE_NAME = os.environ['NODE_NAME'] if 'NODE_NAME' in os.environ else "node"
CLUSTER_MIN_SIZE = os.environ['CLUSTER_MIN_SIZE'] 



def log_error(error_msg, fatal=False):
    print(error_msg)
    if check_error_threshold():
        log_error_date()
        send_discord.send(error_msg)
    if fatal:
        sys.exit(1)



def log_error_date():
        with open('last_error_date.log', 'w') as f:
            f.write(f"{datetime.now()}")

def get_last_error_date():
    with open('last_error_date.log', 'r') as f:
        return f.read()

def check_error_threshold():
    last_error_date = get_last_error_date()
    last_error_date = datetime.strptime(last_error_date, '%Y-%m-%d %H:%M:%S.%f')
    current_date = datetime.now()
    if (current_date - last_error_date).seconds > ERROR_THRESHOLD:
        return True
    return False


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
            log_error(f"Error connecting to MariaDB : {e}. Docker is going to restart now.", True)


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
        except Exception as e:
            log_error(f"Error while fetching data: {e}", fatal=True)


        vars_to_check = [
                {"var" :size,
                 "name": "cluster size (nodes)",
                  "expected": CLUSTER_MIN_SIZE},
                {"var" :quorum_status,
                    "name": "quorum status" },
                {"var" :node_status,
                    "name": f"{NODE_NAME} status" },
                {"var" :node_connectivity,
                    "name": f"{NODE_NAME} connectivity with other nodes" },
                {"var" :incoming_addresses,
                    "name": "incoming addresses" }
        ]
        for to_check in vars_to_check:
            if to_check["var"] is None:
                log_error(f"Could not verify the following information : {to_check['name']}. Got None")
            else:


            

        sleep(int(CHECK_INTERVAL))
    
    
    



