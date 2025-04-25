dependencies:
	npm install

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
	protoc --plugin=./node_modules/.bin/protoc-gen-ts_proto.cmd --ts_proto_out=./src/app/services -I ../proto/ services.proto

build: generate
