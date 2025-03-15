# Project Swift

## About the Project

The project is written in **Go** using the **Gin** web framework and utilizes **MongoDB** as the database. The application is containerized with **Docker**. It parses data from a Google Spreadsheet into the database and provides the specified endpoints through a **RESTful API**. The project also includes a database integration test using **Testcontainers**. The application structure is based on the **go-blueprint**.

## Getting Started

To run the app, you need to have Docker installed. This project was tested with Docker version 27.4.0, build bde2b89.


## How to run (MakeFile Commands)


Build and run containers
```
make docker-run
```

Shutdown containers
```
make docker-down
```

Run tests
```
make test
```

Fetch data from <i>"Interns_2025_SWIFT_CODES"</i>
```
make fetch-data
```
