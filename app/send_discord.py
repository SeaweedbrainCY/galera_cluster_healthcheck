from discord_webhook import DiscordWebhook, DiscordEmbed
import os
from datetime import datetime
import sys

try:
    WEBHOOK_URL = os.environ['WEBHOOK_URL']
except Exception as e :
    print("Error while fetching environment variables. Exiting.")
    print("Required variables are : DB_USER, DB_PASSWORD, DB_HOST, NODE_NAME, CLUSTER_MIN_SIZE")
    print(f"Error : {e}")
    sys.exit(1)


def send(error_msg, include_status, status, is_error):
    if is_error:
        color="f38989"
        title="The galera cluster is in error state"
    else:
        color="78c1a3"
        title="The galera cluster is back to a normal state"
    
    if include_status:
        description = f"Current cluster state at {datetime.now()}"
    else:
        description = f"{error_msg} at {datetime.now()}"
    webhook = DiscordWebhook(url=WEBHOOK_URL)
    embed = DiscordEmbed(title=title, description=description, color=color)
    if include_status:
        for key in status:
            embed.add_embed_field(name=key, value=status[key]['msg'], inline=False)
    webhook.add_embed(embed)
    response = webhook.execute()
    print(f"{datetime.now()} : Discord answered {response}")
    return response.status_code