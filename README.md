# User Photo Api
REST API that handles CRUD for users' profile picture, made using Go, Gin(Go Framework), Gorm(Go ORM), MySQL/MariaDB, and JWT.
 [Thunder Client VS Code Extension](https://www.thunderclient.com/) can be used to access the API. Before accessing API using Thunder CLient, we must create a new local environment in `Env` tab of Thunder Client. Then we create a new variable named `token` which we must fill with token that we get after successful login. `token` is needed to create, edit, and delete profile picture. This token must be used in `Auth: Bearer -> Bearer Token` in Thunder Client.

## How to run the code
1. Create a new database for this API. No need to manually create other tables in the new database because the tables will be created automatically after executing `go run .`(Step 5).
2. Copy and rename file `example.env` into `.env`.
3. Open file `.env` and change `PORT`, `SECRET`, `DATABASE_URL`, and `API_PREFIX` into the appropriate port, secret for generating JWT token, database url, and api prefix.
4. Open terminal, go into root directory of this code, and run `go mod tidy`.
5. Then run `go run .`.
6. Press `Ctrl + C` to terminate the API.