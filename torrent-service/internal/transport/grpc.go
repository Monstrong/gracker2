package grpc

import (
	"context"

	// при смене версии изменить
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1"
)

type Server struct {
	pb.UnimplementedTorrentServiceServer
}

func (s *Server) GetTorrent(ctx context.Context, in *pb.GetTorrentRequest) (*pb.GetTorrentResponse, error) {
	
}

func (s *Server) ListTorrents(ctx context.Context, in *pb.ListTorrentsRequest) (*pb.ListTorrentsResponse, error) {

}

func (s *Server) CreateTorrent(ctx context.Context, in *pb.CreateTorrentRequest) (*pb.CreateTorrentResponse, error) {

}

func (s *Server) DeleteTorrent(ctx context.Context, in *pb.DeleteTorrentRequest) (*pb.DeleteTorrentResponse, error) {

}
