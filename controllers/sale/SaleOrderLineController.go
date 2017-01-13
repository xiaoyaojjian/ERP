package sale

import (
	"encoding/json"
	"goERP/controllers/base"
	md "goERP/models"
	"strconv"
	"strings"
)

// SaleOrderLineController
type SaleOrderLineController struct {
	base.BaseController
}

// Post request
func (ctl *SaleOrderLineController) Post() {
	action := ctl.Input().Get("action")
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}

// Put request
func (ctl *SaleOrderLineController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/sale/order/line/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if orderLine, err := md.GetSaleOrderLineByID(idInt64); err == nil {
			if err := ctl.ParseForm(&orderLine); err == nil {

				if err := md.UpdateSaleOrderLineByID(orderLine); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)

}

// Get request
func (ctl *SaleOrderLineController) Get() {
	ctl.PageName = "销售订单明细管理"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()

	}
	ctl.Data["PageName"] = ctl.PageName + "\\" + ctl.PageAction
	ctl.URL = "/sale/order/line/"
	ctl.Data["URL"] = ctl.URL

	ctl.Data["MenuSaleOrderLineActive"] = "active"
}

// Edit edit sale orde line
func (ctl *SaleOrderLineController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	orderLineInfo := make(map[string]interface{})
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if orderLine, err := md.GetSaleOrderLineByID(idInt64); err == nil {
				ctl.PageAction = orderLine.Name
				orderLineInfo["name"] = orderLine.Name
			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["orderLine"] = orderLineInfo
	ctl.Layout = "base/base.html"
	ctl.TplName = "sale/sale_order_ine_form.html"
}

// Create display create page
func (ctl *SaleOrderLineController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"
	ctl.Layout = "base/base.html"
	ctl.TplName = "sale/sale_order_line_form.html"
}

// Detail display sale order line info
func (ctl *SaleOrderLineController) Detail() {
	//获取信息一样，直接调用Edit
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

// PostCreate post request create sale order line
func (ctl *SaleOrderLineController) PostCreate() {
	orderLine := new(md.SaleOrderLine)
	if err := ctl.ParseForm(orderLine); err == nil {

		if id, err := md.AddSaleOrderLine(orderLine); err == nil {
			ctl.Redirect("/sale/order/line/"+strconv.FormatInt(id, 10)+"?action=detail", 302)
		} else {
			ctl.Get()
		}
	} else {
		ctl.Get()
	}
}

// Validator js valid
func (ctl *SaleOrderLineController) Validator() {
	name := ctl.GetString("name")
	name = strings.TrimSpace(name)
	recordID, _ := ctl.GetInt64("recordId")
	result := make(map[string]bool)
	obj, err := md.GetSaleOrderLineByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.ID {
				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

//SaleOrderLineList 获得符合要求的数据
func (ctl *SaleOrderLineController) SaleOrderLineList(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.SaleOrderLine
	paginator, arrs, err := md.GetAllSaleOrderLine(query, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		//使用多线程来处理数据，待修改
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["name"] = line.Name
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID
			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}

// PostList post request list sale order line
func (ctl *SaleOrderLineController) PostList() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	if result, err := ctl.SaleOrderLineList(query, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

// GetList display sale order line with list
func (ctl *SaleOrderLineController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-sale-order-line"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "sale/sale_order_line_list_search.html"
}