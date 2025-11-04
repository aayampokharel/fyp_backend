package authentication

import (
	"project/internals/domain/entity"
	"project/internals/usecase"
	err "project/package/errors"
	"project/package/utils/common"
	"time"
)

type Controller struct {
	useCase usecase.SqlUseCase
}

func NewController(useCase usecase.SqlUseCase) *Controller {
	return &Controller{useCase: useCase}
}

func (c *Controller) HandleCreateNewInstitution(institution CreateInstitutionRequest) entity.Response {
	// var institution CreateInstitutionRequest

	// fmt.Print("hello hello")
	// //! to be done in middleware .
	// if er := json.NewDecoder(r.Body).Decode(&institution); er != nil {
	// 	c.useCase.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
	// 	return common.HandleErrorResponse(500, err.ErrDecodingJSONString, er)

	// }

	if institution.InstitutionName == "" || institution.ToleAddress == "" || institution.DistrictAddress == "" || institution.WardNumber == 0 {
		c.useCase.Service.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", err.ErrEmptyInstitutionInfo)
		return common.HandleErrorResponse(500, err.ErrCreatingInstitutionString, err.ErrEmptyInstitutionInfo)

	}
	institutionEntity := institution.ToEntity()

	if er := c.useCase.CheckDuplicationByInstitutionInfoUseCase(institutionEntity); er != nil {
		c.useCase.Service.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		return common.HandleErrorResponse(500, err.ErrDuplicateString, er)

	}

	institutionId, er := c.useCase.InsertInstitutionsUseCase(institutionEntity)
	if er != nil {
		c.useCase.Service.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		return common.HandleErrorResponse(500, err.ErrCreatingInstitutionString, er)

	}
	return common.HandleSuccessResponse(CreateInstutionResponse{
		InstitutionID: institutionId,
		IsActive:      false,
	})

}

func (c *Controller) HandleCreateNewUserAccount(newUserAccount CreateUserAccountRequest) entity.Response {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
	// var newUserAccount CreateUserAccountRequest
	// if er := json.NewDecoder(r.Body).Decode(&newUserAccount); er != nil {
	// 	return common.HandleErrorResponse(500, err.ErrDecodingJSONString, er)

	// }

	newUserAccountEntity := newUserAccount.ToEntity()

	createdAtStr, userAccountID, er := c.useCase.InsertUserAccountUseCase(newUserAccountEntity, newUserAccount.InstitutionID, newUserAccount.InstitutionLogoBase64)

	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCreatingUserAccountString, er)

	}
	createdAt, er := time.Parse(time.RFC3339, createdAtStr)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCreatingUserAccountString, er)

	}
	return common.HandleSuccessResponse(CreateUserAccountResponse{
		UserAccountID: userAccountID,
		CreatedAt:     createdAt.Format(time.RFC3339),
	})

}
func (c *Controller) HandleCreateNewFaculty(newFaculty CreateFacultyRequest) (*entity.Institution, entity.Response) {
	newFacultyEntity := newFaculty.ToEntity()
	facultyID, institutionInfo, er := c.useCase.InsertFacultyAndRetrieveInstitutionUseCase(newFacultyEntity)

	if er != nil {
		return nil, common.HandleErrorResponse(500, err.ErrCreatingInstitutionFacultyString, er)

	}

	return institutionInfo, common.HandleSuccessResponse(CreateFacultyResponse{
		InstitutionFacultyID: facultyID,
	})

}

func (c *Controller) HandleCheckInstitutionIsActive(request map[string]string) entity.Response {
	//esma , first email ra password college ko manche le insert garcha , I will compulsorily return list of all institutions associated to the user  .BUT THAT IS HANDLED BY ANOTHER ENDPOINT , as after this that person can choose from multiple list 1 institution , ani tyo select garepachi we hit this associated endpoint and use this controller .
	requestMap, er := common.CheckMapKeysReturnValues(request, CheckInstitutionIsActiveQuery)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrParsingQueryParametersString, er)
	}
	institutionID := requestMap[InstitutionID]

	institutionInfo, er := c.useCase.GetInstitutionInfoUseCase(institutionID)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCheckingIsActiveString, er)
	}
	var checkInstitutionIsActiveResponse CheckInstitutionIsActiveResponse

	return common.HandleSuccessResponse(checkInstitutionIsActiveResponse.FromEntity(institutionInfo))
}
