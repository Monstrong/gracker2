package clients

import (
	"context"

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
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}

	client := torrentv1.NewTorrentServiceClient(conn)

	return &TorrentClient{conn: conn, client: client}, nil
}

func (c *TorrentClient) GetTorrent(
	ctx context.Context,
	req *torrentv1.GetTorrentRequest,
) (*torrentv1.GetTorrentResponse, error) {
	return c.client.GetTorrent(ctx, req)
}

func (c *TorrentClient) CreateTorrent(
	ctx context.Context,
	req *torrentv1.CreateTorrentRequest,
) (*torrentv1.CreateTorrentResponse, error) {
	return c.client.CreateTorrent(ctx, req)
}

func (c *TorrentClient) DeleteTorrent(
	ctx context.Context,
	req *torrentv1.DeleteTorrentRequest,
) (*torrentv1.DeleteTorrentResponse, error) {
	return c.client.DeleteTorrent(ctx, req)
}

func (c *TorrentClient) ListTorrents(
	ctx context.Context,
	req *torrentv1.ListTorrentsRequest,
) (*torrentv1.ListTorrentsResponse, error) {
	return c.client.ListTorrents(ctx, req)
}

