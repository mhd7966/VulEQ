package gopkgs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type response struct {
	Data struct {
		UID uint `json:"id"`
	} `json:"data"`
}

func Auth(c *fiber.Ctx) error {
	client := http.Client{}

	baseurl, exist := os.LookupEnv("DEL_BASE_URL")
	if !exist {
		baseurl = "https://api.abr.ooo"
	}

	req, err := http.NewRequest("GET", baseurl+"/v1/user", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", c.Get("Authorization"))
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Auth Server is Down", "data": nil})
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Authorization Header Not Found or Incorrect", "data": nil})
	}

	var res response
	if err = json.Unmarshal(body, &res); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Authorization Header Not Found or Incorrect", "data": nil})
	}

	if res.Data.UID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Authorization Header Not Found or Incorrect", "data": nil})
	}

	c.Locals("uid", res.Data.UID)
	return c.Next()
}

func UID(c *fiber.Ctx) uint {
	uid := c.Locals("uid")
	if uid == nil {
		panic("Use Auth Middleware Before Calling This Method")
	}
	return uid.(uint)
}
