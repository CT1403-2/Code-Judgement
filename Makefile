dependencies:
	go mod download

	@if command -v mockery >/dev/null 2>&1; then \
		echo "">/dev/null; \
	else \
		go install github.com/vektra/mockery/v2@v2.52.4; \
	fi

	@if command -v protoc >/dev/null 2>&1; then \
		echo "">/dev/null; \
	else \
		if [ -f /etc/arch-release ]; then \
			sudo pacman -Sy --noconfirm protobuf; \
		elif [ -f /etc/debian_version ] || grep -q "debian\|ubuntu" /etc/os-release 2>/dev/null; then \
			apt update && apt install -y protobuf-compiler; \
		else \
			echo "Unsupported distribution. Please install protoc manually."; \
			exit 1; \
		fi; \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
		echo "protoc installed: $$(protoc --version)"; \
	fi

	@if command -v migrate >/dev/null 2>&1; then \
		echo "">/dev/null; \
	else \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3; \
	fi

database:
	PGPASSWORD=$(password) psql -U $(username) -f ./manager/internal/database/migrations/create_table.sql; \
    migrate -database postgres://$(username):$(password)@localhost:5432/judge_db?sslmode=disable -path ./manager/internal/database/migrations up

generate:
	protoc -I ./proto/ services.proto --go_out=./ --go-grpc_out=./
	go generate ./judge/...

build: generate
	go build -gcflags="all=-N -l" -o ./build/bin/$(target) ./$(target)/cmd/