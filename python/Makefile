PACKAGE_NAME := import \"wedge-api/

all: requirements pre-build build post-build

requirements:
	python3 -m pip install --upgrade pip
	python3 -m pip install --upgrade protobuf
	python3 -m pip install grpcio-tools

pre-build:
	mkdir -p src
	mkdir -p proto/wedge-api
	cp -r ../protobuf/* proto/wedge-api/

build:
	python3 -m grpc_tools.protoc  --python_out=src/wedge_api  --grpc_python_out=./src/wedge_api -I=proto/wedge-api proto/wedge-api/slx.proto proto/wedge-api/wedge.proto

post-build:
	cp src/wedge_api/slx_pb2.py src/examples
	cp src/wedge_api/wedge_pb2.py src/examples
	cp src/wedge_api/wedge_pb2_grpc.py src/examples

