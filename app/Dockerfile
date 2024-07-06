FROM python:3.12-slim

# Install dependencies
RUN apt-get update && apt-get install -y \
    libmariadb3 \
    libmariadb-dev \
    gcc 

RUN apt-get clean


# Install python dependencies
RUN pip install discord-webhook  mariadb

# Move scripts
RUN mkdir -p /app/logs  

COPY ./main.py /app/main.py
COPY ./send_discord.py /app/send_discord.py

WORKDIR /app

# Run the script
ENTRYPOINT ["python3" ,"-u", "main.py"]