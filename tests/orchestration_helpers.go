package tests

import (
	"bytes"
	"mime/multipart"

	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Role
	role_cache "github.com/MamangRust/monolith-ecommerce-grpc-role/cache"
	role_handler "github.com/MamangRust/monolith-ecommerce-grpc-role/handler"
	role_repo "github.com/MamangRust/monolith-ecommerce-grpc-role/repository"
	role_service "github.com/MamangRust/monolith-ecommerce-grpc-role/service"

	// User
	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-user/cache"
	user_handler "github.com/MamangRust/monolith-ecommerce-grpc-user/handler"
	user_repo "github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	user_service "github.com/MamangRust/monolith-ecommerce-grpc-user/service"

	// Auth
	auth_cache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	auth_handler "github.com/MamangRust/monolith-ecommerce-auth/handler"
	auth_repo "github.com/MamangRust/monolith-ecommerce-auth/repository"
	auth_service "github.com/MamangRust/monolith-ecommerce-auth/service"

	// Banner
	banner_cache "github.com/MamangRust/monolith-ecommerce-grpc-banner/cache"
	banner_handler "github.com/MamangRust/monolith-ecommerce-grpc-banner/handler"
	banner_repo "github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	banner_service "github.com/MamangRust/monolith-ecommerce-grpc-banner/service"

	// Slider
	slider_cache "github.com/MamangRust/monolith-ecommerce-grpc-slider/cache"
	slider_handler "github.com/MamangRust/monolith-ecommerce-grpc-slider/handler"
	slider_repo "github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
	slider_service "github.com/MamangRust/monolith-ecommerce-grpc-slider/service"

	// Category
	category_cache "github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	category_handler "github.com/MamangRust/monolith-ecommerce-grpc-category/handler"
	category_repo "github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	category_service "github.com/MamangRust/monolith-ecommerce-grpc-category/service"

	// Product
	product_cache "github.com/MamangRust/monolith-ecommerce-grpc-product/cache"
	product_handler "github.com/MamangRust/monolith-ecommerce-grpc-product/handler"
	product_repo "github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	product_service "github.com/MamangRust/monolith-ecommerce-grpc-product/service"

	// Cart
	cart_cache "github.com/MamangRust/monolith-ecommerce-grpc-cart/cache"
	cart_handler "github.com/MamangRust/monolith-ecommerce-grpc-cart/handler"
	cart_repo "github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	cart_service "github.com/MamangRust/monolith-ecommerce-grpc-cart/service"

	// Merchant
	merchant_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	merchant_handler "github.com/MamangRust/monolith-ecommerce-grpc-merchant/handler"
	merchant_repo "github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	merchant_service "github.com/MamangRust/monolith-ecommerce-grpc-merchant/service"

	// Order
	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	order_handler "github.com/MamangRust/monolith-ecommerce-grpc-order/handler"
	order_repo "github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	order_service "github.com/MamangRust/monolith-ecommerce-grpc-order/service"

	// Merchant Award
	merchant_award_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/cache"
	merchant_award_handler "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/handler"
	merchant_award_repo "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/repository"
	merchant_award_service "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/service"

	// Merchant Business
	merchant_business_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/cache"
	merchant_business_handler "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/handler"
	merchant_business_repo "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	merchant_business_service "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/service"

	// Transaction
	transaction_cache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	transaction_handler "github.com/MamangRust/monolith-ecommerce-grpc-transaction/handler"
	transaction_repo "github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	transaction_service "github.com/MamangRust/monolith-ecommerce-grpc-transaction/service"

	// Merchant Detail
	merchant_detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/cache"
	merchant_detail_handler "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/handler"
	merchant_detail_repo "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/repository"
	merchant_detail_service "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/service"

	// Merchant Policy
	merchant_policy_cache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/cache"
	merchant_policy_handler "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/handler"
	merchant_policy_repo "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/repository"
	merchant_policy_service "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/service"

	// Shipping Address
	shipping_address_cache "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/cache"
	shipping_address_handler "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/handler"
	shipping_address_repo "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	shipping_address_service "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"

	// Order Item
	order_item_cache "github.com/MamangRust/monolith-ecommerce-grpc-order-item/cache"
	order_item_handler "github.com/MamangRust/monolith-ecommerce-grpc-order-item/handler"
	order_item_repo "github.com/MamangRust/monolith-ecommerce-grpc-order-item/repository"
	order_item_service "github.com/MamangRust/monolith-ecommerce-grpc-order-item/service"

	// Review
	review_cache "github.com/MamangRust/monolith-ecommerce-grpc-review/cache"
	review_handler "github.com/MamangRust/monolith-ecommerce-grpc-review/handler"
	review_repo "github.com/MamangRust/monolith-ecommerce-grpc-review/repository"
	review_service "github.com/MamangRust/monolith-ecommerce-grpc-review/service"

	// Review Detail
	review_detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/cache"
	review_detail_handler "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/handler"
	review_detail_repo "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/repository"
	review_detail_service "github.com/MamangRust/monolith-ecommerce-grpc-review-detail/service"
)

