#!/bin/bash

# Set the base URL for your application
BASE_URL="http://localhost:3000"

# Create a Ticket Option
echo "Creating a Ticket Option..."
CREATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/ticket_options" \
     -H "Content-Type: application/json" \
     -d '{"name":"Concert Tickets","desc":"Amazing concert event","allocation":1000}')
echo $CREATE_RESPONSE

# Extract the ticket option ID from the response
TICKET_ID=$(echo $CREATE_RESPONSE | jq -r '.id')

# Get the Ticket Option
echo "Retrieving the Ticket Option..."
curl -s -X GET "${BASE_URL}/ticket_options/${TICKET_ID}"

# Purchase Tickets
echo "Purchasing Tickets..."
curl -s -X POST "${BASE_URL}/ticket_options/${TICKET_ID}/purchases" \
     -H "Content-Type: application/json" \
     -d "{\"quantity\":2,\"user_id\":\"$(uuidgen)\"}"

# Get the Ticket Option again to check updated allocation
echo "Retrieving the Ticket Option after purchase..."
curl -s -X GET "${BASE_URL}/ticket_options/${TICKET_ID}"

# Try to get a non-existent Ticket Option
echo "Attempting to retrieve a non-existent Ticket Option..."
curl -s -X GET "${BASE_URL}/ticket_options/00000000-0000-0000-0000-000000000000"

# Try to purchase more tickets than available
echo "Attempting to purchase more tickets than available..."
curl -s -X POST "${BASE_URL}/ticket_options/${TICKET_ID}/purchases" \
     -H "Content-Type: application/json" \
     -d "{\"quantity\":1001,\"user_id\":\"$(uuidgen)\"}"

# Try to create an invalid Ticket Option
echo "Attempting to create an invalid Ticket Option..."
curl -s -X POST "${BASE_URL}/ticket_options" \
     -H "Content-Type: application/json" \
     -d '{"name":"","desc":"Invalid option","allocation":-10}'
