

route:
	bash scripts/route.sh ${name}

test:
	DATABASE_URL=${DATABASE_URL} go test ./...