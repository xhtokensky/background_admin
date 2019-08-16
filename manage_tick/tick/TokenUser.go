package tick

import (
	"tokensky_bg_admin/conf"
	"tokensky_bg_admin/models"
)

//维护用户地址维护
var tickTokenskyUserAddressUpSign bool = true
func TickTokenskyUserAddressUp() error {
	if tickTokenskyUserAddressUpSign{
		tickTokenskyUserAddressUpSign = false
		defer func() {tickTokenskyUserAddressUpSign=true}()
		for _, coinType := range conf.TOKENSKY_ADDRESS_COIN_TYPES {
			total := models.TokenskyUserAddressGetNotUsedCont(coinType)
			if total < conf.TBI_SERVER_ADDRESS_MAX {
				models.TokenskyUserAddressAddNum(coinType, conf.TBI_SERVER_ADDRESS_MAX-total)
			}
		}
	}
	return nil
}