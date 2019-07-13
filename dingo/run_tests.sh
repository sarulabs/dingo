#/bin/sh

set -e

baseDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
testsDir="$baseDir/tests/"

# app tests
cd $testsDir/app

echo ">>> GENERATING CODE..."
go run main.go

echo ">>> RUNNING TESTS..."
go test -v .

echo ">>> REMOVING GENERATED CODE..."
rm -rf generated_services
echo "directory $testsDir/app/generated_services was removed"

