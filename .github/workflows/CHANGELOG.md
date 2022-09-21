# k8f version <>
## Changes
## Bugfix
## Braking changes


  env GOOS=windows GOARCH=amd64 go build . -o k8f.exe
  VERSION="$(git describe --tags --always --abbrev=0 --match='[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
