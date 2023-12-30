package choice

import (
	context "context"
	"log"
	"net"

	"github.com/tthung1997/buddy/core/random"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

var port = ":8080"

type ChoiceServer struct {
	repository random.IChoiceListRepository
}

func NewChoiceServer(repository random.IChoiceListRepository) *ChoiceServer {
	return &ChoiceServer{repository: repository}
}

func (s *ChoiceServer) Run() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listen.Close()

	grpcServer := grpc.NewServer()
	RegisterChoiceServiceServer(grpcServer, s)

	log.Printf("Starting gRPC server on port %s", port)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Printf("Stopping gRPC server on port %s", port)
}

func (s *ChoiceServer) GetChoiceListById(ctx context.Context, request *GetByIdRequest) (*ChoiceList, error) {
	log.Printf("GetChoiceListById: %v", request.Id)

	choiceList, err := s.repository.GetChoiceList(request.Id)
	if err != nil {
		log.Printf("GetChoiceListById: %v", err)

		if err.Error() == "ChoiceList with ID "+request.Id+" not found" {
			return &ChoiceList{}, status.Error(codes.NotFound, err.Error())
		}

		return &ChoiceList{}, status.Error(codes.Internal, err.Error())
	}

	log.Printf("GetChoiceListById: %v", choiceList)

	updatedDateTime := timestamppb.New(choiceList.UpdatedDateTime)

	choices := make([]*Choice, len(choiceList.Choices))
	for i, c := range choiceList.Choices {
		choices[i] = &Choice{
			Id:              c.Id,
			Value:           c.Value,
			Weight:          c.Weight,
			Color:           c.Color,
			UpdatedDateTime: timestamppb.New(c.UpdatedDateTime),
		}
	}

	return &ChoiceList{
		Id:              choiceList.Id,
		Choices:         choices,
		UpdatedDateTime: updatedDateTime,
	}, nil
}

func (s *ChoiceServer) UpsertChoiceList(ctx context.Context, choiceList *ChoiceList) (*UpsertResponse, error) {
	log.Printf("UpsertChoiceList: %v", choiceList)

	choices := make([]random.Choice, len(choiceList.Choices))
	for i, c := range choiceList.Choices {
		choices[i] = random.Choice{
			Id:              c.Id,
			Value:           c.Value,
			Weight:          c.Weight,
			Color:           c.Color,
			UpdatedDateTime: c.UpdatedDateTime.AsTime(),
		}
	}

	err := s.repository.CreateOrUpdateChoiceList(random.ChoiceList{
		Id:              choiceList.Id,
		Choices:         choices,
		UpdatedDateTime: choiceList.UpdatedDateTime.AsTime(),
	})
	if err != nil {
		log.Printf("UpsertChoiceList: %v", err)
		return &UpsertResponse{Success: false, Error: err.Error()}, status.Error(codes.Internal, err.Error())
	}

	return &UpsertResponse{Success: true}, nil
}
