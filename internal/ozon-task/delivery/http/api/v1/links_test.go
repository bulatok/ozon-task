package v1

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	"github.com/bulatok/ozon-task/internal/ozon-task/config"
	"github.com/bulatok/ozon-task/internal/ozon-task/models"
	api_models "github.com/bulatok/ozon-task/internal/ozon-task/models/api/v1"
	"github.com/bulatok/ozon-task/internal/ozon-task/store"
	"github.com/bulatok/ozon-task/internal/ozon-task/usecase"
)

const (
	publicUrlLinksTest = "http://127.0.0.1"

	originalLinkTest = "https://ya.ru/"
)

var (
	linkTest              *models.Link
	linkTestOtherOriginal *models.Link
)

func prepareLinkTest(t *testing.T) {
	linkTest = new(models.Link)
	linkTestOtherOriginal = new(models.Link)

	(*linkTest).Original = originalLinkTest
	(*linkTestOtherOriginal).Original = originalLinkTest + "blabla"

	err := linkTest.SetShortLink(publicUrlLinksTest)
	utils.AssertEqual(t, err, nil, "link inti-lazing")

	(*linkTestOtherOriginal).Short = linkTest.Short

}

type linksSetup struct {
	app  *fiber.App
	l    *zap.Logger
	repo *store.MockLinksRepo
}

func newLinksSetup(t *testing.T) linksSetup {
	ctrl := gomock.NewController(t)
	repo := store.NewMockLinksRepo(ctrl)

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	l, _ := zap.NewDevelopment()
	conf := &config.Config{
		Service: &config.Service{
			PublicUrl: publicUrlLinksTest,
		},
	}

	useCases := UseCases{usecase.ProvideLinks(repo, l)}

	srv := &ApiServer{
		uc:   useCases,
		conf: conf,
		l:    l,
	}

	srv.Init(app)

	return linksSetup{
		repo: repo,
		app:  app,
		l:    l,
	}
}

