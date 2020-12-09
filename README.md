# go_server

protoc -I=/Users/apple/Projects/go_server/goprotobuf --go_out=/Users/apple/Projects/go_server /Users/apple/Projects/go_server/goprotobuf/addressbook.proto 

sprotogen --go_out=addressbook.go --package=ab addressbook.sp