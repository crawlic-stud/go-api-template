

route:
	bash scripts/route.sh ${name}

db-up:
	export GOOSE_DBSTRING=$(DATABASE_URL) & goose up
	export GOOSE_DBSTRING=$(DATABASE_TEMPLATE_URL) & goose up