func (s *BaseTestSuite) SetupRoleService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(queries)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repository:    roleRepos,
		Logger:        s.Log,
		Cache:         roleMencache,
		Observability: s.Obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterRoleQueryServiceServer(server, roleGapi.RoleQuery)
	pb.RegisterRoleCommandServiceServer(server, roleGapi.RoleCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["role"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupUserService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())
	hasher := hash.NewHashingPassword()

	userMencache := user_cache.NewMencache(cacheStore)
	roleQueryClient := pb.NewRoleQueryServiceClient(s.Conns["role"])
	userRepos := user_repo.NewRepositories(queries, roleQueryClient)
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        s.Log,
		Hash:          hasher,
		Cache:         userMencache,
		Observability: s.Obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterUserQueryServiceServer(server, userGapi.UserQuery)
	pb.RegisterUserCommandServiceServer(server, userGapi.UserCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["user"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupAuthService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())
	hasher := hash.NewHashingPassword()
	tokenManager, _ := auth.NewManager("mysecret")

	userQueryClient := pb.NewUserQueryServiceClient(s.Conns["user"])
	userCommandClient := pb.NewUserCommandServiceClient(s.Conns["user"])
	roleQueryClient := pb.NewRoleQueryServiceClient(s.Conns["role"])
	roleCommandClient := pb.NewRoleCommandServiceClient(s.Conns["role"])

	authRepos := auth_repo.NewRepositories(queries, userQueryClient, userCommandClient, roleQueryClient, roleCommandClient)
	authMencache := auth_cache.NewMencache(cacheStore)
	authSvc := auth_service.NewService(&auth_service.Deps{
		Repositories:  authRepos,
		Logger:        s.Log,
		Mencache:      authMencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: s.Obs,
	})
	authGapi := auth_handler.NewAuthHandleGrpc(authSvc, s.Log)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["auth"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupBannerService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	bannerMencache := banner_cache.NewMencache(cacheStore)
	bannerRepos := banner_repo.NewRepositories(queries)
	bannerSvc := banner_service.NewService(&banner_service.Deps{
		Cache:         bannerMencache,
		Repository:    bannerRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	bannerGapi := banner_handler.NewHandler(&banner_handler.Deps{
		Service: bannerSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterBannerQueryServiceServer(server, bannerGapi.BannerQuery)
	pb.RegisterBannerCommandServiceServer(server, bannerGapi.BannerCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["banner"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupSliderService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	sliderMencache := slider_cache.NewMencache(cacheStore)
	sliderRepos := slider_repo.NewRepositories(queries)
	sliderSvc := slider_service.NewService(&slider_service.Deps{
		Mencache:      sliderMencache,
		Repositories:  sliderRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	sliderGapi := slider_handler.NewHandler(&slider_handler.Deps{
		Service: sliderSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterSliderQueryServiceServer(server, sliderGapi.SliderQuery)
	pb.RegisterSliderCommandServiceServer(server, sliderGapi.SliderCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["slider"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCategoryService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	catMencache := category_cache.NewMencache(cacheStore)
	catRepos := category_repo.NewRepositories(queries)
	catSvc := category_service.NewService(&category_service.Deps{
		Cache:         catMencache,
		Repositories:  catRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	catGapi := category_handler.NewHandler(&category_handler.Deps{
		Service: catSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterCategoryQueryServiceServer(server, catGapi.CategoryQuery)
	pb.RegisterCategoryCommandServiceServer(server, catGapi.CategoryCommand)
	pb.RegisterCategoryStatsServiceServer(server, catGapi.CategoryStats)
	pb.RegisterCategoryStatsByIdServiceServer(server, catGapi.CategoryStatsById)
	pb.RegisterCategoryStatsByMerchantServiceServer(server, catGapi.CategoryStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["category"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupProductService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	prodMencache := product_cache.NewMencache(cacheStore)
	catQueryClient := pb.NewCategoryQueryServiceClient(s.Conns["category"])
	merchantQueryClient := pb.NewMerchantQueryServiceClient(s.Conns["merchant"])
	prodRepos := product_repo.NewRepositories(queries, catQueryClient, merchantQueryClient)
	prodSvc := product_service.NewService(&product_service.Deps{
		Cache:         prodMencache,
		Repository:    prodRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	prodGapi := product_handler.NewHandler(&product_handler.Deps{
		Service: prodSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterProductQueryServiceServer(server, prodGapi.ProductQuery)
	pb.RegisterProductCommandServiceServer(server, prodGapi.ProductCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["product"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCartService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	cartMencache := cart_cache.NewMencache(cacheStore)
	userQueryClient := pb.NewUserQueryServiceClient(s.Conns["user"])
	productQueryClient := pb.NewProductQueryServiceClient(s.Conns["product"])
	cartRepos := cart_repo.NewRepositories(queries, userQueryClient, productQueryClient)
	cartSvc := cart_service.NewService(&cart_service.Deps{
		Cache:         cartMencache,
		Repositories:  cartRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	cartGapi := cart_handler.NewHandler(&cart_handler.Deps{
		Service: cartSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterCartQueryServiceServer(server, cartGapi.CartQuery)
	pb.RegisterCartCommandServiceServer(server, cartGapi.CartCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["cart"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	merchantMencache := merchant_cache.NewMencache(cacheStore)
	userQueryClient := pb.NewUserQueryServiceClient(s.Conns["user"])
	merchantRepos := merchant_repo.NewRepositories(queries, userQueryClient)
	merchantSvc := merchant_service.NewService(&merchant_service.Deps{
		Mencache:      merchantMencache,
		Repositories:  merchantRepos,
		Logger:        s.Log,
		Observability: s.Obs,
		Kafka:         nil,
	})
	merchantGapi := merchant_handler.NewHandler(&merchant_handler.Deps{
		Service: merchantSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterMerchantQueryServiceServer(server, merchantGapi.MerchantQuery)
	pb.RegisterMerchantCommandServiceServer(server, merchantGapi.MerchantCommandHandler)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	orderMencache := order_cache.NewMencache(cacheStore)
	orderRepos := order_repo.NewRepositories(&order_repo.Deps{
		DB:                 queries,
		MerchantQuery:      pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		ProductQuery:       pb.NewProductQueryServiceClient(s.Conns["product"]),
		ProductCommand:     pb.NewProductCommandServiceClient(s.Conns["product"]),
		OrderItemQuery:     pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		OrderItemCommand:   pb.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
		UserQuery:          pb.NewUserQueryServiceClient(s.Conns["user"]),
		ShippingCommand:    pb.NewShippingCommandServiceClient(s.Conns["shipping-address"]),
		ShippingQuery:      pb.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
		TransactionCommand: pb.NewTransactionCommandServiceClient(s.Conns["transaction"]),
	})
	orderSvc := order_service.NewService(&order_service.Deps{
		Cache:         orderMencache,
		Repositories:  orderRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	orderGapi := order_handler.NewHandler(&order_handler.Deps{
		Service: orderSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterOrderQueryServiceServer(server, orderGapi.OrderQuery)
	pb.RegisterOrderStatsServiceServer(server, orderGapi.OrderStats)
	pb.RegisterOrderCommandServiceServer(server, orderGapi.OrderCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantAwardService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	awardMencache := merchant_award_cache.NewMencache(cacheStore)
	awardRepos := merchant_award_repo.NewRepositories(queries, pb.NewMerchantQueryServiceClient(s.Conns["merchant"]))
	awardSvc := merchant_award_service.NewService(&merchant_award_service.Deps{
		Cache:         awardMencache,
		Repository:    awardRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	awardGapi := merchant_award_handler.NewHandler(&merchant_award_handler.Deps{
		Service: awardSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterMerchantAwardQueryServiceServer(server, awardGapi.MerchantAwardQuery)
	pb.RegisterMerchantAwardCommandServiceServer(server, awardGapi.MerchantAwardCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant_award"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantBusinessService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	businessMencache := merchant_business_cache.NewMencache(cacheStore)
	businessRepos := merchant_business_repo.NewRepositories(queries, pb.NewMerchantQueryServiceClient(s.Conns["merchant"]))
	businessSvc := merchant_business_service.NewService(&merchant_business_service.Deps{
		Cache:         businessMencache,
		Repository:    businessRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	businessGapi := merchant_business_handler.NewHandler(&merchant_business_handler.Deps{
		Service: businessSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterMerchantBusinessQueryServiceServer(server, businessGapi.MerchantBusinessQuery)
	pb.RegisterMerchantBusinessCommandServiceServer(server, businessGapi.MerchantBusinessCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant_business"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupTransactionService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	transactionMencache := transaction_cache.NewMencache(cacheStore)
	transactionRepos := transaction_repo.NewRepositories(&transaction_repo.Deps{
		DB:             queries,
		UserQuery:      pb.NewUserQueryServiceClient(s.Conns["user"]),
		MerchantQuery:  pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		OrderQuery:     pb.NewOrderQueryServiceClient(s.Conns["order"]),
		OrderItemQuery: pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		ShippingQuery:  pb.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
	})
	transactionSvc := transaction_service.NewService(&transaction_service.Deps{
		Cache:         transactionMencache,
		Repositories:  transactionRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	transactionGapi := transaction_handler.NewHandler(&transaction_handler.Deps{
		Service: transactionSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterTransactionQueryServiceServer(server, transactionGapi.TransactionQuery)
	pb.RegisterTransactionCommandServiceServer(server, transactionGapi.TransactionCommand)
	pb.RegisterTransactionStatsServiceServer(server, transactionGapi.TransactionStats)
	pb.RegisterTransactionStatsByMerchantServiceServer(server, transactionGapi.TransactionStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["transaction"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantDetailService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	detailMencache := merchant_detail_cache.NewMencache(cacheStore)
	detailRepos := merchant_detail_repo.NewRepositories(queries, pb.NewMerchantQueryServiceClient(s.Conns["merchant"]))
	detailSvc := merchant_detail_service.NewService(&merchant_detail_service.Deps{
		Cache:         detailMencache,
		Repository:    detailRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	detailGapi := merchant_detail_handler.NewHandler(&merchant_detail_handler.Deps{
		Service: detailSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterMerchantDetailQueryServiceServer(server, detailGapi.MerchantDetailQuery)
	pb.RegisterMerchantDetailCommandServiceServer(server, detailGapi.MerchantDetailCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant_detail"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantPolicyService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	policyMencache := merchant_policy_cache.NewMencache(cacheStore)
	policyRepos := merchant_policy_repo.NewRepositories(queries, pb.NewMerchantQueryServiceClient(s.Conns["merchant"]))
	policySvc := merchant_policy_service.NewService(&merchant_policy_service.Deps{
		Cache:         policyMencache,
		Repository:    policyRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	policyGapi := merchant_policy_handler.NewHandler(&merchant_policy_handler.Deps{
		Service: policySvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterMerchantPolicyQueryServiceServer(server, policyGapi.MerchantPolicyQuery)
	pb.RegisterMerchantPolicyCommandServiceServer(server, policyGapi.MerchantPolicyCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant_policy"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupShippingAddressService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	addrMencache := shipping_address_cache.NewMencache(cacheStore)
	addrRepos := shipping_address_repo.NewRepositories(queries)
	addrSvc := shipping_address_service.NewService(&shipping_address_service.Deps{
		Mencache:      addrMencache,
		Repositories:  addrRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	addrGapi := shipping_address_handler.NewHandler(&shipping_address_handler.Deps{
		Service: addrSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterShippingQueryServiceServer(server, addrGapi.ShippingQuery)
	pb.RegisterShippingCommandServiceServer(server, addrGapi.ShippingCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["shipping-address"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderItemService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	itemMencache := order_item_cache.NewMencache(cacheStore)
	itemRepos := order_item_repo.NewRepositories(queries)
	itemSvc := order_item_service.NewService(&order_item_service.Deps{
		Cache:         itemMencache,
		Repository:    itemRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	itemGapi := order_item_handler.NewHandler(&order_item_handler.Deps{
		Service: itemSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterOrderItemQueryServiceServer(server, itemGapi.OrderItemQuery)
	pb.RegisterOrderItemCommandServiceServer(server, itemGapi.OrderItemCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order-item"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupReviewService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	reviewMencache := review_cache.NewMencache(cacheStore)
	userQueryClient := pb.NewUserQueryServiceClient(s.Conns["user"])
	productQueryClient := pb.NewProductQueryServiceClient(s.Conns["product"])
	reviewRepos := review_repo.NewRepositories(queries, userQueryClient, productQueryClient)
	reviewSvc := review_service.NewService(&review_service.Deps{
		Cache:         reviewMencache,
		Repositories:  reviewRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	reviewGapi := review_handler.NewHandler(&review_handler.Deps{
		Service: reviewSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterReviewQueryServiceServer(server, reviewGapi.ReviewQuery)
	pb.RegisterReviewCommandServiceServer(server, reviewGapi.ReviewCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["review"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupReviewDetailService() {
	cacheStore := s.GetCacheStore()
	queries := db.New(s.ts.DBPool())

	detailMencache := review_detail_cache.NewMencache(cacheStore)
	detailRepos := review_detail_repo.NewRepositories(queries)
	detailSvc := review_detail_service.NewService(&review_detail_service.Deps{
		Cache:         detailMencache,
		Repositories:  detailRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	detailGapi := review_detail_handler.NewHandler(&review_detail_handler.Deps{
		Service: detailSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pb.RegisterReviewDetailQueryServiceServer(server, detailGapi.ReviewDetailQuery)
	pb.RegisterReviewDetailCommandServiceServer(server, detailGapi.ReviewDetailCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["review-detail"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) GetCacheStore() *cache.CacheStore {
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	return cache.NewCacheStore(s.ts.RedisClient(), s.Log, cacheMetrics)
}

func (s *BaseTestSuite) BuildMultipartRequestBody(fields map[string]string, fieldName, fileName string) ([]byte, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range fields {
		fw, _ := w.CreateFormField(key)
		fw.Write([]byte(r))
	}
	fw, _ := w.CreateFormFile(fieldName, fileName)
	fw.Write([]byte("dummy image content"))
	w.Close()
	return b.Bytes(), w.FormDataContentType()
}
