package grpc_reqAuth

type Options struct {
	HashSecret      string
	MetadataKeyList []string
	MetadataHashKey string
}
