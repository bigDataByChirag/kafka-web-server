# Use the official Kafka image from Docker Hub
FROM bitnami/kafka:latest

# Expose Kafka and Zookeeper ports
EXPOSE 9092
EXPOSE 2181

# Use the official Redis image from Docker Hub
FROM redis:latest

# Expose Redis port
EXPOSE 6379
