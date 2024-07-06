FROM python:3.12-slim

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmariadb3 \
    libmariadb-dev \
    gcc 

RUN apt-get clean

WORKDIR /app
RUN COPY ./requirements.txt /app/requirements.txt
# Install python dependencies
RUN pip install -r requirements.txt

# Define non-root user
USER nobody

# Move scripts
RUN mkdir -p /app/logs  

COPY ./main.py /app/main.py
COPY ./send_discord.py /app/send_discord.py



# Run the script
ENTRYPOINT ["python3" ,"-u", "main.py"]