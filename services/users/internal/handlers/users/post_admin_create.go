package users

// import (
// 	"context"
// 	"log"

// )

// func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
// 	id, err := i.noteService.Create(ctx, converter.ToNoteInfoFromDesc(req.GetInfo()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("inserted note with id: %d", id)

// 	return &desc.CreateResponse{
// 		Id: id,
// 	}, nil
// }
