#!/bin/bash
# startup.sh

# Function to check if Cassandra is up and accepting cqlsh connections
wait_for_cassandra() {
    echo "Waiting for Cassandra to be available..."
    until cqlsh cassandra 9042 -e "describe cluster"; do
        sleep 1
        echo "Waiting for Cassandra..."
    done
    echo "Cassandra is up and running."
}

# Function to create the keyspace in Cassandra
create_keyspace() {
    echo "Creating keyspace if not exists..."
    CQL="CREATE KEYSPACE IF NOT EXISTS nihongowa WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };"
    cqlsh cassandra 9042 -e "$CQL"
    echo "Keyspace created or already exists."
}

# Main execution
# 1. Wait for Cassandra
wait_for_cassandra

# 2. Create the keyspace
create_keyspace

# 3. Start the Go application
echo "Starting the Go application..."
exec ./server
