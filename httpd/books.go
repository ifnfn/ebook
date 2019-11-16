package httpd

import (
	"ebook/util"
	"net/http"

	"github.com/ifnfn/util/system"
	"github.com/labstack/echo"
)

func (h *Httpd) getBook(c echo.Context) error {
	bookID := c.Param("bookID")
	userID := c.QueryParam("id")

	println(bookID, userID)

	if book, found := h.kindle.BooksRobot.Get(bookID); found {
		pargs := make(map[string]interface{})
		pargs["book"] = book
		pargs["userID"] = userID
		if account := h.kindle.GetAcccount(userID); account != nil {
			pargs["favor"] = account.GetFavor(bookID)
			pargs["download"] = account.GetHistory(bookID)
			pargs["email"] = account.Email != ""
		}

		// system.PrintInterface(pargs)
		html := system.GetViewHTML("template/book.html", pargs)

		return c.HTML(http.StatusOK, html)
	}

	return c.NoContent(http.StatusNotFound)
}

func (h *Httpd) pushKindle(c echo.Context) error {
	bookID := c.QueryParam("book_id")
	userID := c.QueryParam("user_id")

	println(bookID, userID)
	if account := h.kindle.GetAcccount(userID); account != nil {
		println(account.Email)
		if account.Email != "" && account.GetHistory(bookID) == false {
			if bok, found := h.kindle.BooksRobot.Get(bookID); found {
				go func() {
					if util.Sendkindle(account.Email, bok) {
						account.AddHistory(bookID)
						h.kindle.Save()
					}
				}()
			}
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *Httpd) favorite(c echo.Context) error {
	bookID := c.QueryParam("book_id")
	userID := c.QueryParam("user_id")

	println(bookID, userID)
	if account := h.kindle.GetAcccount(userID); account != nil {
		if account.GetFavor(bookID) == false {
			if _, found := h.kindle.BooksRobot.Get(bookID); found {
				account.AddFavor(bookID)
				h.kindle.Save()
			}
		}
	}

	return c.NoContent(http.StatusOK)
}
