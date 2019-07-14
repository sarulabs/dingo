#/bin/sh

set -e

testsDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo ">>> GENERATING CODE..."
cd $testsDir/app
go run main.go

echo ">>> RUNNING TESTS..."
cd $testsDir/app/tests
go test -v .

# echo ">>> REMOVING GENERATED CODE..."
cd $testsDir/app
rm -rf generated_services
echo "directory $testsDir/app/generated_services was removed"

