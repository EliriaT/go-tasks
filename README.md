# go-tasks

To run the tasks:
`docker compose build`
`docker compose up`

Running `docker compose up` will run all the migrations, seeding and sql queries.

Database available at: user:password@tcp(sources_db:3306)/sources

Useful commands:

To start the mariadb client inside the sources_db container:
`docker exec -it sources_db mariadb -uusername -ppassword`

To run the migrations:
`make migrateup`

To rollback the migrations:
`make migratedown`

To build the go image:
`docker compose build`

To run the small environment:
`docker compose up --remove-orphans`

To create several migration files:
`make create_migration create_sources_associated_campaigns_table create_sources_table`