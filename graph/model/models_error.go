package model

import "errors"

var (
	ErrNodeNotFound           = errors.New("node not found")
	ErrReviewNotFound         = errors.New("review not found")
	ErrLikeNotFound           = errors.New("like not found")
	ErrAddressNotFound        = errors.New("address not found")
	ErrTagNotFound            = errors.New("tag not found")
	ErrAmenityNotFound        = errors.New("amenity not found")
	ErrAmenityGroupNotFound   = errors.New("amenityGroup not found")
	ErrTagOptNotFound         = errors.New("tagopt not found")
	ErrNodedataExistValue     = errors.New("nodedata already exist")
	ErrNodedataVoteExistValue = errors.New("nodedata vote already exist")

	ErrLikeExist        = errors.New("like already exist")
	ErrActionExistValue = errors.New("action already exist")

	ErrTagOptExistValue = errors.New("exists already option")
)
