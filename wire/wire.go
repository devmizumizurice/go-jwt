package wire

// import (
// 	"github.com/devmizumizurice/go-jwt/controllers"
// 	"github.com/devmizumizurice/go-jwt/repositories"
// 	"github.com/devmizumizurice/go-jwt/services"
// 	"github.com/google/wire"
// 	"gorm.io/gorm"
// )

// func InitializeAuthController(db *gorm.DB) controllers.AuthControllerInterface {
// 	wire.Build(
// 		repositories.NewUserRepository,
// 		services.NewAuthService,
// 		controllers.NewAuthController,
// 	)
// 	return nil
// }

// func InitializeUserController(db *gorm.DB) controllers.UserControllerInterface {
// 	wire.Build(
// 		repositories.NewUserRepository,
// 		services.NewUserService,
// 		controllers.NewUserController,
// 	)
// 	return nil
// }
