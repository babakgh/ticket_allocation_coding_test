# Ticket Allocation System

This project is a ticket allocation system built with Go, using PostgreSQL as the database.

## Development Notes

### Testing Approach
To meet the project deadline, unit tests were not implemented in this version. Instead, a bash script has been provided to test the main functionalities of the system. This approach allowed for quicker development and still provides a way to verify the system's core features.

### Concurrency Control
This system uses a Pessimistic Concurrency Control (PCC) approach for handling concurrent ticket purchases. This decision was made due to uncertainty about the client-side implementation that would be using this API.

Alternative Approach: An Optimistic Concurrency Control (OCC) could potentially be more efficient in some scenarios, but it would require a retrying mechanism on the client side. The PCC approach was chosen for its simplicity and to avoid potential issues with clients that might not handle retries correctly.

## Prerequisites

- Docker

## Starting the Application

To start the application, follow these steps:

1. Clone the repository:
   ```
   git clone git@github.com:babakgh/ticket_allocation_coding_test.git
   cd ticket-allocation-system
   ```

2. Start the application using Docker Compose:
   ```
   docker compose up --build
   ```

   This command will build the Docker images and start the containers for the application and the database.

3. The application should now be running and accessible at `http://localhost:3000`.

## Running Integration Tests

To run the integration tests, follow these steps:

1. Ensure docker compose is running:
   ```
   docker compose up --build
   ```

2. Run the integration test script:
   ```
   ./tests/simple_integration_test.sh
   ```

## API Endpoints

The following endpoints are available:

- `POST /ticket_options`: Create a new ticket option
- `GET /ticket_options/:id`: Get a specific ticket option
- `POST /ticket_options/:id/purchases`: Purchase tickets for a specific option

## Troubleshooting

If you encounter any issues:

1. Verify that the necessary ports (3000 for the app, 5432 for PostgreSQL) are not in use by other applications
2. Stop the docker compose process and run it again:
   ```
   docker compose down -v
   ```
