package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"

	"github.com/NdoleStudio/http-sms-manager/pkg/repositories"
	"github.com/palantir/stacktrace"

	"github.com/NdoleStudio/http-sms-manager/pkg/entities"
	"github.com/NdoleStudio/http-sms-manager/pkg/telemetry"
)

// PhoneService is handles phone requests
type PhoneService struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	repository repositories.PhoneRepository
}

// NewPhoneService creates a new PhoneService
func NewPhoneService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	repository repositories.PhoneRepository,
) (s *PhoneService) {
	return &PhoneService{
		logger:     logger.WithService(fmt.Sprintf("%T", s)),
		tracer:     tracer,
		repository: repository,
	}
}

// Index fetches the heartbeats for a phone number
func (service *PhoneService) Index(ctx context.Context, authUser entities.AuthUser, params repositories.IndexParams) (*[]entities.Phone, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	phones, err := service.repository.Index(ctx, authUser.ID, params)
	if err != nil {
		msg := fmt.Sprintf("could not fetch phones with parms [%+#v]", params)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	ctxLogger.Info(fmt.Sprintf("fetched [%d] phones with prams [%+#v]", len(*phones), params))
	return phones, nil
}

// PhoneUpsertParams are parameters for creating a new entities.Phone
type PhoneUpsertParams struct {
	PhoneNumber phonenumbers.PhoneNumber
	FcmToken    string
	UserID      entities.UserID
}

// Upsert a new entities.Phone
func (service *PhoneService) Upsert(ctx context.Context, params PhoneUpsertParams) (*entities.Phone, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	phone, err := service.repository.Load(ctx, params.UserID, phonenumbers.Format(&params.PhoneNumber, phonenumbers.E164))
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		return service.createPhone(ctx, params)
	}

	phone.FcmToken = &params.FcmToken
	if err = service.repository.Save(ctx, phone); err != nil {
		msg := fmt.Sprintf("cannot update phone with id [%s] and number [%s]", phone.ID, phone.PhoneNumber)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	ctxLogger.Info(fmt.Sprintf("phone saved with id [%s] in the userRepository", phone.ID))
	return phone, nil
}

func (service *PhoneService) createPhone(ctx context.Context, params PhoneUpsertParams) (*entities.Phone, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	phone := &entities.Phone{
		ID:          uuid.New(),
		UserID:      params.UserID,
		FcmToken:    &params.FcmToken,
		PhoneNumber: phonenumbers.Format(&params.PhoneNumber, phonenumbers.E164),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := service.repository.Save(ctx, phone); err != nil {
		msg := fmt.Sprintf("cannot create phone with id [%s] and number [%s]", phone.ID, phone.PhoneNumber)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return phone, nil
}
