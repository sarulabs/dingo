#/bin/sh

set -e

testsDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo ">>> GENERATING CODE ..."
go run "${testsDir}/app/main.go" "${testsDir}/app/generated"

echo ">>> RUNNING TESTS ..."
go test -v "${testsDir}/app/tests"

echo ">>> REMOVING GENERATED CODE from ${testsDir}/app/generated ..."
rm -rf "${testsDir}/app/generated"
