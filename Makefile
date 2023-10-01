build: build-agent build-server

build-server:
	go build -o cmd/server/server cmd/server/*.go

build-agent:
	go build -o cmd/agent/agent cmd/agent/*.go

run-server: build-server
	./cmd/server/server -a="localhost:8080" -i=0

run-agent: build-agent
	./cmd/agent/agent -a="localhost:8080" -r=10 -p=2

stattest:
	go vet -vettool=statictest ./...
	
autotests: build autotests9
	
autotests1:
	./metricstest -test.v -test.run=^TestIteration1$$ -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

autotests2: autotests1
	./metricstest -test.v -test.run=^TestIteration2[AB]*$$ -source-path=. -agent-binary-path=cmd/agent/agent

autotests3: autotests2
	./metricstest -test.v -test.run=^TestIteration3[AB]*$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server

autotests4: autotests3 
	./metricstest -test.v -test.run=^TestIteration4$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008"

autotests5: autotests4
	./metricstest -test.v -test.run=^TestIteration5$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008"

autotests6: autotests5
	./metricstest -test.v -test.run=^TestIteration6$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008" 

autotests7: autotests6
	./metricstest -test.v -test.run=^TestIteration7$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008"

autotests8: autotests7
	./metricstest -test.v -test.run=^TestIteration8$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008"

autotests9: autotests8
	./metricstest -test.v -test.run=^TestIteration9$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008" -file-storage-path=/tmp/metrics-tests-db.json
