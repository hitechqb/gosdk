package sdk

import (
	"fmt"
	"github.com/hitech/gosdk/sdkcm"
	"github.com/labstack/echo/v4"
 	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	echo.Context
	logger *logrus.Logger
 	cf     sdkcm.SDKConfig
}

func NewSDKContext(
	ctx echo.Context,
	logger *logrus.Logger,
	cf sdkcm.SDKConfig,
) *Context {
	return &Context{
		Context: ctx,
		logger:  logger,
		cf:      cf,
	}
}

func GetSDKContext(c echo.Context) *Context {
	return c.(*Context)
}

// Get logger from context.
func (s *Context) GetLogger() *logrus.Logger {
	return s.logger
}

// Get current_page from query param.
func (s *Context) GetCurrentPage() int {
	currentPage, _ := strconv.Atoi(s.QueryParam("page"))
	if currentPage <= 0 {
		currentPage = 1
	}
	return currentPage
}

// Get per_page from Env and can be override by query param.
func (s *Context) GetPerPage() int {

	perPage, _ := strconv.Atoi(s.QueryParam("limit"))
	if perPage <= 0 {
		perPage = s.cf.ItemsPerPage()
		if perPage <= 0 {
			perPage = 10
		}
	}

	return perPage
}

// Calculate paging from.
func (s *Context) GetPagingFrom(total int) int {
	from := s.GetSkip() + 1
	if from > total {
		return total
	}

	return from
}

// Calculate paging to.
func (s *Context) GetPagingTo(total int) int {
	if total == 0 {
		return 0
	}

	from := s.GetPagingFrom(total)
	perPage := s.GetPerPage()

	to := from + perPage - 1

	if to > total {
		return total
	}

	if to < 0 {
		return 0
	}

	return to
}

// Calculate paging totalPages.
func (s *Context) GetTotalPages(total int) int {
	perPage := s.GetPerPage()
	return (total + perPage - 1) / perPage
}

func (s *Context) GetSkip() int {
	currentPage := s.GetCurrentPage()
	perPage := s.GetPerPage()
	skip, _ := strconv.Atoi(s.QueryParam("skip"))

	if currentPage >= 1 && skip > 0 {
		x := (currentPage - 1) * perPage
		return x + skip
	}

	return (currentPage - 1) * perPage
}

func (s *Context) SetPagingFilter(filter sdkcm.IPagingFilter) {
	filter.SetPerPage(s.GetPerPage())
	filter.SetSkip(s.GetSkip())
}

// Success response with paging object.
func (s *Context) SuccessWithPagingJSON(input, data interface{}, total int) error {
	currentPage := s.GetCurrentPage()
	perPage := s.GetPerPage()
	totalPages := s.GetTotalPages(total)

	link, err := url.Parse(fmt.Sprintf(`%s://localhost:%s%s`, s.cf.AppURLScheme(), s.cf.AppPort(), s.Context.Request().URL.Path))
	if err != nil {
		return err
	}

	var firstPageURL *string
	var nextPageURL *string
	var prevPageURL *string

	if totalPages > 0 {
		firstPageQuery := link.Query()
		firstPageQuery.Set("page", "1")
		firstPage := link
		firstPage.RawQuery = firstPageQuery.Encode()
		firstPageURLString := firstPage.String()
		firstPageURL = &firstPageURLString
	}

	lastPageQuery := link.Query()
	lastPageQuery.Set("page", strconv.Itoa(totalPages))
	lastPage := link
	lastPage.RawQuery = lastPageQuery.Encode()
	lastPageURLString := lastPage.String()
	lastPageURL := &lastPageURLString

	if currentPage < totalPages {
		nextPage := link
		nextPageQuery := nextPage.Query()
		nextPageQuery.Set("page", strconv.Itoa(currentPage+1))
		nextPage.RawQuery = nextPageQuery.Encode()
		nextPageString := nextPage.String()
		nextPageURL = &nextPageString
	}

	if currentPage > 1 {
		prevPage := link
		prevPageQuery := prevPage.Query()
		prevPageQuery.Set("page", strconv.Itoa(currentPage-1))
		prevPage.RawQuery = prevPageQuery.Encode()
		prevPageString := prevPage.String()
		prevPageURL = &prevPageString
	}

	from := s.GetPagingFrom(total)
	to := s.GetPagingTo(total)

	return s.JSON(http.StatusOK, sdkcm.SuccessResponse(input, data, &sdkcm.Paging{
		Total:        total,
		CurrentPage:  currentPage,
		PerPage:      perPage,
		LastPage:     totalPages,
		FirstPageURL: firstPageURL,
		LastPageURL:  lastPageURL,
		NextPageURL:  nextPageURL,
		PrevPageURL:  prevPageURL,
		From:         from,
		To:           to,
	}))
}

// Success response with data
func (s *Context) SuccessJSON(input, data interface{}) error {
	return s.JSON(http.StatusOK, sdkcm.SuccessResponse(input, data, nil))
}

// Error response with messages...
func (s *Context) ErrorJSON(statusCode int, input interface{}, rootCause error, messages ...string) error {
	return s.JSON(statusCode, sdkcm.ErrorResponse(statusCode, input, rootCause, messages...))
}

// Error bad request response with messages...
func (s *Context) ErrorBadRequestJSON(input interface{}, rootCause error, messages ...string) error {
	return s.JSON(http.StatusBadRequest, sdkcm.ErrorResponse(http.StatusBadRequest, input, rootCause, messages...))
}

// Error notfound response with messages...
func (s *Context) ErrorNotFoundJSON(input interface{}, rootCause error, messages ...string) error {
	return s.JSON(http.StatusNotFound, sdkcm.ErrorResponse(http.StatusNotFound, input, rootCause, messages...))
}

// Error conflict response with messages...
func (s *Context) ErrorConflictJSON(input interface{}, rootCause error, messages ...string) error {
	return s.JSON(http.StatusConflict, sdkcm.ErrorResponse(http.StatusConflict, input, rootCause, messages...))
}

// Error conflict response with messages...
func (s *Context) ErrorUnAuthorizedJSON(input interface{}, rootCause error, messages ...string) error {
	return s.JSON(http.StatusUnauthorized, sdkcm.ErrorResponse(http.StatusUnauthorized, input, rootCause, messages...))
}

func (s *Context) ErrorTokenInvalidOrExpired(rootCause error) error {
	return s.ErrorUnAuthorizedJSON(http.StatusUnauthorized, rootCause, "Token is invalid or expired")
}

func (s *Context) HandleError(input interface{}, err error) error {
	if e, ok := err.(sdkcm.SDKError); ok {
		switch e.StatusCode {
		case http.StatusBadRequest:
			return s.ErrorBadRequestJSON(input, e.RootCause, e.Error())
		case http.StatusNotFound:
			return s.ErrorNotFoundJSON(input, e.RootCause, e.Error())
		case http.StatusConflict:
			return s.ErrorConflictJSON(input, e.RootCause, e.Error())
		case http.StatusUnauthorized:
			return s.ErrorUnAuthorizedJSON(input, e.RootCause, e.Error())
		case http.StatusUnprocessableEntity:
			return s.ErrorJSON(http.StatusUnprocessableEntity, input, e.RootCause, e.Error())
		default:
			return s.ErrorJSON(http.StatusInternalServerError, input, e.RootCause, e.Error())
		}
	}

	return s.ErrorJSON(http.StatusInternalServerError, input, err, err.Error())
}
