package usecase

import (
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/jung-kurt/gofpdf"
)

func OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderBody.UserID = userID
	cartExist, err := repository.DoesCartExist(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := repository.AddressExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}
	PaymentExist, err := repository.PaymentExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !PaymentExist {
		return domain.OrderSuccessResponse{}, errors.New("paymentmethod does not exist")
	}
	cartItems, err := repository.GetAllItemsFromCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	total, err := repository.TotalAmountInCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	discount_price, err := repository.GetCouponDiscountPrice(int(orderBody.UserID), total)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	err = repository.UpdateCouponDetails(discount_price, orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	FinalPrice := total - discount_price
	if orderBody.PaymentID == 3 {
		wallectAmount, err := repository.WallectAmount(userID)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
		if FinalPrice >= wallectAmount {
			return domain.OrderSuccessResponse{}, errors.New("this much of amount not available in wallet")
		}
	}

	order_id, err := repository.OrderItems(orderBody, FinalPrice)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if orderBody.PaymentID == 3 {
		if err := repository.UpdateWallectAfterOrder(userID, FinalPrice); err != nil {
			return domain.OrderSuccessResponse{}, err
		}
		reason := "Amount debited for purchasing products"
		err = repository.UpdateHistoryForDebit(userID, order_id, FinalPrice, reason)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
	}
	if err := repository.AddOrderProducts(order_id, cartItems); err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderSuccessResponse, err := repository.GetBriefOrderDetails(order_id, orderBody.PaymentID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	// Send Order Confirmation Email (asynchronous)
	go func() {
		orderDetails, err := repository.GetDetailedOrderThroughId(order_id)
		if err == nil {
			_ = helper.SendOrderConfirmationEmail(orderDetails.Email, orderDetails.Firstname, order_id, FinalPrice)
		}
	}()

	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := repository.UpdateCartAfterOrder(userID, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
		// Decrement stock when order is placed
		if err := repository.DecreaseStock(int(c.ProductID), c.Size, int(c.Quantity)); err != nil {
			return domain.OrderSuccessResponse{}, err
		}
	}
	return orderSuccessResponse, nil
}

func GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := repository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

func CancelOrders(orderID int, userId int) error {
	userTest, err := repository.UserOrderRelationship(orderID, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}
	orderProductDetails, err := repository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := repository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item already delivered, cannot cancel")
	}

	if shipmentStatus != "order placed" {
		return errors.New("order cannot be cancelled once it is confirmed or shipped by the admin")
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled")
	}
	err = repository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	payment_status, err := repository.PaymentStatus(orderID)
	if err != nil {
		return err
	}
	err = repository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}
	amount, err := repository.TotalAmountFromOrder(orderID)
	if err != nil {
		return err
	}
	if payment_status == "refunded" {
		err = repository.UpdateAmountToWallet(userId, amount)
		if err != nil {
			return err
		}
		reason := "Amount credited for cancellation of order by user"
		err := repository.UpdateHistory(userId, orderID, amount, reason)
		if err != nil {
			return err
		}
	}
	return nil
}

