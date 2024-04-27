package repository

import (
	"context"
	"github.com/pkg/errors"

	"github.com/link-identity/app/domain"
	"github.com/link-identity/app/infrastructure/mysql"
)

type ContactRepository interface {
	GetContactByEmail(ctx context.Context, email string) (*domain.Contact, error)
	GetContactByPhone(ctx context.Context, phone string) (*domain.Contact, error)
	GetAllContacts(ctx context.Context) ([]*domain.Contact, error)
	GetAllSecondaryContacts(ctx context.Context, linkedID uint) ([]*domain.Contact, error)
	GetPrimaryContactFromLinkedID(ctx context.Context, linkedId uint) (*domain.Contact, error)
	CreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error)
}

type contactDBRepo struct {
	db *mysql.DbConn
}

func NewContactRepository(db *mysql.DbConn) ContactRepository {
	return &contactDBRepo{
		db: db,
	}
}

func (r *contactDBRepo) GetContactByEmail(ctx context.Context, email string) (*domain.Contact, error) {
	db := r.db.GormConn
	contact := &domain.Contact{}
	rows := db.WithContext(ctx).Where("email = ?", email).Find(contact)
	if rows.Error != nil {
		return nil, errors.Wrapf(rows.Error, "[Repository] error while getting contacts by email")
	}
	if rows.RowsAffected == 0 {
		return nil, nil
	}
	return contact, nil
}

func (r *contactDBRepo) GetContactByPhone(ctx context.Context, phone string) (*domain.Contact, error) {
	db := r.db.GormConn
	contact := &domain.Contact{}
	rows := db.WithContext(ctx).Where("phone = ?", phone).Find(contact)
	if rows.Error != nil {
		return nil, errors.Wrapf(rows.Error, "[Repository] error while getting contacts by phone")
	}
	if rows.RowsAffected == 0 {
		return nil, nil
	}
	return contact, nil
}

func (r *contactDBRepo) GetAllContacts(ctx context.Context) ([]*domain.Contact, error) {
	db := r.db.GormConn
	var contacts []*domain.Contact
	rows := db.WithContext(ctx).Find(contacts)
	if rows.Error != nil {
		return nil, errors.Wrapf(rows.Error, "[Repository] error while getting all contacts")
	}
	if rows.RowsAffected == 0 {
		return nil, nil
	}
	return contacts, nil
}

func (r *contactDBRepo) GetAllSecondaryContacts(ctx context.Context, linkedID uint) ([]*domain.Contact, error) {
	db := r.db.GormConn
	var contacts []*domain.Contact
	rows := db.WithContext(ctx).Where("linked_id = ? OR contact_id = ?", linkedID, linkedID).Find(&contacts)
	if rows != nil && rows.Error != nil {
		return nil, errors.Wrapf(rows.Error, "[Repository] error while getting contacts by linked_id")
	}
	if rows.RowsAffected == 0 {
		return nil, nil
	}
	return contacts, nil
}

func (r *contactDBRepo) GetPrimaryContactFromLinkedID(ctx context.Context, linkedID uint) (*domain.Contact, error) {
	db := r.db.GormConn
	var contact *domain.Contact
	rows := db.WithContext(ctx).Where("contact_id = ?", linkedID).Find(contact)
	if rows != nil && rows.Error != nil {
		return nil, errors.Wrapf(rows.Error, "[Repository] error while getting contacts by contact_id")
	}
	if rows.RowsAffected == 0 {
		return nil, nil
	}
	return contact, nil
}

func (r *contactDBRepo) CreateContact(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
	db := r.db.GormConn
	rows := db.WithContext(ctx).Create(contact)
	if rows != nil && rows.Error != nil {
		return nil, errors.Wrapf(rows.Error, "[Repository] error while creating a contact")
	}
	return contact, nil
}
