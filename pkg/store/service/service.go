package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	customError "github.com/ashwingopalsamy/backend-services/pkg/errors"
	"github.com/ashwingopalsamy/backend-services/pkg/store/proto"
	"github.com/ashwingopalsamy/backend-services/pkg/store/validator"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type StoreServiceServer struct {
	proto.UnimplementedStoreServiceServer
}

func (s *StoreServiceServer) CreateStore(ctx context.Context, req *proto.CreateStoreRequest) (*proto.CreateStoreResponse, error) {
	// Validate the request payload
	if validationErrs := validateCreateStoreRequest(req); validationErrs != nil {
		return nil, customError.BadRequest(validationErrs.Error(), validationErrs.Errors)
	}

	// TODO: Database Interactions
	// Check for duplicate store name
	//if s.isDuplicateStore(req.BasicInformation.StoreName) {
	//	return s.buildErrorResponse(409, "A store with this name already exists.", nil)
	//}

	storeID, _ := uuid.NewV7()

	// TODO: Database Interactions
	// Insert store data into the database
	//storeID, err := s.insertStore(req)
	//if err != nil {
	//	return &proto.CreateStoreResponse{
	//		Status:  "error",
	//		Message: fmt.Sprintf("Failed to create store: %v", err),
	//	}, nil
	//}

	// Return the response
	response := &proto.CreateStoreResponse{
		Status:  "success",
		Message: "Store created successfully",
		Data: &proto.StoreData{
			StoreId:   storeID.String(),
			StoreName: req.BasicInformation.StoreName,
			CreatedAt: timestamppb.New(time.Now()),
		},
	}

	return response, nil
}

// Helper to write an error response
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, details map[string]string) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"status":  "error",
		"message": message,
	}
	if details != nil {
		response["details"] = details
	}
	json.NewEncoder(w).Encode(response)
}

// Handle errors returned by the service layer
func handleServiceError(w http.ResponseWriter, err error) {
	var customErr *customError.CustomError
	if errors.As(err, &customErr) {
		writeErrorResponse(w, customErr.StatusCode, customErr.Message, customErr.Details)
	}
}

func validateCreateStoreRequest(req *proto.CreateStoreRequest) *validator.ValidationErrors {
	errors := validator.NewValidationErrors()

	// Validate each section
	validateBasicInformation(req.BasicInformation, errors)
	validateLocation(req.Location, errors)
	validateOperationalHours(req.OperationalHours, errors)
	validateTaxAndPayment(req.TaxAndPayment, errors)
	validateDeliveryConfiguration(req.DeliveryConfiguration, errors)

	// Return errors if any
	if errors.HasErrors() {
		return errors
	}
	return nil
}

func validateBasicInformation(basic *proto.BasicInformation, errors *validator.ValidationErrors) {
	if basic == nil {
		errors.Add("basic_information", "Basic information is required.")
		return
	}

	if len(strings.TrimSpace(basic.StoreName)) == 0 {
		errors.Add("basic_information.store_name", "Store name is required.")
	} else if len(basic.StoreName) > 255 {
		errors.Add("basic_information.store_name", "Store name must be at most 255 characters.")
	}

	if len(strings.TrimSpace(basic.StoreType)) == 0 {
		errors.Add("basic_information.store_type", "Store type is required.")
	}

	if basic.StoreImage != "" && !isValidURLOrBase64(basic.StoreImage) {
		errors.Add("basic_information.store_image", "Store image must be a valid URL or base64 string.")
	}

	if len(strings.TrimSpace(basic.ContactNumber)) == 0 {
		errors.Add("basic_information.contact_number", "Contact number is required.")
	} else if !isValidE164(basic.ContactNumber) {
		errors.Add("basic_information.contact_number", "Contact number must be in E.164 format.")
	}

	if !basic.IsManagerSameAsOwner && len(strings.TrimSpace(basic.ShopManagerId)) == 0 {
		errors.Add("basic_information.shop_manager_id", "Manager ID is required when not the same as owner.")
	}
}

func validateLocation(location *proto.Location, errors *validator.ValidationErrors) {
	if location == nil {
		errors.Add("location", "Location is required.")
		return
	}

	if location.GpsCoordinates == nil {
		errors.Add("location.gps_coordinates", "GPS coordinates are required.")
	} else {
		if location.GpsCoordinates.Latitude < -90 || location.GpsCoordinates.Latitude > 90 {
			errors.Add("location.gps_coordinates.latitude", "Latitude must be between -90 and 90.")
		}
		if location.GpsCoordinates.Longitude < -180 || location.GpsCoordinates.Longitude > 180 {
			errors.Add("location.gps_coordinates.longitude", "Longitude must be between -180 and 180.")
		}
	}

	if len(strings.TrimSpace(location.Area)) == 0 {
		errors.Add("location.area", "Area is required.")
	}

	if len(strings.TrimSpace(location.City)) == 0 {
		errors.Add("location.city", "City is required.")
	}

	if len(strings.TrimSpace(location.Pincode)) == 0 {
		errors.Add("location.pincode", "Pincode is required.")
	}
}

