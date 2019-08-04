# Coupons service

## Running the project

Requirements:
```sh
docker --version # Docker version 19.03.1-ce, build 74b1e89e8a
```

```sh
# Start the database
docker-compose up -d db

# Create the database schema, insert some data
docker-compose up migrate

# Run the server
docker-compose up -d coupons

# You can watch the logs in another terminal if you want:
docker logs -f coupons_coupons_1
```

Ready to go! You can start using the service:
```sh
curl '127.0.0.1:4000/coupons' | python -mjson.tool
```

## Building manually

Requirements:

You will need a running Postgres database, and the following:
```sh
go version
# go version go1.12.1 linux/amd64

go get
go build

export DATABASE_HOST=[host]
export DATABASE_PORT=[port]
export DATABASE_USERNAME=[username]
export DATABASE_PASSWORD=[password]
export DATABASE_DATABASE=[database]

# create the db structure, insert some data
go run migration/migrate.go migration/creation.sql migration/seed.sql

./coupons
```

## Filtering

Available filters:
- min: minimum discount (min >= 0)
- max: maximum discount (max >= 0, 0 = disabled)
- brand: brand the coupons can be used with (case insensitive, wildcard)
- take, skip: pagination (1 <= take <= 20, skip >= 0)
- exp: expiration date (2019-09-16T11:45:26.371Z format, display only coupons  
available after this date)

Omit any filter you don't want to apply.

Example usage:
```sh
# Filter on brands: try Tesco, Asda, LIDL
curl '127.0.0.1:4000/coupons?brand=Tesco' | python -mjson.tool
curl '127.0.0.1:4000/coupons?brand=Tesco&min=25' | python -mjson.tool

# paginate data
curl '127.0.0.1:4000/coupons?take=5' | python -mjson.tool
curl '127.0.0.1:4000/coupons?take=5&skip=5' | python -mjson.tool

# expires after a certain date
curl '127.0.0.1:4000/coupons?exp=2019-09-16T11:45:26.371Z' | python -mjson.tool

# all filters at once: 
# first 5 Tesco coupons that are worth between £15.99 and £23.55
# and expire the 1st of January 2019 or later
curl '127.0.0.1:4000/coupons?take=5&skip=0&brand=Tesco&min=15.99&max=23.55&exp=2019-01-01T11:45:26.371Z' \
  | python -mjson.tool
```

## Update data
Update an existing coupon:
```sh
curl \
  -X PUT \
  -H "Content-Type: application/json" \
  -d '{ "id": 1, "brand": "Tesco", "value": 34, "name": "Save £34 at Tesco", "expiryUtc": "2019-09-16T11:45:26.371Z" }' \
  '127.0.0.1:4000/coupons' \
  | python -mjson.tool
```

Create a new coupon:
```sh
curl \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{ "brand": "Mark & Spencers", "value": 0.01, "name": "Save a penny at Mark & Spencers", "expiryUtc": "2019-09-16T11:45:26.371Z" }' \
  '127.0.0.1:4000/coupons' \
  | python -mjson.tool
```

## Precisions

### Why PostgreSQL?

The data that needs to be stored to serve coupons is following a relational schema,
know in advance, and it is unlikely to vary.  
Unless we are dealing with millions of coupons, a PostgreSQL database
is performant enough to retrieve the coupons effectively.  
When the application scales, we could add an ElasticSearch engine to make
the search faster and more easily configurable.

### Why run the migrations in Docker Compose?

For simplicity of use and not to expose the PostgreSQL database outside of the
Docker Compose network.
At the moment, you need to rebuild the image every time you change a migration
file, but this could be improved by using `docker cp` or Docker volumes for example.

### Why store Brands in a separate SQL table?

This will simplify search, avoid duplicates, make renaming easy, can help eliminating typos... 
