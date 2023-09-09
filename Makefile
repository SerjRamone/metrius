build: build-agent build-server

build-server:
	go build -o cmd/server/server cmd/server/*.go

build-agent:
	go build -o cmd/agent/agent cmd/agent/*.go

run-server:
	go run cmd/server/main.go

run-agent:
	go run cmd/agent/main.go

stattest:
	go vet -vettool=statictest ./...
	
autotests: build autotests3
	./metricstest
	
autotests1:
	./metricstest -test.v -test.run=^TestIteration1$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

autotests2: autotests1
	./metricstest -test.v -test.run=^TestIteration2[AB]*$$ -source-path=. -agent-binary-path=cmd/agent/agent

autotests3: autotests2
	./metricstest -test.v -test.run=^TestIteration3[AB]*$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

