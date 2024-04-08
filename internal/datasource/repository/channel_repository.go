package repository

import (
	"fund-o/api-server/internal/entity"
	"gorm.io/gorm"
)

type ChannelRepository interface {
	Create(channel *entity.Channel) (*entity.Channel, error)
	GetExistingChannel(userId string, memberId string) (*entity.Channel, error)
	GetByUserID(userId string) ([]entity.Channel, error)
}

type channelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) Create(channel *entity.Channel) (*entity.Channel, error) {
	result := r.db.
		Preload("Members").
		Preload("Messages").
		Create(&channel).
		First(&channel)
	if result.Error != nil {
		return nil, result.Error
	}

	return channel, nil
}

func (r *channelRepository) GetExistingChannel(userId string, memberId string) (*entity.Channel, error) {
	var channel entity.Channel
	result := r.db.Raw(`
		SELECT c.* 
		FROM channels c
		JOIN channel_members cm1 ON c.id = cm1.channel_id AND cm1.user_id = ?
		JOIN channel_members cm2 ON c.id = cm2.channel_id AND cm2.user_id = ?
	`, userId, memberId).Scan(&channel)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := r.db.Model(&channel).Association("Members").Find(&channel.Members); err != nil {
		return nil, err
	}

	if err := r.db.Preload("Author").Model(&channel).Association("Messages").Find(&channel.Messages); err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *channelRepository) GetByUserID(userId string) ([]entity.Channel, error) {
	var channels []entity.Channel
	result := r.db.
		Preload("Members").
		Preload("Messages").
		Preload("Messages.Author").
		Joins("JOIN channel_members ON channels.id = channel_members.channel_id").
		Where("channel_members.user_id = ?", userId).
		Find(&channels)
	if result.Error != nil {
		return nil, result.Error
	}

	return channels, nil
}
