FROM node:18.16-alpine

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install application dependencies
RUN npm install

# Copy the application code into the container
COPY . .

# # Expose the port the app runs on
# EXPOSE 3000

# Define the command to run your application
CMD ["node", "server.js"]