func TestApiServer_links(t *testing.T) {
	prepareLinkTest(t)

	type response struct {
		StatusCode    int
		NeedCheckBody bool
		Body          string
	}

	type request struct {
		Target string
		Body   io.Reader
		Method string
	}

	tests := []struct {
		name        string
		prepareRepo func(r *store.MockLinksRepo)
		request     request
		response    response
	}{
		{
			request: request{
				Body:   nil,
				Target: "/" + linkTest.GetUnderlineHash(),
				Method: fiber.MethodGet,
			},
			response: response{
				StatusCode: fiber.StatusOK,
				Body:       testLinkToGetLinkResponse(t),
			},
			name: "GET. /:linkHash. OK",
			prepareRepo: func(r *store.MockLinksRepo) {
				r.EXPECT().
					Get(publicUrlLinksTest+"/"+linkTest.GetUnderlineHash()).
					Return(linkTest, nil)
			},
		},
		{
			request: request{
				Body:   nil,
				Target: "/" + linkTest.GetUnderlineHash(),
				Method: fiber.MethodGet,
			},
			response: response{
				StatusCode: fiber.StatusInternalServerError,
			},
			name: "GET. /:linkHash. database error",
			prepareRepo: func(r *store.MockLinksRepo) {
				r.EXPECT().
					Get(publicUrlLinksTest+"/"+linkTest.GetUnderlineHash()).
					Return(nil, assert.AnError)
			},
		},
		{
			request: request{
				Body:   nil,
				Target: "/",
				Method: fiber.MethodGet,
			},
			response: response{
				StatusCode: fiber.StatusNotFound,
			},
			name:        "GET. /. not found",
			prepareRepo: func(r *store.MockLinksRepo) {},
		},
		{
			request: request{
				Body:   nil,
				Target: "/" + linkTest.GetUnderlineHash(),
				Method: fiber.MethodGet,
			},
			response: response{
				StatusCode: fiber.StatusGone,
			},
			name: "GET. /:linkHash. not found in database error",
			prepareRepo: func(r *store.MockLinksRepo) {
				r.EXPECT().
					Get(publicUrlLinksTest+"/"+linkTest.GetUnderlineHash()).
					Return(nil, models.ErrNotFound)
			},
		},
		{
			request: request{
				Body:   strings.NewReader(testLinkToNewLinkRequestRequest(t)),
				Target: "/new",
				Method: fiber.MethodPost,
			},
			response: response{
				StatusCode:    fiber.StatusOK,
				Body:          testLinkToNewLinkRequestResponse(t),
				NeedCheckBody: true,
			},
			name: "POST. /new. OK",
			prepareRepo: func(r *store.MockLinksRepo) {
				gomock.InOrder(
					r.EXPECT().
						Get(publicUrlLinksTest+"/"+linkTest.GetUnderlineHash()).
						Return(nil, models.ErrNotFound),
					r.EXPECT().
						Save(linkTest).
						Return(nil),
				)
			},
		},
		{
			request: request{
				Body:   strings.NewReader(testLinkToNewLinkRequestRequest(t)),
				Target: "/new",
				Method: fiber.MethodPost,
			},
			response: response{
				StatusCode: fiber.StatusInternalServerError,
			},
			name: "POST. /new. internal error",
			prepareRepo: func(r *store.MockLinksRepo) {
				gomock.InOrder(
					r.EXPECT().
						Get(publicUrlLinksTest+"/"+linkTest.GetUnderlineHash()).
						Return(nil, assert.AnError),
				)
			},
		},
		{
			request: request{
				Body:   strings.NewReader(testLinkToNewLinkRequestRequest(t) + "blabla"),
				Target: "/new",
				Method: fiber.MethodPost,
			},
			response: response{
				StatusCode: fiber.StatusBadRequest,
			},
			name: "POST. /new. invalid body of request",
			prepareRepo: func(r *store.MockLinksRepo) {
			},
		},
		{
			request: request{
				Body:   strings.NewReader(testLinkToNewLinkRequestRequest(t)),
				Target: "/new",
				Method: fiber.MethodPost,
			},
			response: response{
				StatusCode:    fiber.StatusOK,
				Body:          testLinkToNewLinkRequestResponse(t),
				NeedCheckBody: true,
			},
			name: "POST. /new. the same original link",
			prepareRepo: func(r *store.MockLinksRepo) {
				gomock.InOrder(
					r.EXPECT().
						Get(publicUrlLinksTest+"/"+linkTest.GetUnderlineHash()).
						Return(linkTest, nil),
				)
			},
		},
		{
			request: request{
				Body:   strings.NewReader(testLinkToNewLinkRequestRequest(t)),
				Target: "/new",
				Method: fiber.MethodPost,
			},
			response: response{
				StatusCode: fiber.StatusInternalServerError,
			},
			name: "POST. /new. the same equal hash condition",
			prepareRepo: func(r *store.MockLinksRepo) {
				gomock.InOrder(
					r.EXPECT().
						Get(publicUrlLinksTest+"/"+linkTest.GetUnderlineHash()).
						Return(linkTestOtherOriginal, nil),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := newLinksSetup(t)
			tt.prepareRepo(setup.repo)

			req := httptest.NewRequest(tt.request.Method, tt.request.Target, tt.request.Body)

			resp, err := setup.app.Test(req)
			if err != nil {
				t.Fatal(tt.name, err)
			}
			defer resp.Body.Close()
			respBody, _ := io.ReadAll(resp.Body)

			utils.AssertEqual(t, err, nil, "request error")
			utils.AssertEqual(t, tt.response.StatusCode, resp.StatusCode, "response status code")
			if tt.response.NeedCheckBody {
				utils.AssertEqual(t, tt.response.Body, string(respBody), "response body")
			}

		})
	}
}

func testLinkToGetLinkResponse(t *testing.T) string {
	resp := api_models.ApiGetOriginalLinkResponse{
		OriginalLink: linkTest.Original,
	}

	d, err := json.Marshal(resp)
	utils.AssertEqual(t, err, nil, "test link marshalling")

	return string(d)
}

func testLinkToNewLinkRequestRequest(t *testing.T) string {
	resp := api_models.ApiNewLinkRequest{
		OriginalLink: linkTest.Original,
	}

	d, err := json.Marshal(resp)
	utils.AssertEqual(t, err, nil, "test link marshalling")

	return string(d)
}

func testLinkToNewLinkRequestResponse(t *testing.T) string {
	resp := api_models.ApiNewLinkResponse{
		ShortLink: linkTest.Short,
	}

	d, err := json.Marshal(resp)
	utils.AssertEqual(t, err, nil, "test link marshalling")

	return string(d)
}
