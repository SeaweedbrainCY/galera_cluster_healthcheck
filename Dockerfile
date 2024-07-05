FROM python:3.12-slim

# Install dependencies
RUN apt-get update && apt-get install -y \
    mysql-client 

RUN apt-get clean

# Install dependencies
RUN pip install discord-webhook mariadb


# Run the script
CMD ["/app/main.sh"]