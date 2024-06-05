Kafka-Redis-Gin Application
===========================

This application demonstrates a basic setup for consuming messages from a Kafka topic, storing them in Redis, and providing a REST API using the Gin framework to interact with the stored data.

Prerequisites
-------------

-   Go 1.16+
-   Kafka
-   Redis
-   Make sure to install the following Go packages:

   `go get github.com/IBM/sarama
    go get github.com/gin-gonic/gin
    go get github.com/redis/go-redis/v9`

Setup
-----

1.  **Kafka Setup**

    Ensure that Kafka is running on `localhost:9092`. Create a topic named `data`:

    `kafka-topics.sh --create --topic data --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1`

2.  **Redis Setup**

    Ensure that Redis is running on `localhost:6379`.

Running the Application
-----------------------

1.  **Build and Run**

    `go build -o app
    ./app`

2.  **Using Docker (Optional)**

    If you prefer using Docker, create a Dockerfile with the necessary instructions to run this application, and build and run the Docker image.

API Endpoints
-------------

-   **POST /api/post**

    Publish data to Kafka.

    **Request:**

    json

    `{
      "userId": "user1",
      "data": "some data"
    }`

    **Response:**

    json

    `{
      "status": "success"
    }`

-   **GET /api/get/**

    Retrieve data from Redis.

    **Response:**

    json

    `{
      "data": "some data"
    }`

How It Works
------------

1.  **Publishing to Kafka:**

    -   The POST endpoint accepts a JSON payload with `userId` and `data`.
    -   The data is published to the Kafka topic `data` with `userId` as the key.
2.  **Consuming from Kafka:**

    -   A background goroutine continuously consumes messages from the Kafka topic `data`.
    -   The consumed messages are stored in Redis with `userId` as the key.
3.  **Fetching from Redis:**

    -   The GET endpoint retrieves the stored data from Redis using the provided `userId`.

Logging
-------

The application logs important events, such as successful publishing to Kafka and storing data in Redis, as well as any errors encountered during these operations.

Error Handling
--------------

The application provides appropriate HTTP status codes and error messages for various error conditions, such as invalid request payloads and data not found in Redis.

License
-------

This project is licensed under the MIT License.
