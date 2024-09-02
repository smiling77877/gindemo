package grpc

import (
	"context"
	pmtv1 "gindemo/webook/api/proto/gen/payment/v1"
	"gindemo/webook/payment/domain"
	"gindemo/webook/payment/service/wechat"
	"google.golang.org/grpc"
)

type WechatServiceServer struct {
	pmtv1.UnimplementedWechatPaymentServiceServer
	svc *wechat.NativePaymentService
}

func NewWechatServiceServer(svc *wechat.NativePaymentService) *WechatServiceServer {
	return &WechatServiceServer{svc: svc}
}

func (s *WechatServiceServer) Register(server *grpc.Server) {
	pmtv1.RegisterWechatPaymentServiceServer(server, s)
}

func (s *WechatServiceServer) NativePrePay(ctx context.Context, req *pmtv1.PrePayRequest) (*pmtv1.NativePrePayResponse, error) {
	codeURL, err := s.svc.Prepay(ctx, domain.Payment{
		Amt: domain.Amount{
			Currency: req.Amt.Currency,
			Total:    req.Amt.Total,
		},
		BizTradeNO:  req.BizTradeNo,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pmtv1.NativePrePayResponse{
		CodeUrl: codeURL,
	}, nil
}

func (s *WechatServiceServer) GetPayment(ctx context.Context, req *pmtv1.GetPaymentRequest) (*pmtv1.GetPaymentResponse, error) {
	p, err := s.svc.GetPayment(ctx, req.GetBizTradeNo())
	if err != nil {
		return nil, err
	}
	return &pmtv1.GetPaymentResponse{
		Status: pmtv1.PaymentStatus(p.Status),
	}, nil
}
