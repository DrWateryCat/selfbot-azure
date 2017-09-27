package learning

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

var address string

func SendMessageToLearning(m string, length int) (ret string, err error) {
	address = "localhost:69420"
	conn, err := grpc.Dial(address)
	if err != nil {
		return
	}

	defer conn.Close()
	client := NewLearningClient(conn)
	req := LearningRequest{Input: m, Length: int32(length)}
	res, err := client.Train(context.Background(), &req)
	if err != nil {
		return
	}

	ret = res.GetResponse()
	return
}
