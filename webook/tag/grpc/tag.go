package grpc

import (
	"context"
	tagv1 "gindemo/webook/api/proto/gen/tag/v1"
	"gindemo/webook/tag/domain"
	"gindemo/webook/tag/service"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

var _ tagv1.TagServiceServer = (*TagServiceServer)(nil)

type TagServiceServer struct {
	tagv1.UnimplementedTagServiceServer
	service service.TagService
}

func NewTagServiceServer(svc service.TagService) *TagServiceServer {
	return &TagServiceServer{
		service: svc,
	}
}

func (t *TagServiceServer) Register(server grpc.ServiceRegistrar) {
	tagv1.RegisterTagServiceServer(server, t)
}

func (t *TagServiceServer) CreateTag(ctx context.Context, req *tagv1.CreateTagRequest) (*tagv1.CreateTagResponse, error) {
	id, err := t.service.CreateTag(ctx, req.Uid, req.Name)
	return &tagv1.CreateTagResponse{
		Tag: &tagv1.Tag{
			Id:   id,
			Name: req.Name,
			Uid:  req.Uid,
		},
	}, err
}

func (t *TagServiceServer) AttachTags(ctx context.Context, req *tagv1.AttachTagsRequest) (*tagv1.AttachTagsResponse, error) {
	err := t.service.AttachTags(ctx, req.Uid, req.Biz, req.BizId, req.Tids)
	return &tagv1.AttachTagsResponse{}, err
}

func (t *TagServiceServer) GetTags(ctx context.Context, req *tagv1.GetTagsRequest) (*tagv1.GetTagsResponse, error) {
	tags, err := t.service.GetTags(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return &tagv1.GetTagsResponse{
		Tags: slice.Map(tags, func(idx int, src domain.Tag) *tagv1.Tag {
			return t.toDTO(src)
		}),
	}, nil
}

func (t *TagServiceServer) GetBizTags(ctx context.Context, req *tagv1.GetBizTagsRequest) (*tagv1.GetBizTagsResponse, error) {
	res, err := t.service.GetBizTags(ctx, req.Uid, req.Biz, req.BizId)
	if err != nil {
		return nil, err
	}
	return &tagv1.GetBizTagsResponse{
		Tags: slice.Map(res, func(idx int, src domain.Tag) *tagv1.Tag {
			return t.toDTO(src)
		}),
	}, nil
}

func (t *TagServiceServer) toDTO(tag domain.Tag) *tagv1.Tag {
	return &tagv1.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}
