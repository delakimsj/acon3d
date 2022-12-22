package framework

type RBACItem struct {
	FullPath string
	Role     string
}

func GetRBACMatrix() *map[RBACItem]bool {
	matrix := make(map[RBACItem]bool)

	matrix[RBACItem{FullPath: "GET /admin/product", Role: "author"}] = true
	matrix[RBACItem{FullPath: "POST /admin/product", Role: "author"}] = true
	matrix[RBACItem{FullPath: "GET /admin/product/:product_id", Role: "author"}] = true
	matrix[RBACItem{FullPath: "PUT /admin/product/:product_id", Role: "author"}] = true
	matrix[RBACItem{FullPath: "PATCH /admin/product/:product_id", Role: "author"}] = false // editor's approval
	matrix[RBACItem{FullPath: "GET /admin/deal", Role: "author"}] = true
	matrix[RBACItem{FullPath: "POST /admin/deal", Role: "author"}] = false           // only editor can
	matrix[RBACItem{FullPath: "PUT /admin/deal/:deal_id", Role: "author"}] = false   // only editor can
	matrix[RBACItem{FullPath: "PATCH /admin/deal/:deal_id", Role: "author"}] = false // only editor can
	matrix[RBACItem{FullPath: "POST /admin/request_deal_modification", Role: "author"}] = true

	matrix[RBACItem{FullPath: "GET /admin/product", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "POST /admin/product", Role: "editor"}] = false // only author create new product
	matrix[RBACItem{FullPath: "GET /admin/product/:product_id", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "PUT /admin/product/:product_id", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "PATCH /admin/product/:product_id", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "GET /admin/deal", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "POST /admin/deal", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "PUT /admin/deal/:deal_id", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "PATCH /admin/deal/:deal_id", Role: "editor"}] = true
	matrix[RBACItem{FullPath: "POST /admin/request_deal_modification", Role: "editor"}] = true

	return &matrix
}
