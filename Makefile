build: build-agent build-server

build-server:
	go build -o cmd/server/server cmd/server/*.go

build-agent:
	go build -o cmd/agent/agent cmd/agent/*.go

autotests:
	./metricstest
	
autotests1:
	./metricstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent
