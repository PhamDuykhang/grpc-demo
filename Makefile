 gen_file:
	protoc -I=.proto --go_out=.proto/no_touch proto/example.proto
