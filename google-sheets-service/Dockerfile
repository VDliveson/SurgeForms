FROM python:3.12

# Set the working directory in the container
WORKDIR /main

# Copy the local code to the container
COPY . /main

# Install any dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Command to run on container start
CMD ["python", "app.py"]