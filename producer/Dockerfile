# Use an official Node.js runtime as a parent image
FROM node:18.16.0-alpine

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy package.json and package-lock.json to the container
COPY package*.json ./

# Install app dependencies
RUN npm install

# Bundle app source
COPY . .

# Expose the port the app runs on
EXPOSE 3010

# Define environment variables
ENV MONGO_ATLAS_PW='A4sFBBa26o09FB8E'
ENV PORT=3010

# Wait for MongoDB to be ready before starting the application
CMD ["./wait-for-mongo.sh", "npm", "start"]
