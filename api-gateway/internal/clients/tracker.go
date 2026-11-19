package clients

// import (
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"

// 	trackerv1 "github.com/monstrong/gracker2/proto/gen/go/tracker/v1"
// )

// type TrackerClient struct {
// 	conn   *grpc.ClientConn
// 	client trackerv1.TorrentServiceClient
// }

// func (c *TrackerClient) Close() error {
// 	return c.conn.Close()
// }

// func NewTrackerClient(addr string) (*TorrentClient, error) {

// 	conn, err := grpc.NewClient(
// 		addr,
// 		grpc.WithTransportCredentials(
// 			insecure.NewCredentials(),
// 		),
// 		grpc.WithUnaryInterceptor(clientinterceptors.TraceIDClientInterceptor()),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client := trackerv1.NewTorrentServiceClient(conn)

// 	return &TorrentClient{conn: conn, client: client}, nil
// }

