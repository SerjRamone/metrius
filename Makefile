DSN = 'postgres://postgres:postgres@localhost:5432/metrius?sslmode=disable'
DATE = $(shell date +'%Y/%m/%d %H:%M:%S')
COMMIT = $(shell git rev-parse --short HEAD)
	
# f.e. make build VERSION=0.0.1
build: build-agent build-server

# f.e. make build-server VERSION=0.0.1
build-server:
	go build \
		-ldflags "-X main.buildVersion=$(VERSION) -X 'main.buildDate=$(DATE)' -X 'main.buildCommit=$(COMMIT)'" \
		-o cmd/server/server \
		cmd/server/*.go

# f.e. make build-agent VERSION=0.0.1
build-agent:
	go build \
		-ldflags "-X main.buildVersion=$(VERSION) -X 'main.buildDate=$(DATE)' -X 'main.buildCommit=$(COMMIT)'" \
		-o cmd/agent/agent \
		cmd/agent/*.go

build-linter:
	go build -o staticlint cmd/staticlint/*.go

run-server: build-server
	./cmd/server/server -a="localhost:8080" -i=0 -d=$(DSN) -k=testkey

run-agent: build-agent
	./cmd/agent/agent -a="localhost:8080" -r=10 -p=2 -k=testkey -l=2

stattest:
	go vet -vettool=statictest ./...
	./staticlint ./...

test:
	go test -v -race ./...
	
autotests: build autotest_up_to24
	
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

autotests10: autotests9
	./metricstest -test.v -test.run=^TestIteration10[AB]$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8080" -database-dsn=$(DSN)

autotests11: autotests10
	./metricstest -test.v -test.run=^TestIteration11$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008" -database-dsn=$(DSN)

autotests12: autotests11
	./metricstest -test.v -test.run=^TestIteration12$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008" -database-dsn=$(DSN)
	
autotests13: autotests12
	./metricstest -test.v -test.run=^TestIteration13$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008" -database-dsn=$(DSN)

autotests14: autotests13
	./metricstest -test.v -test.run=^TestIteration14$$ -source-path=. -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -server-port="8008" -database-dsn=$(DSN) -key="testkey"

autotest_up_to24: autotests14
	 go test -v -race ./...
