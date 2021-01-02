package local

import (
	"bufio"
	"github.com/PODO/datepath-server/pkg/local"
	"github.com/labstack/echo/v4"
	"io"
)

// Local Keyword Search godoc
// @Summary 키워드 장소 검색
// @Description 키워드 장소 검색
// @Name 추교정
// @Accept x-www-form-urlencoded
// @Produce json
// @Param query query string true "검색을 원하는 질의어" minlength(1) maxlength(15)
// @Param x query number false "중심 좌표의 X값 혹은 longitude. 특정 지역을 중심으로 검색하려고 할 경우 radius와 함께 사용 가능"
// @Param y query number false "중심 좌표의 Y값 혹은 latitude. 특정 지역을 중심으로 검색하려고 할 경우 radius와 함께 사용 가능"
// @Param radius query integer false "중심 좌표부터의 반경거리. 특정 지역을 중심으로 검색하려고 할 경우 중심좌표로 쓰일 x,y와 함께 사용. 단위 meter, 0~20000 사이의 값"
// @Success 200 {object} local.Response
// @Router /local/search/keyword [post]
func SearchKeywordHandler(c echo.Context) error {
	cli := local.NewClient()

	r := &local.Request{}
	if err := c.Bind(r); err != nil {
		return err
	}

	resp, err := cli.SearchKeyword(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.Response().Header().Set("Content-type", "application/json")
	br := bufio.NewReader(resp.Body)
	_, err = io.Copy(c.Response().Writer, br)

	return err
}
