# Stage 1: Build the React SPA and install dependencies
FROM node:20-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Serve the static files
FROM node:20-alpine
WORKDIR /app

# Install the 'serve' package as a production dependency
RUN npm install -g serve

# Copy the static assets from the build stage to the final image
COPY --from=build /app/build/client ./

# The command to run when the container starts
CMD ["serve", "-s", "-l", "8080"]