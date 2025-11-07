package category

import (
	"log"
	"project/internals/domain/entity"
	"project/internals/usecase"
	err "project/package/errors"
	"project/package/utils/common"
)

type Controller struct {
	sqlUseCase *usecase.SqlUseCase
}

func NewController(sqlUseCase *usecase.SqlUseCase) *Controller {
	return &Controller{sqlUseCase: sqlUseCase}
}

func (c *Controller) HandleCreatePDFCategory(request CreatePDFCategoryDto) entity.Response {
	pdfFileCategory, er := request.ToPdfFileCategoryEntity()
	if er != nil {
		log.Println(er)
		return common.HandleErrorResponse(401, er.Error(), er)
	}
	insertedpdfFileCategory, er := c.sqlUseCase.InsertAndGetPDFCategoryUseCase(pdfFileCategory)
	if er != nil {
		return common.HandleErrorResponse(401, er.Error(), er)
	}
	return common.HandleSuccessResponse(CreatePDFCategoryResponseDto{
		CategoryID:   insertedpdfFileCategory.CategoryID,
		CategoryName: insertedpdfFileCategory.CategoryName,
	})
}

func (c *Controller) HandleGetPDFCategoriesList(request map[string]string) entity.Response {
	requestMap, er := common.CheckMapKeysReturnValues(request, GetPDFCategoryRequestDtoQuery)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrParsingQueryParametersString, er)
	}
	institutionFacultyID := requestMap[InstitutionFacultyID]
	institutionID := requestMap[InstitutionID]
	pdfFileCategories, er := c.sqlUseCase.GetPDFCategoriesListUseCase(institutionID, institutionFacultyID)
	if er != nil {
		return common.HandleErrorResponse(401, er.Error(), er)
	}
	return common.HandleSuccessResponse(pdfFileCategories)

}
