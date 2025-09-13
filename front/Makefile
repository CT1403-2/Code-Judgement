dependencies:
	npm install --force

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
		echo "protoc installed: $$(protoc --version)"; \
	fi

generate:
	protoc -I=../proto services.proto \
    --js_out=import_style=commonjs,binary:./src/app/services/proto \
    --grpc-web_out=import_style=typescript,mode=grpcwebtext:./src/app/services/proto

build: generate
	ng build --output-path ../build --configuration production
