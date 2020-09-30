gen_grpc:
	protoc --go_out=plugins=grpc:proto proto/*.proto

gen_cert_file:
	openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr

