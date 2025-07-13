#!bin/bash

name=$1
camelName=$(echo "$name" | sed -r 's/(.)_+(.)/\1\U\2/g;s/^[a-z]/\U&/')

touch internal/server/router/$name.go
echo "" > internal/server/router/$name.go

echo "package router" >> internal/server/router/$name.go
echo "" >> internal/server/router/$name.go

echo "import (" >> internal/server/router/$name.go
echo "	\"net/http\"" >> internal/server/router/$name.go
echo ")" >> internal/server/router/$name.go
echo "" >> internal/server/router/$name.go

echo "func (api *router) $camelName(w http.ResponseWriter, r *http.Request) {" >> internal/server/router/$name.go
echo "" >> internal/server/router/$name.go
echo "}" >> internal/server/router/$name.go