func Checkout(userID int) (models.CheckoutDetails, error) {
	allUserAddress, err := repository.GetAllAddresses(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	paymentDetails, err := repository.GetAllPaymentOption(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	cartItems, err := repository.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	grandTotal, err := repository.GetTotalPrice(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	return models.CheckoutDetails{
		AddressInfoResponse: allUserAddress,
		Payment_Method:      paymentDetails,
		Cart:                cartItems,
		Total_Price:         grandTotal.FinalPrice,
	}, nil
}

func PaymentMethodID(order_id int) (int, error) {
	id, err := repository.PaymentMethodID(order_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetPaymentMethodName(orderID int) (string, error) {
	name, err := repository.GetPaymentMethodName(orderID)
	if err != nil {
		return "", err
	}
	return name, nil
}

func ExecutePurchaseCOD(orderID int) error {
	err := repository.OrderExist(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := repository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item  delivered, cannot pay")
	}
	if shipmentStatus == "order placed" {
		return errors.New("item placed, cannot pay")
	}
	if shipmentStatus == "cancelled" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + "so can't paid")
	}
	if shipmentStatus == "processing" {
		return errors.New("the order is already paid")
	}
	err = repository.UpdateOrder(orderID)
	if err != nil {
		return err
	}
	return nil
}

func PrintInvoice(orderId int) (*gofpdf.Fpdf, error) {

	if orderId < 1 {
		return nil, errors.New("enter a valid order id")
	}

	order, err := repository.GetDetailedOrderThroughId(orderId)
	if err != nil {
		return nil, err
	}

	items, err := repository.GetItemsByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	fmt.Println("order details ", order)
	fmt.Println("itemssss", items)
	if order.OrderId == "" {
		return nil, errors.New("order not found or invalid details")
	}

	// Allowing most statuses for invoice printing to ensure user accessibility
	restrictedStatuses := map[string]bool{
		"cancelled": true,
		"returned":  true,
	}

	if restrictedStatuses[order.ShipmentStatus] {
		return nil, fmt.Errorf("invoice is not available for %s orders", order.ShipmentStatus)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Colors
	brandGreenR, brandGreenG, brandGreenB := 27, 67, 50
	textDarkR, textDarkG, textDarkB := 26, 26, 26
	textLightR, textLightG, textLightB := 120, 120, 120
	lineGrayR, lineGrayG, lineGrayB := 230, 230, 230

	// Brand Header
	pdf.SetFont("Arial", "B", 32)
	pdf.SetTextColor(brandGreenR, brandGreenG, brandGreenB)
	pdf.Text(10, 25, "SOXENLY")

	// Invoice Label
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(textLightR, textLightG, textLightB)
	pdf.Text(150, 25, "TAX INVOICE")

	// Company Info
	pdf.SetY(32)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(textLightR, textLightG, textLightB)
	pdf.Cell(0, 5, "123 Soxenly Plaza, Urban Sector, Kerala, India")
	pdf.Ln(5)
	pdf.Cell(0, 5, "GSTIN: 32AAAAA0000A1Z5 | Support: hello@soxenly.com")
	pdf.Ln(12)

	pdf.SetDrawColor(lineGrayR, lineGrayG, lineGrayB)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(10)

	// Bill To & Order Info Grid
	yBeforeGrid := pdf.GetY()

	// Left Column: Bill To
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(brandGreenR, brandGreenG, brandGreenB)
	pdf.Text(10, yBeforeGrid+5, "BILL TO")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(textDarkR, textDarkG, textDarkB)
	pdf.Text(10, yBeforeGrid+12, order.Firstname)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(textLightR, textLightG, textLightB)
	pdf.Text(10, yBeforeGrid+17, order.HouseName)
	pdf.Text(10, yBeforeGrid+22, order.Street+", "+order.City)
	pdf.Text(10, yBeforeGrid+27, order.State+" - "+order.Pin)
	pdf.Text(10, yBeforeGrid+32, "Phone: "+order.Phone)

	// Right Column: Order Info
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(brandGreenR, brandGreenG, brandGreenB)
	pdf.Text(120, yBeforeGrid+5, "ORDER DETAILS")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(textDarkR, textDarkG, textDarkB)
	pdf.Text(120, yBeforeGrid+12, "Order ID:")
	pdf.SetFont("Arial", "B", 9)
	pdf.Text(150, yBeforeGrid+12, "#"+order.OrderId)

	pdf.SetFont("Arial", "", 9)
	pdf.Text(120, yBeforeGrid+17, "Date:")
	pdf.Text(150, yBeforeGrid+17, time.Now().Format("Jan 02, 2006"))

	pdf.Text(120, yBeforeGrid+22, "Payment:")
	pdf.Text(150, yBeforeGrid+22, order.PaymentStatus)

	pdf.Text(120, yBeforeGrid+27, "Shipment:")
	pdf.Text(150, yBeforeGrid+27, order.ShipmentStatus)

	pdf.SetY(yBeforeGrid + 45)

	// Table Header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(brandGreenR, brandGreenG, brandGreenB)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(90, 10, "  PRODUCT DESCRIPTION", "0", 0, "L", true, 0, "")
	pdf.CellFormat(30, 10, "UNIT PRICE", "0", 0, "C", true, 0, "")
	pdf.CellFormat(20, 10, "QTY", "0", 0, "C", true, 0, "")
	pdf.CellFormat(50, 10, "TOTAL", "0", 0, "C", true, 0, "")
	pdf.Ln(10)

	// Table Body
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(textDarkR, textDarkG, textDarkB)
	for i, item := range items {
		// Alternate background color for rows
		if i%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		name := item.ProductName
		if name == "" {
			name = "Product ID: " + strconv.Itoa(int(item.ProductID))
		}

		pdf.CellFormat(90, 10, "  "+name, "0", 0, "L", true, 0, "")
		pdf.CellFormat(30, 10, "INR "+strconv.FormatFloat(item.TotalPrice/item.Quantity, 'f', 2, 64), "0", 0, "C", true, 0, "")
		pdf.CellFormat(20, 10, strconv.Itoa(int(item.Quantity)), "0", 0, "C", true, 0, "")
		pdf.CellFormat(50, 10, "INR "+strconv.FormatFloat(item.TotalPrice, 'f', 2, 64), "0", 0, "C", true, 0, "")
		pdf.Ln(10)
	}

	pdf.Ln(5)

	// Summary Section
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(textLightR, textLightG, textLightB)
	pdf.CellFormat(140, 8, "Subtotal", "0", 0, "R", false, 0, "")

	var subTotal float64
	for _, item := range items {
		subTotal += item.TotalPrice
	}
	pdf.SetTextColor(textDarkR, textDarkG, textDarkB)
	pdf.CellFormat(50, 8, "INR "+strconv.FormatFloat(subTotal, 'f', 2, 64), "0", 0, "C", false, 0, "")
	pdf.Ln(8)

	pdf.SetTextColor(textLightR, textLightG, textLightB)
	pdf.CellFormat(140, 8, "Discount", "0", 0, "R", false, 0, "")
	offerApplied := subTotal - order.FinalPrice
	pdf.SetTextColor(200, 0, 0)
	pdf.CellFormat(50, 8, "- INR "+strconv.FormatFloat(offerApplied, 'f', 2, 64), "0", 0, "C", false, 0, "")
	pdf.Ln(12)

	// Final Amount
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(248, 249, 245)
	pdf.SetTextColor(brandGreenR, brandGreenG, brandGreenB)
	pdf.CellFormat(140, 12, "TOTAL AMOUNT PAYABLE", "0", 0, "R", true, 0, "")
	pdf.CellFormat(50, 12, "INR "+strconv.FormatFloat(order.FinalPrice, 'f', 2, 64), "0", 0, "C", true, 0, "")
	pdf.Ln(30)

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(textLightR, textLightG, textLightB)
	pdf.Cell(0, 5, "This is a computer generated invoice and does not require a physical signature.")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(0, 5, "Thank you for choosing Soxenly - Walk with Nature.")
	return pdf, nil
}
