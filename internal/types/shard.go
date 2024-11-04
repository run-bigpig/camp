package types

type ListRoomRequest struct {
	NtwNum           string `json:"ntwNum"`
	OpenID           string `json:"openId"`
	PromotionChannel int    `json:"promotionChannel"`
	RoomStartDate    string `json:"roomStartDate"`
	RoomEndDate      string `json:"roomEndDate"`
	RoomPageNum      int    `json:"roomPageNum"`
	RoomPageSize     int    `json:"roomPageSize"`
	ID               string `json:"id"`
}

type RoomList struct {
	Name       string `json:"name"`
	RoomTypeID string `json:"roomTypeId"`
}
type RoomData struct {
	PageNum  int        `json:"pageNum"`
	PageSize int        `json:"pageSize"`
	Total    int        `json:"total"`
	List     []RoomList `json:"list"`
}

type ListRoomResponse struct {
	Code string   `json:"code"`
	Msg  string   `json:"msg"`
	Data RoomData `json:"data"`
}

type CommitRequest struct {
	NtwNum                  string        `json:"ntwNum"`
	OpenID                  string        `json:"openId"`
	PromotionChannel        int           `json:"promotionChannel"`
	BoolPublish             bool          `json:"boolPublish"`
	RoomTypeID              string        `json:"roomTypeId"`
	Count                   int           `json:"count"`
	StartDate               string        `json:"startDate"`
	EndDate                 string        `json:"endDate"`
	CustomerPhone           string        `json:"customerPhone"`
	CustomerName            string        `json:"customerName"`
	Remark                  string        `json:"remark"`
	BookingType             int           `json:"bookingType"`
	AppType                 int           `json:"appType"`
	BoolUseIntegral         bool          `json:"boolUseIntegral"`
	DiscountInstList        []interface{} `json:"discountInstList"`
	BoolUseContinueLiveInst bool          `json:"boolUseContinueLiveInst"`
}

type OrderData struct {
	OrderID string `json:"orderId"`
}

type CommitResponse struct {
	Code string    `json:"code"`
	Msg  string    `json:"msg"`
	Data OrderData `json:"data"`
}
