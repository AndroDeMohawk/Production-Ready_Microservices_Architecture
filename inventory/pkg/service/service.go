package service

import (
	"cmp"
	"context"
	"slices"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
)

// Part представляет деталь космического корабля
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
}

// server реализует gRPC сервис
type server struct {
	inventoryv1.UnimplementedInventoryServiceServer
	parts map[uuid.UUID]Part
}

// NewServer создаёт сервер с предзагруженными seed-данными
func NewServer() *server {
	now := timestamppb.Now()

	return &server{
		parts: map[uuid.UUID]Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 5,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 8,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 3,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      inventoryv1.PartType_PART_TYPE_SHIELD,
				StockQuantity: 6,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      inventoryv1.PartType_PART_TYPE_WEAPON,
				StockQuantity: 7,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Экспериментальный корпус (нет на складе)",
				Price:         2000000, // 20000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

// GetPart возвращает деталь по UUID
func (s *server) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	id := req.GetUuid()
	if id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "uuid не может быть пустым")
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid некорректный")
	}

	part, exists := s.parts[idParsed]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "деталь не найдена")
	}
	finalPart := &inventoryv1.Part{
		Uuid:          id,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
		PartType:      part.PartType,
	}
	return &inventoryv1.GetPartResponse{
		Part: finalPart,
	}, nil

}

// ListParts возвращает список деталей с опциональной фильтрацией по типу
func (s *server) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {

	uuids := req.GetUuids()

	// 1. Если передан список uuids → найти детали по UUID (сохраняя порядок запроса)
	if len(uuids) > 0 {
		parts := make([]*inventoryv1.Part, 0, len(uuids))

		for _, idStr := range uuids {
			idStr := uuid.MustParse(idStr)
			// Ищем деталь в нашей базе данных s.parts
			val, exists := s.parts[idStr] // Предполагается, что s.parts это map[string]ВашаСтруктура
			if !exists {
				// Если хоть один UUID не найден → NOT_FOUND
				return nil, status.Error(codes.NotFound, "part not found")
			}

			// Маппим в protobuf-структуру (порядок сохраняется за счет цикла по запросу)
			parts = append(parts, &inventoryv1.Part{
				Uuid:          val.UUID,
				Name:          val.Name,
				Description:   val.Description,
				PartType:      val.PartType,
				Price:         val.Price,
				StockQuantity: val.StockQuantity,
				CreatedAt:     val.CreatedAt,
			})
		}

		// Для UUID сортировка НЕ нужна (по пункту 4)
		return &inventoryv1.ListPartsResponse{Parts: parts}, nil
	}

	// 2. Иначе если part_type == UNSPECIFIED (или 0) → вернуть все детали
	// В protobuf enum по умолчанию равен 0 (обычно это _UNSPECIFIED)
	targetType := req.GetPartType()
	if targetType == inventoryv1.PartType_PART_TYPE_UNSPECIFIED { // Замените на точное имя вашего enum по умолчанию

		parts := s.makeResponseSlice(SortByName, targetType, true)
		return parts, nil
	}

	// 3 & 4. Иначе → фильтровать по типу и отсортировать ПО ИМЕНИ
	parts := s.makeResponseSlice(SortByName, targetType, false)
	return parts, nil
}

type SortBy string

const (
	SortByPartType SortBy = "type"
	SortByName     SortBy = "name"
)

// Добавим в метод параметры фильтрации, чтобы делать всё за один проход
func (s *server) makeResponseSlice(sortBy SortBy, filterType inventoryv1.PartType, returnAll bool) *inventoryv1.ListPartsResponse {
	parts := make([]*inventoryv1.Part, 0, len(s.parts))

	for _, val := range s.parts {
		// Если не нужно возвращать всё и тип не совпадает — пропускаем
		if !returnAll && val.PartType != filterType {
			continue
		}

		parts = append(parts, &inventoryv1.Part{
			Uuid:          val.UUID,
			Name:          val.Name,
			Description:   val.Description,
			PartType:      val.PartType,
			Price:         val.Price,
			StockQuantity: val.StockQuantity,
			CreatedAt:     val.CreatedAt,
		})
	}

	// Сортировка
	slices.SortFunc(parts, func(a, b *inventoryv1.Part) int {
		switch sortBy {
		case SortByName:
			return cmp.Compare(a.Name, b.Name)
		case SortByPartType:
			fallthrough
		default:
			return cmp.Compare(a.PartType, b.PartType)
		}
	})

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}
}
