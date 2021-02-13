# Open Social

This is a port of https://github.com/reecerussell/social-media.

## Stack

-   Golang
-   ReactJS/Node
-   SQL Server
-   Google Cloud Platform

## Development

To build and run this locally, you'll need both Docker and NodeJS installed.

### Backend

By running docker-compose, all of the services will be spun up and exposed via the ingress, on port 80.

Before running the docker-compose command, ensure there is a `CONNECTION_STRING` environment variable, set with a connection string to a SQL Server instance. An example of a connection string is `sqlserver://<username>:<password>@<host>?database=<database>`.

```
$ docker-compose -f docker-compose.dev.yaml up
```

### Frontend

NPM is used to start the frontend locally, exposing it on port 3000.

```
cd ui
npm start
```

## Google Cloud Platform

Currently, this is built to use GCP storage buckets as a filestore for the media service. The media service uses an interface which can be implemented for other services, such as AWS or Azure.
