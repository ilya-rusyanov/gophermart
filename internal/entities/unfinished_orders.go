package entities

type UnfinishedOrders map[OrderID]OrderStatus

func NewUnfinishedOrders() UnfinishedOrders {
	return UnfinishedOrders{}
}
