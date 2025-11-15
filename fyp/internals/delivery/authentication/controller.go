package authentication

import (
	"project/internals/domain/entity"
	"project/internals/usecase"
	"project/package/enum"
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
		return common.HandleErrorResponse(400, err.ErrCreatingInstitutionString, err.ErrEmptyInstitutionInfo)

	}
	institutionEntity := institution.ToEntity()

	if er := c.useCase.CheckDuplicationByInstitutionInfoUseCase(institutionEntity); er != nil {
		c.useCase.Service.Logger.Errorln("[authentication_controller] Error: HandleCreateNewInstitution::", er)
		return common.HandleErrorResponse(409, err.ErrDuplicateString, er)

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

	newUserAccountEntity, er := newUserAccount.ToEntity()
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCreatingUserAccountString, er)
	}

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
		if er == err.ErrInstitutionAlreadyVerified {
			//nil because we don't need to broadcast this information.
			return nil, common.HandleSuccessResponse(CreateFacultyResponse{
				InstitutionFacultyID: facultyID,
			})
		}
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
		return common.HandleErrorResponse(400, err.ErrParsingQueryParametersString, er)
	}
	institutionID := requestMap[InstitutionID]

	institutionInfo, er := c.useCase.GetInstitutionInfoUseCase(institutionID)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCheckingIsActiveString, er)
	}
	facultyListInfo, er := c.useCase.GetFacultyListUseCase(institutionID)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrCheckingIsActiveString, er)
	}
	var checkInstitutionIsActiveResponse CheckInstitutionIsActiveResponse

	return common.HandleSuccessResponse(checkInstitutionIsActiveResponse.FromEntity(institutionInfo, facultyListInfo))
}

func (c *Controller) HandleInstitutionsLogin(request InstitutionLoginRequest) entity.Response {

	hashedPassword, _, er := common.HashData(request.Password)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrVerifyingInstituteString, er)
	}
	userID, createdAt, er := c.useCase.VerifyUserLoginUseCase(request.Email, hashedPassword, enum.INSTITUTE)

	if er != nil {
		return common.HandleErrorResponse(401, err.ErrVerifyingInstituteString, er)
	}
	institutionList, er := c.useCase.GetInstitutionsForUserUseCase(userID)
	if er != nil {
		return common.HandleErrorResponse(500, err.ErrVerifyingInstituteString, er)
	}

	return common.HandleSuccessResponse(InstitutionLoginResponse{UserID: userID, CreatedAt: createdAt.Format(time.RFC3339), InstitutionList: institutionList})

}

func (c *Controller) HandleGetFacultiesForInstitutionID(request map[string]string) entity.Response {
	requestMap, er := common.CheckMapKeysReturnValues(request, GetInstitutionFacultiesQuery)
	if er != nil {
		return common.HandleErrorResponse(400, err.ErrParsingQueryParametersString, er)
	}

	institutionID := requestMap[InstitutionID]
	faculties, er := c.useCase.GetFacultiesForInstitutionIDUseCase(institutionID)
	if er != nil {
		return common.HandleErrorResponse(500, er.Error(), er)
	}
	return common.HandleSuccessResponse(faculties)

}
