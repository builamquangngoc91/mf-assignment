# Docker Setup for Project

## Prerequisites

Make sure you have the following tools installed on your machine:

- Docker
- Docker Compose

## Getting Started

1. Clone the repository:

   ```bash
   https://github.com/builamquangngoc91/mf-assignment
   ```
2. Navigate to the project directory:
    ``` 
    cd mf-assignment
    ```
3. Build and start the Docker containers:
    ```
    docker-compose up -d
    ```
4. Run migration files (in imgrations directory) using docker exec:
    ```
    docker exec -it <mf-db> bash
    
    psql -U postgres
    
    CREATE TABLE banking;
    
    \c banking

    <Run file by file in migrations directory>
    ```
5. Access the application in your browser:
    ```
    http://localhost:8081
    ```

### list APIs
- Create User
    ```
    curl --location 'localhost:8081/users' \
    --header 'Content-Type: application/json' \
    --data '{
        "name": "Alice"
    }'
    ```
- Get users
    ```
    curl --location 'localhost:8081/users'
    ```
- Get users/:userID
    ```
    curl --location 'localhost:8081/users/7a6eead1-0d62-41d7-bf51-8984cdb918fc'
    ```
- Create account
    ```
    curl --location 'localhost:8081/accounts' \
    --header 'Content-Type: application/json' \
    --data '{
        "user_id": "7a6eead1-0d62-41d7-bf51-8984cdb918fc",
        "name": "account_2"
    }'
    ```
- Get accounts
    ```
    curl --location 'localhost:8081/accounts'
    ```
- Get accounts/:accountID
    ```
    curl --location 'localhost:8081/accounts/52f3d4fa-87d2-44d1-a181-2dbc567c56f3'
    ```
- Deposit account
    ```
    curl --location 'localhost:8081/accounts/fde7f07a-fd12-493c-83a9-7bec2644c4c2/deposit' \
    --header 'Content-Type: application/json' \
    --data '{
        "amount": 2000
    }'
    ```
- Withdraw account
    ```
    curl --location 'localhost:8081/accounts/fde7f07a-fd12-493c-83a9-7bec2644c4c2/withdraw' \
    --header 'Content-Type: application/json' \
    --data '{
        "amount": 2000
    }'
    ```
- Transfer amount from account to another account
    ```
    curl --location 'localhost:8081/accounts/fde7f07a-fd12-493c-83a9-7bec2644c4c2/transfer' \
    --header 'Content-Type: application/json' \
    --data '{
        "to_account_id": "e66c2ba2-34fd-4801-9650-567e274bf69e",
        "amount": 200
    }'
    ```
- Get transations
    ```
    curl --location 'localhost:8081/accounts/fde7f07a-fd12-493c-83a9-7bec2644c4c2/transactions'
    ```