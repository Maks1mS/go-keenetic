package client

import (
	"errors"
	"fmt"
	"strings"
)

type IProute struct {
	Destination string
	Gateway     string
	Interface   string
	Metric      string
}

func GetIpRoutesCmd() string {
	return "show ip route"
}

func (c *Client) GetIpRoutes() []IProute {
	raw := c.ExecuteCommand(GetIpRoutesCmd())

	prefix := 305
	postfix := 12

	string_result := string(raw[prefix : len(raw)-postfix])
	lines := strings.Split(string_result, "\r\n")

	result := make([]IProute, len(lines))

	for i, x := range lines {
		l := strings.Fields(x)
		result[i] = IProute{Destination: l[0], Gateway: l[1], Interface: l[2], Metric: l[3]}
	}

	return result
}

func AddIpRouteCmd(route IProute, auto bool) string {
	autoString := "auto"

	if !auto {
		autoString = ""
	}

	return fmt.Sprintf("ip route %s %s %s %s %s", route.Destination, route.Gateway, route.Interface, route.Metric, autoString)
}

func (c *Client) AddIpRoute(route IProute, auto bool) error {
	cmd := AddIpRouteCmd(route, auto)
	raw := c.ExecuteCommand(cmd)

	prefix := len(cmd)*4 + 7
	postfix := 11

	result := string(raw[prefix : len(raw)-postfix])

	if strings.HasPrefix(result, "Network::RoutingTable: added static route") || strings.HasPrefix(result, "Network::RoutingTable: renewed static route") {
		return nil
	} else {
		return errors.New(result)
	}
}

func RemoveIpRouteCmd(route IProute) string {
	return fmt.Sprintf("no ip route %s %s %s", route.Destination, route.Interface, route.Metric)
}

func (c *Client) RemoveIpRoute(route IProute) error {
	cmd := RemoveIpRouteCmd(route)
	raw := c.ExecuteCommand(cmd)

	prefix := len(cmd)*4 + 7
	postfix := 11

	result := string(raw[prefix : len(raw)-postfix])

	if strings.HasPrefix(result, "Network::RoutingTable: deleted static route:") {
		return nil
	} else {
		return errors.New(result)
	}
}
