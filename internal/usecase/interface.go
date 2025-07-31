package usecase

// Repository interfaces define contracts for data access in the use case layer.
// These interfaces follow Clean Architecture principles by depending only on
// the entity layer and providing abstractions for outer layers to implement.

// Example repository interface structure:
// type EntityRepository interface {
//     Create(ctx context.Context, entity *entity.Entity) error
//     GetByID(ctx context.Context, id string) (*entity.Entity, error)
//     Update(ctx context.Context, entity *entity.Entity) error
//     Delete(ctx context.Context, id string) error
//     List(ctx context.Context, filter EntityFilter) ([]*entity.Entity, error)
// }

// Example service interface structure:
// type ExternalService interface {
//     ProcessData(ctx context.Context, data interface{}) error
//     ValidateRequest(ctx context.Context, request interface{}) error
// }