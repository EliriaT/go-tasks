# go-tasks

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