func validateOperationalHours(operationalHours *proto.OperationalHours, errors *validator.ValidationErrors) {
	if operationalHours == nil {
		errors.Add("operational_hours", "Operational hours are required.")
		return
	}

	if !operationalHours.IsOpen_24Hours {
		if len(strings.TrimSpace(operationalHours.OpeningTime)) == 0 || !isValidTimeFormat(operationalHours.OpeningTime) {
			errors.Add("operational_hours.opening_time", "Opening time is required and must be in HH:MM format if the store is not open 24 hours.")
		}

		if len(strings.TrimSpace(operationalHours.ClosingTime)) == 0 || !isValidTimeFormat(operationalHours.ClosingTime) {
			errors.Add("operational_hours.closing_time", "Closing time is required and must be in HH:MM format if the store is not open 24 hours.")
		}
	}

	if operationalHours.IsOwnPickupEnabled && !isAllowedPickupReadyTime(operationalHours.OwnPickupReadyTime) {
		errors.Add("operational_hours.own_pickup_ready_time", "Pickup ready time must be one of the allowed options: '15 minutes', '30 minutes', '1 hour', '2 hours'.")
	}
}

func validateTaxAndPayment(taxAndPayment *proto.TaxAndPayment, errors *validator.ValidationErrors) {
	if taxAndPayment == nil {
		errors.Add("tax_and_payment", "Tax and payment information is required.")
		return
	}

	if taxAndPayment.GstRegistered && len(strings.TrimSpace(taxAndPayment.GstNumber)) == 0 {
		errors.Add("tax_and_payment.gst_number", "GST number is required if the store is GST registered.")
	}

	if taxAndPayment.IsPaymentGatewayEnabled {
		if taxAndPayment.BankInformation == nil {
			errors.Add("tax_and_payment.bank_information", "Bank information is required if the payment gateway is enabled.")
			return
		}

		if len(strings.TrimSpace(taxAndPayment.BankInformation.AccountHolderName)) == 0 {
			errors.Add("tax_and_payment.bank_information.account_holder_name", "Account holder name is required in bank information.")
		}

		if len(strings.TrimSpace(taxAndPayment.BankInformation.BankAccountNumber)) == 0 {
			errors.Add("tax_and_payment.bank_information.bank_account_number", "Bank account number is required in bank information.")
		}

		if len(strings.TrimSpace(taxAndPayment.BankInformation.BankIfscCode)) != 11 {
			errors.Add("tax_and_payment.bank_information.bank_ifsc_code", "IFSC code must be exactly 11 characters in bank information.")
		}
	}
}

func validateDeliveryConfiguration(delivery *proto.DeliveryConfiguration, errors *validator.ValidationErrors) {
	if delivery == nil {
		errors.Add("delivery_configuration", "Delivery configuration is required.")
		return
	}

	if delivery.IsDeliveryEnabled {
		switch delivery.DeliveryLocationType {
		case "Radius":
			if delivery.DeliveryRadiusKm <= 0 {
				errors.Add("delivery_configuration.delivery_radius_km", "Delivery radius must be a positive number if delivery location type is 'Radius'.")
			}
		case "City", "International":
			if len(delivery.DeliveryLocations) == 0 {
				errors.Add("delivery_configuration.delivery_locations", "Delivery locations must be provided if delivery location type is 'City' or 'International'.")
			}
		case "PAN India":
			// No specific validation required for PAN India
		default:
			errors.Add("delivery_configuration.delivery_location_type", "Delivery location type must be one of the allowed options: 'PAN India', 'City', 'Radius', 'International'.")
		}

		if delivery.FreeDeliveryMinOrder < 0 {
			errors.Add("delivery_configuration.free_delivery_min_order", "Minimum order amount for free delivery must be a positive number.")
		}

		if delivery.DeliveryFeeIfMinNotMet < 0 {
			errors.Add("delivery_configuration.delivery_fee_if_min_not_met", "Delivery fee must be a positive number.")
		}
	}
}

func isValidURLOrBase64(input string) bool {
	// Basic validation for URL or base64
	return strings.HasPrefix(input, "http") || len(input) > 0 // Extend as needed
}

func isValidE164(phone string) bool {
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}

func isValidTimeFormat(timeStr string) bool {
	re := regexp.MustCompile(`^(2[0-3]|[01]?[0-9]):([0-5]?[0-9])$`)
	return re.MatchString(timeStr)
}

func isAllowedPickupReadyTime(input string) bool {
	allowed := []string{"15 minutes", "30 minutes", "1 hour", "2 hours"}
	for _, v := range allowed {
		if v == input {
			return true
		}
	}
	return false
}

// Helper function to check for duplicate store name
//func (s *StoreServiceServer) isDuplicateStore(storeName string) bool {
//	query := `SELECT COUNT(1) FROM stores WHERE store_name = $1`
//	var count int
//	err := s.DB.QueryRowContext(context.Background(), query, storeName).Scan(&count)
//	if err != nil {
//		return false // Assume no duplicate if query fails
//	}
//	return count > 0
//}

// Helper function to verify shop_manager_id
func (s *StoreServiceServer) verifyManager(managerID string) (valid bool, locked bool) {
	// Placeholder logic for manager verification
	// In production, implement logic to query the manager's account status
	if managerID == "locked" {
		return false, true
	}
	if managerID == "invalid" {
		return false, false
	}
	return true, false
}

// Helper function to build error responses
func (s *StoreServiceServer) buildErrorResponse(statusCode int, message string, _ map[string]string) (*proto.CreateStoreResponse, error) {
	response := &proto.CreateStoreResponse{
		Status:  "error",
		Message: message,
	}
	return response, fmt.Errorf("status_code: %d, message: %s", statusCode, message)
}
