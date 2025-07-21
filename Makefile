

route:
	bash scripts/route.sh ${name}

ws: 
	bash scripts/ws.sh websocket_${name}

test:
	DATABASE_URL=${DATABASE_URL} go test ./...
