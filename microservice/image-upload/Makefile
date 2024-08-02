# Define the paths for proto files and output
PROTO_PATH = proto
PROTO_FILES = $(PROTO_PATH)/*.proto
OUT_DIR = .

# Define the command for protoc
PROTOC = protoc

# Targets
all: pb

# Rule to generate Go code from proto files
pb:
	$(PROTOC) --proto_path=$(PROTO_PATH) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_FILES)


