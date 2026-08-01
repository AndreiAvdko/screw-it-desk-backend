package handlers

// ------------ Пример создания новой записи в БД ------------
// func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
// 	// Делаем запрос на вставку записи в таблицу note
// 	builderInsert := sq.Insert("note").
// 		PlaceholderFormat(sq.Dollar).
// 		Columns("title", "body").
// 		Values(gofakeit.City(), gofakeit.Address().Street).
// 		Suffix("RETURNING id")

// 	query, args, err := builderInsert.ToSql()
// 	if err != nil {
// 		log.Fatalf("failed to build query: %v", err)
// 	}

// 	var noteID int64
// 	err = s.pool.QueryRow(ctx, query, args...).Scan(&noteID)
// 	if err != nil {
// 		log.Fatalf("failed to insert note: %v", err)
// 	}

// 	log.Printf("inserted note with id: %d", noteID)

// 	return &desc.CreateResponse{
// 		Id: noteID,
// 	}, nil
// }

// /
//
