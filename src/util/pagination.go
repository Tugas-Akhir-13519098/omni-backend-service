package util

import (
	"math"
	"strings"
)

// Pagination struct model to construct pagination key value
type Pagination struct {
	// page number of data
	Page int `json:"page"`
	// amount per page of data
	PageSize int `json:"page_size"`
	// numbers of pages
	TotalPage int `json:"total_page"`
	// total existing data in database
	Total int `json:"total"`
	// sort data in database
	Orders []*Order `json:"-"`
}

type Order struct {
	ColumnName string
	Direction  Direction
}

// Direction shows the sort direction
type Direction string

func (d Direction) String() string {
	return string(d)
}

const (
	// DirectionAsc sorts by ascending order
	DirectionAsc Direction = "ASC"
	// DirectionDesc sorts by descending order
	DirectionDesc Direction = "DESC"
)

const (
	// DefaultPage use for default page if user doesn't assign page value
	DefaultPage = 1
	// DefaultPageSize use for default page size if user doesn't assign pageSize value
	DefaultPageSize = 10
	// MaxPageSize use for limit of rows data each page
	MaxPageSize = 100
)

// Parse to pagination type
func Parse(page, pageSize int) *Pagination {
	return &Pagination{
		Page:      page,
		PageSize:  pageSize,
		TotalPage: 0,
		Total:     0,
		Orders:    nil,
	}
}

// SetPagination utils function to set pagination response
func (p *Pagination) SetPagination() {
	if p.Total > 0 && p.Total < p.PageSize {
		p.PageSize = p.Total
	}
	p.TotalPage = int(math.Ceil(float64(p.Total) / float64(p.PageSize)))
}

// ValidatePagination utils function to validate pagination values request
func (p *Pagination) ValidatePagination() {
	// check is inputed page and pageSize acceptable
	if p.Page <= 0 || p.PageSize <= 0 {
		p.SetToDefault()
		return
	}
	// check is inputed pageSize exceed of maximal pageSize
	if p.PageSize > MaxPageSize {
		p.SetToDefault()
		return
	}
	return
}

// SetToDefault to set default pagination
func (p *Pagination) SetToDefault() {
	p.Page, p.PageSize = DefaultPage, DefaultPageSize
}

// GetOffset to get offset value of pagination
// this function will validate pagination value, then calculate offset
func (p *Pagination) GetOffset() int {
	p.ValidatePagination()
	return (p.Page - 1) * p.PageSize
}

// SetOrder single set order of pagination
func (p *Pagination) SetOrder(order *Order) {
	if order == nil {
		return
	}
	p.Orders = append(p.Orders, order)
}

// SetOrders of pagination
func (p *Pagination) SetOrders(orders []*Order) {
	if orders == nil {
		return
	}
	for _, order := range orders {
		p.SetOrder(order)
	}
}

func (p *Pagination) GetOrderQuery() (query string, ok bool) {
	if len(p.Orders) == 0 {
		return
	}
	if p.Orders[0] == nil {
		return
	}
	if p.Orders[0].ColumnName == "" {
		return
	}
	if p.Orders[0].Direction.String() == "" {
		return
	}
	query = p.Orders[0].ColumnName + " " + p.Orders[0].Direction.String()
	ok = true
	return
}

func (p *Pagination) GetOrdersQuery() (query string, ok bool) {
	if len(p.Orders) == 0 {
		return
	}
	typeOrders := make(map[string][]string)
	for _, order := range p.Orders {
		if order == nil {
			return
		}
		typeOrders[order.Direction.String()] = append(typeOrders[order.Direction.String()], order.ColumnName)
	}
	var queries []string
	for direction, orders := range typeOrders {
		queries = append(queries, strings.Join(orders, ",")+" "+direction)
	}
	query = strings.Join(queries, ",")
	return
}
