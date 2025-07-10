FROM python:3.13-slim

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmariadb3 \
    libmariadb-dev \
    gcc 

RUN apt-get clean


RUN mkdir -p /app/logs  
WORKDIR /app

# Install python dependencies
COPY ./requirements.txt /app/requirements.txt
RUN pip install -r requirements.txt

# Define non-root user
USER nobody

# Move scripts
COPY ./main.py /app/main.py
COPY ./send_discord.py /app/send_discord.py



# Run the script
ENTRYPOINT ["python3" ,"-u", "main.py"]