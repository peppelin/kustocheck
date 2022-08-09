set -eu
set -o pipefail

 # Get tools for converting Go's test reports
go get github.com/jstemmer/go-junit-report
go get github.com/axw/gocov/gocov
go get github.com/AlekSi/gocov-xml

go install github.com/jstemmer/go-junit-report
go install github.com/axw/gocov/gocov
go install github.com/AlekSi/gocov-xml

# Run Go tests and turn output into JUnit test result format
go test $(go list ./... | grep -v /cmd/ | grep -v /ci/) -v -coverprofile=coverage.cov -covermode count | $GOPATH/bin/go-junit-report > report.xml

# Convert coverage file into XML
$GOPATH/bin/gocov convert coverage.cov > coverage.json
$GOPATH/bin/gocov-xml < coverage.json > coverage.xml
