from discord_webhook import DiscordWebhook, DiscordEmbed
import os
from datetime import datetime

WEBHOOK_URL = os.environ['WEBHOOK_URL']


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
            embed.add_embed_field(name=key, value=status[key]['msg'])
    webhook.add_embed(embed)
    response = webhook.execute()
    print(f"{datetime.now()} : Discord answered {response}")