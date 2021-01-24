package direction

import (
	"bufio"
	datepath_server "github.com/PODO/datepath-server"
	"github.com/PODO/datepath-server/pkg/direction"
	"github.com/labstack/echo/v4"
	"io"
)

// DirectionCoordinates godoc
// @Summary 길찾기
// @Description 길찾기
// @Name 추교정
// @Produce json
// @Param origin query string true "출발 지점의 위도,경도. default: 범계역" minlength(3) maxlength(40) default(37.389624,126.950745)
// @Param destination query string true "도착 지점의 위도,경도역. default: 판교역" minlength(3) maxlength(40) default(37.394640,127.111135)
// @Param departure_time query string false "출발 시간. UTC 1970 년 1월 1일 0시 0분 0초 부터 경과한 초를 정수로 반환한 값. 현재 시각으로 찾으려면 'now'를 전달한다." default(1614332700)
// @Success 200 {object} direction.Response
// @Router /direction/coordinates [get]
func DirectionHandler(c echo.Context) error {
	cli := c.Get("datepathServiceClients").(*datepath_server.ServiceClients).Direction

	r := &direction.Request{}
	if err := c.Bind(r); err != nil {
		return err
	}

	resp, err := cli.Direction(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.Response().Header().Set("Content-type", "application/json")
	br := bufio.NewReader(resp.Body)
	_, err = io.Copy(c.Response().Writer, br)

	return err
}
