# For simplicity for users just cloning and running this package, I stored the
# configuration in this file rather than using environment variables.

# If you need to run the test database on another port, update the value here.

export DATABASE_HOST=127.0.0.1
export DATABASE_PORT=5432
export DATABASE_USERNAME=perkbox
export DATABASE_PASSWORD=perkbox
export DATABASE_DATABASE=coupons

export POSTGRES_USER=$DATABASE_USERNAME
export POSTGRES_PASSWORD=$DATABASE_PASSWORD
export POSTGRES_DB=$DATABASE_DATABASE

echo "Building"
go get
go build

echo "Running the db"

# If the image is present don't write anything to the console
if [[ -z "$(docker images -q postgres:12-alpine 2> /dev/null)" ]]; then
    docker pull postgres:12-alpine
fi

docker run \
    --rm --name dbtest \
    -d -p $DATABASE_PORT:5432 \
    -e POSTGRES_DB \
    -e POSTGRES_USER \
    -e POSTGRES_PASSWORD \
    postgres:12-alpine

echo "Waiting for the container to start"
sleep 5

# Creating the db structure with every test to avoid race conditions between tests
echo "Testing: database client"
go run ./migration/migrate.go ./migration/creation.sql 2>/dev/null
go test ./coupon
echo

echo "Testing: coupon client"
go run ./migration/migrate.go ./migration/creation.sql 2>/dev/null
go test ./coupon
echo

echo "Testing: server"
go run ./migration/migrate.go ./migration/creation.sql  2>/dev/null
go test .
echo

echo "Cleaning up"
docker kill dbtest > /dev/null