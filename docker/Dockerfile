# Use an official Node.js base image with the desired version
FROM node:19.5.0

# Set the maintainer label
LABEL maintainer="Azraf Al Moznim dev@monzim.com"

# Set the image name and version label
LABEL name="uiudiscordbot" version="1.0.3"

# Set other labels for description, author, and license
LABEL description="This is Bot for UIU made by @monzim"
LABEL author="Azraf Al Monzim"
LABEL license="ISC"

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install app dependencies
RUN npm install

# Copy the rest of the app files to the working directory
COPY . .

# Build the TypeScript files
RUN npm run build

# Start the app
CMD ["node", "-r", "dotenv/config", "./prod/index.js"]
