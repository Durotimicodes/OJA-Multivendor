package router

import (
	"net/http"
	"os"

	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/server/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	ContentType string
	handlers    map[string]func(w http.ResponseWriter, r *http.Request)
}

func SetupRouter(h *handlers.Handler) (*gin.Engine, string) {
	router := gin.Default()

	apirouter := router.Group("/api/v1")
	apirouter.GET("/ping", handlers.PingHandler)
	apirouter.GET("/searchproducts", h.SearchProductHandler)
	apirouter.POST("buyer/forgotpassword", h.BuyerForgotPasswordEMailHandler)
	apirouter.POST("seller/forgotpassword", h.SellerForgotPasswordEMailHandler)
	apirouter.GET("/sellers", h.GetSellers)
	apirouter.GET("/product/:id", h.GetProductById)
	apirouter.GET("/sellers", h.GetSellers)
	apirouter.GET("/product/:id", h.GetProductById)
	apirouter.POST("/buyersignup", h.BuyerSignUpHandler)
	apirouter.POST("/sellersignup", h.SellerSignUpHandler)
	apirouter.POST("/loginbuyer", h.LoginBuyerHandler)
	apirouter.POST("/loginseller", h.LoginSellerHandler)

	apirouter.PUT("/sellerforgotpassword/", h.SellerForgotPasswordResetHandler)
	apirouter.PUT("/buyerforgotpassword/", h.BuyerForgotPasswordResetHandler)

	apirouter.GET("/seller/totalorder/:id", h.SellerTotalOrders)

	//All authorized routes here
	authorizedRoutesBuyer := apirouter.Group("/")
	authorizedRoutesSeller := apirouter.Group("/")
	authorizedRoutesSeller.Use(middleware.AuthorizeSeller(h.DB.FindSellerByEmail, h.DB.TokenInBlacklist))
	authorizedRoutesBuyer.PUT("/buyer/resetpassword/:email", h.BuyerUpdatePassword)
	authorizedRoutesSeller.PUT("/seller/resetpassword/:email", h.SellerUpdatePassword)
	authorizedRoutesBuyer.PUT("/updatebuyerprofile/:id", h.UpdateBuyerProfileHandler)
	authorizedRoutesSeller.PUT("/updatesellerprofile/:id", h.UpdateSellerProfileHandler)

	authorizedRoutesBuyer.Use(middleware.AuthorizeBuyer(h.DB.FindBuyerByEmail, h.DB.TokenInBlacklist))
	{
		authorizedRoutesBuyer.PUT("/updatebuyerprofile", h.UpdateBuyerProfileHandler)
		authorizedRoutesBuyer.GET("/getbuyerprofile", h.GetBuyerProfileHandler)
		authorizedRoutesBuyer.GET("/buyerorders/", h.AllBuyerOrders)
		authorizedRoutesBuyer.PUT("/buyer/updatepassword", h.BuyerUpdatePassword)
	}

	authorizedRoutesBuyer := apirouter.Group("/")
	authorizedRoutesSeller.Use(middleware.AuthorizeSeller(h.DB.FindSellerByEmail, h.DB.TokenInBlacklist))
	{

		authorizedRoutesSeller.GET("/sellerorders/", h.AllSellerOrders)
		authorizedRoutesSeller.GET("/seller/totalorder/", h.SellerTotalOrders)
		authorizedRoutesSeller.PUT("/updatesellerprofile", h.UpdateSellerProfileHandler)
		authorizedRoutesSeller.GET("/getsellerprofile", h.GetSellerProfileHandler)
		authorizedRoutesSeller.PUT("/seller/updatepassword", h.SellerUpdatePassword)
		authorizedRoutesSeller.GET("/seller/shop", h.HandleGetSellerShopByProfileAndProduct())
		authorizedRoutesSeller.GET("/seller/total/product/count", h.GetTotalProductCountForSeller)
		authorizedRoutesSeller.GET("/seller/product", h.SellerIndividualProduct)
		authorizedRoutesSeller.GET("/seller/total/product/sold", h.GetTotalSoldProductCount)
		authorizedRoutesSeller.GET("/seller/allproducts", h.SellerAllProducts)

	}

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":8081"
	}
	return router, port
}
