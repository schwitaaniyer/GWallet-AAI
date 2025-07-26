# Build stage
FROM debian:latest AS build-env

# Install Flutter dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    unzip \
    xz-utils \
    zip \
    libglu1-mesa \
    openjdk-11-jdk \
    && rm -rf /var/lib/apt/lists/*

# Clone Flutter
RUN git clone https://github.com/flutter/flutter.git /flutter
ENV PATH="/flutter/bin:${PATH}"

# Enable web
RUN flutter channel stable
RUN flutter upgrade
RUN flutter config --enable-web

# Copy app files
WORKDIR /app
COPY pubspec.yaml .
COPY lib/ lib/
COPY assets/ assets/

# Get app dependencies
RUN flutter pub get

# Build the app for the web
RUN flutter build web

# Production stage
FROM nginx:alpine

# Copy built app from previous stage
COPY --from=build-env /app/build/web /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port
EXPOSE 8080

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 