package clients

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	torrentv1 "github.com/monstrong/gracker2/proto/gen/go/torrent/v1"
)

type TorrentClient struct {
	conn   *grpc.ClientConn
	client torrentv1.TorrentServiceClient
}

func (c *TorrentClient) Close() error {
	return c.conn.Close()
}

func NewTorrentClient(addr string) (*TorrentClient, error) {

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(), // TODO: add tls
		),
	)
	if err != nil {
		return nil, err
	}

	client := torrentv1.NewTorrentServiceClient(conn)

	return &TorrentClient{conn: conn, client: client}, nil
}
