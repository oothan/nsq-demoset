package server

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"nsq-demoset/app/app-services/cmd/market/data"
	marketpb "nsq-demoset/app/app-services/proto/market/v1"
	"time"
)

type MarketServer struct {
	marketData *data.MarketData
	marketpb.UnimplementedMarketServer
}

func NewMarketServer() *MarketServer {
	m := data.NewMarketData()
	return &MarketServer{
		marketData: m,
	}
}

func (s *MarketServer) Subscribe(in *marketpb.MarketRequest, stream marketpb.Market_SubscribeServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended.")
		default:
			time.Sleep(time.Second)
			for key, val := range s.marketData.List {
				err := stream.SendMsg(&marketpb.MarketResponse{
					Symbol:             key,
					OpenPrice:          val.OpenPrice,
					ClosePrice:         val.HighPrice,
					LastPrice:          val.LastPrice,
					PriceChange:        val.PriceChange,
					PriceChangePercent: val.PriceChangePercent,
				})
				if err != nil {
					return status.Error(codes.Internal, "Stream interrupted.")
				}
			}
		}
	}
}

type responseItem struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (s *MarketServer) GetPrice(ctx context.Context, in *marketpb.PriceRequest) (*marketpb.PriceResponse, error) {
	resp, err := http.Get("https://api3.binance.com/api/v3/ticker/price")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	list := make([]*responseItem, 0)
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, err
	}

	item, ok := s.findSymbol(in.GetSymbol(), list)
	if !ok {
		return nil, status.Error(codes.NotFound, "Symbol not found")
	}

	return &marketpb.PriceResponse{
		Symbol: item.Symbol,
		Price:  item.Price,
	}, nil
}

func (s *MarketServer) findSymbol(symbol string, list []*responseItem) (*responseItem, bool) {
	for _, item := range list {
		if item.Symbol == symbol {
			return item, true
		}
	}
	return nil, false
}
