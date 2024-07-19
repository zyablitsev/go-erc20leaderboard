BUILDCMD=CGO_ENABLED=0 go build -mod=vendor

# load config parameters
ifneq (,$(wildcard ./config))
	include config
endif

v_rpc_url=$(or ${RPC_URL},${rpc_url})
v_dial_timeout=$(or ${DIAL_TIMEOUT},${dial_timeout})
v_tls_handshake_timeout=$(or ${TLS_HANDSHAKE_TIMEOUT},${tls_handshake_timeout})
v_http_client_timeout=$(or ${HTTP_CLIENT_TIMEOUT},${http_client_timeout})
v_depth=$(or ${DEPTH},${depth})

v_env=RPC_URL="${v_rpc_url}" \
	DIAL_TIMEOUT="${v_dial_timeout}" \
	TLS_HANDSHAKE_TIMEOUT="${v_tls_handshake_timeout}" \
	HTTP_CLIENT_TIMEOUT="${v_http_client_timeout}" \
	DEPTH="${v_depth}"

build-dev:
	${BUILDCMD} -o ./bin/leaderboard_dev cmd/main.go

run-dev:
	${v_env} ./bin/leaderboard_dev

build-amd64:
	GOOS=linux GOARCH=amd64 \
		${BUILDCMD} -o ./bin/leaderboard_amd64 cmd/main.go

run-amd64:
	${v_env} ./bin/leaderboard_amd64

docker-build:
	docker build -t zyablitsev/leaderboard .

docker-run:
	docker run \
		-e "RPC_URL=${v_rpc_url}" \
		-e "DIAL_TIMEOUT=${v_dial_timeout}" \
		-e "TLS_HANDSHAKE_TIMEOUT=${v_tls_handshake_timeout}" \
		-e "HTTP_CLIENT_TIMEOUT=${v_http_client_timeout}" \
		-e "DEPTH=${v_depth}" \
		--network host \
		--rm zyablitsev/leaderboard
