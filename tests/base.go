package tests

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type BaseTestSuite struct {
	suite.Suite
	ts      *TestSuite
	Log     logger.LoggerInterface
	Obs     observability.TraceLoggerObservability
	Conns   map[string]*grpc.ClientConn
	Servers []*grpc.Server
	Ctx     context.Context
	Cancel  context.CancelFunc
}

// seedUserCounter and seedUniqueSuffix make seeded values unique across multiple
// seed calls within the same test process (suites share one database per package
// run), so fixed slugs/emails never collide between test methods.
var seedUserCounter int64

func uniqueSuffix() string {
	seedUserCounter++
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), seedUserCounter)
}

func (s *BaseTestSuite) SetupSuite() {
	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	ts, err := SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	logger.ResetInstance()
	lp := sdklog.NewLoggerProvider()
	s.Log, _ = logger.NewLogger("test", lp)

	if s.Log == nil || (reflect.ValueOf(s.Log).Kind() == reflect.Ptr && reflect.ValueOf(s.Log).IsNil()) {
		z, _ := zap.NewDevelopment()
		s.Log = &logger.Logger{Log: z}
	}

	s.Obs, err = observability.NewObservability("test", s.Log)
	s.Require().NoError(err)
	s.Require().NotNil(s.Obs)
	s.Conns = make(map[string]*grpc.ClientConn)
}

func (s *BaseTestSuite) TearDownSuite() {
	for _, conn := range s.Conns {
		conn.Close()
	}
	for _, server := range s.Servers {
		server.GracefulStop()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
	if s.Cancel != nil {
		s.Cancel()
	}
}

func (s *BaseTestSuite) DBPool() *pgxpool.Pool {
	return s.ts.DBPool()
}

func (s *BaseTestSuite) RedisClient() *goredis.Client {
	return s.ts.RedisClient()
}

func (s *BaseTestSuite) RegisterServer(server *grpc.Server) string {
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)
	return addr
}

func (s *BaseTestSuite) GetConnection(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) SeedUser(ctx context.Context) int {
	// Each call seeds a unique email so suites with multiple test methods can
	// share one database without colliding on a fixed address.
	email := fmt.Sprintf("seed.user.%s@example.com", uniqueSuffix())
	res, err := pb.NewUserCommandServiceClient(s.Conns["user"]).Create(ctx, &pb.CreateUserRequest{
		Firstname:       "Seed",
		Lastname:        "User",
		Email:           email,
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedCategory(ctx context.Context) int {
	seedSuffix := uniqueSuffix()
	res, err := pb.NewCategoryCommandServiceClient(s.Conns["category"]).Create(ctx, &pb.CreateCategoryRequest{
		Name:          "Seed Category " + seedSuffix,
		Description:   "Seed Description",
		SlugCategory:  "seed-category-" + seedSuffix,
		ImageCategory: "seed.jpg",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedMerchant(ctx context.Context, userID int) int {
	res, err := pb.NewMerchantCommandServiceClient(s.Conns["merchant"]).Create(ctx, &pb.CreateMerchantRequest{
		UserId:       int32(userID),
		Name:         "Seed Merchant",
		Description:  "Seed Description",
		Address:      "Seed Address",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedProduct(ctx context.Context, merchantID int, categoryID int) int {
	seedSuffix := uniqueSuffix()
	res, err := pb.NewProductCommandServiceClient(s.Conns["product"]).Create(ctx, &pb.CreateProductRequest{
		MerchantId:   int32(merchantID),
		CategoryId:   int32(categoryID),
		Name:         "Seed Product " + seedSuffix,
		Description:  "Seed Description",
		Price:        10000,
		CountInStock: 100,
		Brand:        "Seed Brand",
		Weight:       1000,
		SlugProduct:  "seed-product-" + seedSuffix,
		ImageProduct: "seed.jpg",
		Barcode:      "123456789",
		Rating:       5,
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedShippingAddress(ctx context.Context, orderID int) int {
	res, err := pb.NewShippingCommandServiceClient(s.Conns["shipping-address"]).CreateShipping(ctx, &pb.CreateShippingAddressRequest{
		OrderId:        int32(orderID),
		Alamat:         "Seed Address",
		Provinsi:       "Seed Province",
		Kota:           "Seed City",
		Negara:         "Seed Country",
		Courier:        "Seed Courier",
		ShippingMethod: "Seed Method",
		ShippingCost:   1000,
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrder(ctx context.Context, userID int, merchID int, prodID int) int {
	res, err := pb.NewOrderCommandServiceClient(s.Conns["order"]).Create(ctx, &pb.CreateOrderRequest{
		UserId:     int32(userID),
		MerchantId: int32(merchID),
		TotalPrice: 10000,
		Items: []*pb.CreateOrderItemRequest{
			{
				ProductId: int32(prodID),
				Quantity:  1,
				Price:     10000,
			},
		},
		Shipping: &pb.CreateShippingAddressRequest{
			Alamat:         "Seed Address",
			Provinsi:       "Seed Province",
			Kota:           "Seed City",
			Negara:         "Seed Country",
			Courier:        "Seed Courier",
			ShippingMethod: "Seed Method",
			ShippingCost:   1000,
		},
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedReview(ctx context.Context, userID int, productID int) int {
	res, err := pb.NewReviewCommandServiceClient(s.Conns["review"]).Create(ctx, &pb.CreateReviewRequest{
		UserId:    int32(userID),
		ProductId: int32(productID),
		Rating:    5,
		Comment:   "Seed Review",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrderItem(ctx context.Context, orderID int, productID int) int {
	res, err := pb.NewOrderItemCommandServiceClient(s.Conns["order-item"]).CreateOrderItem(ctx, &pb.CreateOrderItemRecordRequest{
		OrderId:   int32(orderID),
		ProductId: int32(productID),
		Quantity:  1,
		Price:     1000,
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}
