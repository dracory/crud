package crud

import (
	"fmt"
	"strconv"

	"github.com/dracory/hb"
)

func (crud *Crud) renderPagination(currentPage int, totalRows int64, baseURL string) string {
	if crud.pageSize <= 0 {
		return ""
	}

	totalPages := int((totalRows + int64(crud.pageSize) - 1) / int64(crud.pageSize))
	if totalPages <= 1 {
		return ""
	}

	nav := hb.Nav().Attr("aria-label", "Page navigation")
	ul := hb.UL().Class("pagination justify-content-center mt-3")

	// Previous button
	prevDisabled := ""
	prevPage := currentPage - 1
	if currentPage <= 1 {
		prevDisabled = " disabled"
		prevPage = 1
	}
	ul.Child(
		hb.LI().Class("page-item" + prevDisabled).Child(
			hb.Hyperlink().Class("page-link").
				Href(baseURL + "&page=" + strconv.Itoa(prevPage)).
				HTML("&laquo;"),
		),
	)

	// Page numbers
	for i := 1; i <= totalPages; i++ {
		activeClass := ""
		if i == currentPage {
			activeClass = " active"
		}
		ul.Child(
			hb.LI().Class("page-item" + activeClass).Child(
				hb.Hyperlink().Class("page-link").
					Href(baseURL + "&page=" + strconv.Itoa(i)).
					HTML(fmt.Sprintf("%d", i)),
			),
		)
	}

	// Next button
	nextDisabled := ""
	nextPage := currentPage + 1
	if currentPage >= totalPages {
		nextDisabled = " disabled"
		nextPage = totalPages
	}
	ul.Child(
		hb.LI().Class("page-item" + nextDisabled).Child(
			hb.Hyperlink().Class("page-link").
				Href(baseURL + "&page=" + strconv.Itoa(nextPage)).
				HTML("&raquo;"),
		),
	)

	nav.Child(ul)
	return nav.ToHTML()
}
