# go-pkgs
### for installing this package you have to save your github credentials
### there is two way you can do it:
```
git config --global credential.helper store
```
### or
```
git config  --global  url."https://{username}:{token}@github.com".insteadOf "https://github.com"
```
### Then install the private package
```
GOPRIVATE=github.com go get github.com/abr-ooo/go-pkgs
```
## Example
```
package routes

import (
	"github.com/abr-ooo/go-pkgs"
	"github.com/gofiber/fiber/v2"
)

func Router(app fiber.Router) {
	api := app.Group("/uid", gopkgs.Auth)
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(gopkgs.UID(c))
	})
}

```
