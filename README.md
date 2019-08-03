# Coupons service

## Running the project

sh
```
# Start the database
docker-compose up -d db

# Create the database schema, insert some data
docker-compose up migrate

# Run the server
docker-compose up -d coupons

# Ready to go! You can start using the service - try Tesco, Asda, LIDL
curl 127.0.0.1:4000/coupons?brand=Tesco
curl 127.0.0.1:4000/coupons?brand=Tesco&minValue=25

# paginate data
curl 127.0.0.1:4000/coupons?take=5
curl 127.0.0.1:4000/coupons?take=5&skip=5

# expires after a certain date
curl 127.0.0.1:4000/coupons?expiresAfter=2019-09-16T11:45:26.371Z

# update an existing coupon
curl \
  -X PUT \
  -H "Accept: application/json" \
  -d '{ id: 1, brand: "Tesco", value: 34, name: "Save £34 at Tesco", expiryUtc: "2019-09-16T11:45:26.371Z" }' \
  127.0.0.1:4000/coupons/1

# create a new coupon
curl \
  -X POST \
  -H "Accept: application/json" \
  -d '{ brand: "Sainsbury", value: 0.01, name: "Save a penny at Sainsbury", expiryUtc: "2019-09-16T11:45:26.371Z" }' \
  127.0.0.1:4000/coupons
```

## Precisions

Why PostgreSQL?

The data that needs to be stored to serve coupons is following a relational schema,
know in advance, and it is unlikely to vary.  
Unless we are dealing with millions of coupons, a PostgreSQL database
is performant enough to retrieve the coupons effectively.  
When the application scales, we could add an ElasticSearch engine to make
the search faster and more easily configurable.

Why run the migrations in Docker Compose ?

For simplicity of use and not to expose the PostgreSQL database outside of the
Docker Compose network.
At the moment, you need to rebuild the image every time you change a migration
file, but this could be improved by using volumes.