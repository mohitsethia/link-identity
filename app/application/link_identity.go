package application

import (
	"context"
	"database/sql"

	"github.com/link-identity/app/domain"
	"github.com/link-identity/app/infrastructure/repository"

	"github.com/pkg/errors"
)

type LinkIdentityService interface {
	Identify(ctx context.Context, email, phone string) ([]*domain.Contact, error)
}

type service struct {
	repo repository.ContactRepository
}

func NewService(contactRepo repository.ContactRepository) LinkIdentityService {
	return &service{
		repo: contactRepo,
	}
}

func (s *service) Identify(ctx context.Context, email, phone string) ([]*domain.Contact, error) {
	existingContactByEmail, err := s.repo.GetContactByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrapf(err, "[Service][LinkIdentity] error from repo while getting contacts by email")
	}

	existingContactByPhone, err := s.repo.GetContactByPhone(ctx, phone)
	if err != nil {
		return nil, errors.Wrapf(err, "[Service][LinkIdentity] error from repo while getting contacts by phone")
	}

	contact := &domain.Contact{
		Email:            sql.NullString{String: email, Valid: true},
		Phone:            sql.NullString{String: phone, Valid: true},
		LinkedPrecedence: "primary",
	}

	switch {
	case existingContactByEmail != nil && existingContactByPhone != nil:
		{
			contact.LinkedPrecedence = "secondary"
			if existingContactByPhone.ContactId == existingContactByEmail.ContactId {
				contact.LinkedID = existingContactByEmail.ContactId
				if existingContactByEmail.LinkedPrecedence == "secondary" {
					contact.LinkedID = existingContactByEmail.LinkedID
				}
			} else {
				if existingContactByEmail.CreatedAt.After(*existingContactByPhone.CreatedAt) {
					contact.LinkedID = existingContactByPhone.ContactId
					if existingContactByPhone.LinkedPrecedence == "secondary" {
						contact.LinkedID = existingContactByPhone.LinkedID
					}
					if existingContactByEmail.LinkedPrecedence == "primary" || existingContactByEmail.LinkedID != contact.LinkedID {
						existingContactByEmail.LinkedPrecedence = "secondary"
						existingContactByEmail.LinkedID = contact.LinkedID
						existingContactByEmail, err = s.repo.UpdateContact(ctx, existingContactByEmail)
						if err != nil {
							return nil, errors.Wrapf(err, "[Service][LinkIdentity] error while updating contact")
						}
					}
				} else {
					contact.LinkedID = existingContactByEmail.ContactId
					if existingContactByEmail.LinkedPrecedence == "secondary" {
						contact.LinkedID = existingContactByEmail.LinkedID
					}
					if existingContactByPhone.LinkedPrecedence == "primary" || existingContactByPhone.LinkedID != contact.LinkedID {
						existingContactByPhone.LinkedPrecedence = "secondary"
						existingContactByPhone.LinkedID = contact.LinkedID
						existingContactByPhone, err = s.repo.UpdateContact(ctx, existingContactByPhone)
						if err != nil {
							return nil, errors.Wrapf(err, "[Service][LinkIdentity] error while updating contact")
						}
					}
				}
			}
			secondaryContacts, err := s.repo.GetAllSecondaryContacts(ctx, contact.LinkedID)
			if err != nil {
				return nil, errors.Wrapf(err, "[Service][LinkIdentity] error while getting secondary contacts")
			}
			return secondaryContacts, nil
		}
	case existingContactByEmail != nil:
		{
			contact.LinkedPrecedence = "secondary"
			contact.LinkedID = existingContactByEmail.ContactId
			if existingContactByEmail.LinkedPrecedence == "secondary" {
				contact.LinkedID = existingContactByEmail.LinkedID
			}
		}
	case existingContactByPhone != nil:
		{
			contact.LinkedPrecedence = "secondary"
			contact.LinkedID = existingContactByPhone.ContactId
			if existingContactByPhone.LinkedPrecedence == "secondary" {
				contact.LinkedID = existingContactByPhone.LinkedID
			}
		}
	}

	contact, err = s.repo.CreateContact(ctx, contact)
	if err != nil {
		return nil, errors.Wrapf(err, "[Service][LinkIdentity] error while creating contact")
	}

	if contact.LinkedPrecedence == "primary" {
		return []*domain.Contact{contact}, nil
	}

	secondaryContacts, err := s.repo.GetAllSecondaryContacts(ctx, contact.LinkedID)
	if err != nil {
		return nil, errors.Wrapf(err, "[Service][LinkIdentity] error while getting secondary contacts")
	}

	return secondaryContacts, nil
}
