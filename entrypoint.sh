#!/bin/sh

echo "Waiting for Kafka on port 29092..."
while ! nc -z kafka 29092; do
  sleep 1
done

echo "Kafka is ready!"

./go-service