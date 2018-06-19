#/bin/sh

set -e

testsDir="$(dirname $0)/../dingo/tests/"

# app tests
cd $testsDir/app

echo ">>> GENERATING CODE..."
go run main.go

echo ">>> RUNNING TESTS..."
go test -v .

echo ">>> REMOVING GENERATED CODE..."
rm -rf generated_services
echo "directory $testsDir/app/generated_services was removed"

