# Use the official Python image with Alpine Linux
FROM python:3.9-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the local code to the container
COPY . /app

# Install any dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Command to run on container start
CMD ["python", "app.py"]
