package pagination

import (
	"gorm.io/gorm"
	"math"
)

// Page 单个分页元素
type Page struct {
	// 链接
	URL string
	// 页码
	Number uint64
}

// ViewData 同视图渲染的数据
type ViewData struct {
	// 是否需要显示分页
	HasPages bool

	// 下一页
	Next    Page
	HasNext bool

	// 上一页
	Prev    Page
	HasPrev bool

	Current Page

	// 数据库的内容总数量
	TotalCount uint64
	// 总页数
	TotalPage uint64
}

// Pagination 分页对象
type Pagination struct {
	PerPage uint32
	Page    uint64
	Count   uint64
	DB      *gorm.DB
}

// New 分页对象构建器
// db —— GORM 查询句柄，用以查询数据集和获取数据总数
// page —— page
// perPage —— 每页条数，传参为小于或者等于 0 时为默认值  10
func New(db *gorm.DB, page uint64, perPage uint32) *Pagination {
	// 默认每页数量
	if perPage <= 0 {
		perPage = 10
	}

	// 实例对象
	p := &Pagination{
		DB:      db,
		PerPage: perPage,
		Page:    page,
		Count:   0,
	}

	// 设置当前页码
	p.SetPage(page)

	return p
}

// Paging 返回渲染分页所需的数据
func (p *Pagination) Paging() ViewData {

	return ViewData{
		HasPages: p.HasPages(),

		Next:    p.NewPage(p.NextPage()),
		HasNext: p.HasNext(),

		Prev:    p.NewPage(p.PrevPage()),
		HasPrev: p.HasPrev(),

		Current:   p.NewPage(p.CurrentPage()),
		TotalPage: p.TotalPage(),

		TotalCount: p.Count,
	}
}

// NewPage 设置当前页
func (p Pagination) NewPage(page uint64) Page {
	return Page{
		Number: page,
	}
}

// SetPage 设置当前页
func (p *Pagination) SetPage(page uint64) {
	if page <= 0 {
		page = 1
	}

	p.Page = page
}

// CurrentPage 返回当前页码
func (p Pagination) CurrentPage() uint64 {
	totalPage := p.TotalPage()
	if totalPage == 0 {
		return 0
	}

	if p.Page > totalPage {
		return totalPage
	}

	return p.Page
}

// Results 返回请求数据，请注意 data 参数必须为 GROM 模型的 Slice 对象
func (p Pagination) Results(data interface{}) error {
	var err error
	var offset uint64
	page := p.CurrentPage()
	if page == 0 {
		return err
	}

	if page > 1 {
		offset = (page - 1) * uint64(p.PerPage)
	}

	return p.DB.Debug().Limit(int(p.PerPage)).Offset(int(offset)).Find(data).Error
}

// TotalCount 返回的是数据库里的条数
func (p *Pagination) TotalCount() uint64 {
	if p.Count == 0 {
		var count int64
		if err := p.DB.Count(&count).Error; err != nil {
			return 0
		}
		p.Count = uint64(count)
	}

	return p.Count
}

// HasPages 总页数大于 1 时会返回 true
func (p *Pagination) HasPages() bool {
	n := p.TotalCount()
	return n > uint64(p.PerPage)
}

// HasNext returns true if current page is not the last page
func (p Pagination) HasNext() bool {
	totalPage := p.TotalPage()
	if totalPage == 0 {
		return false
	}

	page := p.CurrentPage()
	if page == 0 {
		return false
	}

	return page < totalPage
}

// PrevPage 前一页码，0 意味着这就是第一页
func (p Pagination) PrevPage() uint64 {
	hasPrev := p.HasPrev()

	if !hasPrev {
		return 0
	}

	page := p.CurrentPage()
	if page == 0 {
		return 0
	}

	return page - 1
}

// NextPage 下一页码，0 的话就是最后一页
func (p Pagination) NextPage() uint64 {
	hasNext := p.HasNext()
	if !hasNext {
		return 0
	}

	page := p.CurrentPage()
	if page == 0 {
		return 0
	}

	return page + 1
}

// HasPrev 如果当前页不为第一页，就返回 true
func (p Pagination) HasPrev() bool {
	page := p.CurrentPage()
	if page == 0 {
		return false
	}

	return page > 1
}

// TotalPage 返回总页数
func (p Pagination) TotalPage() uint64 {
	count := p.TotalCount()
	if count == 0 {
		return 0
	}

	nums := int64(math.Ceil(float64(count) / float64(p.PerPage)))
	if nums == 0 {
		nums = 1
	}

	return uint64(nums)
}
