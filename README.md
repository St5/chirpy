# Chirpy Server API tool
Education project for learn build server on GO

## Introduce
This project covers various aspects of building an API server using Go. It includes:

- Authentication and authorization of users
- CRUD operations for Chirpy (small posts similar to tweets)
- Building and managing a web server

This educational project aims to provide a comprehensive understanding of server-side development with Go.

## Docker Setup

This project requires PostgreSQL, which can be easily set up using Docker. To start the PostgreSQL container, run the following command:

```sh
docker-compose --compatibility up -d
```
For more commands and credentials, refer to the docker-compose.yml file. This file contains all the necessary configurations for running the PostgreSQL container and other services required by the project. 

## Instalation
To install and run the Chirpy Server API tool, follow these steps:

1. **Clone the repository:**
    ```sh
    git clone https://github.com/St5/chirpy.git
    cd chirpy
    ```

2. **Install dependencies:**
    Make sure you have Go installed on your machine. Then, run:
    ```sh
    go mod tidy
    ```

3. **Set up environment variables:**
    Create a `.env` file in the root directory and add the necessary environment variables. For example:
    ```sh
    DB_URL="YOUR_CONNECTION_STRING_HERE"
    TOKEN_SECRET="TOKEN_SECRET"
    POLKA_KEY="WEBHOOK KEY"
    ```
    DB_URL is the connection string to PostgreSQL with password and username. 
    TOKEN_SECRET is the secret key for generating JWT tokens. POLKA_KEY is the key for the webhook.

4. **Run the server:**
    ```sh
    go run main.go
    ```

5. **Access the API:**
    Open your browser or API client and navigate to `http://localhost:8585`. You can change port in main.go file.

Now you should have the Chirpy Server API tool up and running on your local machine.

## API Endpoints
The Chirpy Server API tool provides the following endpoints:

- `GET /api/chirps`: Get all chirps
- `GET /api/chirps/:id`: Get a chirp by ID
- `POST /api/chirps`: Create a new chirp
- `PUT /api/chirps/:id`: Update a chirp by ID
- `DELETE /api/chirps/:id`: Delete a chirp by ID
- `POST /api/users`: Register a new user
- `PUT /api/users`: Update a user
- `POST /api/login`: Login a user
- `POST /api/refresh`: Refresh the JWT token by providing a valid refresh token
- `POST /api/revoke`: Revoke refresh tokens
- `POST /api/revopolka/webhooks`: A webhook to mark chirpy red for a user
- `GET /admin/reset`: Reset the database and all entries
- `/app/`: Web interface to return file content from public folder
- `GET /admin/metrics`: Calculate the metrics visiting of the server. Result.
- `GET /admin/healthz`: Calculate the metrics visiting of the server.


## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
```

